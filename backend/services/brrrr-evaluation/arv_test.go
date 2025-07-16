package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculateARV_Basic(t *testing.T) {
	service := NewArvService()

	req := ArvRequest{
		PurchasePrice:  50000,
		RehabCost:      15000,
		HoldingCosts:   3000,
		ClosingCosts:   2000,
		ARV:            120000,
		FinancingCosts: 4000,
		SellingCosts:   6000,
	}

	result := service.CalculateARV(req)

	assert.Equal(t, 74000.00, result.TotalInvestment)
	assert.Equal(t, 40000.00, result.PotentialProfit)
	assert.InDelta(t, 54.05, result.ROI, 0.01)
	assert.Equal(t, true, result.Is70RuleGood)
	assert.Equal(t, "Low", result.RiskLevel)
	assert.NotEmpty(t, result.Recommendations)
}

// Additional tests for current BRRRR functionality edge cases

func TestCalculateARV_HighRehabCosts(t *testing.T) {
	service := NewArvService()

	req := ArvRequest{
		PurchasePrice:  60000,
		RehabCost:      45000, // Very high rehab - 37.5% of ARV
		HoldingCosts:   5000,
		ClosingCosts:   3000,
		ARV:            120000,
		FinancingCosts: 2000,
		SellingCosts:   6000,
	}

	result := service.CalculateARV(req)

	// Should flag high rehab costs in recommendations
	assert.Contains(t, result.Recommendations, "Rehab costs are high (>30% of ARV) - verify estimates with contractors")
	assert.Equal(t, "Very High", result.RiskLevel)
}

func TestCalculateARV_BreaksEven(t *testing.T) {
	service := NewArvService()

	req := ArvRequest{
		PurchasePrice:  50000,
		RehabCost:      20000,
		HoldingCosts:   5000,
		ClosingCosts:   3000,
		ARV:            84000, // Exactly breaks even after selling costs
		FinancingCosts: 0,
		SellingCosts:   6000,
	}

	result := service.CalculateARV(req)

	assert.Equal(t, 0.0, result.PotentialProfit)
	assert.Equal(t, 0.0, result.ROI)
	assert.Equal(t, "Very High", result.RiskLevel)
}

func TestCalculate70Rule(t *testing.T) {
	service := NewArvService()
	arv := 100000.0
	rehab := 20000.0

	expected := 50000.0
	actual := service.Calculate70Rule(arv, rehab)
	assert.Equal(t, expected, actual)
}

func TestCalculateROI_ZeroInvestment(t *testing.T) {
	service := NewArvService()
	roi := service.CalculateROI(10000, 0)
	assert.Equal(t, 0.0, roi)
}

func TestCalculateROI_PositiveInvestment(t *testing.T) {
	service := NewArvService()
	roi := service.CalculateROI(25000, 100000)
	assert.Equal(t, 25.0, roi)
}

func TestEstimateARVFromComps(t *testing.T) {
	service := NewArvService()
	comps := []ComparableProperty{
		{SalePrice: 100000, Distance: 0.5, Bedrooms: 3, Bathrooms: 2.0, SquareFeet: 1200},
		{SalePrice: 95000, Distance: 1.0, Bedrooms: 3, Bathrooms: 1.5, SquareFeet: 1100},
	}
	est := service.EstimateARVFromComps(comps, 3, 2.0, 1200)
	assert.InDelta(t, 100000, est, 5000)
}

func TestEstimateARVFromComps_EmptyComps(t *testing.T) {
	service := NewArvService()
	comps := []ComparableProperty{}
	est := service.EstimateARVFromComps(comps, 3, 2.0, 1200)
	assert.Equal(t, 0.0, est)
}

func TestRiskAssessment(t *testing.T) {
	service := NewArvService()

	lowRisk := service.assessRisk(25, true, 120000, 50000)
	assert.Equal(t, "Low", lowRisk)

	mediumRisk := service.assessRisk(12, false, 120000, 80000)
	assert.Equal(t, "Medium", mediumRisk)

	highRisk := service.assessRisk(6, false, 120000, 110000)
	assert.Equal(t, "High", highRisk)

	veryHighRisk := service.assessRisk(1, false, 120000, 115000)
	assert.Equal(t, "Very High", veryHighRisk)
}

func TestGenerateRecommendations(t *testing.T) {
	service := NewArvService()
	req := ArvRequest{
		PurchasePrice:  90000,
		RehabCost:      40000,
		HoldingCosts:   10000,
		ClosingCosts:   5000,
		ARV:            120000,
		FinancingCosts: 5000,
		SellingCosts:   5000,
	}
	recs := service.generateRecommendations(req, 5, false)
	assert.Greater(t, len(recs), 0)
}

func TestCalculateCashOnCashReturn(t *testing.T) {
	service := NewArvService()

	// Test normal case
	cashReturn := service.CalculateCashOnCashReturn(12000, 100000)
	assert.Equal(t, 12.0, cashReturn)

	// Test zero investment
	cashReturn = service.CalculateCashOnCashReturn(12000, 0)
	assert.Equal(t, 0.0, cashReturn)

	// Test negative investment
	cashReturn = service.CalculateCashOnCashReturn(12000, -100000)
	assert.Equal(t, 0.0, cashReturn)
}

func TestCalculateCapRate(t *testing.T) {
	service := NewArvService()

	// Test normal case
	capRate := service.CalculateCapRate(15000, 200000)
	assert.Equal(t, 7.5, capRate)

	// Test zero property value
	capRate = service.CalculateCapRate(15000, 0)
	assert.Equal(t, 0.0, capRate)

	// Test negative property value
	capRate = service.CalculateCapRate(15000, -200000)
	assert.Equal(t, 0.0, capRate)
}

func TestCalculateARV_EdgeCases(t *testing.T) {
	service := NewArvService()

	// Test with very high purchase price (doesn't meet 70% rule)
	req := ArvRequest{
		PurchasePrice:  90000,
		RehabCost:      15000,
		HoldingCosts:   3000,
		ClosingCosts:   2000,
		ARV:            120000,
		FinancingCosts: 4000,
		SellingCosts:   6000,
	}

	result := service.CalculateARV(req)
	assert.Equal(t, false, result.Is70RuleGood)
	assert.Equal(t, 69000.0, result.MaxOffer70) // (120000 * 0.7) - 15000
}

func TestCalculateComparableAdjustments(t *testing.T) {
	service := NewArvService()

	comp := ComparableProperty{
		Bedrooms:   2,
		Bathrooms:  1.5,
		SquareFeet: 1000,
	}

	// Subject property: 3 bed, 2 bath, 1200 sq ft
	adjustments := service.calculateComparableAdjustments(comp, 3, 2.0, 1200)

	// Expected: +1 bed ($5000) + 0.5 bath ($1500) + 200 sq ft ($10000) = $16500
	expected := 5000 + 1500 + 10000
	assert.Equal(t, float64(expected), adjustments)
}

func TestBRRRRStrategy(t *testing.T) {
	service := NewArvService()

	req := ArvRequest{
		PurchasePrice:  60000,
		RehabCost:      20000,
		HoldingCosts:   5000,
		ClosingCosts:   3000,
		ARV:            120000,
		FinancingCosts: 2000,
		SellingCosts:   6000,
	}

	result := service.CalculateARV(req)

	// BRRRR max offer = (ARV * 0.75) - Rehab = (120000 * 0.75) - 20000 = 70000
	assert.Equal(t, 70000.0, result.BrrrrMaxOffer)
	assert.Equal(t, 24000.0, result.BrrrrProfit) // Same as potential profit in this case
}
