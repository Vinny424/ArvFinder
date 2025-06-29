package handlers

import (
	"io"
	"net/http"
	"arvfinder-backend/services"

	"github.com/gin-gonic/gin"
)

// StripeHandler handles Stripe-related endpoints
type StripeHandler struct {
	stripeService *services.StripeService
}

// NewStripeHandler creates a new Stripe handler
func NewStripeHandler(stripeSecretKey string) *StripeHandler {
	return &StripeHandler{
		stripeService: services.NewStripeService(stripeSecretKey),
	}
}

// GetSubscriptionPlans returns available subscription plans
func (h *StripeHandler) GetSubscriptionPlans(c *gin.Context) {
	plans := h.stripeService.GetSubscriptionPlans()
	reportInfo := h.stripeService.GetReportPaymentInfo()
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"plans": plans,
			"report_payment": reportInfo,
		},
	})
}

// CreateSubscription creates a new subscription for a customer
func (h *StripeHandler) CreateSubscription(c *gin.Context) {
	var req struct {
		Email   string `json:"email" binding:"required,email"`
		Name    string `json:"name" binding:"required"`
		PriceID string `json:"price_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	// Create customer first
	customer, err := h.stripeService.CreateCustomer(req.Email, req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create customer",
			"details": err.Error(),
		})
		return
	}

	// Create subscription
	subscription, err := h.stripeService.CreateSubscription(customer.ID, req.PriceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create subscription",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"customer_id":     customer.ID,
			"subscription_id": subscription.ID,
			"client_secret":   subscription.LatestInvoice.PaymentIntent.ClientSecret,
			"status":          subscription.Status,
		},
	})
}

// CreatePaymentIntent creates a payment intent for one-time payments
func (h *StripeHandler) CreatePaymentIntent(c *gin.Context) {
	var req struct {
		Amount     int64  `json:"amount" binding:"required,min=50"`     // Minimum $0.50
		Currency   string `json:"currency" binding:"required"`
		CustomerID string `json:"customer_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	// Default to USD if not specified
	if req.Currency == "" {
		req.Currency = "usd"
	}

	paymentIntent, err := h.stripeService.CreatePaymentIntent(req.Amount, req.Currency, req.CustomerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create payment intent",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"client_secret": paymentIntent.ClientSecret,
			"payment_intent_id": paymentIntent.ID,
		},
	})
}

// CreateReportPayment creates a payment intent for report generation
func (h *StripeHandler) CreateReportPayment(c *gin.Context) {
	var req struct {
		CustomerEmail string `json:"customer_email" binding:"required,email"`
		CustomerName  string `json:"customer_name" binding:"required"`
		PropertyID    string `json:"property_id" binding:"required"`
		UserTier      string `json:"user_tier"` // starter, professional, enterprise
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	// Check if user gets free reports
	tier := services.SubscriptionTier(req.UserTier)
	if req.UserTier == "" {
		tier = services.TierStarter // Default to starter if not specified
	}

	if h.stripeService.CanGenerateReportForFree(tier) {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data": gin.H{
				"free_report": true,
				"message": "Report generation is included in your subscription",
			},
		})
		return
	}

	// Create customer first if they don't exist
	customer, err := h.stripeService.CreateCustomer(req.CustomerEmail, req.CustomerName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create customer",
			"details": err.Error(),
		})
		return
	}

	// Create payment intent for report
	paymentIntent, err := h.stripeService.CreateReportPaymentIntent(customer.ID, req.PropertyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create payment intent",
			"details": err.Error(),
		})
		return
	}

	reportInfo := h.stripeService.GetReportPaymentInfo()

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"client_secret": paymentIntent.ClientSecret,
			"payment_intent_id": paymentIntent.ID,
			"customer_id": customer.ID,
			"amount": reportInfo.Price,
			"currency": reportInfo.Currency,
			"description": reportInfo.Description,
			"free_report": false,
		},
	})
}

// CancelSubscription cancels a user's subscription
func (h *StripeHandler) CancelSubscription(c *gin.Context) {
	var req struct {
		SubscriptionID string `json:"subscription_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	subscription, err := h.stripeService.CancelSubscription(req.SubscriptionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to cancel subscription",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"subscription_id": subscription.ID,
			"status": subscription.Status,
			"canceled_at": subscription.CanceledAt,
		},
	})
}

// UpdateSubscription updates a subscription to a new plan
func (h *StripeHandler) UpdateSubscription(c *gin.Context) {
	var req struct {
		SubscriptionID string `json:"subscription_id" binding:"required"`
		NewPriceID     string `json:"new_price_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	subscription, err := h.stripeService.UpdateSubscription(req.SubscriptionID, req.NewPriceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update subscription",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"subscription_id": subscription.ID,
			"status": subscription.Status,
		},
	})
}

// GetSubscriptionStatus returns subscription status and usage
func (h *StripeHandler) GetSubscriptionStatus(c *gin.Context) {
	// This would typically get the user's current subscription from the database
	// For now, return a mock response
	tier := services.TierStarter // This would come from the user's database record
	currentUsage := 3            // This would come from usage tracking

	status := h.stripeService.GetSubscriptionStatus(tier, currentUsage)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": status,
	})
}

// HandleWebhook handles Stripe webhooks
func (h *StripeHandler) HandleWebhook(c *gin.Context) {
	const MaxBodyBytes = int64(65536)
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, MaxBodyBytes)

	payload, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "Error reading request body",
		})
		return
	}

	// Get the signature header
	signature := c.GetHeader("Stripe-Signature")
	
	// In production, you would store this in environment variables
	endpointSecret := "whsec_your_webhook_secret_here"

	event, err := h.stripeService.ValidateWebhookSignature(payload, signature, endpointSecret)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid signature",
		})
		return
	}

	// Handle the event
	switch event.Type {
	case "payment_intent.succeeded":
		// Handle successful payment
		// Update user's subscription status in database
		break
	case "invoice.payment_succeeded":
		// Handle successful subscription payment
		// Update user's subscription status and reset usage counters
		break
	case "customer.subscription.deleted":
		// Handle subscription cancellation
		// Update user's subscription status in database
		break
	case "customer.subscription.updated":
		// Handle subscription updates
		// Update user's subscription tier in database
		break
	default:
		// Unexpected event type
		break
	}

	c.JSON(http.StatusOK, gin.H{
		"received": true,
	})
}

// SetupPrices creates the subscription prices in Stripe (for initial setup)
func (h *StripeHandler) SetupPrices(c *gin.Context) {
	err := h.stripeService.CreatePrices()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create prices",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Prices created successfully",
	})
}