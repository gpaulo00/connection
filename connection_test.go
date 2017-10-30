package connection

import (
	"reflect"
	"testing"

	"github.com/Masterminds/squirrel"
)

func TestOpaqueCursor(t *testing.T) {
	sql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Select("image_id").From("images")
	example := CreateCursor("example")

	type args struct {
		params QueryConfig
	}
	tests := []struct {
		name         string
		args         args
		wantResult   string
		wantPageSize int
		wantErr      bool
	}{
		{
			name: "query without pagination, desc order",
			args: args{params: QueryConfig{
				SQL: sql, ID: "image_id", Order: -1,
			}},
			wantResult: "SELECT image_id FROM images ORDER BY image_id DESC",
		},
		{
			name: "query without pagination, asc order",
			args: args{params: QueryConfig{
				SQL: sql, ID: "image_id",
			}},
			wantResult: "SELECT image_id FROM images ORDER BY image_id ASC",
		},
		{
			name: "query only first 10 elements",
			args: args{params: QueryConfig{
				SQL: sql, ID: "image_id", First: 10,
			}},
			wantResult:   "SELECT image_id FROM images ORDER BY image_id ASC LIMIT 11",
			wantPageSize: 11,
		},
		{
			name: "query elements after cursor",
			args: args{params: QueryConfig{
				SQL: sql, ID: "image_id", After: example,
			}},
			wantResult: "SELECT image_id FROM images WHERE image_id > $1 ORDER BY image_id ASC",
		},
		{
			name: "query only first 10 elements after cursor",
			args: args{params: QueryConfig{
				SQL: sql, ID: "image_id", After: example, First: 10,
			}},
			wantResult:   "SELECT image_id FROM images WHERE image_id > $1 ORDER BY image_id ASC LIMIT 11",
			wantPageSize: 11,
		},
		{
			name: "not query with invalid 'first' argument",
			args: args{params: QueryConfig{
				SQL: sql, ID: "image_id", First: -10,
			}},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, gotPageSize, err := OpaqueCursor(tt.args.params)
			var gotResult string
			if (err != nil) != tt.wantErr {
				t.Errorf("OpaqueCursor() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if result != nil {
				gotResult, _, _ = result.ToSql()
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("OpaqueCursor() gotResult = %v, want %v", gotResult, tt.wantResult)
			}
			if gotPageSize != tt.wantPageSize {
				t.Errorf("OpaqueCursor() gotPageSize = %v, want %v", gotPageSize, tt.wantPageSize)
			}
		})
	}
}

func TestBuildConnection(t *testing.T) {
	type args struct {
		nodes    interface{}
		pageSize int
	}
	type node struct {
		ID    string
		Value int
	}
	example := node{ID: "aaa", Value: 1}
	exampleCursor := CreateCursor(example.ID)

	tests := []struct {
		name    string
		args    args
		want    *Connection
		wantErr bool
	}{
		{
			name: "create simple connection",
			args: args{nodes: []node{example}},
			want: &Connection{
				Edges: []Edge{{Node: example, Cursor: exampleCursor}},
				PageInfo: PageInfo{
					HasNextPage: false,
					StartCursor: exampleCursor,
					EndCursor:   exampleCursor,
				},
			},
		},
		{
			name: "create connection without next page",
			args: args{nodes: []node{example}, pageSize: 2},
			want: &Connection{
				Edges: []Edge{{Node: example, Cursor: exampleCursor}},
				PageInfo: PageInfo{
					HasNextPage: false,
					StartCursor: exampleCursor,
					EndCursor:   exampleCursor,
				},
			},
		},
		{
			name: "create connection with next page",
			args: args{
				nodes:    []node{example, node{ID: "aab", Value: 2}},
				pageSize: 2,
			},
			want: &Connection{
				Edges: []Edge{{Node: example, Cursor: exampleCursor}},
				PageInfo: PageInfo{
					HasNextPage: true,
					StartCursor: exampleCursor,
					EndCursor:   exampleCursor,
				},
			},
		},
		{
			name:    "not create connection without node as slice",
			args:    args{nodes: example},
			wantErr: true,
		},
		{
			name: "not create connection without node.ID",
			args: args{nodes: []struct{ Value int }{
				{Value: 100},
				{Value: 10},
			}},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := BuildConnection(tt.args.nodes, tt.args.pageSize)
			if (err != nil) != tt.wantErr {
				t.Errorf("BuildConnection() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BuildConnection() = %v, want %v", got, tt.want)
			}
		})
	}
}
