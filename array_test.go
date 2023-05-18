package gql_auto_test

import (
	"github.com/SbstnErhrdt/gql_auto"
	"github.com/graphql-go/graphql"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestEncoder_ArrayOfStirng(t *testing.T) {
	t.Parallel()
	ass := assert.New(t)
	obj, err := gql_auto.NewEncoder().ArrayOf(reflect.TypeOf(""))
	ass.NoError(err)
	ass.NotNil(obj)
	ass.Equal(graphql.NewList(graphql.String), obj)
	ass.Equal("[String]", obj.String())
}

func TestEncoder_ArrayOfStirngPointer(t *testing.T) {
	t.Parallel()
	s := ""
	ass := assert.New(t)
	obj, err := gql_auto.NewEncoder().ArrayOf(reflect.TypeOf(&s))
	ass.NoError(err)
	ass.NotNil(obj)
	ass.Equal(graphql.NewList(graphql.String), obj)
	ass.Equal("[String]", obj.String())
}

func TestEncoder_ArrayOfNumber(t *testing.T) {
	t.Parallel()
	ass := assert.New(t)
	obj, err := gql_auto.NewEncoder().ArrayOf(reflect.TypeOf(1337))
	ass.NoError(err)
	ass.NotNil(obj)
	ass.Equal(graphql.NewList(graphql.Int), obj)
	ass.Equal("[Int]", obj.String())
}

func TestEncoder_ArrayOfNumberPointer(t *testing.T) {
	t.Parallel()
	n := 1337
	ass := assert.New(t)
	obj, err := gql_auto.NewEncoder().ArrayOf(reflect.TypeOf(&n))
	ass.NoError(err)
	ass.NotNil(obj)
	ass.Equal(graphql.NewList(graphql.Int), obj)
	ass.Equal("[Int]", obj.String())
}

func TestEncoder_ArrayOfFloat(t *testing.T) {
	t.Parallel()
	ass := assert.New(t)
	obj, err := gql_auto.NewEncoder().ArrayOf(reflect.TypeOf(1337.0))
	ass.NoError(err)
	ass.NotNil(obj)
	ass.Equal(graphql.NewList(graphql.Float), obj)
	ass.Equal("[Float]", obj.String())
}

func TestEncoder_ArrayOfFloatPointer(t *testing.T) {
	t.Parallel()
	n := 1337.0
	ass := assert.New(t)
	obj, err := gql_auto.NewEncoder().ArrayOf(reflect.TypeOf(&n))
	ass.NoError(err)
	ass.NotNil(obj)
	ass.Equal(graphql.NewList(graphql.Float), obj)
	ass.Equal("[Float]", obj.String())
}

func TestEncoder_ArrayOfCustomFieldType(t *testing.T) {
	t.Parallel()
	ass := assert.New(t)
	obj, err := gql_auto.NewEncoder().ArrayOf(reflect.TypeOf(CustomFieldType{}))
	ass.NoError(err)
	ass.NotNil(obj)
	ass.Equal("[CustomFieldType]", obj.String())
}

func TestEncoder_ArrayOfCustomFieldPointer(t *testing.T) {
	t.Parallel()
	cf := CustomFieldType{}
	ass := assert.New(t)
	obj, err := gql_auto.NewEncoder().ArrayOf(reflect.TypeOf(&cf))
	ass.NoError(err)
	ass.NotNil(obj)
	ass.Equal("[CustomFieldType]", obj.String())
}

type CustomFieldTypeWithNoGraphQLConverted interface{}

func TestEncoder_ArrayOfNoGraphQLConverted(t *testing.T) {
	t.Parallel()
	ass := assert.New(t)

	type StructExampleNotFound struct {
		Field1 []CustomFieldTypeWithNoGraphQLConverted `graphql:"field1"`
	}

	enc := gql_auto.NewEncoder()
	_, err := enc.ArrayOf(reflect.TypeOf(&StructExampleNotFound{}))
	ass.Error(err)
	ass.ErrorContains(err, "not recognized")
}
