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
	CreateIndex(ctx context.Context, index string)
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
func (cmd *QueryCmd) SetVal(val []interface{}) {
	cmd.SliceCmd.SetVal(val)
}

func (cmd *QueryCmd) Val() []interface{} {
	return cmd.SliceCmd.Val()
}

func (cmd *QueryCmd) Result() ([]interface{}, error) {
	return cmd.SliceCmd.Result()
}

func (cmd *QueryCmd) String() string {
	return cmd.SliceCmd.String()
}

// Drop an index
func (c *Client) DropIndex(ctx context.Context, index string, dropDocuments bool) *redis.BoolCmd {
	args := []interface{}{"ft.dropindex", index}
	if dropDocuments {
		args = append(args, "DD")
	}
	return redis.NewBoolCmd(ctx, args...)
}

// Create an index
func (c *Client) CreateIndex(ctx context.Context, index string, options *IndexOptions) *redis.BoolCmd {
	args := []interface{}{"ft.create", index}
	args = append(args, options.serialize()...)
	cmd := redis.NewBoolCmd(ctx, args...)
	_ = c.Process(ctx, cmd)
	return cmd
}

/* ------------------ Useful internals --------- */

// serializeCountedArgs is used to serialize a string array to
// NAME <count> values. If incZero is true then NAME 0 will be generated
// otherwise empty results will not be generated.
func serializeCountedArgs(name string, incZero bool, args []string) []interface{} {
	if len(args) > 0 || incZero {
		result := make([]interface{}, 2+len(args))

		result[0] = name
		result[1] = len(args)
		for pos, val := range args {
			result[pos+2] = val
		}

		return result
	} else {
		return nil
	}
}
