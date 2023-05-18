# Golang GrahpQL Automatic Schema Generation   

This repository contains a library to auto generate a GraphQL schemas based on golang structs.

It based on [github.com/lab259/go-graphql-struct](https://github.com/lab259/go-graphql-struct)
and [github.com/graphql-go/graphql](https://github.com/graphql-go/graphql).

Usually, building the schema is a one time task, and it is done
statically. So, this library does not degrade the performance, not even
a little, but in that one-time initialization.

## Usage

```go
type Person struct {
	Name    string   `graphql:"!name"`
	Age     int      `graphql:"age"`
	Friends []Person `graphql:"friends"`
}
```

## Custom Types

The default data types of the GraphQL can be count in one hand, which is
not a bad thing. However, that means that you may need to implement some
scalar types (or event complexes types) yourself.

In order to provide custom types for the fields the `GraphqlTyped`
interface was defined:

```go
type GraphqlTyped interface {
    GraphqlType() graphql.Type
}
```

An example:

```go
type TypeA string

func (*TypeA) GraphqlType() graphql.Type {
    return graphql.Int
}

```

Remember, this library is all about declaring the schema. If you need
marshalling/unmarshaling a custom type to another, use the implementation
of the [github.com/graphql-go/graphql](https://github.com/graphql-go/graphql)
library (check on the `graphql.NewScalar` and `graphql.ScalarConfig`).

## Resolver

To implement resolvers over a Custom Type, you will implement the
interface `GraphqlResolver`:

```go
type GraphqlResolver interface {
    GraphqlResolve(p graphql.ResolveParams) (interface{}, error)
}
```

**IMPORTANT**: Although the method `GraphqlResolve` is a member of a struct, it is
called statically. So, do not make any references of the struct itself,
inside of this method.

An example:

```go
type TypeA string

func (*TypeA) GraphqlType() graphql.Type {
    return graphql.Int
}
```

## License

MIT