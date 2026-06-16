package prometheus

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type TimeseriesPoint struct {
	Timestamp time.Time
	Value     float64
}

type Series struct {
	Labels map[string]string
	Points []TimeseriesPoint
}

type Sample struct {
	Labels map[string]string
	Value  float64
}

type Client struct {
	baseURL    string
	token      string
	httpClient *http.Client
}

func New(baseURL, token string) *Client {
	return &Client{
		baseURL:    baseURL,
		token:      token,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// QueryRange executes a PromQL range query via the Prometheus HTTP API v1.
func (c *Client) QueryRange(ctx context.Context, query string, start, end time.Time, step time.Duration) ([]Series, error) {
	params := url.Values{
		"query": {query},
		"start": {formatUnix(start)},
		"end":   {formatUnix(end)},
		"step":  {strconv.Itoa(int(step.Seconds())) + "s"},
	}

	var resp struct {
		Status string `json:"status"`
		Data   struct {
			Result []struct {
				Metric map[string]string `json:"metric"`
				Values [][]any   `json:"values"`
			} `json:"result"`
		} `json:"data"`
		Error string `json:"error"`
	}

	if err := c.get(ctx, "/api/v1/query_range", params, &resp); err != nil {
		return nil, err
	}
	if resp.Status != "success" {
		return nil, fmt.Errorf("prometheus query_range: %s", resp.Error)
	}

	series := make([]Series, 0, len(resp.Data.Result))
	for _, r := range resp.Data.Result {
		points := make([]TimeseriesPoint, 0, len(r.Values))
		for _, v := range r.Values {
			ts, _ := v[0].(float64)
			valStr, _ := v[1].(string)
			val, _ := strconv.ParseFloat(valStr, 64)
			points = append(points, TimeseriesPoint{
				Timestamp: time.Unix(int64(ts), 0),
				Value:     val,
			})
		}
		series = append(series, Series{Labels: r.Metric, Points: points})
	}
	return series, nil
}

// Query executes an instant PromQL query.
func (c *Client) Query(ctx context.Context, query string, at time.Time) ([]Sample, error) {
	params := url.Values{
		"query": {query},
		"time":  {formatUnix(at)},
	}

	var resp struct {
		Status string `json:"status"`
		Data   struct {
			Result []struct {
				Metric map[string]string `json:"metric"`
				Value  []any     `json:"value"`
			} `json:"result"`
		} `json:"data"`
		Error string `json:"error"`
	}

	if err := c.get(ctx, "/api/v1/query", params, &resp); err != nil {
		return nil, err
	}
	if resp.Status != "success" {
		return nil, fmt.Errorf("prometheus query: %s", resp.Error)
	}

	samples := make([]Sample, 0, len(resp.Data.Result))
	for _, r := range resp.Data.Result {
		if len(r.Value) < 2 {
			continue
		}
		valStr, _ := r.Value[1].(string)
		val, _ := strconv.ParseFloat(valStr, 64)
		samples = append(samples, Sample{Labels: r.Metric, Value: val})
	}
	return samples, nil
}

func (c *Client) get(ctx context.Context, path string, params url.Values, out any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+path+"?"+params.Encode(), nil)
	if err != nil {
		return fmt.Errorf("prometheus: build request: %w", err)
	}
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("prometheus: request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		return fmt.Errorf("prometheus: HTTP %d", res.StatusCode)
	}
	return json.NewDecoder(res.Body).Decode(out)
}

func formatUnix(t time.Time) string {
	return strconv.FormatFloat(float64(t.UnixNano())/1e9, 'f', 3, 64)
}
