package connection

import (
	"github.com/graphql-go/graphql"
)

// PageInfoType show information about pagination in the connection in a GraphQL API.
var PageInfoType = graphql.NewObject(graphql.ObjectConfig{
	Name:        "PageInfo",
	Description: "Information about pagination in a connection.",
	Fields: graphql.Fields{
		"endCursor": &graphql.Field{
			Type:        graphql.String,
			Description: "The last cursor in the page.",
		},
		"startCursor": &graphql.Field{
			Type:        graphql.String,
			Description: "The first cursor in the page",
		},
		"hasNextPage": &graphql.Field{
			Type:        graphql.NewNonNull(graphql.Boolean),
			Description: "Are there more items?",
		},
	},
})

// ConnectionArguments describes the arguments that must have the connections.
var ConnectionArguments = graphql.FieldConfigArgument{
	"first": &graphql.ArgumentConfig{
		Type:        graphql.Int,
		Description: "Returns the first *n* elements from the list.",
	},
	"after": &graphql.ArgumentConfig{
		Type:        graphql.String,
		Description: "Returns the elements in the list that come after the specified cursor.",
	},
}

// GenerateConnectionSchema create the GraphQL definition for the connection.
func GenerateConnectionSchema(name string, nodeType graphql.Type) (connectionType graphql.Type, edgeType graphql.Type) {
	// create edge
	edgeType = graphql.NewObject(graphql.ObjectConfig{
		Name: name + "Edge",
		Fields: graphql.Fields{
			"cursor": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "A cursor for use in pagination.",
			},
			"node": &graphql.Field{
				Type:        graphql.NewNonNull(nodeType),
				Description: "The item at the end of the edge.",
			},
		},
	})

	// create connection
	connectionType = graphql.NewObject(graphql.ObjectConfig{
		Name: name + "Connection",
		Fields: graphql.Fields{
			"edges": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.NewList(edgeType)),
				Description: "A list of edges.",
			},
			"pageInfo": &graphql.Field{
				Type:        graphql.NewNonNull(PageInfoType),
				Description: "Information to aid in pagination.",
			},
		},
	})
	return
}
