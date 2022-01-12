package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/MAKLs/nextcloud-exporter/metrics"
	"github.com/MAKLs/nextcloud-exporter/models"
)

const (
	authHeader = "NC-Token"
	ncApi      = "/ocs/v2.php/apps/serverinfo/api/v1/info?format=json"
)

type Client interface {
	FetchNCServerInfo() (*models.NCServerInfo, error)
}

type NCClient struct {
	httpClient *http.Client
	url        *url.URL
	token      string
}

func NewNCClient(baseUrl *url.URL, token string) *NCClient {
	client := http.DefaultClient
	if apiUrl, err := baseUrl.Parse(ncApi); err != nil {
		panic(fmt.Sprintf("failed to parse URL: %v", err))
	} else {
		return &NCClient{httpClient: client, url: apiUrl, token: token}
	}
}

func (c *NCClient) prepareRequest() (*http.Request, error) {
	req, err := http.NewRequest(http.MethodGet, c.url.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set(authHeader, c.token)

	return req, nil
}

func (c *NCClient) FetchNCServerInfo() (*models.NCServerInfo, error) {
	req, err := c.prepareRequest()
	if err != nil {
		return nil, err
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	decodedBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	// Handle response
	var (
		result    *models.NCServerInfo
		errResult error
	)

	metrics.ScrapeCount.WithLabelValues(strconv.FormatUint(uint64(res.StatusCode), 10)).Inc()

	switch res.StatusCode {
	case http.StatusOK:
		var ncMetrics models.NCServerInfo
		err = json.Unmarshal(decodedBody, &ncMetrics)
		if err != nil {
			result = nil
			errResult = err
		} else {
			result = &ncMetrics
			errResult = nil
		}
	default:
		var ncError models.NCError
		result = nil
		err = json.Unmarshal(decodedBody, &ncError)
		if err != nil {
			errResult = err
		} else {
			errResult = fmt.Errorf("error fetching NC metrics: %s", ncError.Ocs.Meta.Message)
		}
	}

	return result, errResult
}
