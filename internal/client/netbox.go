package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/jaepetto/nbjrnlctl/internal/models"
)

// NetboxClient handles communication with the Netbox API
type NetboxClient struct {
	BaseURL    string
	APIToken   string
	HTTPClient *http.Client
}

// extractKindFromDisplay extracts the kind/type from the display field
// The display field in Netbox journal entries typically follows patterns like:
// "2025-07-29 08:00 (Info)" - where the kind is in parentheses
// This function attempts to extract the kind/type portion
func extractKindFromDisplay(display string) string {
	if display == "" {
		return ""
	}

	// Pattern: "date time (Kind)" - extract content between parentheses
	if openParen := strings.LastIndex(display, "("); openParen != -1 {
		if closeParen := strings.LastIndex(display, ")"); closeParen != -1 && closeParen > openParen {
			return strings.TrimSpace(display[openParen+1 : closeParen])
		}
	}

	// Fallback patterns for other formats
	// Pattern 1: "Kind: description" or "Kind - description"
	if colonIndex := strings.Index(display, ":"); colonIndex != -1 {
		return strings.TrimSpace(display[:colonIndex])
	}

	if dashIndex := strings.Index(display, " - "); dashIndex != -1 {
		// Check if this looks like a date pattern (likely not the kind)
		parts := strings.Split(display[:dashIndex], "-")
		if len(parts) == 3 && len(parts[0]) == 4 {
			// This looks like a date, so skip it and look for kind in the rest
			rest := strings.TrimSpace(display[dashIndex+3:])
			if colonIndex := strings.Index(rest, ":"); colonIndex != -1 {
				return strings.TrimSpace(rest[:colonIndex])
			}
			return rest
		}
		return strings.TrimSpace(display[:dashIndex])
	}

	// If no clear separator, return the first word as the kind
	parts := strings.Fields(display)
	if len(parts) > 0 {
		return parts[0]
	}

	return display
}

// NewNetboxClient creates a new client for Netbox API
func NewNetboxClient(baseURL, apiToken string) *NetboxClient {
	return &NetboxClient{
		BaseURL:  baseURL,
		APIToken: apiToken,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// GetDeviceIDByName retrieves a device ID by its name
func (c *NetboxClient) GetDeviceIDByName(name string) (int, error) {
	var nextURL string = fmt.Sprintf("%s/api/dcim/devices/", c.BaseURL)

	for nextURL != "" {
		req, err := http.NewRequest("GET", nextURL, nil)
		if err != nil {
			return 0, err
		}

		req.Header.Set("Authorization", fmt.Sprintf("Token %s", c.APIToken))

		resp, err := c.HTTPClient.Do(req)
		if err != nil {
			return 0, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return 0, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		}

		var result struct {
			Count    int    `json:"count"`
			Next     string `json:"next"`
			Previous string `json:"previous"`
			Results  []struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
			} `json:"results"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return 0, err
		}

		// Check for a device match in current page
		for _, device := range result.Results {
			if device.Name == name {
				return device.ID, nil
			}
		}

		// Set next URL for pagination
		nextURL = result.Next
	}

	return 0, fmt.Errorf("no device found with name: %s", name)
}

// CreateJournalEntry creates a new journal entry for a device
func (c *NetboxClient) CreateJournalEntry(deviceID int, entry models.JournalEntry) (*models.JournalEntry, error) {
	url := fmt.Sprintf("%s/api/dcim/devices/%d/journal/", c.BaseURL, deviceID)

	reqBody, err := json.Marshal(entry)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Token %s", c.APIToken))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result models.JournalEntry
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetJournalEntry retrieves a specific journal entry
func (c *NetboxClient) GetJournalEntry(journalID int) (*models.JournalEntry, error) {
	url := fmt.Sprintf("%s/api/extras/journal-entries/%d/", c.BaseURL, journalID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Token %s", c.APIToken))

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result models.JournalEntry
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// ListJournalEntries lists all journal entries for a device
func (c *NetboxClient) ListJournalEntries(deviceID int) ([]models.JournalEntry, error) {
	url := fmt.Sprintf("%s/graphql/", c.BaseURL)

	// Create GraphQL query to fetch device journal entries by device name
	query := `
	{
		device_list(filters: {id: {exact: $deviceID}}) {
			id
			name
			journal_entries {
				id
				created
				display
				comments
				created_by {
					username
				}
			}
		}
	}
	`
	// Replace placeholder with actual device ID
	query = strings.Replace(query, "$deviceID", fmt.Sprintf("%d", deviceID), 1)

	// Create request body
	reqBody, err := json.Marshal(map[string]string{
		"query": query,
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Token %s", c.APIToken))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Read the full response body into a string for inspection
	var bodyBytes []byte
	bodyBytes, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	var result struct {
		Data struct {
			DeviceList []struct {
				ID             string `json:"id"`
				Name           string `json:"name"`
				JournalEntries []struct {
					ID        string `json:"id"`
					Created   string `json:"created"`
					Display   string `json:"display"`
					Comments  string `json:"comments"`
					CreatedBy struct {
						Username string `json:"username"`
					} `json:"created_by"`
				} `json:"journal_entries"`
			} `json:"device_list"`
		} `json:"data"`
	}

	err = json.Unmarshal(bodyBytes, &result)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %w", err)
	}

	if len(result.Data.DeviceList) == 0 {
		return []models.JournalEntry{}, nil
	}

	// Convert the parsed journal entries to the models.JournalEntry type
	entries := make([]models.JournalEntry, len(result.Data.DeviceList[0].JournalEntries))
	for i, entry := range result.Data.DeviceList[0].JournalEntries {
		// Parse the created time string into time.Time
		createdTime, err := time.Parse(time.RFC3339, entry.Created)
		if err != nil {
			// If there's an error parsing time, use current time as fallback
			createdTime = time.Now()
		}

		// Convert ID from string to int
		var entryID int
		fmt.Sscanf(entry.ID, "%d", &entryID)

		// Extract kind from display field
		// The display field typically contains the kind information in Netbox
		kind := extractKindFromDisplay(entry.Display)

		entries[i] = models.JournalEntry{
			ID:        entryID,
			Comments:  entry.Comments,
			CreatedBy: entry.CreatedBy.Username,
			Created:   createdTime,
			Kind:      kind,
		}
	}

	return entries, nil
}

// UpdateJournalEntry updates a journal entry
func (c *NetboxClient) UpdateJournalEntry(journalID int, entry models.JournalEntry) (*models.JournalEntry, error) {
	url := fmt.Sprintf("%s/api/extras/journal-entries/%d/", c.BaseURL, journalID)

	reqBody, err := json.Marshal(entry)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Token %s", c.APIToken))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result models.JournalEntry
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// DeleteJournalEntry deletes a journal entry
func (c *NetboxClient) DeleteJournalEntry(journalID int) error {
	url := fmt.Sprintf("%s/api/extras/journal-entries/%d/", c.BaseURL, journalID)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Token %s", c.APIToken))

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
