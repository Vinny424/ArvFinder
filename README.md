# ArvFinder

A comprehensive rental property scouting web application that specializes in ARV (After Repair Value) calculations and investment analysis.

## Features

- **ARV Calculator**: Calculate accurate After Repair Values using the 70% rule
- **Property Portfolio**: Manage and track multiple investment properties
- **Comparable Analysis**: Find and analyze comparable properties
- **Investment Metrics**: ROI, cash-on-cash return, and profit margin calculations
- **Subscription Management**: Stripe-powered billing with three-tier pricing
- **Multi-tenant Architecture**: Support for multiple users and organizations
- **PDF Reports**: Generate professional property analysis reports

## Tech Stack

- **Frontend**: SvelteKit with TypeScript and Tailwind CSS
- **Backend**: Go with Gin framework
- **Database**: PostgreSQL
- **Authentication**: JWT tokens
- **Containerization**: Docker and Docker Compose

## Getting Started

### Prerequisites

- Docker and Docker Compose
- Node.js 18+ (for local development)
- Go 1.21+ (for local development)

### Quick Start with Docker

1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd ArvFinder
   ```

2. Start the application:
   ```bash
   docker-compose up -d
   ```

3. Access the application:
   - Frontend: http://localhost:5173
   - Backend API: http://localhost:8080
   - Database: localhost:5432

### Local Development

#### Backend

1. Install Go dependencies:
   ```bash
   cd backend
   go mod tidy
   ```

2. Set up environment variables:
   ```bash
   export DATABASE_URL="postgres://arvfinder:arvfinder_dev@localhost:5432/arvfinder?sslmode=disable"
   export JWT_SECRET="your-super-secret-jwt-key"
   ```

3. Run the backend:
   ```bash
   go run main.go
   ```

#### Frontend

1. Install dependencies:
   ```bash
   cd frontend/frontend
   npm install
   ```

2. Set up environment variables:
   ```bash
   echo "VITE_API_URL=http://localhost:8080/api/v1" > .env
   ```

3. Run the frontend:
   ```bash
   npm run dev
   ```

## API Endpoints

### Authentication
- `POST /api/v1/auth/login` - User login
- `POST /api/v1/auth/register` - User registration
- `POST /api/v1/auth/refresh` - Refresh JWT token

### Properties
- `GET /api/v1/properties` - Get all properties
- `POST /api/v1/properties` - Create new property
- `GET /api/v1/properties/:id` - Get property by ID
- `PUT /api/v1/properties/:id` - Update property
- `DELETE /api/v1/properties/:id` - Delete property

### ARV Calculations
- `POST /api/v1/arv/calculate` - Calculate ARV
- `POST /api/v1/arv/70-rule` - Calculate 70% rule
- `POST /api/v1/arv/roi` - Calculate ROI
- `POST /api/v1/arv/cash-on-cash` - Calculate cash-on-cash return
- `POST /api/v1/arv/cap-rate` - Calculate cap rate
- `POST /api/v1/arv/estimate-from-comps` - Estimate ARV from comparables

### Stripe Payments
- `GET /api/v1/payments/plans` - Get subscription plans
- `POST /api/v1/payments/create-subscription` - Create subscription
- `POST /api/v1/payments/cancel-subscription` - Cancel subscription
- `POST /api/v1/payments/webhook` - Handle Stripe webhooks

## Database Schema

The application uses a multi-tenant PostgreSQL database with the following main tables:

- `tenants` - Organization/tenant information
- `users` - User accounts linked to tenants
- `properties` - Property information and basic metrics
- `arv_calculations` - ARV calculation results
- `comparables` - Comparable property data

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

[License details here]

## Support

For support, please contact [support contact information]