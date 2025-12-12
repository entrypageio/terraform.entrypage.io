package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type apiClient struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

type apiError struct {
	StatusCode int
	Body       string
}

func (e *apiError) Error() string {
	return fmt.Sprintf("entrypage API error: status %d: %s", e.StatusCode, e.Body)
}

func newClient(apiKey, baseURL string) *apiClient {
	return &apiClient{
		apiKey:  apiKey,
		baseURL: strings.TrimRight(baseURL, "/"),
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

type userRequest struct {
	PreferredUsername string `json:"preferredUsername"`
	Name              string `json:"name"`
}

type userResponse struct {
	Sub string `json:"sub"`
}

type userSummaryResponse struct {
	Sub               string  `json:"sub"`
	Active            *bool   `json:"active,omitempty"`
	PreferredUsername string  `json:"preferredUsername"`
	Email             *string `json:"email,omitempty"`
	Picture           *string `json:"picture,omitempty"`
	IdentityProvider  string  `json:"identityProvider"`
}

// CreateUser issues a create call and returns the subject identifier.
func (c *apiClient) CreateUser(ctx context.Context, domain, preferredUsername, name string) (string, error) {
	payload := userRequest{
		PreferredUsername: preferredUsername,
		Name:              name,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	path := fmt.Sprintf("%s/v1/domain/%s/user", c.baseURL, url.PathEscape(domain))
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, path, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("x-api-key", c.apiKey)

	respBody, status, err := c.doRequest(req)
	if status == http.StatusConflict {
		// If the user already exists, try to find and adopt it to keep creates idempotent.
		existing, findErr := c.FindUserByPreferredUsername(ctx, domain, preferredUsername)
		if findErr == nil && existing != nil && existing.Sub != "" {
			return existing.Sub, nil
		}
		if findErr != nil {
			return "", fmt.Errorf("user already exists but lookup failed: %w", findErr)
		}
		return "", &apiError{StatusCode: status, Body: string(respBody)}
	}
	if err != nil {
		return "", err
	}
	if status != http.StatusOK && status != http.StatusCreated {
		return "", &apiError{StatusCode: status, Body: string(respBody)}
	}

	var resp userResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return "", err
	}
	if resp.Sub == "" {
		return "", fmt.Errorf("entrypage API did not return a user id")
	}
	return resp.Sub, nil
}

// FindUser queries users using the search filter and returns the user with the requested sub.
func (c *apiClient) FindUser(ctx context.Context, domain, sub string) (*userSummaryResponse, error) {
	params := url.Values{}
	params.Set("search", sub)
	params.Set("limit", "200")

	path := fmt.Sprintf("%s/v1/domain/%s/user?%s", c.baseURL, url.PathEscape(domain), params.Encode())
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("x-api-key", c.apiKey)

	body, status, err := c.doRequest(req)
	if err != nil {
		if apiErr, ok := err.(*apiError); ok && apiErr.StatusCode == http.StatusNotFound {
			return nil, nil
		}
		return nil, err
	}
	if status == http.StatusNotFound {
		return nil, nil
	}
	if status != http.StatusOK {
		return nil, &apiError{StatusCode: status, Body: string(body)}
	}

	var users []userSummaryResponse
	if err := json.Unmarshal(body, &users); err != nil {
		return nil, err
	}
	for _, u := range users {
		if u.Sub == sub {
			return &u, nil
		}
	}
	return nil, nil
}

// FindUserByPreferredUsername searches by preferred username and returns the first match.
func (c *apiClient) FindUserByPreferredUsername(ctx context.Context, domain, preferredUsername string) (*userSummaryResponse, error) {
	params := url.Values{}
	params.Set("search", preferredUsername)
	params.Set("limit", "200")

	path := fmt.Sprintf("%s/v1/domain/%s/user?%s", c.baseURL, url.PathEscape(domain), params.Encode())
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("x-api-key", c.apiKey)

	body, status, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	if status == http.StatusNotFound {
		return nil, nil
	}
	if status != http.StatusOK {
		return nil, &apiError{StatusCode: status, Body: string(body)}
	}

	var users []userSummaryResponse
	if err := json.Unmarshal(body, &users); err != nil {
		return nil, err
	}
	for _, u := range users {
		if strings.EqualFold(u.PreferredUsername, preferredUsername) {
			return &u, nil
		}
	}
	return nil, nil
}

// DeleteUser removes a user by subject identifier.
func (c *apiClient) DeleteUser(ctx context.Context, domain, sub string) error {
	path := fmt.Sprintf("%s/v1/domain/%s/user/%s", c.baseURL, url.PathEscape(domain), url.PathEscape(sub))
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("x-api-key", c.apiKey)

	body, status, err := c.doRequest(req)
	if err != nil {
		if apiErr, ok := err.(*apiError); ok && apiErr.StatusCode == http.StatusNotFound {
			return nil
		}
		return err
	}
	if status == http.StatusNotFound {
		return nil
	}
	if status != http.StatusNoContent {
		return &apiError{StatusCode: status, Body: string(body)}
	}
	return nil
}

func (c *apiClient) doRequest(req *http.Request) ([]byte, int, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return nil, resp.StatusCode, readErr
	}

	if resp.StatusCode >= 400 {
		return body, resp.StatusCode, &apiError{StatusCode: resp.StatusCode, Body: string(body)}
	}

	return body, resp.StatusCode, nil
}
