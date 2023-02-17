// create provides an interface to RedisSearch's create index functionality.
package ftsearch

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type (
	CreateOptions struct {
		Index   string
		On      string
		Schemas []*SchemaOptions
	}
	// SCHEMA {identifier} AS {attribute} {attribute type} {options...}:
	SchemaOptions struct {
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
func NewCreate() *CreateOptions {
	return &CreateOptions{
		On: "HASH", // since it is the default
	}
}

// NewSchema creates a new schema with defaults set
func NewSchema() *SchemaOptions {
	return &SchemaOptions{}
}
func (s *SchemaOptions) WithIdentifier(identifier string) *SchemaOptions {
	s.identifier = identifier
	return s
}
func (s *CreateOptions) serialize() []interface{} {
	var args = []interface{}{"FT.CREATE", s.Index, "ON", s.On}
	// SCHEMA
	args = append(args, "SCHEMA")
	for _, schema := range s.Schemas {
		schemaArgs := schema.serialize()
		args = append(args, schemaArgs...)
	}
	return args
}
func (s *SchemaOptions) serialize() []interface{} {
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
func (s *SchemaOptions) AsAttribute(attribute string) *SchemaOptions {
	s.attribute = attribute
	return s
}
func (s *SchemaOptions) AttributeType(attributeType string) *SchemaOptions {
	s.attributeType = attributeType
	return s
}
func (q *CreateOptions) WithSchema(s *SchemaOptions) *CreateOptions {
	q.Schemas = append(q.Schemas, s)
	return q
}

// WithIndex sets the index to be search on a create, returning the
// udpated create for chaining
func (q *CreateOptions) WithIndex(index string) *CreateOptions {
	q.Index = index
	return q
}

func (q *CreateOptions) OnJSON() *CreateOptions {
	q.On = "JSON"
	return q
}
func (q *CreateOptions) OnHASH() *CreateOptions {
	q.On = "HASH"
	return q
}
func (q *CreateOptions) String() string {
	return fmt.Sprintf("%v", q.serialize())
}

type (
	CreateIndexResults struct {
		RawResults interface{}
	}
)

func (c *Client) CreateIndex(ctx context.Context, qry *Create) (*CreateIndexResults, error) {
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
