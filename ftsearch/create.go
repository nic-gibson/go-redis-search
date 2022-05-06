// create provides an interface to RedisSearch's create index functionality.
package ftsearch

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

type (
	create struct {
		Index   string
		On      string
		schemas []*schema
	}
	// SCHEMA {identifier} AS {attribute} {attribute type} {options...}:
	schema struct {
		identifier    string
		attribute     string
		attributeType string
	}
)

/*
FT.CREATE echoTokenStoreIdx ON JSON SCHEMA $.metadata.type AS type TEXT $.metadata.client_id AS client_id TEXT $.metadata.subject AS subject TEXT
*/
// NewCreate creates a new create with defaults set
// https://redis.io/commands/ft.create/
func NewCreate() *create {
	return &create{
		On: "HASH", // since it is the default
	}
}

// NewSchema creates a new schema with defaults set
func NewSchema() *schema {
	return &schema{}
}
func (s *schema) WithIdentifier(identifier string) *schema {
	s.identifier = identifier
	return s
}
func (s *create) serialize() []interface{} {
	var args = []interface{}{"FT.CREATE", s.Index, "ON", s.On}
	// SCHEMA
	args = append(args, "SCHEMA")
	for _, schema := range s.schemas {
		schemaArgs := schema.serialize()
		args = append(args, schemaArgs...)
	}
	return args
}
func (s *schema) serialize() []interface{} {
	var args = []interface{}{s.identifier}

	if s.attribute != "" {
		args = append(args, "AS")
		args = append(args, s.attribute)
	}
	if s.attributeType != "" {
		args = append(args, s.attributeType)
	}
	return args
}
func (s *schema) AsAttribute(attribute string) *schema {
	s.attribute = attribute
	return s
}
func (s *schema) AttributeType(attributeType string) *schema {
	s.attributeType = attributeType
	return s
}
func (q *create) WithSchema(s *schema) *create {
	q.schemas = append(q.schemas, s)
	return q
}

// WithIndex sets the index to be search on a create, returning the
// udpated create for chaining
func (q *create) WithIndex(index string) *create {
	q.Index = index
	return q
}

func (q *create) OnJSON() *create {
	q.On = "JSON"
	return q
}
func (q *create) OnHASH() *create {
	q.On = "HASH"
	return q
}
func (q *create) String() string {
	return fmt.Sprintf("%v", q.serialize())
}

type (
	CreateIndexResults struct {
		RawResults interface{}
	}
)

func (c *Client) ReIndex(ctx context.Context, index string, qry *create) (*CreateIndexResults, error) {
	dropIndex := NewDropIndex().WithIndex(index)
	c.DropIndex(ctx, dropIndex)

	serialized := qry.serialize()
	cmd := redis.NewCmd(ctx, serialized...)
	if err := c.client.Process(ctx, cmd); err != nil {
		return nil, err
	} else if rawResults, err := cmd.Result(); err != nil {
		return nil, err
	} else {
		return &CreateIndexResults{
			RawResults: rawResults,
		}, nil
	}
}

func (c *Client) CreateIndex(ctx context.Context, qry *create) (*CreateIndexResults, error) {
	serialized := qry.serialize()
	cmd := redis.NewCmd(ctx, serialized...)
	if err := c.client.Process(ctx, cmd); err != nil {
		return nil, err
	} else if rawResults, err := cmd.Result(); err != nil {
		return nil, err
	} else {
		return &CreateIndexResults{
			RawResults: rawResults,
		}, nil
	}
}
