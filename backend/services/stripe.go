package services

import (
	"fmt"
	"log"

	"github.com/stripe/stripe-go/v79"
	"github.com/stripe/stripe-go/v79/customer"
	"github.com/stripe/stripe-go/v79/paymentintent"
	"github.com/stripe/stripe-go/v79/price"
	"github.com/stripe/stripe-go/v79/product"
	"github.com/stripe/stripe-go/v79/subscription"
	"github.com/stripe/stripe-go/v79/webhook"
)

// StripeService handles all Stripe-related operations
type StripeService struct {
	secretKey string
}

// NewStripeService creates a new Stripe service instance
func NewStripeService(secretKey string) *StripeService {
	stripe.Key = secretKey
	return &StripeService{
		secretKey: secretKey,
	}
}

// SubscriptionTier represents a subscription plan
type SubscriptionTier string

const (
	TierStarter      SubscriptionTier = "starter"
	TierProfessional SubscriptionTier = "professional"
	TierEnterprise   SubscriptionTier = "enterprise"
)

// SubscriptionPlan represents pricing and features for each tier
type SubscriptionPlan struct {
	Name        string  `json:"name"`
	Price       int64   `json:"price"`        // Price in cents
	PriceID     string  `json:"price_id"`     // Stripe Price ID
	Features    []string `json:"features"`
	ArvLimit    int     `json:"arv_limit"`    // -1 for unlimited
	Popular     bool    `json:"popular"`
}

// GetSubscriptionPlans returns all available subscription plans
func (s *StripeService) GetSubscriptionPlans() map[SubscriptionTier]SubscriptionPlan {
	return map[SubscriptionTier]SubscriptionPlan{
		TierStarter: {
			Name:     "Starter",
			Price:    0, // Free
			PriceID:  "", // No Stripe price for free tier
			ArvLimit: 10,
			Features: []string{
				"10 ARV calculations per month",
				"Basic property analysis",
				"Pay $9.99 per report generation",
				"Email support",
			},
			Popular: false,
		},
		TierProfessional: {
			Name:     "Professional",
			Price:    2900, // $29.00
			PriceID:  "price_professional_monthly", // Will be created in Stripe
			ArvLimit: -1, // Unlimited
			Features: []string{
				"Unlimited ARV calculations",
				"Advanced property analysis",
				"FREE report generation",
				"Custom reports with branding",
				"Mobile app access",
				"Priority support",
				"BRRRR strategy analysis",
				"Portfolio dashboard",
			},
			Popular: true,
		},
		TierEnterprise: {
			Name:     "Enterprise",
			Price:    5900, // $59.00
			PriceID:  "price_enterprise_monthly", // Will be created in Stripe
			ArvLimit: -1, // Unlimited
			Features: []string{
				"Everything in Professional",
				"FREE report generation",
				"API access",
				"Batch property processing",
				"White-label reports",
				"Dedicated support",
				"Advanced analytics",
				"Team collaboration",
				"Custom integrations",
			},
			Popular: false,
		},
	}
}

// ReportPaymentInfo contains information about report payments
type ReportPaymentInfo struct {
	Price       int64  `json:"price"`        // Price in cents
	Currency    string `json:"currency"`
	Description string `json:"description"`
}

// GetReportPaymentInfo returns information about one-time report payments
func (s *StripeService) GetReportPaymentInfo() ReportPaymentInfo {
	return ReportPaymentInfo{
		Price:       999, // $9.99
		Currency:    "usd",
		Description: "Professional ARV Analysis Report",
	}
}

// CreateCustomer creates a new Stripe customer
func (s *StripeService) CreateCustomer(email, name string) (*stripe.Customer, error) {
	params := &stripe.CustomerParams{
		Email: stripe.String(email),
		Name:  stripe.String(name),
	}

	return customer.New(params)
}

// CreateSubscription creates a new subscription for a customer
func (s *StripeService) CreateSubscription(customerID, priceID string) (*stripe.Subscription, error) {
	params := &stripe.SubscriptionParams{
		Customer: stripe.String(customerID),
		Items: []*stripe.SubscriptionItemsParams{
			{
				Price: stripe.String(priceID),
			},
		},
		PaymentBehavior: stripe.String("default_incomplete"),
		PaymentSettings: &stripe.SubscriptionPaymentSettingsParams{
			SaveDefaultPaymentMethod: stripe.String("on_subscription"),
			PaymentMethodOptions: &stripe.SubscriptionPaymentSettingsPaymentMethodOptionsParams{
				Card: &stripe.SubscriptionPaymentSettingsPaymentMethodOptionsCardParams{
					RequestThreeDSecure: stripe.String("automatic"),
				},
			},
		},
		// Ensure subscription is set to automatically collect payment
		CollectionMethod: stripe.String("charge_automatically"),
	}

	params.AddExpand("latest_invoice.payment_intent")
	params.AddExpand("customer")

	return subscription.New(params)
}

// CreatePaymentIntent creates a payment intent for one-time payments
func (s *StripeService) CreatePaymentIntent(amount int64, currency, customerID string) (*stripe.PaymentIntent, error) {
	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(amount),
		Currency: stripe.String(currency),
		Customer: stripe.String(customerID),
		AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
			Enabled: stripe.Bool(true),
		},
	}

	return paymentintent.New(params)
}

// CreateReportPaymentIntent creates a payment intent specifically for report generation
func (s *StripeService) CreateReportPaymentIntent(customerID, propertyID string) (*stripe.PaymentIntent, error) {
	reportInfo := s.GetReportPaymentInfo()
	
	params := &stripe.PaymentIntentParams{
		Amount:      stripe.Int64(reportInfo.Price),
		Currency:    stripe.String(reportInfo.Currency),
		Customer:    stripe.String(customerID),
		Description: stripe.String(reportInfo.Description),
		AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
			Enabled: stripe.Bool(true),
		},
		Metadata: map[string]string{
			"type":        "report_generation",
			"property_id": propertyID,
		},
	}

	return paymentintent.New(params)
}

// CancelSubscription cancels a subscription
func (s *StripeService) CancelSubscription(subscriptionID string) (*stripe.Subscription, error) {
	params := &stripe.SubscriptionCancelParams{}
	return subscription.Cancel(subscriptionID, params)
}

// UpdateSubscription updates a subscription to a new price
func (s *StripeService) UpdateSubscription(subscriptionID, newPriceID string) (*stripe.Subscription, error) {
	// Get current subscription to get the subscription item ID
	sub, err := subscription.Get(subscriptionID, nil)
	if err != nil {
		return nil, err
	}

	if len(sub.Items.Data) == 0 {
		return nil, fmt.Errorf("subscription has no items")
	}

	params := &stripe.SubscriptionParams{
		Items: []*stripe.SubscriptionItemsParams{
			{
				ID:    stripe.String(sub.Items.Data[0].ID),
				Price: stripe.String(newPriceID),
			},
		},
	}

	return subscription.Update(subscriptionID, params)
}

// GetSubscription retrieves subscription details
func (s *StripeService) GetSubscription(subscriptionID string) (*stripe.Subscription, error) {
	return subscription.Get(subscriptionID, nil)
}

// GetCustomer retrieves customer details
func (s *StripeService) GetCustomer(customerID string) (*stripe.Customer, error) {
	return customer.Get(customerID, nil)
}

// ValidateWebhookSignature validates Stripe webhook signatures
func (s *StripeService) ValidateWebhookSignature(payload []byte, signature, endpointSecret string) (stripe.Event, error) {
	return webhook.ConstructEvent(payload, signature, endpointSecret)
}

// CreatePrices creates subscription prices in Stripe (run this once during setup)
func (s *StripeService) CreatePrices() error {
	plans := s.GetSubscriptionPlans()

	// Create Professional plan price
	if professionalPlan, exists := plans[TierProfessional]; exists && professionalPlan.Price > 0 {
		// First create the product
		productParams := &stripe.ProductParams{
			Name:        stripe.String("ArvFinder Professional"),
			Description: stripe.String("Professional plan with unlimited ARV calculations and advanced features"),
		}
		prod, err := product.New(productParams)
		if err != nil {
			log.Printf("Error creating professional product: %v", err)
			return err
		}

		// Then create the price
		params := &stripe.PriceParams{
			UnitAmount: stripe.Int64(professionalPlan.Price),
			Currency:   stripe.String("usd"),
			Recurring: &stripe.PriceRecurringParams{
				Interval: stripe.String("month"),
			},
			Product: stripe.String(prod.ID),
		}

		professionalPrice, err := price.New(params)
		if err != nil {
			log.Printf("Error creating professional price: %v", err)
			return err
		}
		log.Printf("Created professional price: %s", professionalPrice.ID)
	}

	// Create Enterprise plan price
	if enterprisePlan, exists := plans[TierEnterprise]; exists && enterprisePlan.Price > 0 {
		// First create the product
		productParams := &stripe.ProductParams{
			Name:        stripe.String("ArvFinder Enterprise"),
			Description: stripe.String("Enterprise plan with API access, white-label reports, and team features"),
		}
		prod, err := product.New(productParams)
		if err != nil {
			log.Printf("Error creating enterprise product: %v", err)
			return err
		}

		// Then create the price
		params := &stripe.PriceParams{
			UnitAmount: stripe.Int64(enterprisePlan.Price),
			Currency:   stripe.String("usd"),
			Recurring: &stripe.PriceRecurringParams{
				Interval: stripe.String("month"),
			},
			Product: stripe.String(prod.ID),
		}

		enterprisePrice, err := price.New(params)
		if err != nil {
			log.Printf("Error creating enterprise price: %v", err)
			return err
		}
		log.Printf("Created enterprise price: %s", enterprisePrice.ID)
	}

	return nil
}

// Usage tracking for subscription limits
func (s *StripeService) TrackUsage(subscriptionTier SubscriptionTier, currentUsage int) bool {
	plans := s.GetSubscriptionPlans()
	plan, exists := plans[subscriptionTier]
	if !exists {
		return false
	}

	// Unlimited usage for paid plans
	if plan.ArvLimit == -1 {
		return true
	}

	// Check if within limits
	return currentUsage < plan.ArvLimit
}

// GetSubscriptionStatus returns the status and limits for a subscription
type SubscriptionStatus struct {
	Tier               SubscriptionTier `json:"tier"`
	ArvLimit           int             `json:"arv_limit"`
	ArvUsed            int             `json:"arv_used"`
	IsActive           bool            `json:"is_active"`
	NextBilling        string          `json:"next_billing,omitempty"`
	FreeReports        bool            `json:"free_reports"`
	ReportPrice        int64           `json:"report_price,omitempty"` // in cents
}

func (s *StripeService) GetSubscriptionStatus(tier SubscriptionTier, currentUsage int) SubscriptionStatus {
	plans := s.GetSubscriptionPlans()
	plan, exists := plans[tier]
	if !exists {
		// Default to starter if tier not found
		plan = plans[TierStarter]
		tier = TierStarter
	}

	// Professional and Enterprise get free reports, Starter pays $9.99
	freeReports := tier == TierProfessional || tier == TierEnterprise
	var reportPrice int64 = 0
	if !freeReports {
		reportPrice = s.GetReportPaymentInfo().Price
	}

	return SubscriptionStatus{
		Tier:        tier,
		ArvLimit:    plan.ArvLimit,
		ArvUsed:     currentUsage,
		IsActive:    true, // This would be determined by actual subscription status
		FreeReports: freeReports,
		ReportPrice: reportPrice,
	}
}

// CanGenerateReportForFree checks if user can generate reports without payment
func (s *StripeService) CanGenerateReportForFree(tier SubscriptionTier) bool {
	return tier == TierProfessional || tier == TierEnterprise
}