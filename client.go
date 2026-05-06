package lumina

import (
	"net/http"

	"github.com/nexula-rg/go-lumina/internal/config"
)

type Client struct {
	client *http.Client
	apiKey string
}

func NewClient(apiKey string) (*Client, error) {
	config.Url = config.Getenv(config.LuminaURL, config.DefaultUrl)

	client := Client{
		client: &http.Client{},
		apiKey: apiKey,
	}

	return &client, client.Check(apiKey)
}
