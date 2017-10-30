package connection

import (
	"reflect"
	"testing"
)

func TestParseCursor(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "parse simple cursor",
			args: args{input: "Z3BhdWxvOmFhLWJiYi1jYw=="},
			want: "aa-bbb-cc",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseCursor(tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseCursor() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseCursor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreateCursor(t *testing.T) {
	type args struct {
		id string
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "create simple cursor",
			args: args{id: "aa-bbb-cc"},
			want: "Z3BhdWxvOmFhLWJiYi1jYw==",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CreateCursor(tt.args.id); got != tt.want {
				t.Errorf("CreateCursor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreateEdge(t *testing.T) {
	type args struct {
		node interface{}
		id   string
	}
	tests := []struct {
		name string
		args args
		want *Edge
	}{
		{
			name: "create a simple Edge",
			args: args{node: "xyz", id: "aa-bbb-cc"},
			want: &Edge{Node: "xyz", Cursor: "Z3BhdWxvOmFhLWJiYi1jYw=="},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CreateEdge(tt.args.node, tt.args.id); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateEdge() = %v, want %v", got, tt.want)
			}
		})
	}
}
