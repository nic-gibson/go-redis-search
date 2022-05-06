package ftsearch

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateIndexJSON(t *testing.T) {
	const (
		expected = `[FT.CREATE test ON JSON SCHEMA $.metadata.type AS type TEXT $.metadata.client_id AS client_id TEXT $.metadata.subject AS subject TEXT]`
	)
	create := NewCreate().WithIndex("test").OnJSON().
		WithSchema(NewSchema().
			WithIdentifier("$.metadata.type").AsAttribute("type").AttributeType("TEXT")).
		WithSchema(NewSchema().
			WithIdentifier("$.metadata.client_id").AsAttribute("client_id").AttributeType("TEXT")).
		WithSchema(NewSchema().
			WithIdentifier("$.metadata.subject").AsAttribute("subject").AttributeType("TEXT"))

	createCmd := create.String()
	require.Equal(t, expected, createCmd)
	require.Nil(t, nil)
}
