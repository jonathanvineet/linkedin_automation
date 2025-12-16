package search

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"github.com/jonathanvineet/linkedin-automation/internal/logger"
)

// PeopleSearch handles LinkedIn people search
type PeopleSearch struct {
	page *rod.Page
}

// SearchCriteria defines search parameters
type SearchCriteria struct {
	Keywords string
	Location string
	Company  string
	JobTitle string
}

// ProfileResult represents a search result
type ProfileResult struct {
	Name       string
	Title      string
	Company    string
	Location   string
	ProfileURL string
}

// NewPeopleSearch creates a new search handler
func NewPeopleSearch(page *rod.Page) *PeopleSearch {
	return &PeopleSearch{
		page: page,
	}
}

// Search performs a people search with given criteria
func (ps *PeopleSearch) Search(criteria SearchCriteria, maxResults int) ([]ProfileResult, error) {
	logger.Log.WithField("criteria", fmt.Sprintf("%+v", criteria)).Info("Starting LinkedIn people search")
	
	// Build search URL
	searchURL := ps.buildSearchURL(criteria)
	
	// Navigate to search page
	err := ps.page.Navigate(searchURL)
	if err != nil {
		return nil, fmt.Errorf("failed to navigate to search: %w", err)
	}
	
	// Wait for page load
	time.Sleep(3 * time.Second)
	err = ps.page.WaitLoad()
	if err != nil {
		return nil, fmt.Errorf("failed to load search page: %w", err)
	}
	
	logger.Log.Info("Search page loaded successfully")
	
	// Extract profile results
	results, err := ps.extractProfiles(maxResults)
	if err != nil {
		return nil, fmt.Errorf("failed to extract profiles: %w", err)
	}
	
	logger.Log.WithField("count", len(results)).Info("Search completed")
	return results, nil
}

// buildSearchURL constructs LinkedIn search URL from criteria
func (ps *PeopleSearch) buildSearchURL(criteria SearchCriteria) string {
	baseURL := "https://www.linkedin.com/search/results/people/?"
	
	params := []string{}
	
	if criteria.Keywords != "" {
		params = append(params, fmt.Sprintf("keywords=%s", criteria.Keywords))
	}
	
	if criteria.Location != "" {
		params = append(params, fmt.Sprintf("location=%s", criteria.Location))
	}
	
	if criteria.Company != "" {
		params = append(params, fmt.Sprintf("company=%s", criteria.Company))
	}
	
	if criteria.JobTitle != "" {
		params = append(params, fmt.Sprintf("title=%s", criteria.JobTitle))
	}
	
	return baseURL + strings.Join(params, "&")
}

// extractProfiles extracts profile information from search results
func (ps *PeopleSearch) extractProfiles(maxResults int) ([]ProfileResult, error) {
	results := make([]ProfileResult, 0)
	
	// Find all search result containers
	elements, err := ps.page.Timeout(10 * time.Second).Elements(".reusable-search__result-container")
	if err != nil {
		logger.Log.Warn("Could not find search results")
		return results, nil
	}
	
	count := 0
	for _, element := range elements {
		if count >= maxResults {
			break
		}
		
		profile := ps.extractProfileFromElement(element)
		if profile.ProfileURL != "" {
			results = append(results, profile)
			count++
		}
	}
	
	return results, nil
}

// extractProfileFromElement extracts profile data from a single result element
func (ps *PeopleSearch) extractProfileFromElement(element *rod.Element) ProfileResult {
	profile := ProfileResult{}
	
	// Extract name
	nameElement, err := element.Element(".entity-result__title-text a")
	if err == nil {
		profile.Name, _ = nameElement.Text()
		href, _ := nameElement.Attribute("href")
		if href != nil {
			profile.ProfileURL = *href
		}
	}
	
	// Extract title
	titleElement, err := element.Element(".entity-result__primary-subtitle")
	if err == nil {
		profile.Title, _ = titleElement.Text()
	}
	
	// Extract company
	companyElement, err := element.Element(".entity-result__secondary-subtitle")
	if err == nil {
		profile.Company, _ = companyElement.Text()
	}
	
	// Extract location
	locationElement, err := element.Element(".entity-result__location")
	if err == nil {
		profile.Location, _ = locationElement.Text()
	}
	
	// Clean profile URL
	if profile.ProfileURL != "" {
		// Remove query parameters
		if idx := strings.Index(profile.ProfileURL, "?"); idx > 0 {
			profile.ProfileURL = profile.ProfileURL[:idx]
		}
	}
	
	return profile
}

// ScrollToLoadMore scrolls down to load more results
func (ps *PeopleSearch) ScrollToLoadMore() error {
	// Scroll down gradually
	for i := 0; i < 5; i++ {
		err := ps.page.Mouse.Scroll(0, 500, 10)
		if err != nil {
			return err
		}
		time.Sleep(time.Duration(500+i*200) * time.Millisecond)
	}
	return nil
}

// HasNextPage checks if there are more pages
func (ps *PeopleSearch) HasNextPage() bool {
	_, err := ps.page.Timeout(2 * time.Second).Element("button[aria-label='Next']")
	return err == nil
}

// GoToNextPage navigates to the next page of results
func (ps *PeopleSearch) GoToNextPage() error {
	nextButton, err := ps.page.Element("button[aria-label='Next']")
	if err != nil {
		return fmt.Errorf("no next page available")
	}
	
	// Human-like delay before clicking
	time.Sleep(time.Duration(1000+time.Now().UnixNano()%1000) * time.Millisecond)
	
	err = nextButton.Click(proto.InputMouseButtonLeft, 1)
	if err != nil {
		return err
	}
	
	// Wait for page load
	time.Sleep(3 * time.Second)
	return ps.page.WaitLoad()
}
