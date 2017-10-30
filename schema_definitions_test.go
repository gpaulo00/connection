package connection

import (
	"reflect"
	"testing"

	"github.com/graphql-go/graphql"
)

func TestGenerateConnectionSchema(t *testing.T) {
	type args struct {
		name     string
		nodeType graphql.Type
	}
	edgeType := graphql.NewObject(graphql.ObjectConfig{
		Name: "StringEdge",
		Fields: graphql.Fields{
			"cursor": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "A cursor for use in pagination.",
			},
			"node": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "The item at the end of the edge.",
			},
		},
	})
	connType := graphql.NewObject(graphql.ObjectConfig{
		Name: "StringConnection",
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

	tests := []struct {
		name               string
		args               args
		wantConnectionType graphql.Type
		wantEdgeType       graphql.Type
	}{
		{
			name:               "generate regular connection schema",
			args:               args{name: "String", nodeType: graphql.String},
			wantEdgeType:       edgeType,
			wantConnectionType: connType,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotConnectionType, gotEdgeType := GenerateConnectionSchema(tt.args.name, tt.args.nodeType)
			if !reflect.DeepEqual(gotConnectionType, tt.wantConnectionType) {
				t.Errorf("GenerateConnectionSchema() gotConnectionType = %v, want %v", gotConnectionType, tt.wantConnectionType)
			}
			if !reflect.DeepEqual(gotEdgeType, tt.wantEdgeType) {
				t.Errorf("GenerateConnectionSchema() gotEdgeType = %v, want %v", gotEdgeType, tt.wantEdgeType)
			}
		})
	}
}
