package util

import (
	"testing"

	"cloud.google.com/go/bigquery"
	"github.com/cpurta/harmony-one-to-bigquery/internal/schema"
	"github.com/stretchr/testify/assert"
)

func TestSchemasEqual(t *testing.T) {
	testcases := []struct {
		name          string
		schemaA       bigquery.Schema
		schemaB       bigquery.Schema
		expectedEqual bool
	}{
		{
			name:          "test empty schemas",
			schemaA:       nil,
			schemaB:       nil,
			expectedEqual: true,
		},
		{
			name:          "one schema is nil",
			schemaA:       make(bigquery.Schema, 0),
			schemaB:       nil,
			expectedEqual: true,
		},
		{
			name:          "empty schemas are equal",
			schemaA:       make(bigquery.Schema, 0),
			schemaB:       make(bigquery.Schema, 0),
			expectedEqual: true,
		},
		{
			name:          "mismatched lengths",
			schemaA:       schema.BlocksTableSchema,
			schemaB:       schema.TransactionsTableSchema,
			expectedEqual: false,
		},
		{
			name:          "equal schemas",
			schemaA:       schema.BlocksTableSchema,
			schemaB:       schema.BlocksTableSchema,
			expectedEqual: true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			equal := SchemasEqual(tc.schemaA, tc.schemaB)

			assert.Equal(t, tc.expectedEqual, equal)
		})
	}
}
