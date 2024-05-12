package news

import (
	"errors"
	"net/http"
	"os"

	"github.com/dghubble/sling"
)

var apiURL = os.Getenv("NEWS_API_URL")

type NewsClient struct {
	sling *sling.Sling
}

func NewNewsClient(client *http.Client) *NewsClient {
	return &NewsClient{
		sling: sling.New().Client(client).Base(apiURL),
	}
}

type NewsError struct {
	Code string
}

type MetaResponse struct {
	Found    int `json:"found"`
	Returned int `json:"returned"`
	Limit    int `json:"limit"`
	Page     int `json:"page"`
}

type Source struct {
	SourceID   string   `json:"source_id"`
	Domain     string   `json:"domain"`
	Language   string   `json:"language"`
	Categories []string `json:"categories"`
}

type SourcesResponse struct {
	Meta    MetaResponse `json:"meta"`
	Sources []Source     `json:"data"`
}

type News struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Desc  string `json:"desc"`
	Link  string `json:"link"`
	Img   string `json:"img"`
}

type SourcesParams struct {
	ApiToken   string `url:"api_token"`
	Categories string `url:"categories"`
	Language   string `url:"language"`
}

func (c *NewsClient) sources(languages string, categories string) ([]Source, error) {
	sourcesResp := new(SourcesResponse)
	apiError := new(NewsError)
	path := "/v1/news/sources"
	params := &SourcesParams{
		Categories: categories,
		Language:   languages,
		ApiToken:   os.Getenv("NEWS_API_TOKEN"),
	}
    c.sling.New().Get(path).QueryStruct(params).Receive(sourcesResp, apiError)
	if apiError.Code != "" {
		return nil, errors.New(apiError.Code)
	}

	return sourcesResp.Sources, nil
}
