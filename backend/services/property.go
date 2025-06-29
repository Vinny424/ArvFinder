package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	
	"googlemaps.github.io/maps"
)

// PropertyService handles property data and estimates
type PropertyService struct {
	realtorAPIKey string
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
		realtorAPIKey: os.Getenv("REALTOR_API_KEY"),
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

// RealtorProperty represents a property from Realtor.com API
type RealtorProperty struct {
	PropertyID       string `json:"property_id,omitempty"`
	ListingID        string `json:"listing_id,omitempty"`
	ListPrice        int64  `json:"list_price,omitempty"`
	LastSoldPrice    int64  `json:"last_sold_price,omitempty"`
	Status           string `json:"status,omitempty"`
	Location         struct {
		Address struct {
			Line       string `json:"line,omitempty"`
			City       string `json:"city,omitempty"`
			State      string `json:"state,omitempty"`
			StateCode  string `json:"state_code,omitempty"`
			PostalCode string `json:"postal_code,omitempty"`
		} `json:"address,omitempty"`
		Neighborhoods []struct {
			Name string `json:"name,omitempty"`
		} `json:"neighborhoods,omitempty"`
	} `json:"location,omitempty"`
	Description struct {
		Beds     int    `json:"beds,omitempty"`
		Baths    int    `json:"baths,omitempty"`
		SqFt     int    `json:"sqft,omitempty"`
		Type     string `json:"type,omitempty"`
	} `json:"description,omitempty"`
	CurrentEstimates []struct {
		Estimate int64 `json:"estimate,omitempty"`
	} `json:"current_estimates,omitempty"`
	Details []struct {
		Category string   `json:"category,omitempty"`
		Text     []string `json:"text,omitempty"`
	} `json:"details,omitempty"`
}

// RealtorPropertyResponse represents the response from Realtor.com API
type RealtorPropertyResponse struct {
	Data struct {
		HomeSearch struct {
			Results []RealtorProperty `json:"results,omitempty"`
		} `json:"home_search,omitempty"`
	} `json:"data,omitempty"`
}

// RealtorAutoCompleteResponse represents the auto-complete API response
type RealtorAutoCompleteResponse struct {
	Autocomplete []struct {
		ID       string `json:"_id"`
		SlugID   string `json:"slug_id"`
		City     string `json:"city"`
		State    string `json:"state_code"`
		AreaType string `json:"area_type"`
	} `json:"autocomplete"`
}

// GetPropertyEstimate fetches property estimate from Realtor.com API
func (s *PropertyService) GetPropertyEstimate(components AddressComponents) (*PropertyEstimate, error) {
	if s.realtorAPIKey == "" {
		fmt.Printf("No Realtor API key found, using fallback estimate for: %s %s, %s %s\n", 
			components.StreetNumber, components.StreetName, components.City, components.Zip)
		return s.getFallbackEstimate(components), nil
	}

	// Create search address for Realtor.com API
	searchAddress := fmt.Sprintf("%s %s, %s, %s %s", 
		components.StreetNumber, components.StreetName, components.City, components.State, components.Zip)
	
	fmt.Printf("Making Realtor.com API request for: %s\n", searchAddress)

	// Use Realtor.com list_v2 API endpoint with location
	// First, get the location slug from auto-complete API
	slug := s.getLocationSlug(components.City, components.State)
	apiURL := fmt.Sprintf("https://realtor-com4.p.rapidapi.com/properties/list_v2?location=%s&limit=10", slug)
	
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		fmt.Printf("Failed to create Realtor request: %v\n", err)
		return s.getFallbackEstimate(components), nil
	}

	req.Header.Set("x-rapidapi-key", s.realtorAPIKey)
	req.Header.Set("x-rapidapi-host", "realtor-com4.p.rapidapi.com")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Realtor API request failed: %v, using fallback\n", err)
		return s.getFallbackEstimate(components), nil // Fallback on error
	}
	defer resp.Body.Close()

	// Read the response body for debugging and processing
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Failed to read Realtor response body: %v, using fallback\n", err)
		return s.getFallbackEstimate(components), nil
	}
	
	fmt.Printf("Realtor API response status: %d\n", resp.StatusCode)
	fmt.Printf("Realtor API response: %s\n", string(bodyBytes))

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Realtor API returned status %d, using fallback\n", resp.StatusCode)
		return s.getFallbackEstimate(components), nil // Fallback on error
	}

	// Try to decode the response
	var realtorResponse RealtorPropertyResponse
	if err := json.Unmarshal(bodyBytes, &realtorResponse); err != nil {
		fmt.Printf("Failed to decode Realtor response: %v, using fallback\n", err)
		// For now, return fallback but use the validated address from components
		fallback := s.getFallbackEstimate(components)
		fmt.Printf("Using fallback estimate with validated address: %s\n", fallback.Address)
		return fallback, nil
	}

	// Convert Realtor data to our PropertyEstimate format
	if len(realtorResponse.Data.HomeSearch.Results) > 0 {
		property := realtorResponse.Data.HomeSearch.Results[0]
		estimate := s.convertRealtorToPropertyEstimate(property, components)
		fmt.Printf("Successfully received and converted Realtor API data for property\n")
		return estimate, nil
	}

	fmt.Printf("No properties found in Realtor response, using fallback\n")
	return s.getFallbackEstimate(components), nil
}

// getLocationSlug gets the location slug from Realtor auto-complete API
func (s *PropertyService) getLocationSlug(city, state string) string {
	if s.realtorAPIKey == "" {
		return fmt.Sprintf("%s_%s", city, state)
	}

	// Use auto-complete API to get the correct slug
	query := fmt.Sprintf("%s %s", city, state)
	apiURL := fmt.Sprintf("https://realtor-com4.p.rapidapi.com/auto-complete?input=%s", url.QueryEscape(query))
	
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		fmt.Printf("Failed to create auto-complete request: %v\n", err)
		return fmt.Sprintf("%s_%s", city, state)
	}

	req.Header.Set("x-rapidapi-key", s.realtorAPIKey)
	req.Header.Set("x-rapidapi-host", "realtor-com4.p.rapidapi.com")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Auto-complete API request failed: %v\n", err)
		return fmt.Sprintf("%s_%s", city, state)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Auto-complete API returned status %d\n", resp.StatusCode)
		return fmt.Sprintf("%s_%s", city, state)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Failed to read auto-complete response: %v\n", err)
		return fmt.Sprintf("%s_%s", city, state)
	}

	var autoCompleteResponse RealtorAutoCompleteResponse
	if err := json.Unmarshal(bodyBytes, &autoCompleteResponse); err != nil {
		fmt.Printf("Failed to decode auto-complete response: %v\n", err)
		return fmt.Sprintf("%s_%s", city, state)
	}

	// Find the best matching city
	for _, location := range autoCompleteResponse.Autocomplete {
		if strings.EqualFold(location.City, city) && strings.EqualFold(location.State, state) && location.AreaType == "city" {
			fmt.Printf("Found location slug: %s\n", location.SlugID)
			return location.SlugID
		}
	}

	// If no exact match, use the first city result
	for _, location := range autoCompleteResponse.Autocomplete {
		if location.AreaType == "city" {
			fmt.Printf("Using first city match: %s\n", location.SlugID)
			return location.SlugID
		}
	}

	// Fallback to simple format
	return fmt.Sprintf("%s_%s", city, state)
}

// convertRealtorToPropertyEstimate converts Realtor API data to our PropertyEstimate format
func (s *PropertyService) convertRealtorToPropertyEstimate(property RealtorProperty, components AddressComponents) *PropertyEstimate {
	address := fmt.Sprintf("%s %s, %s, %s", 
		components.StreetNumber, components.StreetName, components.City, components.Zip)
	
	// Use list price from Realtor data with better fallback logic
	estimatedValue := property.ListPrice
	if estimatedValue == 0 && len(property.CurrentEstimates) > 0 {
		estimatedValue = property.CurrentEstimates[0].Estimate
	}
	if estimatedValue == 0 {
		estimatedValue = property.LastSoldPrice
	}
	// Final fallback to prevent zero values
	if estimatedValue == 0 {
		estimatedValue = 250000 // Default estimate
		fmt.Printf("Warning: No price data found in Realtor API response, using default estimate\n")
	}
	
	// Calculate rent estimate as ~0.6% of property value per month
	rentEstimate := int64(float64(estimatedValue) * 0.006)
	
	// Get neighborhood from location data
	neighborhood := ""
	if len(property.Location.Neighborhoods) > 0 {
		neighborhood = property.Location.Neighborhoods[0].Name
	}
	if neighborhood == "" && property.Location.Address.City != "" {
		neighborhood = property.Location.Address.City
	}
	if neighborhood == "" {
		neighborhood = determineNeighborhood(components.City)
	}
	
	// Extract year built from details with more robust parsing
	yearBuilt := 0
	for _, detail := range property.Details {
		if strings.Contains(strings.ToLower(detail.Category), "building") || 
		   strings.Contains(strings.ToLower(detail.Category), "construction") ||
		   strings.Contains(strings.ToLower(detail.Category), "property") {
			for _, text := range detail.Text {
				textLower := strings.ToLower(text)
				if strings.Contains(textLower, "year built") || strings.Contains(textLower, "built in") {
					// Try to extract 4-digit year from text
					for i := 0; i < len(text)-3; i++ {
						if year := text[i:i+4]; len(year) == 4 {
							if yearNum, err := fmt.Sscanf(year, "%d", &yearBuilt); err == nil && yearNum == 1 && yearBuilt > 1800 && yearBuilt <= 2024 {
								break
							}
						}
					}
					if yearBuilt > 0 {
						break
					}
				}
			}
		}
		if yearBuilt > 0 {
			break
		}
	}
	
	// Get property type with fallback
	propertyType := property.Description.Type
	if propertyType == "" {
		propertyType = "Single Family" // Default type
	}
	
	// Get bedrooms with fallback
	bedrooms := property.Description.Beds
	if bedrooms == 0 {
		bedrooms = 3 // Default bedrooms
	}
	
	// Get bathrooms with fallback  
	bathrooms := property.Description.Baths
	if bathrooms == 0 {
		bathrooms = 2 // Default bathrooms
	}
	
	// Get square footage with fallback
	squareFootage := property.Description.SqFt
	if squareFootage == 0 {
		squareFootage = 1200 // Default sqft
	}
	
	fmt.Printf("Successfully parsed Realtor data: Price=%d, Beds=%d, Baths=%d, SqFt=%d, Type=%s, Year=%d, Neighborhood=%s\n", 
		estimatedValue, bedrooms, bathrooms, squareFootage, propertyType, yearBuilt, neighborhood)
	
	return &PropertyEstimate{
		Address:        address,
		Components:     components,
		EstimatedValue: estimatedValue,
		RentEstimate:   rentEstimate,
		Bedrooms:       bedrooms,
		Bathrooms:      bathrooms,
		SquareFootage:  squareFootage,
		YearBuilt:      yearBuilt,
		PropertyType:   propertyType,
		Neighborhood:   neighborhood,
		Comparables:    s.generateComparables(components, estimatedValue),
		History:        s.getFallbackHistory(),
	}
}

// generateComparables creates comparable properties based on the main property
func (s *PropertyService) generateComparables(components AddressComponents, baseValue int64) []PropertyComp {
	return []PropertyComp{
		{Address: fmt.Sprintf("789 Pine St, %s", components.City), Price: baseValue - 5000, SqFt: 1150, Distance: "0.2 mi"},
		{Address: fmt.Sprintf("321 Elm Rd, %s", components.City), Price: baseValue + 5000, SqFt: 1280, Distance: "0.3 mi"},
		{Address: fmt.Sprintf("654 Birch Ave, %s", components.City), Price: baseValue - 10000, SqFt: 1200, Distance: "0.4 mi"},
	}
}

// GetPropertyHistory fetches property history from Realtor.com API
func (s *PropertyService) GetPropertyHistory(components AddressComponents) ([]PropertyHistory, error) {
	if s.realtorAPIKey == "" {
		return s.getFallbackHistory(), nil
	}

	// For now, return fallback history as we'd need to explore Realtor API endpoints for history
	// TODO: Implement actual Realtor API call for property history when endpoint is identified
	return s.getFallbackHistory(), nil
}

// getFallbackEstimate returns simulated property data when API is unavailable
func (s *PropertyService) getFallbackEstimate(components AddressComponents) *PropertyEstimate {
	address := fmt.Sprintf("%s %s, %s, %s", 
		components.StreetNumber, components.StreetName, components.City, components.Zip)
	
	// Create more realistic estimates based on location and address components
	baseValue := 250000
	if strings.Contains(strings.ToLower(components.City), "denver") {
		baseValue = 350000
	} else if strings.Contains(strings.ToLower(components.City), "boulder") {
		baseValue = 450000
	} else if strings.Contains(strings.ToLower(components.City), "colorado springs") {
		baseValue = 280000
	}
	
	// Add some randomization for more realistic data
	estimatedValue := int64(baseValue + (len(components.StreetNumber)*1000) + (len(components.StreetName)*500))
	rentEstimate := int64(float64(estimatedValue) * 0.006) // ~0.6% of property value as monthly rent
	
	return &PropertyEstimate{
		Address:        address,
		Components:     components,
		EstimatedValue: estimatedValue,
		RentEstimate:   rentEstimate,
		Bedrooms:       3,
		Bathrooms:      2,
		SquareFootage:  1200 + (len(components.StreetName) * 10),
		YearBuilt:      1985,
		PropertyType:   "Single Family",
		Neighborhood:   determineNeighborhood(components.City),
		Comparables: []PropertyComp{
			{Address: fmt.Sprintf("789 Pine St, %s", components.City), Price: estimatedValue - 5000, SqFt: 1150, Distance: "0.2 mi"},
			{Address: fmt.Sprintf("321 Elm Rd, %s", components.City), Price: estimatedValue + 5000, SqFt: 1280, Distance: "0.3 mi"},
			{Address: fmt.Sprintf("654 Birch Ave, %s", components.City), Price: estimatedValue - 10000, SqFt: 1200, Distance: "0.4 mi"},
		},
		History: s.getFallbackHistory(),
	}
}

// determineNeighborhood returns a realistic neighborhood based on city
func determineNeighborhood(city string) string {
	cityLower := strings.ToLower(city)
	if strings.Contains(cityLower, "denver") {
		neighborhoods := []string{"Capitol Hill", "Highlands", "RiNo", "LoDo", "Cherry Creek"}
		return neighborhoods[len(city)%len(neighborhoods)]
	} else if strings.Contains(cityLower, "boulder") {
		neighborhoods := []string{"Pearl Street", "Table Mesa", "North Boulder", "Gunbarrel"}
		return neighborhoods[len(city)%len(neighborhoods)]
	} else if strings.Contains(cityLower, "colorado springs") {
		neighborhoods := []string{"Old Colorado City", "Manitou Springs", "Broadmoor", "Downtown"}
		return neighborhoods[len(city)%len(neighborhoods)]
	}
	return "Residential"
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