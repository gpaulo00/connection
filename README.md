
# Connections (using opaque-cursors)
This package is used to generate *Connection* boilerplate in **Go** (golang) without
using Relay tools.

## Pagination in SQL
This package uses some *SQL Statements* to create the pagination, avoiding to read **all**
the database (and parsing it with **ConnectionFromArray**) and requesting the data that
it actually needs. Example (SQL):
```sql
-- input query
SELECT * FROM images;

-- output (first 10 after "XYZ")
SELECT * FROM images WHERE id < 'XYZ' LIMIT 10;
```

## Build the GraphQL boilerplate
With the **GenerateConnectionSchema** method you can automatically generate the schema of
the connection (**Connection** and **Edge** types).

## Example
```go
import (
  // schema
  "log"
  "github.com/graphql-go/graphql"

  // resolver
  conn "github.com/gpaulo00/connection"
  sq "github.com/Masterminds/squirrel"
)

// ImagesResolver is a field resolver that returns a ImageConnection type.
func ImagesResolver(params graphql.ResolveParams) (interface{}, error) {
  // select only first 10 images after "XYZ"
  builder := sq.Select("*").From("images")
  builder, pageSize, err := conn.OpaqueCursor(conn.QueryConfig{
    SQL: builder,
    ID: "image_id",
    First: 10,
    After: "XYZ",
  })
  if err != nil {
    return nil, err
  }
  query, args, _ := builder.ToSql()

  // query using a "jmoiron/sqlx" instance (map to example Image struct)
  images := []Image{}
  err = db.Select(&images, query, args...)
  if err != nil {
    return nil, err
  }

  // build the connection with the result (using pageSize of OpaqueCursor)
  result, err := conn.BuildConnection(images, pageSize)
  if err != nil {
    return nil, err
  }
  
  return result, nil
}

schema, _ := graphql.NewSchema(graphql.SchemaConfig{
  Query: graphql.NewObject(graphql.ObjectConfig{
    Name:   "RootQuery",
    Fields: graphql.Fields{
      "images": &graphql.Field{
        // generate ImageConnection using ImageType (graphql definition type)
        Type: conn.GenerateConnectionSchema("Image", ImageType),
        // using the resolver defined above
        Resolve: ImagesResolver,
      },
    },
  }),
})
```

## Limitations
This package requires that you use the [Masterminds/squirrel](https://github.com/Masterminds/squirrel)
package to build your queries, and this does not support **backwards pagination**.

## License
**MIT**