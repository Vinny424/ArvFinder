package main

import (
	"log"
	"net/http"
	"os"
	"arvfinder-backend/database"
	"arvfinder-backend/handlers"
	"arvfinder-backend/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize database
	db, err := database.InitDB()
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer database.CloseDB()

	// Run database migrations
	err = database.RunMigrations(db)
	if err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	r := gin.Default()

	// Get Stripe secret key from environment or use default for development
	stripeSecretKey := os.Getenv("STRIPE_SECRET_KEY")
	if stripeSecretKey == "" {
		stripeSecretKey = "sk_test_51Rf9L600n2nnxa7pNjxkeVUzm8I54V9VZO1gg4P5iDckkGJzZegdbzyGMMHz7RzeocEequ2Ah1Wtb3Ru73Q8ES4m0041YIezPX"
	}

	// Initialize handlers
	arvHandler := handlers.NewArvHandler()
	stripeHandler := handlers.NewStripeHandler(stripeSecretKey)
	propertyHandler := handlers.NewPropertyHandler()
	authHandler := handlers.NewAuthHandler()

	// Security middleware
	r.Use(middleware.SecurityHeadersMiddleware())
	r.Use(middleware.RateLimitMiddleware())

	// CORS middleware
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
			"service": "arvfinder-backend",
		})
	})

	// API routes
	api := r.Group("/api/v1")
	{
		// Authentication routes
		auth := api.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.POST("/register", authHandler.Register)
			auth.POST("/verify-2fa", authHandler.Verify2FA)
			auth.POST("/refresh", refreshTokenHandler) // TODO: Implement
			auth.POST("/logout", logoutHandler)        // TODO: Implement
			auth.POST("/forgot-password", forgotPasswordHandler) // TODO: Implement
			auth.POST("/reset-password", resetPasswordHandler)   // TODO: Implement
		}

		// Property routes (protected)
		properties := api.Group("/properties")
		properties.Use(middleware.AuthMiddleware())
		{
			properties.GET("/", getPropertiesHandler)
			properties.POST("/", createPropertyHandler)
			properties.GET("/:id", getPropertyHandler)
			properties.PUT("/:id", updatePropertyHandler)
			properties.DELETE("/:id", deletePropertyHandler)
		}

		// ARV calculation routes (protected - disabled for now)
		arv := api.Group("/arv")
		// arv.Use(authMiddleware()) // Disable auth for now to test functionality
		{
			arv.POST("/calculate", arvHandler.CalculateARV)
			arv.POST("/70-rule", arvHandler.Calculate70Rule)
			arv.POST("/roi", arvHandler.CalculateROI)
			arv.POST("/cash-on-cash", arvHandler.CalculateCashOnCash)
			arv.POST("/cap-rate", arvHandler.CalculateCapRate)
			arv.POST("/estimate-from-comps", arvHandler.EstimateARVFromComps)
		}

		// Property estimate routes
		api.POST("/property-estimate", propertyHandler.GetPropertyEstimate)
		api.POST("/property-history", propertyHandler.GetPropertyHistory)
		api.POST("/address-suggestions", propertyHandler.GetAddressSuggestions)
		api.POST("/geocode-address", propertyHandler.GeocodeAddress)
		api.GET("/property-search", propertyHandler.SearchProperties)

		// Stripe payment routes
		payments := api.Group("/payments")
		{
			payments.GET("/plans", stripeHandler.GetSubscriptionPlans)
			payments.POST("/create-subscription", stripeHandler.CreateSubscription)
			payments.POST("/create-payment-intent", stripeHandler.CreatePaymentIntent)
			payments.POST("/create-report-payment", stripeHandler.CreateReportPayment)
			payments.POST("/cancel-subscription", stripeHandler.CancelSubscription)
			payments.POST("/update-subscription", stripeHandler.UpdateSubscription)
			payments.GET("/subscription-status", stripeHandler.GetSubscriptionStatus)
			payments.POST("/webhook", stripeHandler.HandleWebhook)
			payments.POST("/setup-prices", stripeHandler.SetupPrices) // For initial setup only
		}
	}

	log.Println("Server starting on :8080")
	log.Fatal(r.Run(":8080"))
}

// TODO: Implement these handlers
func refreshTokenHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Refresh token endpoint - to be implemented"})
}

func logoutHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Logout endpoint - to be implemented"})
}

func forgotPasswordHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Forgot password endpoint - to be implemented"})
}

func resetPasswordHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Reset password endpoint - to be implemented"})
}

func getPropertiesHandler(c *gin.Context) {
	// Sample data for now
	properties := []map[string]interface{}{
		{
			"id":      "1",
			"address": "123 Main St, Denver, CO",
			"price":   180000,
			"arv":     250000,
			"roi":     15.8,
		},
		{
			"id":      "2",
			"address": "456 Oak Ave, Boulder, CO",
			"price":   220000,
			"arv":     300000,
			"roi":     12.4,
		},
	}
	c.JSON(http.StatusOK, properties)
}

func createPropertyHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Create property endpoint - to be implemented"})
}

func getPropertyHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Get property endpoint - to be implemented"})
}

func updatePropertyHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Update property endpoint - to be implemented"})
}

func deletePropertyHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Delete property endpoint - to be implemented"})
}


