package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	
	"googlemaps.github.io/maps"
)

// PropertyService handles property data and estimates
type PropertyService struct {
	repliersAPIKey string
	googleMapsClient *maps.Client
}

// NewPropertyService creates a new property service instance
func NewPropertyService() *PropertyService {
	googleAPIKey := os.Getenv("GOOGLE_MAPS_API_KEY")
	var googleClient *maps.Client
	
	if googleAPIKey != "" {
		client, err := maps.NewClient(maps.WithAPIKey(googleAPIKey))
		if err == nil {
			googleClient = client
		}
	}
	
	return &PropertyService{
		repliersAPIKey: os.Getenv("REPLIERS_API_KEY"),
		googleMapsClient: googleClient,
	}
}

// AddressComponents represents the components of an address
type AddressComponents struct {
	StreetNumber string `json:"streetNumber"`
	StreetName   string `json:"streetName"`
	City         string `json:"city"`
	Zip          string `json:"zip"`
	State        string `json:"state,omitempty"`
}

// PropertyEstimate represents property estimate data
type PropertyEstimate struct {
	Address       string             `json:"address"`
	Components    AddressComponents  `json:"components"`
	EstimatedValue int64             `json:"estimatedValue,omitempty"`
	RentEstimate   int64             `json:"rentEstimate,omitempty"`
	Bedrooms       int               `json:"bedrooms,omitempty"`
	Bathrooms      int               `json:"bathrooms,omitempty"`
	SquareFootage  int               `json:"squareFootage,omitempty"`
	YearBuilt      int               `json:"yearBuilt,omitempty"`
	PropertyType   string            `json:"propertyType,omitempty"`
	Neighborhood   string            `json:"neighborhood,omitempty"`
	Comparables    []PropertyComp    `json:"comparables,omitempty"`
	History        []PropertyHistory `json:"history,omitempty"`
}

// PropertyComp represents comparable property data
type PropertyComp struct {
	Address   string `json:"address"`
	Price     int64  `json:"price"`
	SqFt      int    `json:"sqft"`
	Distance  string `json:"distance,omitempty"`
	SoldDate  string `json:"soldDate,omitempty"`
}

// PropertyHistory represents historical property data
type PropertyHistory struct {
	Date  string `json:"date"`
	Price int64  `json:"price"`
	Event string `json:"event"` // "sold", "listed", "pending", etc.
}

// RepliersEstimateRequest represents the request payload for Repliers API
type RepliersEstimateRequest struct {
	Address struct {
		StreetNumber string `json:"streetNumber"`
		StreetName   string `json:"streetName"`
		City         string `json:"city"`
		Zip          string `json:"zip"`
	} `json:"address"`
	Details struct {
		Bedrooms      int    `json:"bedrooms,omitempty"`
		Bathrooms     int    `json:"bathrooms,omitempty"`
		SquareFootage int    `json:"squareFootage,omitempty"`
		YearBuilt     int    `json:"yearBuilt,omitempty"`
		PropertyType  string `json:"propertyType,omitempty"`
	} `json:"details,omitempty"`
}

// GetPropertyEstimate fetches property estimate from Repliers API
func (s *PropertyService) GetPropertyEstimate(components AddressComponents) (*PropertyEstimate, error) {
	if s.repliersAPIKey == "" {
		return s.getFallbackEstimate(components), nil
	}

	// Create estimate request for Repliers API
	requestBody := RepliersEstimateRequest{}
	requestBody.Address.StreetNumber = components.StreetNumber
	requestBody.Address.StreetName = components.StreetName
	requestBody.Address.City = components.City
	requestBody.Address.Zip = components.Zip

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Make request to Repliers API
	req, err := http.NewRequest("POST", "https://api.repliers.io/estimates", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("REPLIERS-API-KEY", s.repliersAPIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return s.getFallbackEstimate(components), nil // Fallback on error
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return s.getFallbackEstimate(components), nil // Fallback on error
	}

	var estimate PropertyEstimate
	if err := json.NewDecoder(resp.Body).Decode(&estimate); err != nil {
		return s.getFallbackEstimate(components), nil // Fallback on error
	}

	return &estimate, nil
}

// GetPropertyHistory fetches property history from Repliers API
func (s *PropertyService) GetPropertyHistory(components AddressComponents) ([]PropertyHistory, error) {
	if s.repliersAPIKey == "" {
		return s.getFallbackHistory(), nil
	}

	url := fmt.Sprintf("https://api.repliers.io/listings/history?streetNumber=%s&streetName=%s&city=%s&zip=%s",
		components.StreetNumber, components.StreetName, components.City, components.Zip)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("REPLIERS-API-KEY", s.repliersAPIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return s.getFallbackHistory(), nil // Fallback on error
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return s.getFallbackHistory(), nil // Fallback on error
	}

	var history []PropertyHistory
	if err := json.NewDecoder(resp.Body).Decode(&history); err != nil {
		return s.getFallbackHistory(), nil // Fallback on error
	}

	return history, nil
}

// getFallbackEstimate returns simulated property data when API is unavailable
func (s *PropertyService) getFallbackEstimate(components AddressComponents) *PropertyEstimate {
	address := fmt.Sprintf("%s %s, %s, %s", 
		components.StreetNumber, components.StreetName, components.City, components.Zip)
	
	return &PropertyEstimate{
		Address:        address,
		Components:     components,
		EstimatedValue: 280000, // Simulated ARV
		RentEstimate:   1800,   // Simulated rent
		Bedrooms:       3,
		Bathrooms:      2,
		SquareFootage:  1200,
		YearBuilt:      1985,
		PropertyType:   "Single Family",
		Neighborhood:   "Residential",
		Comparables: []PropertyComp{
			{Address: "789 Pine St", Price: 275000, SqFt: 1150, Distance: "0.2 mi"},
			{Address: "321 Elm Rd", Price: 285000, SqFt: 1280, Distance: "0.3 mi"},
			{Address: "654 Birch Ave", Price: 270000, SqFt: 1200, Distance: "0.4 mi"},
		},
		History: s.getFallbackHistory(),
	}
}

// getFallbackHistory returns simulated historical data
func (s *PropertyService) getFallbackHistory() []PropertyHistory {
	return []PropertyHistory{
		{Date: "2023-08-15", Price: 250000, Event: "sold"},
		{Date: "2021-03-22", Price: 230000, Event: "sold"},
		{Date: "2019-07-10", Price: 210000, Event: "sold"},
		{Date: "2017-11-05", Price: 195000, Event: "sold"},
	}
}

// ValidateAddress performs basic address validation
func (s *PropertyService) ValidateAddress(components AddressComponents) bool {
	return components.StreetNumber != "" && 
		   components.StreetName != "" && 
		   components.City != "" && 
		   components.Zip != ""
}

// FormatAddress creates a formatted address string
func (s *PropertyService) FormatAddress(components AddressComponents) string {
	address := fmt.Sprintf("%s %s", components.StreetNumber, components.StreetName)
	if components.City != "" {
		address += fmt.Sprintf(", %s", components.City)
	}
	if components.State != "" {
		address += fmt.Sprintf(", %s", components.State)
	}
	if components.Zip != "" {
		address += fmt.Sprintf(" %s", components.Zip)
	}
	return address
}

// AddressSuggestion represents an address suggestion
type AddressSuggestion struct {
	Description   string `json:"description"`
	PlaceID       string `json:"place_id"`
	MainText      string `json:"main_text,omitempty"`
	SecondaryText string `json:"secondary_text,omitempty"`
}

// GetAddressSuggestions gets address autocomplete suggestions using Google Places API
func (s *PropertyService) GetAddressSuggestions(input string) ([]AddressSuggestion, error) {
	if s.googleMapsClient == nil {
		return s.getFallbackSuggestions(input), nil
	}

	request := &maps.PlaceAutocompleteRequest{
		Input:    input,
		Language: "en",
		Components: map[maps.Component][]string{
			maps.ComponentCountry: {"us"}, // Restrict to US
		},
	}

	response, err := s.googleMapsClient.PlaceAutocomplete(context.Background(), request)
	if err != nil {
		return s.getFallbackSuggestions(input), nil // Fallback on error
	}

	suggestions := make([]AddressSuggestion, 0, len(response.Predictions))
	for _, prediction := range response.Predictions {
		suggestion := AddressSuggestion{
			Description: prediction.Description,
			PlaceID:     prediction.PlaceID,
		}
		
		if len(prediction.StructuredFormatting.MainText) > 0 {
			suggestion.MainText = prediction.StructuredFormatting.MainText
			suggestion.SecondaryText = prediction.StructuredFormatting.SecondaryText
		}
		
		suggestions = append(suggestions, suggestion)
	}

	return suggestions, nil
}

// GeocodeAddress uses Google Geocoding API to get detailed address information
func (s *PropertyService) GeocodeAddress(address string) (*AddressComponents, error) {
	if s.googleMapsClient == nil {
		return s.parseAddressFallback(address), nil
	}

	request := &maps.GeocodingRequest{
		Address: address,
	}

	response, err := s.googleMapsClient.Geocode(context.Background(), request)
	if err != nil || len(response) == 0 {
		return s.parseAddressFallback(address), nil // Fallback on error
	}

	result := response[0]
	components := &AddressComponents{}

	// Parse address components
	for _, component := range result.AddressComponents {
		for _, componentType := range component.Types {
			switch componentType {
			case "street_number":
				components.StreetNumber = component.LongName
			case "route":
				components.StreetName = component.LongName
			case "locality":
				components.City = component.LongName
			case "postal_code":
				components.Zip = component.LongName
			case "administrative_area_level_1":
				components.State = component.ShortName
			}
		}
	}

	return components, nil
}

// SearchNearbyProperties searches for properties near a given location
func (s *PropertyService) SearchNearbyProperties(lat, lng float64, radius int) ([]PropertyEstimate, error) {
	if s.googleMapsClient == nil {
		return s.getFallbackNearbyProperties(), nil
	}

	request := &maps.NearbySearchRequest{
		Location: &maps.LatLng{Lat: lat, Lng: lng},
		Radius:   uint(radius),
		Type:     maps.PlaceTypeRealEstateAgency,
	}

	response, err := s.googleMapsClient.NearbySearch(context.Background(), request)
	if err != nil {
		return s.getFallbackNearbyProperties(), nil
	}

	properties := make([]PropertyEstimate, 0, len(response.Results))
	for _, place := range response.Results {
		// This would be enhanced with actual property data
		property := PropertyEstimate{
			Address: place.FormattedAddress,
			EstimatedValue: 250000 + int64(place.Rating*50000), // Simulated
			Neighborhood: place.Name,
		}
		properties = append(properties, property)
	}

	return properties, nil
}

// getFallbackSuggestions returns fallback suggestions when Google API is unavailable
func (s *PropertyService) getFallbackSuggestions(input string) []AddressSuggestion {
	return []AddressSuggestion{
		{Description: input + ", Denver, CO, USA", PlaceID: "fake1"},
		{Description: input + ", Boulder, CO, USA", PlaceID: "fake2"},
		{Description: input + ", Colorado Springs, CO, USA", PlaceID: "fake3"},
	}
}

// parseAddressFallback provides basic address parsing when Google API is unavailable
func (s *PropertyService) parseAddressFallback(address string) *AddressComponents {
	parts := strings.Split(address, ",")
	components := &AddressComponents{}
	
	if len(parts) > 0 {
		streetParts := strings.Fields(strings.TrimSpace(parts[0]))
		if len(streetParts) > 0 {
			components.StreetNumber = streetParts[0]
			if len(streetParts) > 1 {
				components.StreetName = strings.Join(streetParts[1:], " ")
			}
		}
	}
	
	if len(parts) > 1 {
		components.City = strings.TrimSpace(parts[1])
	}
	
	if len(parts) > 2 {
		stateZip := strings.TrimSpace(parts[2])
		stateZipParts := strings.Fields(stateZip)
		if len(stateZipParts) > 0 {
			components.State = stateZipParts[0]
		}
		if len(stateZipParts) > 1 {
			components.Zip = stateZipParts[1]
		}
	}
	
	return components
}

// getFallbackNearbyProperties returns fallback nearby properties
func (s *PropertyService) getFallbackNearbyProperties() []PropertyEstimate {
	return []PropertyEstimate{
		{Address: "123 Sample St, Denver, CO", EstimatedValue: 275000, Neighborhood: "Downtown"},
		{Address: "456 Example Ave, Denver, CO", EstimatedValue: 285000, Neighborhood: "Highlands"},
		{Address: "789 Demo Dr, Denver, CO", EstimatedValue: 265000, Neighborhood: "Capitol Hill"},
	}
}