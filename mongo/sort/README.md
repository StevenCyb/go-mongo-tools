# Query for MongoDB-Sort

MongoDB has a client that allows you to **sort** the result of a **find** request.
This parser support a simple syntax to write sort expressions e.g. by the requester of an API (see example below).

## The language

The syntax of this language is simple: `field_name_to_sort_by = sort_order`.
You can chain multiple sort criteria with with the separator `,` e.g. `last_name=asc,first_name=asc`.
There are two ways to sort:

1. `ASC` or `1` to sort ascending
2. `DESC` or `-1` to sort descending

## Example

### For API

```golang
import (
	"github.com/StevenCyb/go-mongo-tools/mongo/sort"

	"go.mongodb.org/mongo-driver/mongo/options"
)
// ...

  sortExpressionString := r.URL.Query().Get("sort")

  parser := sort.NewParser(nil)
  sortExpression, err := parser.Parse(sortExpressionString)
  // ...

  opts := options.Find()
  opts.SetSort(sortExpression)
  // ...

  coll.Find(r.Context(), filter, opts...)
  // ...
```

### For API with policy

This parser supports two types of policies:

1. `WHITELIST_POLICY` -> disallow everything except given fields
2. `BLACKLIST_POLICY` -> allow everything except given fields

```golang
import (
	"github.com/StevenCyb/go-mongo-tools/mongo/sort"
	"github.com/StevenCyb/go-mongo-tools/tokenizer"

	"go.mongodb.org/mongo-driver/mongo/options"
)

func ListHandler(w http.ResponseWriter, r *http.Request) {
  sortExpressionString := r.URL.Query().Get("sort")

  parser := sort.NewParser(
    // just allow the requester to sort by
    // "first_name", "last_name" and "age"
    tokenizer.NewPolicy(
      tokenizer.WHITELIST_POLICY,
      "first_name", "last_name", "age",
    )
  )
  sortExpression, err := parser.Parse(sortExpressionString)
  // ...

  opts := options.Find()
  opts.SetSort(sortExpression)
  // ...

  coll.Find(r.Context(), filter, opts...)
  // ...
}
```
