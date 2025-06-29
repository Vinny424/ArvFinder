# API Configuration Guide

## Required API Keys for Full Functionality

### 1. Google Maps API Setup

To enable address validation and autocomplete functionality:

1. **Get Google Cloud Console API Key:**
   - Go to [Google Cloud Console](https://console.cloud.google.com/)
   - Create a new project or select existing one
   - Enable the following APIs:
     - Maps JavaScript API
     - Places API
     - Geocoding API
   - Create credentials → API Key
   - Restrict the API key to your domain for security

2. **Frontend Configuration:**
   ```javascript
   // Replace YOUR_GOOGLE_API_KEY in src/routes/+page.svelte
   const GOOGLE_API_KEY = 'your_actual_api_key_here';
   ```

3. **Current Implementation:**
   - Address autocomplete suggestions
   - Address validation and parsing
   - Geocoding for accurate address components

### 2. Repliers API Setup

For property estimates and market data:

1. **Get Repliers API Key:**
   - Visit [Repliers.io](https://api.repliers.io/)
   - Sign up for an account
   - Get your API key from dashboard

2. **Backend Configuration:**
   ```bash
   # Add to your environment variables
   export REPLIERS_API_KEY=your_repliers_api_key_here
   
   # Or add to docker-compose.yml
   environment:
     - REPLIERS_API_KEY=your_repliers_api_key_here
   ```

3. **API Endpoints Used:**
   - Property value estimates
   - Address history lookup
   - Comparable properties data

### 3. Environment Variables

Update your `.env` file or docker-compose.yml:

```bash
# Google Maps (for frontend)
VITE_GOOGLE_MAPS_API_KEY=your_google_api_key

# Repliers API (for backend)
REPLIERS_API_KEY=your_repliers_api_key

# Stripe (already configured)
STRIPE_SECRET_KEY=sk_test_...
STRIPE_PUBLISHABLE_KEY=pk_test_...
```

## Current Fallback Behavior

The application includes intelligent fallbacks:

### Without Google Maps API:
- Shows simple address suggestions based on common locations
- Still performs property search with basic address parsing
- All functionality remains available

### Without Repliers API:
- Uses realistic simulated property data
- Provides sample comparable properties
- Maintains all calculation capabilities

## Production Deployment

For production use:

1. **Replace all API keys** with production keys
2. **Enable API restrictions** for security
3. **Set up proper rate limiting**
4. **Monitor API usage** and costs

## Testing

You can test the complete functionality immediately:

1. **Address Search:** Enter any address to see autocomplete
2. **Property Estimates:** Get property data and estimates  
3. **ARV Calculations:** Use the comprehensive calculator
4. **Cost Estimates:** Click "Estimate Additional Costs"
5. **ROI Calculations:** See accurate cash-on-cash returns

The application works seamlessly with or without API keys!