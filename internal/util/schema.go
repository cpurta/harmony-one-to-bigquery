package util

import "cloud.google.com/go/bigquery"

func SchemasEqual(a, b bigquery.Schema) bool {
	if len(a) != len(b) {
		return false
	}

	for _, aFieldSchema := range a {
		var field *bigquery.FieldSchema

		for _, bFieldSchema := range b {
			if aFieldSchema.Name == bFieldSchema.Name {
				field = bFieldSchema
			}
		}

		if field == nil {
			return false
		}

		if aFieldSchema.Type != field.Type {
			return false
		}
	}

	return true
}
