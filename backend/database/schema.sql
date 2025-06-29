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

-- Create users table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    role VARCHAR(50) NOT NULL DEFAULT 'user',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
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

-- Create indexes for performance
CREATE INDEX idx_users_tenant_id ON users(tenant_id);
CREATE INDEX idx_users_email ON users(email);
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

-- Insert sample data for development
INSERT INTO tenants (id, name, subscription_tier) VALUES 
    ('00000000-0000-0000-0000-000000000001', 'Demo Tenant', 'professional');

INSERT INTO users (id, tenant_id, email, password_hash, first_name, last_name, role) VALUES 
    ('00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001', 'demo@arvfinder.com', '$2a$10$dummy.hash.for.demo.user', 'Demo', 'User', 'admin');

INSERT INTO properties (tenant_id, address, city, state, zip_code, price, arv, bedrooms, bathrooms, square_feet, property_type) VALUES 
    ('00000000-0000-0000-0000-000000000001', '123 Main St', 'Denver', 'CO', '80202', 180000, 250000, 3, 2, 1200, 'Single Family'),
    ('00000000-0000-0000-0000-000000000001', '456 Oak Ave', 'Boulder', 'CO', '80301', 220000, 300000, 4, 3, 1800, 'Single Family');