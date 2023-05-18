package gql_auto_test

import (
	"github.com/SbstnErhrdt/gql_auto"
	"github.com/graphql-go/graphql"
	"github.com/stretchr/testify/assert"
	"testing"
)

type CustomFieldType struct {
	Value string
}

func (t *CustomFieldType) GraphqlType() graphql.Type {
	return graphql.Float
}

type CustomFieldTypeWithResolver struct {
	Value string
}

func (t *CustomFieldTypeWithResolver) GraphqlResolve(p graphql.ResolveParams) (interface{}, error) {
	panic("only to catch")
}

func (t *CustomFieldTypeWithResolver) GraphqlType() graphql.Type {
	return graphql.Float
}

func TestStructExample(t *testing.T) {
	t.Parallel()
	ass := assert.New(t)
	type StructExample struct {
		Field1 string `graphql:"field1"`
		Field2 int
		Field3 bool `graphql:"field3"`
		Field4 bool `graphql:"-"`
	}

	obj, err := gql_auto.NewEncoder().Struct(&StructExample{})
	ass.NoError(err)
	ass.NotNil(obj)
	ass.Equal("StructExample", obj.Name())
	fields := obj.Fields()
	ass.NotNil(fields)
	ass.Len(fields, 3)
	ass.Contains(fields, "field1")
	ass.Contains(fields, "field2")
	ass.Contains(fields, "field3")
	ass.NotContains(fields, "field4")

	ass.Equal(graphql.String, fields["field1"].Type)
	ass.Equal("field1", fields["field1"].Name)
	ass.Equal(graphql.Int, fields["field2"].Type)
	ass.Equal("field2", fields["field2"].Name)
	ass.Equal(graphql.Boolean, fields["field3"].Type)
	ass.Equal("field3", fields["field3"].Name)
	ass.Equal("Boolean", fields["field3"].Type.String())
}

func TestWithDescription(t *testing.T) {
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

func TestTypeNotRecognizedWithStructError_Error(t *testing.T) {
	t.Parallel()
	ass := assert.New(t)

	type StructExample struct {
		Field1 string `graphql:"field1"`
		Field2 string `graphql:"field2"`
	}

	_, err := gql_auto.NewEncoder().Field(&StructExample{}, &erroredOption{})
	ass.Error(err)
}

func TestTypeNotRecognizedWithStructError_Error2(t *testing.T) {
	t.Parallel()
	ass := assert.New(t)

	type StructExample struct {
		Field1 []interface{} `graphql:"field1"`
	}

	_, err := gql_auto.NewEncoder().Field(&StructExample{}, &erroredOption{})
	ass.Error(err)
	ass.ErrorContains(err, "interface {}")
}

func TestTypeNotRecognizedWithStructError_Error3(t *testing.T) {
	t.Parallel()
	ass := assert.New(t)

	type StructExample struct {
		Field1 []*interface{} `graphql:"field1"`
	}

	_, err := gql_auto.NewEncoder().Field(&StructExample{}, &erroredOption{})
	ass.Error(err)
	ass.ErrorContains(err, "interface {}")
}
