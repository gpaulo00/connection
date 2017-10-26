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
			Type: graphql.String,
		},
		"startCursor": &graphql.Field{
			Type: graphql.String,
		},
		"hasNextPage": &graphql.Field{
			Type: graphql.NewNonNull(graphql.Boolean),
		},
	},
})

// GenerateConnectionSchema create the GraphQL definition for the connection.
func GenerateConnectionSchema(name string, nodeType graphql.Type) (connectionType graphql.Type) {
	// create edge
	edgeType := graphql.NewObject(graphql.ObjectConfig{
		Name: name + "Edge",
		Fields: graphql.Fields{
			"cursor": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
			"node": &graphql.Field{
				Type: graphql.NewNonNull(nodeType),
			},
		},
	})

	// create connection
	connectionType = graphql.NewObject(graphql.ObjectConfig{
		Name: name + "Connection",
		Fields: graphql.Fields{
			"edges": &graphql.Field{
				Type: graphql.NewNonNull(graphql.NewList(edgeType)),
			},
			"pageInfo": &graphql.Field{
				Type: graphql.NewNonNull(PageInfoType),
			},
		},
	})
	return
}
