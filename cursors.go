package connection

import (
	"encoding/base64"
	"errors"
	"strings"
)

var cursorPrefix = "gpaulo"
var errCursor = errors.New("this cursor is invalid")

// ParseCursor gets information about an opaque cursor.
func ParseCursor(input string) (string, error) {
	decode, err := base64.URLEncoding.DecodeString(input)
	if err != nil {
		return "", errCursor
	}

	splitted := strings.Split(string(decode), ":")
	prefix, cursor := splitted[0], splitted[1]
	if prefix != cursorPrefix {
		return "", errCursor
	}

	return cursor, nil
}

// CreateCursor creates a opaque cursor for a connection (base64 encoded).
func CreateCursor(id string) string {
	cursor := []byte(cursorPrefix + ":" + id)
	return base64.URLEncoding.EncodeToString(cursor)
}

// CreateEdge creates an edge of connection from a node.
func CreateEdge(node interface{}, id string) *Edge {
	return &Edge{
		Node:   node,
		Cursor: CreateCursor(id),
	}
}
