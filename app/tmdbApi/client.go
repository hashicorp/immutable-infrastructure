package tmdbApi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

//take requests and data deserialization,
//and put it all together in our own
// api.Client type.

type Client struct {
	APIKey string
}

func (c *Client) GetTrendingMovies(ctx context.Context) (*GetTrendingMoviesResult, error) {
	if c.APIKey == "" {
		return nil, errors.New("API key must be set")
	}
	//must set your own api key
	requestURL := "https://api.themoviedb.org/3/trending/movie/week?api_key=" + c.APIKey
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating HTTP request: %w", err)
	}
	req.Header.Add("Accept", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending HTTP request: %w", err)
	}

	var response GetTrendingMoviesResult
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}

func (c *Client) GetTopRatedTVShows(ctx context.Context) (*ListTVResult, error) {
	if c.APIKey == "" {
		return nil, errors.New("API key must be set")
	}
	//must set your own api key
	requestURL := "https://api.themoviedb.org/3/tv/top_rated?api_key=" + c.APIKey
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating HTTP request: %w", err)
	}
	req.Header.Add("Accept", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending HTTP request: %w", err)
	}

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading HTTP body: %w", err)
	}
	switch res.StatusCode {
	case http.StatusTooManyRequests:
		return nil, fmt.Errorf("we got rate limited!")
	case http.StatusOK:
	default:
		return nil, fmt.Errorf("failed to get TV season, status code %d, body %q", res.StatusCode, string(bodyBytes))
	}

	var response ListTVResult
	err = json.Unmarshal(bodyBytes, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}

func (c *Client) GetTVDetails(ctx context.Context, tvID int) (*GetTVDetailsResult, error) {
	if c.APIKey == "" {
		return nil, errors.New("API key must be set")
	}
	//must set your own api key
	requestURL := "https://api.themoviedb.org/3/tv/" + strconv.Itoa(tvID) + "?api_key=" + c.APIKey
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating HTTP request: %w", err)
	}
	req.Header.Add("Accept", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending HTTP request: %w", err)
	}

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading HTTP body: %w", err)
	}
	switch res.StatusCode {
	case http.StatusTooManyRequests:
		return nil, fmt.Errorf("we got rate limited!")
	case http.StatusOK:
	default:
		return nil, fmt.Errorf("failed to get TV details, status code %d, body %q", res.StatusCode, string(bodyBytes))
	}

	var response GetTVDetailsResult
	err = json.Unmarshal(bodyBytes, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}

func (c *Client) GetTVSeason(ctx context.Context, tvID int, seasonNumber int) (*GetTVSeasonResult, error) {
	if c.APIKey == "" {
		return nil, errors.New("API key must be set")
	}
	//must set your own api key
	requestURL := "https://api.themoviedb.org/3/tv/" + strconv.Itoa(tvID) + "/season/" + strconv.Itoa(seasonNumber) + "?api_key=" + c.APIKey
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating HTTP request: %w", err)
	}
	req.Header.Add("Accept", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending HTTP request: %w", err)
	}

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading HTTP body: %w", err)
	}
	switch res.StatusCode {
	case http.StatusTooManyRequests:
		return nil, fmt.Errorf("we got rate limited!")
	case http.StatusOK:
	default:
		return nil, fmt.Errorf("failed to get TV season, status code %d, body %q", res.StatusCode, string(bodyBytes))
	}

	var response GetTVSeasonResult
	err = json.Unmarshal(bodyBytes, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}
