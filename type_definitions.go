package connection

// PageInfo contains information about the pagination in a connection.
type PageInfo struct {
	HasNextPage bool   `json:"hasNextPage"`
	StartCursor string `json:"startCursor"`
	EndCursor   string `json:"endCursor"`
}

// Edge is an edge in a connection.
type Edge struct {
	Node   interface{} `json:"node"`
	Cursor string      `json:"cursor"`
}

// Connection is a paginated list of nodes.
type Connection struct {
	PageInfo PageInfo `json:"pageInfo"`
	Edges    []Edge   `json:"edges"`
}