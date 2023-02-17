// ftsearch main module - defines the client class
package ftsearch

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type SearchCmdAble interface {
	redis.Cmdable
	FTSearch(ctx context.Context, index string, query string, options *QueryOptions) *QueryCmd
	DropIndex(ctx context.Context, index string, dropDocuments bool) *redis.BoolCmd
}

type QueryCmd struct {
	redis.SliceCmd
}

type Client struct {
	redis.Client
}

// NewClient returns a new search client
func NewClient(options *redis.Options) *Client {
	return &Client{Client: *redis.NewClient(options)}
}

// NewQueryCmd returns an initialised query command.
func NewQueryCmd(ctx context.Context, args ...interface{}) *QueryCmd {
	return &QueryCmd{
		SliceCmd: *redis.NewSliceCmd(ctx, args...),
	}
}

// Drop an index
func (c *Client) DropIndex(ctx context.Context, index string, dropDocuments bool) *redis.BoolCmd {
	args := []interface{}{"ft.dropindex", index}
	if dropDocuments {
		args = append(args, "DD")
	}
	return redis.NewBoolCmd(ctx, args...)
}
