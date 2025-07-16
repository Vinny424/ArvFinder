package services

import (
	"math"
)

// ArvRequest represents the input data for ARV calculation
type ArvRequest struct {
	PurchasePrice  float64 `json:"purchase_price" binding:"required,min=1"`
	RehabCost      float64 `json:"rehab_cost" binding:"min=0"`
	HoldingCosts   float64 `json:"holding_costs" binding:"min=0"`
	ClosingCosts   float64 `json:"closing_costs" binding:"min=0"`
	ARV            float64 `json:"arv" binding:"required,min=1"`
	FinancingCosts float64 `json:"financing_costs" binding:"min=0"`
	SellingCosts   float64 `json:"selling_costs" binding:"min=0"`
}

// ArvResult represents the calculated ARV analysis results
type ArvResult struct {
	// Input values
	PurchasePrice  float64 `json:"purchase_price"`
	RehabCost      float64 `json:"rehab_cost"`
	HoldingCosts   float64 `json:"holding_costs"`
	ClosingCosts   float64 `json:"closing_costs"`
	ARV            float64 `json:"arv"`
	FinancingCosts float64 `json:"financing_costs"`
	SellingCosts   float64 `json:"selling_costs"`

	// 70% Rule calculations
	MaxOffer70   float64 `json:"max_offer_70"`
	Is70RuleGood bool    `json:"is_70_rule_good"`

	// Investment analysis
	TotalInvestment float64 `json:"total_investment"`
	PotentialProfit float64 `json:"potential_profit"`
	ProfitMargin    float64 `json:"profit_margin"`
	ROI             float64 `json:"roi"`

	// BRRRR Strategy metrics
	BrrrrMaxOffer float64 `json:"brrrr_max_offer"`
	BrrrrProfit   float64 `json:"brrrr_profit"`

	// Risk assessment
	RiskLevel       string   `json:"risk_level"`
	Recommendations []string `json:"recommendations"`
}

// ArvService handles ARV calculations and analysis
type ArvService struct{}

// NewArvService creates a new ARV service instance
func NewArvService() *ArvService {
	return &ArvService{}
}

// CalculateARV performs comprehensive ARV analysis
func (s *ArvService) CalculateARV(req ArvRequest) ArvResult {
	// Calculate total investment
	totalInvestment := req.PurchasePrice + req.RehabCost + req.HoldingCosts +
		req.ClosingCosts + req.FinancingCosts

	// 70% Rule: Max offer = (ARV * 0.70) - Rehab costs
	maxOffer70 := (req.ARV * 0.70) - req.RehabCost
	is70RuleGood := req.PurchasePrice <= maxOffer70

	// Calculate potential profit (including selling costs)
	potentialProfit := req.ARV - totalInvestment - req.SellingCosts

	// Calculate profit margin and ROI
	var profitMargin, roi float64
	if totalInvestment > 0 {
		profitMargin = (potentialProfit / totalInvestment) * 100
		roi = (potentialProfit / totalInvestment) * 100
	}

	// BRRRR Strategy calculations (75% of ARV for refinancing)
	brrrrMaxOffer := (req.ARV * 0.75) - req.RehabCost
	brrrrProfit := req.ARV - totalInvestment - req.SellingCosts

	// Risk assessment
	riskLevel := s.assessRisk(profitMargin, is70RuleGood, req.ARV, req.PurchasePrice)
	recommendations := s.generateRecommendations(req, profitMargin, is70RuleGood)

	return ArvResult{
		PurchasePrice:   req.PurchasePrice,
		RehabCost:       req.RehabCost,
		HoldingCosts:    req.HoldingCosts,
		ClosingCosts:    req.ClosingCosts,
		ARV:             req.ARV,
		FinancingCosts:  req.FinancingCosts,
		SellingCosts:    req.SellingCosts,
		MaxOffer70:      math.Round(maxOffer70*100) / 100,
		Is70RuleGood:    is70RuleGood,
		TotalInvestment: math.Round(totalInvestment*100) / 100,
		PotentialProfit: math.Round(potentialProfit*100) / 100,
		ProfitMargin:    math.Round(profitMargin*100) / 100,
		ROI:             math.Round(roi*100) / 100,
		BrrrrMaxOffer:   math.Round(brrrrMaxOffer*100) / 100,
		BrrrrProfit:     math.Round(brrrrProfit*100) / 100,
		RiskLevel:       riskLevel,
		Recommendations: recommendations,
	}
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
