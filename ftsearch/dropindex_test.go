package ftsearch

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDropIndex1(t *testing.T) {
	const (
		expected = `[FT.DROPINDEX test DD]`
	)
	dropIndex := NewDropIndex().WithIndex("test").WithDD()

	createCmd := dropIndex.String()
	require.Equal(t, expected, createCmd)
	require.Nil(t, nil)
}
func TestDropIndex2(t *testing.T) {
	const (
		expected = `[FT.DROPINDEX test]`
	)
	dropIndex := NewDropIndex().WithIndex("test")

	createCmd := dropIndex.String()
	require.Equal(t, expected, createCmd)
	require.Nil(t, nil)
}
