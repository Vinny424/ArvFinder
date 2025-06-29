package models

import (
	"time"
)

type Property struct {
	ID           string    `json:"id" db:"id"`
	TenantID     string    `json:"tenant_id" db:"tenant_id"`
	Address      string    `json:"address" db:"address"`
	City         string    `json:"city" db:"city"`
	State        string    `json:"state" db:"state"`
	ZipCode      string    `json:"zip_code" db:"zip_code"`
	Price        float64   `json:"price" db:"price"`
	ARV          float64   `json:"arv" db:"arv"`
	RehabCost    float64   `json:"rehab_cost" db:"rehab_cost"`
	HoldingCosts float64   `json:"holding_costs" db:"holding_costs"`
	ClosingCosts float64   `json:"closing_costs" db:"closing_costs"`
	Bedrooms     int       `json:"bedrooms" db:"bedrooms"`
	Bathrooms    float64   `json:"bathrooms" db:"bathrooms"`
	SquareFeet   int       `json:"square_feet" db:"square_feet"`
	LotSize      float64   `json:"lot_size" db:"lot_size"`
	YearBuilt    int       `json:"year_built" db:"year_built"`
	PropertyType string    `json:"property_type" db:"property_type"`
	Notes        string    `json:"notes" db:"notes"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

type ArvCalculation struct {
	ID           string    `json:"id" db:"id"`
	PropertyID   string    `json:"property_id" db:"property_id"`
	TenantID     string    `json:"tenant_id" db:"tenant_id"`
	PurchasePrice float64  `json:"purchase_price" db:"purchase_price"`
	RehabCost    float64   `json:"rehab_cost" db:"rehab_cost"`
	HoldingCosts float64   `json:"holding_costs" db:"holding_costs"`
	ClosingCosts float64   `json:"closing_costs" db:"closing_costs"`
	ARV          float64   `json:"arv" db:"arv"`
	MaxOffer     float64   `json:"max_offer" db:"max_offer"`
	PotentialProfit float64 `json:"potential_profit" db:"potential_profit"`
	ProfitMargin float64   `json:"profit_margin" db:"profit_margin"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

type Comparable struct {
	ID           string    `json:"id" db:"id"`
	PropertyID   string    `json:"property_id" db:"property_id"`
	Address      string    `json:"address" db:"address"`
	SalePrice    float64   `json:"sale_price" db:"sale_price"`
	SaleDate     time.Time `json:"sale_date" db:"sale_date"`
	Distance     float64   `json:"distance" db:"distance"`
	Bedrooms     int       `json:"bedrooms" db:"bedrooms"`
	Bathrooms    float64   `json:"bathrooms" db:"bathrooms"`
	SquareFeet   int       `json:"square_feet" db:"square_feet"`
	PricePerSqFt float64   `json:"price_per_sq_ft" db:"price_per_sq_ft"`
	Adjustments  float64   `json:"adjustments" db:"adjustments"`
	AdjustedValue float64  `json:"adjusted_value" db:"adjusted_value"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}