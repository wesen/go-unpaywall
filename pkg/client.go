package pkg

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"net/url"
)

// See the API endpoints at https://unpaywall.org/products/api

type Client struct {
	BaseURL string
	Email   string
}

type ClientOption func(*Client)

func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) {
		c.BaseURL = baseURL
	}
}

func WithEmail(email string) ClientOption {
	return func(c *Client) {
		c.Email = email
	}
}

func NewClient(opts ...ClientOption) *Client {
	c := &Client{
		BaseURL: "https://api.unpaywall.org",
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (c *Client) GetDOI(ctx context.Context, doi string) (*DOI, error) {
	if c.Email == "" {
		return nil, errors.New("email is required")
	}

	req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		fmt.Sprintf("%s/v2/%s?email=%s", c.BaseURL, doi, c.Email), nil)
	if err != nil {
		return nil, err
	}
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// read body to string
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	reader := bytes.NewReader(body)

	// Decode the JSON response into a DOIObject struct.
	var doiObject DOI
	if err := json.NewDecoder(reader).Decode(&doiObject); err != nil {
		return nil, err
	}

	return &doiObject, nil
}

type SearchResult struct {
	Response interface{} `json:"response"`
	Score    float64     `json:"score"`
	Snippet  string      `json:"snippet"`
}

type SearchResults []SearchResult

type SearchResponse struct {
	Results SearchResults `json:"results"`
}

type SearchRequest struct {
	Query string `json:"query"`
	IsOA  *bool  `json:"is_oa"`
	Page  *int   `json:"page"`
}

func (c *Client) Search(ctx context.Context, searchRequest SearchRequest) (SearchResults, error) {
	// create the API endpoint URL
	u, err := url.Parse(fmt.Sprintf("%s/v2/search", c.BaseURL))
	if err != nil {
		return nil, err
	}

	// set query parameters
	q := u.Query()
	q.Set("query", searchRequest.Query)
	q.Set("email", c.Email)
	if searchRequest.IsOA != nil {
		q.Set("is_oa", fmt.Sprintf("%t", *searchRequest.IsOA))
	}
	if searchRequest.Page != nil {
		q.Set("page", fmt.Sprintf("%d", *searchRequest.Page))
	}
	u.RawQuery = q.Encode()

	// create a new HTTP client
	client := &http.Client{}

	// create a new HTTP request with the URL and context
	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	// send the request and retrieve the response
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// read body to string
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// decode the JSON response into a SearchResults struct
	var response SearchResponse
	reader := bytes.NewReader(body)
	if err := json.NewDecoder(reader).Decode(&response); err != nil {
		return nil, err
	}

	return response.Results, nil
}
