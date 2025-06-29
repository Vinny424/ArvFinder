package handlers

import (
	"net/http"
	"arvfinder-backend/services"
	
	"github.com/gin-gonic/gin"
)

// ArvHandler handles ARV-related endpoints
type ArvHandler struct {
	arvService *services.ArvService
}

// NewArvHandler creates a new ARV handler
func NewArvHandler() *ArvHandler {
	return &ArvHandler{
		arvService: services.NewArvService(),
	}
}

// CalculateARV handles ARV calculation requests
func (h *ArvHandler) CalculateARV(c *gin.Context) {
	var req services.ArvRequest
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request data",
			"details": err.Error(),
		})
		return
	}
	
	// Perform ARV calculation
	result := h.arvService.CalculateARV(req)
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": result,
	})
}

// Calculate70Rule handles 70% rule calculation requests
func (h *ArvHandler) Calculate70Rule(c *gin.Context) {
	var req struct {
		ARV       float64 `json:"arv" binding:"required,min=1"`
		RehabCost float64 `json:"rehab_cost" binding:"min=0"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request data",
			"details": err.Error(),
		})
		return
	}
	
	maxOffer := h.arvService.Calculate70Rule(req.ARV, req.RehabCost)
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"arv": req.ARV,
			"rehab_cost": req.RehabCost,
			"max_offer": maxOffer,
			"rule": "70% Rule: Max offer = (ARV × 0.70) - Rehab costs",
		},
	})
}

// CalculateROI handles ROI calculation requests
func (h *ArvHandler) CalculateROI(c *gin.Context) {
	var req struct {
		Profit     float64 `json:"profit" binding:"required"`
		Investment float64 `json:"investment" binding:"required,min=1"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request data",
			"details": err.Error(),
		})
		return
	}
	
	roi := h.arvService.CalculateROI(req.Profit, req.Investment)
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"profit": req.Profit,
			"investment": req.Investment,
			"roi": roi,
			"roi_formatted": roi,
		},
	})
}

// CalculateCashOnCash handles cash-on-cash return calculation requests
func (h *ArvHandler) CalculateCashOnCash(c *gin.Context) {
	var req struct {
		AnnualCashFlow     float64 `json:"annual_cash_flow" binding:"required"`
		TotalCashInvested  float64 `json:"total_cash_invested" binding:"required,min=1"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request data",
			"details": err.Error(),
		})
		return
	}
	
	cashOnCashReturn := h.arvService.CalculateCashOnCashReturn(req.AnnualCashFlow, req.TotalCashInvested)
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"annual_cash_flow": req.AnnualCashFlow,
			"total_cash_invested": req.TotalCashInvested,
			"cash_on_cash_return": cashOnCashReturn,
		},
	})
}

// CalculateCapRate handles cap rate calculation requests
func (h *ArvHandler) CalculateCapRate(c *gin.Context) {
	var req struct {
		NetOperatingIncome float64 `json:"net_operating_income" binding:"required"`
		PropertyValue      float64 `json:"property_value" binding:"required,min=1"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request data",
			"details": err.Error(),
		})
		return
	}
	
	capRate := h.arvService.CalculateCapRate(req.NetOperatingIncome, req.PropertyValue)
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"net_operating_income": req.NetOperatingIncome,
			"property_value": req.PropertyValue,
			"cap_rate": capRate,
		},
	})
}

// EstimateARVFromComps handles ARV estimation from comparable properties
func (h *ArvHandler) EstimateARVFromComps(c *gin.Context) {
	var req struct {
		Comparables       []services.ComparableProperty `json:"comparables" binding:"required,dive"`
		SubjectBedrooms   int                          `json:"subject_bedrooms" binding:"required,min=0"`
		SubjectBathrooms  float64                      `json:"subject_bathrooms" binding:"required,min=0"`
		SubjectSquareFeet int                          `json:"subject_square_feet" binding:"required,min=1"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request data",
			"details": err.Error(),
		})
		return
	}
	
	estimatedARV := h.arvService.EstimateARVFromComps(
		req.Comparables,
		req.SubjectBedrooms,
		req.SubjectBathrooms,
		req.SubjectSquareFeet,
	)
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"estimated_arv": estimatedARV,
			"comparables_used": len(req.Comparables),
			"subject_property": gin.H{
				"bedrooms": req.SubjectBedrooms,
				"bathrooms": req.SubjectBathrooms,
				"square_feet": req.SubjectSquareFeet,
			},
		},
	})
}