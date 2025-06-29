package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"arvfinder-backend/services"
)

// PropertyHandler handles property-related HTTP requests
type PropertyHandler struct {
	propertyService *services.PropertyService
}

// NewPropertyHandler creates a new property handler
func NewPropertyHandler() *PropertyHandler {
	return &PropertyHandler{
		propertyService: services.NewPropertyService(),
	}
}

// PropertyEstimateRequest represents the request payload for property estimates
type PropertyEstimateRequest struct {
	StreetNumber string `json:"streetNumber" binding:"required"`
	StreetName   string `json:"streetName" binding:"required"`
	City         string `json:"city" binding:"required"`
	Zip          string `json:"zip" binding:"required"`
	State        string `json:"state"`
}

// PropertyEstimateResponse represents the response for property estimates
type PropertyEstimateResponse struct {
	Success bool                        `json:"success"`
	Data    *services.PropertyEstimate  `json:"data,omitempty"`
	Error   string                      `json:"error,omitempty"`
}

// PropertyHistoryResponse represents the response for property history
type PropertyHistoryResponse struct {
	Success bool                        `json:"success"`
	Data    []services.PropertyHistory  `json:"data,omitempty"`
	Error   string                      `json:"error,omitempty"`
}

// GetPropertyEstimate handles property estimate requests
func (h *PropertyHandler) GetPropertyEstimate(c *gin.Context) {
	var req PropertyEstimateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, PropertyEstimateResponse{
			Success: false,
			Error:   "Invalid request data: " + err.Error(),
		})
		return
	}

	// Convert to address components
	components := services.AddressComponents{
		StreetNumber: req.StreetNumber,
		StreetName:   req.StreetName,
		City:         req.City,
		Zip:          req.Zip,
		State:        req.State,
	}

	// Validate address
	if !h.propertyService.ValidateAddress(components) {
		c.JSON(http.StatusBadRequest, PropertyEstimateResponse{
			Success: false,
			Error:   "Invalid address components",
		})
		return
	}

	// Get property estimate
	estimate, err := h.propertyService.GetPropertyEstimate(components)
	if err != nil {
		c.JSON(http.StatusInternalServerError, PropertyEstimateResponse{
			Success: false,
			Error:   "Failed to get property estimate: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, PropertyEstimateResponse{
		Success: true,
		Data:    estimate,
	})
}

// GetPropertyHistory handles property history requests
func (h *PropertyHandler) GetPropertyHistory(c *gin.Context) {
	var req PropertyEstimateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, PropertyHistoryResponse{
			Success: false,
			Error:   "Invalid request data: " + err.Error(),
		})
		return
	}

	// Convert to address components
	components := services.AddressComponents{
		StreetNumber: req.StreetNumber,
		StreetName:   req.StreetName,
		City:         req.City,
		Zip:          req.Zip,
		State:        req.State,
	}

	// Validate address
	if !h.propertyService.ValidateAddress(components) {
		c.JSON(http.StatusBadRequest, PropertyHistoryResponse{
			Success: false,
			Error:   "Invalid address components",
		})
		return
	}

	// Get property history
	history, err := h.propertyService.GetPropertyHistory(components)
	if err != nil {
		c.JSON(http.StatusInternalServerError, PropertyHistoryResponse{
			Success: false,
			Error:   "Failed to get property history: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, PropertyHistoryResponse{
		Success: true,
		Data:    history,
	})
}

// AddressSuggestionsRequest represents the request for address suggestions
type AddressSuggestionsRequest struct {
	Input string `json:"input" binding:"required"`
}

// AddressSuggestionsResponse represents the response for address suggestions
type AddressSuggestionsResponse struct {
	Success bool                           `json:"success"`
	Data    []services.AddressSuggestion   `json:"data,omitempty"`
	Error   string                         `json:"error,omitempty"`
}

// GeocodeRequest represents the request for geocoding
type GeocodeRequest struct {
	Address string `json:"address" binding:"required"`
}

// GeocodeResponse represents the response for geocoding
type GeocodeResponse struct {
	Success bool                           `json:"success"`
	Data    *services.AddressComponents    `json:"data,omitempty"`
	Error   string                         `json:"error,omitempty"`
}

// GetAddressSuggestions handles address autocomplete requests
func (h *PropertyHandler) GetAddressSuggestions(c *gin.Context) {
	var req AddressSuggestionsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, AddressSuggestionsResponse{
			Success: false,
			Error:   "Invalid request data: " + err.Error(),
		})
		return
	}

	suggestions, err := h.propertyService.GetAddressSuggestions(req.Input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, AddressSuggestionsResponse{
			Success: false,
			Error:   "Failed to get address suggestions: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, AddressSuggestionsResponse{
		Success: true,
		Data:    suggestions,
	})
}

// GeocodeAddress handles address geocoding requests
func (h *PropertyHandler) GeocodeAddress(c *gin.Context) {
	var req GeocodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, GeocodeResponse{
			Success: false,
			Error:   "Invalid request data: " + err.Error(),
		})
		return
	}

	components, err := h.propertyService.GeocodeAddress(req.Address)
	if err != nil {
		c.JSON(http.StatusInternalServerError, GeocodeResponse{
			Success: false,
			Error:   "Failed to geocode address: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, GeocodeResponse{
		Success: true,
		Data:    components,
	})
}

// SearchProperties handles property search requests (placeholder)
func (h *PropertyHandler) SearchProperties(c *gin.Context) {
	// This could be extended to search multiple properties
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Property search endpoint - coming soon",
	})
}