package util

import "strings"

// GenerateDynamoID returns a dynamoDB id for Adjacency List strategy
func GenerateDynamoID(table, id string) string {
	return strings.Title(table) + "#" + id
}

// FromDynamoID returns a key from a composite dynamoDB key
func FromDynamoID(key string) string {
	keys := strings.Split(key, "#")
	if len(keys) < 2 {
		return ""
	}

	return keys[1]
}
