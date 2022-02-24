// ftsearch main module - defines the client class
package ftsearch

import "github.com/go-redis/redis/v8"

type Client struct {
	client *redis.Client
}

// NewClient returns a new search client
func NewClient(c *redis.Client) *Client {
	return &Client{
		client: c,
	}
}
