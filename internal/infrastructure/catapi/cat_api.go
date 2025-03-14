package catapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// CatAPI describes methods of interaction with thecatapi.
type CatAPI interface {
	IsBreedValid(ctx context.Context, breedName string) (bool, error)
}

type catAPI struct {
	httpClient *http.Client
	apiURL     string
	apiKey     string
}

// NewCatAPI Creates an instance for working with thecatapi.
func NewCatAPI(apiURL, apiKey string) CatAPI {
	return &catAPI{
		httpClient: &http.Client{Timeout: 10 * time.Second},
		apiURL:     apiURL,
		apiKey:     apiKey,
	}
}

func (c *catAPI) IsBreedValid(ctx context.Context, breedName string) (bool, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.apiURL+"/v1/breeds", nil)
	if err != nil {
		return false, err
	}

	if c.apiKey != "" {
		req.Header.Set("x-api-key", c.apiKey)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("catapi: unexpected status code %d", resp.StatusCode)
	}

	var data []struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return false, err
	}

	for _, b := range data {
		if b.Name == breedName {
			return true, nil
		}
	}
	return false, nil
}
