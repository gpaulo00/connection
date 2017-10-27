package connection

import (
	"errors"
	"fmt"
	"github.com/Masterminds/squirrel"
	"reflect"
)

// QueryConfig defines the arguments in a OpaqueCursor call.
type QueryConfig struct {
	SQL   squirrel.SelectBuilder
	ID    string
	First int
	After string
}

// OpaqueCursor applies pagination in a regular query.
func OpaqueCursor(params QueryConfig) (result *squirrel.SelectBuilder, pageSize int, err error) {
	builder := params.SQL

	if params.First == 0 {
		// pass
	} else if params.First < 1 {
		return nil, 0, errors.New("the 'first' argument cannot be less than 1")
	} else {
		// add 1 more element (for hasNextPage)
		pageSize = params.First + 1
		builder = builder.Limit(uint64(pageSize))
	}
	if params.After != "" {
		cursor, err := ParseCursor(params.After)
		if err != nil {
			return nil, 0, err
		}
		builder = builder.Where(params.ID+" < ?", cursor)
	}

	return &builder, pageSize, nil
}

// BuildConnection creates a connection type from a slice of nodes.
func BuildConnection(nodes interface{}, pageSize int) (*Connection, error) {
	// check argument
	slice := reflect.ValueOf(nodes)
	if slice.Kind() != reflect.Slice {
		return nil, errors.New("the argument is not an slice")
	}

	// check page size
	count := slice.Len()
	result := Connection{}
	if pageSize != 0 && count == pageSize {
		result.PageInfo.HasNextPage = true
		count--
	}

	// create the connection
	for i := 0; i < count; i++ {
		// parse a node
		node := slice.Index(i).Interface()
		id := reflect.ValueOf(node).FieldByName("ID")
		if !id.IsValid() {
			return nil, errors.New("cannot create cursor of node: does not have ID field")
		}

		// create edge
		result.Edges = append(result.Edges, Edge{
			Node:   node,
			Cursor: CreateCursor(fmt.Sprintf("%v", id.Interface())),
		})
	}

	// complete page info
	if count > 0 {
		result.PageInfo.StartCursor = result.Edges[0].Cursor
		result.PageInfo.EndCursor = result.Edges[count-1].Cursor
	}
	return &result, nil
}
