package services

import (
	"math"
)

// ArvRequest represents the input data for ARV calculation
type ArvRequest struct {
	PurchasePrice    float64 `json:"purchase_price" binding:"required,min=1"`
	RehabCost        float64 `json:"rehab_cost" binding:"min=0"`
	HoldingCosts     float64 `json:"holding_costs" binding:"min=0"`
	ClosingCosts     float64 `json:"closing_costs" binding:"min=0"`
	ARV              float64 `json:"arv" binding:"required,min=1"`
	FinancingCosts   float64 `json:"financing_costs" binding:"min=0"`
	SellingCosts     float64 `json:"selling_costs" binding:"min=0"`

	// BRRRR-specific fields
	MonthlyRent      float64 `json:"monthly_rent" binding:"min=0"`
	VacancyRate      float64 `json:"vacancy_rate" binding:"min=0,max=100"` // percentage
	PropertyTaxes    float64 `json:"property_taxes" binding:"min=0"`       // annual
	Insurance        float64 `json:"insurance" binding:"min=0"`            // annual
	Maintenance      float64 `json:"maintenance" binding:"min=0"`          // annual
	PropertyMgmt     float64 `json:"property_mgmt" binding:"min=0"`        // annual or percentage
	CapEx            float64 `json:"capex" binding:"min=0"`                // annual capital expenditures
	OtherExpenses    float64 `json:"other_expenses" binding:"min=0"`       // annual
	RefinanceLTV     float64 `json:"refinance_ltv" binding:"min=0,max=100"` // percentage, default 75%
	InterestRate     float64 `json:"interest_rate" binding:"min=0,max=30"`  // percentage for refinance loan
	LoanTerm         int     `json:"loan_term" binding:"min=1,max=50"`      // years, default 30
}

// ArvResult represents the calculated ARV analysis results
type ArvResult struct {
	// Input values
	PurchasePrice    float64 `json:"purchase_price"`
	RehabCost        float64 `json:"rehab_cost"`
	HoldingCosts     float64 `json:"holding_costs"`
	ClosingCosts     float64 `json:"closing_costs"`
	ARV              float64 `json:"arv"`
	FinancingCosts   float64 `json:"financing_costs"`
	SellingCosts     float64 `json:"selling_costs"`

	// Income Analysis
	MonthlyRent      float64 `json:"monthly_rent"`
	AnnualGrossIncome float64 `json:"annual_gross_income"`
	EffectiveIncome  float64 `json:"effective_income"` // after vacancy

	// Expense Analysis
	AnnualExpenses   float64 `json:"annual_expenses"`
	ExpenseRatio     float64 `json:"expense_ratio"` // expenses as % of gross income

	// NOI and Cash Flow
	NOI              float64 `json:"noi"` // Net Operating Income

	// 70% Rule calculations (kept for comparison)
	MaxOffer70       float64 `json:"max_offer_70"`
	Is70RuleGood     bool    `json:"is_70_rule_good"`

	// Investment analysis
	TotalInvestment  float64 `json:"total_investment"`
	PotentialProfit  float64 `json:"potential_profit"`
	ProfitMargin     float64 `json:"profit_margin"`
	ROI              float64 `json:"roi"`

	// BRRRR Strategy metrics
	BrrrrMaxOffer    float64 `json:"brrrr_max_offer"`
	BrrrrProfit      float64 `json:"brrrr_profit"`
	RefinanceAmount  float64 `json:"refinance_amount"`  // 75% of ARV by default
	CashRecovered    float64 `json:"cash_recovered"`    // cash pulled out in refinance
	CashLeftIn       float64 `json:"cash_left_in"`      // remaining cash investment
	MonthlyDebtService float64 `json:"monthly_debt_service"` // P&I payment
	MonthlyCashFlow  float64 `json:"monthly_cash_flow"`
	AnnualCashFlow   float64 `json:"annual_cash_flow"`

	// Returns
	CashOnCashReturn float64 `json:"cash_on_cash_return"` // based on cash left in deal
	CapRate          float64 `json:"cap_rate"`            // NOI / ARV
	DSCR             float64 `json:"dscr"`                // Debt Service Coverage Ratio

	// Analysis flags
	IsInfiniteReturn bool    `json:"is_infinite_return"` // all cash recovered
	IsCashFlowPositive bool  `json:"is_cash_flow_positive"`

	// Risk assessment
	RiskLevel        string   `json:"risk_level"`
	Recommendations  []string `json:"recommendations"`

	// Validation warnings
	Warnings         []string `json:"warnings"`
}

// ArvService handles ARV calculations and analysis
type ArvService struct{}

// NewArvService creates a new ARV service instance
func NewArvService() *ArvService {
	return &ArvService{}
}

// CalculateARV performs comprehensive BRRRR analysis with income-based calculations
func (s *ArvService) CalculateARV(req ArvRequest) ArvResult {
	result := ArvResult{
		PurchasePrice:  req.PurchasePrice,
		RehabCost:      req.RehabCost,
		HoldingCosts:   req.HoldingCosts,
		ClosingCosts:   req.ClosingCosts,
		ARV:            req.ARV,
		FinancingCosts: req.FinancingCosts,
		SellingCosts:   req.SellingCosts,
		Warnings:       []string{},
	}

	// Set defaults and validate inputs
	s.setDefaultsAndValidate(&req, &result)

	// Calculate total investment
	result.TotalInvestment = req.PurchasePrice + req.RehabCost + req.HoldingCosts +
		req.ClosingCosts + req.FinancingCosts

	// Income calculations
	result.MonthlyRent = req.MonthlyRent
	result.AnnualGrossIncome = req.MonthlyRent * 12

	// Apply vacancy rate
	vacancyLoss := result.AnnualGrossIncome * (req.VacancyRate / 100)
	result.EffectiveIncome = result.AnnualGrossIncome - vacancyLoss

	// Calculate total annual expenses
	result.AnnualExpenses = req.PropertyTaxes + req.Insurance + req.Maintenance +
		req.CapEx + req.OtherExpenses

	// Add property management (could be flat fee or percentage)
	if req.PropertyMgmt > 0 {
		if req.PropertyMgmt < 1000 { // Assume percentage if less than $1000
			result.AnnualExpenses += result.AnnualGrossIncome * (req.PropertyMgmt / 100)
		} else { // Assume flat annual fee
			result.AnnualExpenses += req.PropertyMgmt
		}
	}

	// Calculate expense ratio
	if result.AnnualGrossIncome > 0 {
		result.ExpenseRatio = (result.AnnualExpenses / result.AnnualGrossIncome) * 100
	}

	// Calculate NOI
	result.NOI = result.EffectiveIncome - result.AnnualExpenses

	// BRRRR refinance calculations
	result.RefinanceAmount = req.ARV * (req.RefinanceLTV / 100)
	result.CashRecovered = math.Min(result.RefinanceAmount, result.TotalInvestment)
	result.CashLeftIn = math.Max(0, result.TotalInvestment - result.CashRecovered)

	// Calculate monthly debt service for refinance loan
	if req.InterestRate > 0 && req.LoanTerm > 0 {
		result.MonthlyDebtService = s.calculateMonthlyPayment(
			result.RefinanceAmount, req.InterestRate, req.LoanTerm)
	}

	// Calculate monthly and annual cash flow
	result.MonthlyCashFlow = (result.EffectiveIncome / 12) - (result.AnnualExpenses / 12) - result.MonthlyDebtService
	result.AnnualCashFlow = result.MonthlyCashFlow * 12

	// Calculate returns
	if result.CashLeftIn > 0 {
		result.CashOnCashReturn = (result.AnnualCashFlow / result.CashLeftIn) * 100
	} else if result.AnnualCashFlow > 0 {
		result.IsInfiniteReturn = true
		result.CashOnCashReturn = 999.99 // Represent infinite return
	}

	// Calculate cap rate
	if req.ARV > 0 {
		result.CapRate = (result.NOI / req.ARV) * 100
	}

	// Calculate DSCR (Debt Service Coverage Ratio)
	annualDebtService := result.MonthlyDebtService * 12
	if annualDebtService > 0 {
		result.DSCR = result.NOI / annualDebtService
	}

	// Set analysis flags
	result.IsCashFlowPositive = result.MonthlyCashFlow > 0

	// Keep 70% rule for comparison
	result.MaxOffer70 = (req.ARV * 0.70) - req.RehabCost
	result.Is70RuleGood = req.PurchasePrice <= result.MaxOffer70

	// Keep legacy calculations for backward compatibility
	result.PotentialProfit = req.ARV - result.TotalInvestment - req.SellingCosts
	if result.TotalInvestment > 0 {
		result.ProfitMargin = (result.PotentialProfit / result.TotalInvestment) * 100
		result.ROI = (result.PotentialProfit / result.TotalInvestment) * 100
	}
	result.BrrrrMaxOffer = (req.ARV * 0.75) - req.RehabCost
	result.BrrrrProfit = req.ARV - result.TotalInvestment - req.SellingCosts

	// Risk assessment and recommendations - use legacy for backward compatibility
	result.RiskLevel = s.assessRisk(result.ProfitMargin, result.Is70RuleGood, req.ARV, req.PurchasePrice)
	result.Recommendations = s.generateRecommendations(req, result.ProfitMargin, result.Is70RuleGood)

	// Round all financial values
	s.roundFinancialValues(&result)

	return result
}

// Calculate70Rule specifically calculates the 70% rule
func (s *ArvService) Calculate70Rule(arv, rehabCost float64) float64 {
	return (arv * 0.70) - rehabCost
}

// CalculateROI calculates return on investment
func (s *ArvService) CalculateROI(profit, investment float64) float64 {
	if investment <= 0 {
		return 0
	}
	return (profit / investment) * 100
}

// CalculateCashOnCashReturn calculates cash-on-cash return for rental properties
func (s *ArvService) CalculateCashOnCashReturn(annualCashFlow, totalCashInvested float64) float64 {
	if totalCashInvested <= 0 {
		return 0
	}
	return (annualCashFlow / totalCashInvested) * 100
}

// CalculateCapRate calculates capitalization rate
func (s *ArvService) CalculateCapRate(netOperatingIncome, propertyValue float64) float64 {
	if propertyValue <= 0 {
		return 0
	}
	return (netOperatingIncome / propertyValue) * 100
}

// assessRisk determines the risk level of the investment
func (s *ArvService) assessRisk(profitMargin float64, meets70Rule bool, arv, purchasePrice float64) string {
	// Calculate equity percentage
	equityPercent := ((arv - purchasePrice) / arv) * 100

	if profitMargin >= 20 && meets70Rule && equityPercent >= 25 {
		return "Low"
	} else if profitMargin >= 10 && (meets70Rule || equityPercent >= 15) {
		return "Medium"
	} else if profitMargin >= 5 {
		return "High"
	}
	return "Very High"
}

// generateRecommendations provides investment recommendations based on analysis
func (s *ArvService) generateRecommendations(req ArvRequest, profitMargin float64, meets70Rule bool) []string {
	var recommendations []string

	if !meets70Rule {
		recommendations = append(recommendations, "Property does not meet the 70% rule - consider negotiating a lower purchase price")
	}

	if profitMargin < 10 {
		recommendations = append(recommendations, "Low profit margin - consider reducing rehab costs or finding a lower purchase price")
	}

	if req.RehabCost > req.ARV*0.3 {
		recommendations = append(recommendations, "Rehab costs are high (>30% of ARV) - verify estimates with contractors")
	}

	if req.HoldingCosts > req.ARV*0.05 {
		recommendations = append(recommendations, "Holding costs seem high - consider faster renovation timeline")
	}

	if profitMargin >= 20 && meets70Rule {
		recommendations = append(recommendations, "Excellent investment opportunity with strong profit potential")
	}

	// Market-based recommendations
	equityPercent := ((req.ARV - req.PurchasePrice) / req.ARV) * 100
	if equityPercent >= 30 {
		recommendations = append(recommendations, "High equity position - good for BRRRR strategy")
	}

	if len(recommendations) == 0 {
		recommendations = append(recommendations, "Moderate investment opportunity - proceed with careful due diligence")
	}

	return recommendations
}

// ComparableProperty represents a comparable property for ARV estimation
type ComparableProperty struct {
	Address       string  `json:"address"`
	SalePrice     float64 `json:"sale_price"`
	SaleDate      string  `json:"sale_date"`
	Bedrooms      int     `json:"bedrooms"`
	Bathrooms     float64 `json:"bathrooms"`
	SquareFeet    int     `json:"square_feet"`
	Distance      float64 `json:"distance"`
	Adjustments   float64 `json:"adjustments"`
	AdjustedValue float64 `json:"adjusted_value"`
}

// EstimateARVFromComps estimates ARV based on comparable properties
func (s *ArvService) EstimateARVFromComps(comps []ComparableProperty, subjectBedrooms int, subjectBathrooms float64, subjectSquareFeet int) float64 {
	if len(comps) == 0 {
		return 0
	}

	totalAdjustedValue := 0.0
	weightedTotal := 0.0

	for _, comp := range comps {
		// Calculate adjustments based on differences
		adjustments := s.calculateComparableAdjustments(comp, subjectBedrooms, subjectBathrooms, float64(subjectSquareFeet))
		adjustedValue := comp.SalePrice + adjustments

		// Weight by distance (closer properties have more weight)
		weight := 1.0 / (1.0 + comp.Distance)

		totalAdjustedValue += adjustedValue * weight
		weightedTotal += weight
	}

	if weightedTotal == 0 {
		return 0
	}

	estimatedArv := totalAdjustedValue / weightedTotal
	return math.Round(estimatedArv*100) / 100
}

// calculateComparableAdjustments calculates adjustments for comparable properties
func (s *ArvService) calculateComparableAdjustments(comp ComparableProperty, subjectBeds int, subjectBaths, subjectSqFt float64) float64 {
	adjustments := 0.0

	// Bedroom adjustment (~$5,000 per bedroom difference)
	bedroomDiff := subjectBeds - comp.Bedrooms
	adjustments += float64(bedroomDiff) * 5000

	// Bathroom adjustment (~$3,000 per bathroom difference)
	bathroomDiff := subjectBaths - comp.Bathrooms
	adjustments += bathroomDiff * 3000

	// Square footage adjustment (~$50 per sq ft difference)
	sqFtDiff := subjectSqFt - float64(comp.SquareFeet)
	adjustments += sqFtDiff * 50

	return adjustments
}

// setDefaultsAndValidate sets reasonable defaults and validates inputs
func (s *ArvService) setDefaultsAndValidate(req *ArvRequest, result *ArvResult) {
	// Set default refinance LTV if not provided
	if req.RefinanceLTV == 0 {
		req.RefinanceLTV = 75.0
	}

	// Set default loan term if not provided
	if req.LoanTerm == 0 {
		req.LoanTerm = 30
	}

	// Set default vacancy rate if not provided (market average)
	if req.VacancyRate == 0 {
		req.VacancyRate = 8.0 // 8% is reasonable default
		result.Warnings = append(result.Warnings, "Using default vacancy rate of 8%")
	}

	// Estimate monthly rent if not provided (1% rule as fallback)
	if req.MonthlyRent == 0 {
		req.MonthlyRent = req.ARV * 0.01 // 1% rule
		result.Warnings = append(result.Warnings, "Monthly rent estimated using 1% rule - verify with market data")
	}

	// Estimate expenses if not provided
	if req.PropertyTaxes == 0 {
		req.PropertyTaxes = req.ARV * 0.015 // 1.5% of ARV annually
		result.Warnings = append(result.Warnings, "Property taxes estimated at 1.5% of ARV")
	}

	if req.Insurance == 0 {
		req.Insurance = req.ARV * 0.005 // 0.5% of ARV annually
		result.Warnings = append(result.Warnings, "Insurance estimated at 0.5% of ARV")
	}

	if req.Maintenance == 0 {
		req.Maintenance = req.MonthlyRent * 12 * 0.10 // 10% of gross rent
		result.Warnings = append(result.Warnings, "Maintenance estimated at 10% of gross rent")
	}

	if req.CapEx == 0 {
		req.CapEx = req.MonthlyRent * 12 * 0.05 // 5% of gross rent
		result.Warnings = append(result.Warnings, "CapEx estimated at 5% of gross rent")
	}

	// Set default interest rate if not provided
	if req.InterestRate == 0 {
		req.InterestRate = 7.0 // Current market rate
		result.Warnings = append(result.Warnings, "Using default interest rate of 7%")
	}

	// Validate critical inputs
	if req.ARV <= req.PurchasePrice + req.RehabCost {
		result.Warnings = append(result.Warnings, "WARNING: ARV may be too low compared to total acquisition costs")
	}

	if req.VacancyRate > 20 {
		result.Warnings = append(result.Warnings, "WARNING: Vacancy rate seems unusually high")
	}
}

// calculateMonthlyPayment calculates monthly P&I payment
func (s *ArvService) calculateMonthlyPayment(principal, annualRate float64, years int) float64 {
	if annualRate == 0 {
		return principal / float64(years*12)
	}

	monthlyRate := annualRate / 100 / 12
	numPayments := float64(years * 12)

	// Standard mortgage payment formula
	payment := principal * (monthlyRate * math.Pow(1+monthlyRate, numPayments)) /
		(math.Pow(1+monthlyRate, numPayments) - 1)

	return payment
}

// assessBRRRRisk provides more sophisticated risk assessment for BRRRR strategy
func (s *ArvService) assessBRRRRisk(result ArvResult) string {
	riskScore := 0

	// Cash flow risk
	if result.MonthlyCashFlow < 0 {
		riskScore += 3
	} else if result.MonthlyCashFlow < 100 {
		riskScore += 2
	} else if result.MonthlyCashFlow < 200 {
		riskScore += 1
	}

	// DSCR risk
	if result.DSCR < 1.0 {
		riskScore += 3
	} else if result.DSCR < 1.25 {
		riskScore += 2
	} else if result.DSCR < 1.5 {
		riskScore += 1
	}

	// Cap rate risk
	if result.CapRate < 4 {
		riskScore += 2
	} else if result.CapRate < 6 {
		riskScore += 1
	}

	// Cash-on-cash return risk
	if result.CashOnCashReturn < 8 && !result.IsInfiniteReturn {
		riskScore += 2
	} else if result.CashOnCashReturn < 12 && !result.IsInfiniteReturn {
		riskScore += 1
	}

	// Expense ratio risk
	if result.ExpenseRatio > 60 {
		riskScore += 2
	} else if result.ExpenseRatio > 50 {
		riskScore += 1
	}

	// Determine risk level
	if riskScore >= 8 {
		return "Very High"
	} else if riskScore >= 6 {
		return "High"
	} else if riskScore >= 3 {
		return "Medium"
	} else {
		return "Low"
	}
}

// generateBRRRRRecommendations provides specific BRRRR strategy recommendations
func (s *ArvService) generateBRRRRRecommendations(req ArvRequest, result ArvResult) []string {
	var recommendations []string

	// Cash flow recommendations
	if result.MonthlyCashFlow < 0 {
		recommendations = append(recommendations, "CRITICAL: Negative cash flow - property will require monthly contributions")
	} else if result.MonthlyCashFlow < 100 {
		recommendations = append(recommendations, "Low cash flow - consider higher rent or lower expenses")
	}

	// DSCR recommendations
	if result.DSCR < 1.0 {
		recommendations = append(recommendations, "CRITICAL: DSCR below 1.0 - property cannot service debt from income")
	} else if result.DSCR < 1.25 {
		recommendations = append(recommendations, "Low DSCR - lender may require higher down payment or reject loan")
	}

	// Refinance recommendations
	if result.CashRecovered >= result.TotalInvestment * 0.9 {
		recommendations = append(recommendations, "Excellent BRRRR opportunity - can recover most/all invested capital")
	} else if result.CashRecovered < result.TotalInvestment * 0.5 {
		recommendations = append(recommendations, "Limited cash recovery in refinance - consider if BRRRR is optimal strategy")
	}

	// Cap rate recommendations
	if result.CapRate < 4 {
		recommendations = append(recommendations, "Low cap rate - property may be overvalued for rental income")
	} else if result.CapRate > 10 {
		recommendations = append(recommendations, "High cap rate - verify income and expense estimates for accuracy")
	}

	// Expense ratio recommendations
	if result.ExpenseRatio > 60 {
		recommendations = append(recommendations, "High expense ratio - review all expense categories for accuracy")
	} else if result.ExpenseRatio < 30 {
		recommendations = append(recommendations, "Low expense ratio - ensure all expenses are accounted for")
	}

	// 70% rule comparison
	if !result.Is70RuleGood {
		recommendations = append(recommendations, "Property fails 70% rule - higher risk flip/BRRRR deal")
	}

	// Positive recommendations
	if result.IsInfiniteReturn && result.IsCashFlowPositive {
		recommendations = append(recommendations, "EXCELLENT: Infinite return with positive cash flow - ideal BRRRR deal")
	} else if result.CashOnCashReturn > 15 && result.IsCashFlowPositive {
		recommendations = append(recommendations, "Strong BRRRR opportunity with good returns and cash flow")
	}

	if len(recommendations) == 0 {
		recommendations = append(recommendations, "Moderate BRRRR opportunity - perform detailed due diligence")
	}

	return recommendations
}

// roundFinancialValues rounds all financial values to 2 decimal places
func (s *ArvService) roundFinancialValues(result *ArvResult) {
	result.AnnualGrossIncome = math.Round(result.AnnualGrossIncome*100) / 100
	result.EffectiveIncome = math.Round(result.EffectiveIncome*100) / 100
	result.AnnualExpenses = math.Round(result.AnnualExpenses*100) / 100
	result.ExpenseRatio = math.Round(result.ExpenseRatio*100) / 100
	result.NOI = math.Round(result.NOI*100) / 100
	result.MaxOffer70 = math.Round(result.MaxOffer70*100) / 100
	result.TotalInvestment = math.Round(result.TotalInvestment*100) / 100
	result.PotentialProfit = math.Round(result.PotentialProfit*100) / 100
	result.ProfitMargin = math.Round(result.ProfitMargin*100) / 100
	result.ROI = math.Round(result.ROI*100) / 100
	result.BrrrrMaxOffer = math.Round(result.BrrrrMaxOffer*100) / 100
	result.BrrrrProfit = math.Round(result.BrrrrProfit*100) / 100
	result.RefinanceAmount = math.Round(result.RefinanceAmount*100) / 100
	result.CashRecovered = math.Round(result.CashRecovered*100) / 100
	result.CashLeftIn = math.Round(result.CashLeftIn*100) / 100
	result.MonthlyDebtService = math.Round(result.MonthlyDebtService*100) / 100
	result.MonthlyCashFlow = math.Round(result.MonthlyCashFlow*100) / 100
	result.AnnualCashFlow = math.Round(result.AnnualCashFlow*100) / 100
	result.CashOnCashReturn = math.Round(result.CashOnCashReturn*100) / 100
	result.CapRate = math.Round(result.CapRate*100) / 100
	result.DSCR = math.Round(result.DSCR*100) / 100
}

// CalculateEnhancedBRRRR performs enhanced BRRRR analysis with new risk assessment
func (s *ArvService) CalculateEnhancedBRRRR(req ArvRequest) ArvResult {
	result := s.CalculateARV(req)
	
	// Apply enhanced risk assessment and recommendations
	result.RiskLevel = s.assessBRRRRisk(result)
	result.Recommendations = s.generateBRRRRRecommendations(req, result)
	
	return result
}
