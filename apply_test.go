package gql_auto_test

import (
	"errors"
	"github.com/SbstnErhrdt/gql_auto"
	"github.com/graphql-go/graphql"
	"github.com/stretchr/testify/assert"
	"testing"
)

type erroredOption struct{}

func (*erroredOption) Apply(dst interface{}) error {
	return errors.New("forced error")
}

func TestApplyWithDescription(t *testing.T) {
	t.Parallel()
	ass := assert.New(t)

	type StructExample struct {
		Field1 string `graphql:"field1"`
		Field2 string `graphql:"field2"`
	}

	obj, err := gql_auto.NewEncoder().Field(&StructExample{}, gql_auto.WithDescription("Description 1"))
	ass.NoError(err)
	ass.Equal("Description 1", obj.Description)
}

func TestApplyWithArgumentDescription(t *testing.T) {
	t.Parallel()
	ass := assert.New(t)
	argument := graphql.ArgumentConfig{}
	err := gql_auto.WithDescription("Description 1").Apply(&argument)
	ass.NoError(err)
	ass.Equal("Description 1", argument.Description)
}

func TestApplyWithDescriptionToObject(t *testing.T) {
	t.Parallel()
	ass := assert.New(t)
	obj := graphql.ObjectConfig{}
	err := gql_auto.WithDescription("Description 1").Apply(&obj)
	ass.NoError(err)
	ass.Equal("Description 1", obj.Description)
}

func TestWithDescriptionApplyError(t *testing.T) {
	t.Parallel()
	ass := assert.New(t)
	obj := map[string]interface{}{}
	err := gql_auto.WithDescription("Description 1").Apply(&obj)
	ass.Error(err)
}

func TestApplyDefaultValue(t *testing.T) {
	t.Parallel()
	ass := assert.New(t)
	argument := graphql.ArgumentConfig{}
	err := gql_auto.WithDefaultValue("default value").Apply(&argument)
	ass.NoError(err)
	ass.Equal("default value", argument.DefaultValue)
}

func TestApplyDefaultValueError(t *testing.T) {
	t.Parallel()
	ass := assert.New(t)
	field := graphql.Field{}
	err := gql_auto.WithDefaultValue("default value").Apply(&field)
	ass.Error(err)
	ass.ErrorContains(err, "is not supported")
	ass.ErrorContains(err, "Field")
}

func TestApplyDeprecationReason(t *testing.T) {
	t.Parallel()
	ass := assert.New(t)
	field := graphql.Field{}
	err := gql_auto.WithDeprecationReason("deprecation reason").Apply(&field)
	ass.NoError(err)
	ass.Equal("deprecation reason", field.DeprecationReason)
}

func TestApplyDeprecationReasonError(t *testing.T) {
	t.Parallel()
	ass := assert.New(t)
	argument := graphql.ArgumentConfig{}
	err := gql_auto.WithDeprecationReason("default value").Apply(&argument)
	ass.Error(err)
	ass.ErrorContains(err, "is not supported")
	ass.ErrorContains(err, "ArgumentConfig")
}

func TestWithResolver_Apply(t *testing.T) {
	t.Parallel()
	ass := assert.New(t)
	field := graphql.Field{}
	err := gql_auto.WithResolver(func(p graphql.ResolveParams) (interface{}, error) {
		return nil, nil
	}).Apply(&field)
	ass.NoError(err)
	ass.NotNil(field.Resolve)
}

func TestWithResolver_ApplyError(t *testing.T) {
	t.Parallel()
	ass := assert.New(t)
	argument := graphql.ArgumentConfig{}
	err := gql_auto.WithResolver(func(p graphql.ResolveParams) (interface{}, error) {
		return nil, nil
	}).Apply(&argument)
	ass.Error(err)
	ass.ErrorContains(err, "is not supported")
	ass.ErrorContains(err, "ArgumentConfig")
}

func TestWithArgs_Apply(t *testing.T) {
	t.Parallel()
	ass := assert.New(t)
	type Args struct {
		Name string `graphql:"name"`
		Age  int    `graphql:"age"`
		No   int    `graphql:"-"`
	}

	field := graphql.Field{}
	err := gql_auto.WithArgs(Args{}).Apply(&field)
	ass.NoError(err)
	ass.NotNil(field.Args)

	ass.Len(field.Args, 2)
	ass.Contains(field.Args, "name")
	ass.Contains(field.Args, "age")
}

func TestWithArgs_ApplyEncoder(t *testing.T) {
	t.Parallel()
	ass := assert.New(t)

	enc := gql_auto.NewEncoder()

	type Args struct {
		Name string `graphql:"name"`
		Age  int    `graphql:"age"`
	}

	field := graphql.Field{}
	err := gql_auto.WithArgs(enc, Args{}).Apply(&field)
	ass.NoError(err)

	ass.NotNil(field.Args)
	ass.Len(field.Args, 2)
	ass.Contains(field.Args, "name")
	ass.Contains(field.Args, "age")
}

func TestWithArgs_ApplyPanic(t *testing.T) {
	t.Parallel()
	ass := assert.New(t)

	type Args struct {
		Name string `graphql:"name"`
		Age  int    `graphql:"age"`
	}

	field := graphql.Field{}
	ass.Panics(func() {
		_ = gql_auto.WithArgs(Args{}, Args{}).Apply(&field)
	})
}

func TestWithArgs_ApplyPanic2(t *testing.T) {
	t.Parallel()
	ass := assert.New(t)

	field := graphql.Field{}
	ass.Panics(func() {
		_ = gql_auto.WithArgs(1, 2, 3).Apply(&field)
	})
}

func TestWithArgs_ApplyError(t *testing.T) {
	t.Parallel()
	ass := assert.New(t)
	field := graphql.Field{}
	err := gql_auto.WithArgs([]interface{}{}).Apply(&field)
	ass.Error(err)
	ass.ErrorContains(err, "cannot build args from a non struct")
}

func TestWithArgs_ApplyError2(t *testing.T) {
	t.Parallel()
	ass := assert.New(t)
	type Args struct {
		Field1 []interface{} `graphql:"field1"`
	}

	field := graphql.Field{}
	err := gql_auto.WithArgs(Args{}).Apply(&field)
	ass.Error(err)
	ass.ErrorContains(err, "not recognized")
	ass.ErrorContains(err, "interface {}")
}

func TestWithArgs_ApplyError3(t *testing.T) {
	t.Parallel()
	ass := assert.New(t)
	type Args struct {
		Field1 []interface{} `graphql:"field1"`
	}

	obj := graphql.ObjectConfig{}
	err := gql_auto.WithArgs(Args{}).Apply(&obj)
	ass.Error(err)
	ass.ErrorContains(err, "not supported")
	ass.ErrorContains(err, "ObjectConfig")
}

type ThisIsAType struct {
	Name string `graphql:"name"`
}

var tt = gql_auto.Struct(ThisIsAType{})

func TestWithType_Apply(t *testing.T) {
	t.Parallel()
	ass := assert.New(t)
	field := graphql.Field{}
	err := gql_auto.WithType(tt).Apply(&field)
	ass.NoError(err)
	ass.Equal(tt, field.Type)
}

func TestWithType_ApplyError(t *testing.T) {
	t.Parallel()
	ass := assert.New(t)
	argument := graphql.ArgumentConfig{}
	err := gql_auto.WithType(tt).Apply(&argument)
	ass.Error(err)
	ass.ErrorContains(err, "is not supported")
	ass.ErrorContains(err, "ArgumentConfig")
}
