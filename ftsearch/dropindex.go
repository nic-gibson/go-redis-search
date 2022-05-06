// create provides an interface to RedisSearch's create index functionality.
package ftsearch

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

type (
	dropindex struct {
		Index string
		DD    bool
	}
)

/*
FT.DROPINDEX echoTokenStoreIdx
*/
// NewDropIndex creates a new dropindex with defaults set
// https://redis.io/commands/ft.dropindex/
func NewDropIndex() *dropindex {
	return &dropindex{}
}

func (s *dropindex) serialize() []interface{} {
	var args = []interface{}{"FT.DROPINDEX", s.Index}
	if s.DD {
		args = append(args, "DD")
	}
	return args
}

// WithIndex sets the index to be search on a create, returning the
// udpated create for chaining
func (q *dropindex) WithIndex(index string) *dropindex {
	q.Index = index
	return q
}

func (q *dropindex) WithDD() *dropindex {
	q.DD = true
	return q
}

func (q *dropindex) String() string {
	return fmt.Sprintf("%v", q.serialize())
}

type (
	DropIndexResults struct {
		RawResults interface{}
	}
)

func (c *Client) DropIndex(ctx context.Context, qry *dropindex) (*DropIndexResults, error) {
	serialized := qry.serialize()
	cmd := redis.NewCmd(ctx, serialized...)
	if err := c.client.Process(ctx, cmd); err != nil {
		return nil, err
	} else if rawResults, err := cmd.Result(); err != nil {
		return nil, err
	} else {
		return &DropIndexResults{
			RawResults: rawResults,
		}, nil
	}
}
