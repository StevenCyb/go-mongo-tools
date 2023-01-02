# Query for JSON patch

With patch operations patches can be specified in detail base on [RFC6902](https://datatracker.ietf.org/doc/html/rfc6902)
(`test` operation is not implemented - use [RSQL Parser](parser/mongo/rsql/README.md) instead)
This features are supported for `MongoDB 4.2+`.

There are five operations available:
| Operation | Description |
|-----------|-------------------------------------------------------------------------------|
| `remove` | remove the value at the target location. |
| `add` | add a value or array to an array at the target location. |
| `replace` | replaces the value at the target location with a new value. |
| `move` | removes the value at a specified location and adds it to the target location. |
| `copy` | copies the value from a specified location to the target location. |

_NOTE_ Mongo Object ID's can be written as 12 bytes long array or 24 character hex string.
In addition - they currently only supported as single field of object or in array (not in map).

There are multiple possibilities on how to use this package:

- Without any restrictions
- With custom manually defined rules (soft restriction)
- Witch reference model that enforce structure and types automatically
- Witch reference model and rule annotations to maximize security

## How to use

### Without any restrictions

```go
import (
  "github.com/StevenCyb/go-mongo-tools/mongo/jsonpatch"

  "go.mongodb.org/mongo-driver/mongo/options"
)
// ...

  var operations []jsonpatch.operation.Spec
  err = json.NewDecoder(req.Body).Decode(&operations)
  // ...

  parser := jsonpatch.NewParser()
  query, err := parser.Parse(operations...)
  // ...

  result := collection.FindOneAndUpdate(ctx, filter, query, updateOptions)
  // ...
```

### With custom manually defined rules

Additionally, simple rules can be set:
| Policy | Description |
|---------------------------------|----------------------------------------------------------|
| `DisallowPathPolicy` | specifies a path that is not allowed. |
| `DisallowOperationOnPathPolicy` | disallows specified operation on path. |
| `ForceTypeOnPathPolicy` | forces the value of a specif path to be from given type. |
| `ForceRegexMatchPolicy` | forces the value of a specif path to match expression. |
| `StrictPathPolicy` | forces path to be strictly one of. |
The path fields can be set to `*` for any field name. E.g. `*.version` will match `product.version` but not `version`.

```go
import (
  "github.com/StevenCyb/go-mongo-tools/mongo/jsonpatch"

  "go.mongodb.org/mongo-driver/mongo/options"
)
// ...

  var operations []jsonpatch.operation.Spec
  err = json.NewDecoder(req.Body).Decode(&operations)
  // ...

  parser := jsonpatch.NewParser(
    DisallowPathPolicy{Details: "illegal ID modification", Path: "_id"},
    ForceTypeOnPathPolicy{Details: "age as number", Path: "user.age", Kind: reflect.Int64},
  )
  query, err := parser.Parse(operations...)
  // ...

  result := collection.FindOneAndUpdate(ctx, filter, query, updateOptions)
  // ...
```

### Witch reference model that enforce structure and types automatically

By using `NewSmartParser` and providing a reference type of the API resource, the parser automatically determine valid patches and types.
Only fields with `bson` tag are considered.
E.g. for the following example a path defined by patch request could be `name` or `metadata.my_map_key.key` and nothing else.
A patch request event could not set a number as `name` since the type of `name` is `string`.

```go
import (
  "github.com/StevenCyb/go-mongo-tools/mongo/jsonpatch"

  "go.mongodb.org/mongo-driver/mongo/options"
)
// ...

type Person struct {
  Name string `bson:"name"`
  Metadata map[string]struct{
    Key string `bson:"key"`
  } `bson:"metadata"`
  IgnoredField string
}

// ...

  var operations []jsonpatch.operation.Spec
  err = json.NewDecoder(req.Body).Decode(&operations)
  // ...

  parser := jsonpatch.NewSmartParser(reflect.TypeOf(Person{}))
  query, err := parser.Parse(operations...)
  // `err` will contain an error if any rule is violated
  // ...

  result := collection.FindOneAndUpdate(ctx, filter, query, updateOptions)
  // ...
```

### Witch reference model and rule annotations to maximize security

The previous approach using `NewSmartParser` could be restricted even more.

#### Disallow any operation on field

Defined by `jp_disallow:"<bool>"`.
No operation with `name` as path is allowed.

```go
type Person struct {
  Name string `bson:"name" jp_disallow:"true"`
}
```

#### Minimum for field

Defined by `jp_min:"<float>"`.
Differentiate by type:

- size on `array`, `map`, `slice` or `string`.
- numeric value on `int` and `float`

```go
type Person struct {
  Age uint `bson:"age" jp_min:"18"`
}
```

#### Maximum for field

Defined by `jp_max:"<float>"`.
Differentiate by type:

- size on `array`, `map`, `slice` or `string`.
- numeric value on `int` and `float`

```go
type Person struct {
  Points uint `bson:"points" jp_max:"100"`
}
```

#### Expression for field value

Defined by `jp_expression:"<escaped_regular_expression>"`.
The value is printed formatted (using `%+v`) before matching.

```go
type Person struct {
  Name string `bson:"name" jp_expression:"^([A-Z][a-z]+ ?){2,}$"`
}
```

#### Allow only specific operations by whitelisting

Defined by `jp_op_allowed:"add,remove,replace,move,copy"` (would allow all).
If the value of a field may only be changed (overwritten):

```go
type Person struct {
  Name string `bson:"name" jp_op_allowed:"replace"`
}
```

#### Allow only specific operations by blacklisting

Defined by `jp_op_disallowed:"add,remove,replace,move,copy"` (same as `jp_disallow`).
If a field should never be deleted:

```go
type Person struct {
  Name string `bson:"name" jp_op_disallowed:"remove"`
}
```
