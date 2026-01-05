package lumina

import (
	"net/http"

	"github.com/nexula-rg/go-lumina/lumina/internal/config"
)

type Client struct {
	client *http.Client
}

func NewClient() (*Client, error) {
	config.Url = config.Getenv(config.GoLuminaURL, config.DefaultUrl)

	client := Client{
		client: &http.Client{},
	}

	return &client, client.Ping()
}
