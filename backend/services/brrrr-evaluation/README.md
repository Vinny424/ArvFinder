# BRRRR Calculation Evaluation

This folder contains the core BRRRR (Buy, Rehab, Rent, Refinance, Repeat) calculation logic and comprehensive tests for evaluating AI-generated improvements.

## Contents

- `arv.go` - Current BRRRR calculation service implementation
- `arv_test.go` - Comprehensive unit tests including edge cases
- `go.mod` - Go module dependencies
- `Dockerfile` - Container setup for testing
- `README.md` - This file

## Current Issues with BRRRR Calculations

The existing implementation has several limitations:

1. **Oversimplified refinancing logic** - Uses fixed 75% rule without considering actual loan requirements
2. **Missing rental income analysis** - No NOI (Net Operating Income) calculations
3. **Poor missing data handling** - Defaults to zero instead of graceful error handling
4. **Static assumptions** - Uses fixed percentages instead of dynamic market-based calculations
5. **No cash flow analysis** - Doesn't account for monthly rent, vacancy, or operating expenses

## Test Coverage

The test suite covers:

- Basic ARV calculations
- 70% rule validation
- ROI and profit margin calculations
- Comparable property analysis
- Risk assessment logic
- Edge cases (high rehab costs, break-even scenarios)
- Missing data scenarios

## Expected Improvements

An ideal refactored implementation should:

1. **Add rental income parameters**: `MonthlyRent`, `VacancyRate`, `AnnualExpenses`
2. **Calculate NOI properly**: Effective Gross Income - Operating Expenses
3. **Implement realistic refinancing**: 75% LTV with debt service coverage ratio (DSCR) validation
4. **Handle missing inputs gracefully**: Warning system instead of silent failures
5. **Time-based calculations**: Annualized returns instead of simple profit formulas

## Docker Usage

```bash
# Build and test
docker build -t brrrr-eval .

# Run tests in container
docker run --rm brrrr-eval

# Or test locally
go test -v
```

## Evaluation Criteria

When evaluating AI-generated improvements, check for:

- [ ] Proper NOI calculation implementation
- [ ] Realistic refinancing constraints (LTV, DSCR)
- [ ] Graceful handling of missing rental data
- [ ] Warning systems for incomplete analysis
- [ ] Time-based ROI calculations
- [ ] Debt service coverage analysis
- [ ] Cash flow projections
- [ ] Risk assessment improvements

## Prompt Used

> The BRRRR-related calculations in this Go service are too simple and not accurate for real-world use. The CalculateARV function uses fixed percentages like the 70% rule and refinance logic, but it doesn't consider important factors like monthly rent, vacancy rate, annual expenses, or refinance limits based on loan-to-value ratios. I want you to rewrite the function so that it calculates return on investment and profit margin using actual income and expenses over time, not just a basic profit formula. The refinance amount should be based on 75% of the ARV. Use rental income minus expenses to calculate net operating income. If any inputs are missing, don't just default them to zero—handle them in a way that avoids errors or misleading results. The goal is to make the BRRRR analysis reflect how investors actually make decisions, not just rough estimates.