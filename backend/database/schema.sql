-- ArvFinder Database Schema
-- Multi-tenant architecture with shared database, separate schemas

-- Create tenants table
CREATE TABLE tenants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    subscription_tier VARCHAR(50) NOT NULL DEFAULT 'starter',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create users table with advanced security features
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    email VARCHAR(255) UNIQUE NOT NULL,
    email_verified BOOLEAN NOT NULL DEFAULT FALSE,
    email_verification_token VARCHAR(255),
    email_verification_expires_at TIMESTAMP WITH TIME ZONE,
    password_hash VARCHAR(512) NOT NULL, -- Increased for Argon2
    password_salt VARCHAR(128) NOT NULL, -- Dedicated salt field
    password_reset_token VARCHAR(255),
    password_reset_expires_at TIMESTAMP WITH TIME ZONE,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    phone_number VARCHAR(20),
    phone_verified BOOLEAN NOT NULL DEFAULT FALSE,
    role VARCHAR(50) NOT NULL DEFAULT 'user',
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    two_factor_enabled BOOLEAN NOT NULL DEFAULT FALSE,
    two_factor_secret VARCHAR(255), -- For TOTP
    backup_codes TEXT[], -- Array of backup codes
    last_login_at TIMESTAMP WITH TIME ZONE,
    last_login_ip INET,
    failed_login_attempts INTEGER NOT NULL DEFAULT 0,
    locked_until TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create user sessions table for JWT token management
CREATE TABLE user_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    refresh_token VARCHAR(512) NOT NULL UNIQUE,
    refresh_token_hash VARCHAR(512) NOT NULL, -- Hashed version for security
    access_token_jti VARCHAR(255) NOT NULL, -- JWT ID for access token
    device_fingerprint VARCHAR(255),
    user_agent TEXT,
    ip_address INET,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    revoked BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create SMS 2FA verification codes table
CREATE TABLE sms_verification_codes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    phone_number VARCHAR(20) NOT NULL,
    code VARCHAR(10) NOT NULL,
    code_hash VARCHAR(255) NOT NULL, -- Hashed version of code
    purpose VARCHAR(50) NOT NULL, -- 'login', 'register', 'password_reset'
    attempts INTEGER NOT NULL DEFAULT 0,
    max_attempts INTEGER NOT NULL DEFAULT 3,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    verified BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create security audit log table
CREATE TABLE security_audit_log (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    event_type VARCHAR(100) NOT NULL, -- 'login', 'logout', 'failed_login', 'password_change', etc.
    event_description TEXT,
    ip_address INET,
    user_agent TEXT,
    additional_data JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create rate limiting table
CREATE TABLE rate_limits (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    identifier VARCHAR(255) NOT NULL, -- IP address or user ID
    action VARCHAR(100) NOT NULL, -- 'login', 'register', 'password_reset'
    attempts INTEGER NOT NULL DEFAULT 1,
    window_start TIMESTAMP WITH TIME ZONE NOT NULL,
    blocked_until TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(identifier, action)
);

-- Create properties table
CREATE TABLE properties (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    address VARCHAR(500) NOT NULL,
    city VARCHAR(100),
    state VARCHAR(50),
    zip_code VARCHAR(20),
    price DECIMAL(12,2),
    arv DECIMAL(12,2),
    rehab_cost DECIMAL(12,2) DEFAULT 0,
    holding_costs DECIMAL(12,2) DEFAULT 0,
    closing_costs DECIMAL(12,2) DEFAULT 0,
    bedrooms INTEGER,
    bathrooms DECIMAL(3,1),
    square_feet INTEGER,
    lot_size DECIMAL(10,2),
    year_built INTEGER,
    property_type VARCHAR(100),
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create arv_calculations table
CREATE TABLE arv_calculations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    property_id UUID REFERENCES properties(id) ON DELETE CASCADE,
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    purchase_price DECIMAL(12,2) NOT NULL,
    rehab_cost DECIMAL(12,2) DEFAULT 0,
    holding_costs DECIMAL(12,2) DEFAULT 0,
    closing_costs DECIMAL(12,2) DEFAULT 0,
    arv DECIMAL(12,2) NOT NULL,
    max_offer DECIMAL(12,2) GENERATED ALWAYS AS (arv * 0.7 - rehab_cost) STORED,
    potential_profit DECIMAL(12,2) GENERATED ALWAYS AS (arv - purchase_price - rehab_cost - holding_costs - closing_costs) STORED,
    profit_margin DECIMAL(5,2) GENERATED ALWAYS AS (
        CASE 
            WHEN (purchase_price + rehab_cost + holding_costs + closing_costs) > 0 
            THEN ((arv - purchase_price - rehab_cost - holding_costs - closing_costs) / (purchase_price + rehab_cost + holding_costs + closing_costs) * 100)
            ELSE 0 
        END
    ) STORED,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create comparables table
CREATE TABLE comparables (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    property_id UUID REFERENCES properties(id) ON DELETE CASCADE,
    address VARCHAR(500) NOT NULL,
    sale_price DECIMAL(12,2) NOT NULL,
    sale_date DATE NOT NULL,
    distance DECIMAL(8,2), -- in miles
    bedrooms INTEGER,
    bathrooms DECIMAL(3,1),
    square_feet INTEGER,
    price_per_sq_ft DECIMAL(8,2) GENERATED ALWAYS AS (
        CASE WHEN square_feet > 0 THEN sale_price / square_feet ELSE 0 END
    ) STORED,
    adjustments DECIMAL(12,2) DEFAULT 0,
    adjusted_value DECIMAL(12,2) GENERATED ALWAYS AS (sale_price + adjustments) STORED,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create indexes for performance and security
CREATE INDEX idx_users_tenant_id ON users(tenant_id);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_email_verification_token ON users(email_verification_token);
CREATE INDEX idx_users_password_reset_token ON users(password_reset_token);
CREATE INDEX idx_users_phone_number ON users(phone_number);
CREATE INDEX idx_users_last_login_at ON users(last_login_at);
CREATE INDEX idx_users_locked_until ON users(locked_until);

-- Session management indexes
CREATE INDEX idx_user_sessions_user_id ON user_sessions(user_id);
CREATE INDEX idx_user_sessions_refresh_token ON user_sessions(refresh_token);
CREATE INDEX idx_user_sessions_access_token_jti ON user_sessions(access_token_jti);
CREATE INDEX idx_user_sessions_expires_at ON user_sessions(expires_at);
CREATE INDEX idx_user_sessions_revoked ON user_sessions(revoked);

-- SMS 2FA indexes
CREATE INDEX idx_sms_verification_codes_user_id ON sms_verification_codes(user_id);
CREATE INDEX idx_sms_verification_codes_phone ON sms_verification_codes(phone_number);
CREATE INDEX idx_sms_verification_codes_expires_at ON sms_verification_codes(expires_at);
CREATE INDEX idx_sms_verification_codes_purpose ON sms_verification_codes(purpose);

-- Security audit indexes
CREATE INDEX idx_security_audit_log_user_id ON security_audit_log(user_id);
CREATE INDEX idx_security_audit_log_event_type ON security_audit_log(event_type);
CREATE INDEX idx_security_audit_log_created_at ON security_audit_log(created_at);
CREATE INDEX idx_security_audit_log_ip_address ON security_audit_log(ip_address);

-- Rate limiting indexes
CREATE INDEX idx_rate_limits_identifier_action ON rate_limits(identifier, action);
CREATE INDEX idx_rate_limits_window_start ON rate_limits(window_start);
CREATE INDEX idx_rate_limits_blocked_until ON rate_limits(blocked_until);

-- Property indexes
CREATE INDEX idx_properties_tenant_id ON properties(tenant_id);
CREATE INDEX idx_properties_address ON properties(address);
CREATE INDEX idx_arv_calculations_tenant_id ON arv_calculations(tenant_id);
CREATE INDEX idx_arv_calculations_property_id ON arv_calculations(property_id);
CREATE INDEX idx_comparables_property_id ON comparables(property_id);
CREATE INDEX idx_comparables_sale_date ON comparables(sale_date);

-- Create updated_at trigger function
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create triggers for updated_at
CREATE TRIGGER update_tenants_updated_at BEFORE UPDATE ON tenants
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_properties_updated_at BEFORE UPDATE ON properties
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_rate_limits_updated_at BEFORE UPDATE ON rate_limits
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Security constraints and checks
ALTER TABLE users ADD CONSTRAINT check_email_format 
    CHECK (email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$');

ALTER TABLE users ADD CONSTRAINT check_phone_format 
    CHECK (phone_number IS NULL OR phone_number ~ '^\+?[1-9]\d{1,14}$');

ALTER TABLE users ADD CONSTRAINT check_failed_login_attempts 
    CHECK (failed_login_attempts >= 0 AND failed_login_attempts <= 100);

ALTER TABLE sms_verification_codes ADD CONSTRAINT check_code_length 
    CHECK (LENGTH(code) >= 4 AND LENGTH(code) <= 10);

ALTER TABLE sms_verification_codes ADD CONSTRAINT check_attempts_range 
    CHECK (attempts >= 0 AND attempts <= max_attempts);

ALTER TABLE rate_limits ADD CONSTRAINT check_attempts_positive 
    CHECK (attempts > 0);

-- Create function to clean up expired records
CREATE OR REPLACE FUNCTION cleanup_expired_records()
RETURNS void AS $$
BEGIN
    -- Clean up expired SMS verification codes
    DELETE FROM sms_verification_codes WHERE expires_at < NOW() - INTERVAL '1 day';
    
    -- Clean up expired user sessions
    DELETE FROM user_sessions WHERE expires_at < NOW();
    
    -- Clean up old audit logs (keep for 1 year)
    DELETE FROM security_audit_log WHERE created_at < NOW() - INTERVAL '1 year';
    
    -- Clean up old rate limit records
    DELETE FROM rate_limits WHERE window_start < NOW() - INTERVAL '1 day' AND blocked_until < NOW();
END;
$$ LANGUAGE plpgsql;

-- Create function to automatically lock users after failed attempts
CREATE OR REPLACE FUNCTION check_and_lock_user()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.failed_login_attempts >= 5 AND (OLD.failed_login_attempts IS NULL OR OLD.failed_login_attempts < 5) THEN
        NEW.locked_until = NOW() + INTERVAL '30 minutes';
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_check_and_lock_user BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION check_and_lock_user();

-- Insert sample data for development
INSERT INTO tenants (id, name, subscription_tier) VALUES 
    ('00000000-0000-0000-0000-000000000001', 'Demo Tenant', 'professional');

INSERT INTO users (id, tenant_id, email, password_hash, password_salt, first_name, last_name, role, email_verified) VALUES 
    ('00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001', 'demo@arvfinder.com', '$2a$10$dummy.hash.for.demo.user', 'dummy_salt_for_demo', 'Demo', 'User', 'admin', TRUE);

INSERT INTO properties (tenant_id, address, city, state, zip_code, price, arv, bedrooms, bathrooms, square_feet, property_type) VALUES 
    ('00000000-0000-0000-0000-000000000001', '123 Main St', 'Denver', 'CO', '80202', 180000, 250000, 3, 2, 1200, 'Single Family'),
    ('00000000-0000-0000-0000-000000000001', '456 Oak Ave', 'Boulder', 'CO', '80301', 220000, 300000, 4, 3, 1800, 'Single Family');