package http

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"tax-helper/internal/logger"
	"time"
)

type Currency struct {
	Code          string    `json:"code"`
	Quantity      int       `json:"quantity"`
	RateFormated  string    `json:"rateFormated"`
	DiffFormated  string    `json:"diffFormated"`
	Rate          float64   `json:"rate"`
	Name          string    `json:"name"`
	Diff          float64   `json:"diff"`
	Date          time.Time `json:"date"`
	ValidFromDate time.Time `json:"validFromDate"`
}

type CurrencyRateResponse struct {
	Date       time.Time  `json:"date"`
	Currencies []Currency `json:"currencies"`
}

type Client struct {
	client  *http.Client
	baseURL string
	logger  logger.Logger
}

func NewClient(baseURL string, timeout time.Duration, l logger.Logger) *Client {
	return &Client{
		client: &http.Client{
			Timeout: timeout,
		},
		baseURL: baseURL,
		logger:  l,
	}
}

func (c *Client) GetCurrencyRates(ctx context.Context, path string) ([]CurrencyRateResponse, error) {
	c.logger.Debugf("Fetching currency rates from path=%q", path)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+path, nil)
	if err != nil {
		c.logger.Error(err)
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		c.logger.Error(err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.logger.Error(err)
		return nil, err
	}

	var rates []CurrencyRateResponse
	if err = json.Unmarshal(body, &rates); err != nil {
		c.logger.Errorf("failed to unmarshal currency rates: %v", err)
		return nil, err
	}
	ratesJSON, _ := json.Marshal(rates)
	c.logger.Debugf("Fetched currency rates: %s", ratesJSON)
	return rates, nil
}
