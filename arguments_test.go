package gql_auto_test

import (
	"github.com/SbstnErhrdt/gql_auto"
	"github.com/graphql-go/graphql"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

type ComplexModel struct {
	FirstName  string `graphql:"!firstName"`
	MiddleName string
	LastName   string `graphql:"lastName"`
	Age        int    `graphql:"age"`
}

func TestComplexModel(t *testing.T) {
	t.Parallel()
	ass := assert.New(t)

	fields, err := gql_auto.NewEncoder().ArgsOf(reflect.TypeOf(ComplexModel{}))
	ass.NoError(err)
	ass.NotNil(fields)
	ass.Len(fields, 4)
	ass.Contains(fields, "firstName")
	ass.Equal(graphql.NewNonNull(graphql.String), fields["firstName"].Type)
	ass.Contains(fields, "lastName")
	ass.Equal(graphql.String, fields["lastName"].Type)
	ass.Contains(fields, "age")
	ass.Equal(graphql.Int, fields["age"].Type)
}

func TestComplexModelPointer(t *testing.T) {
	t.Parallel()
	ass := assert.New(t)

	fields, err := gql_auto.NewEncoder().ArgsOf(reflect.TypeOf(&ComplexModel{}))
	ass.NoError(err)
	ass.NotNil(fields)
	ass.Len(fields, 4)
	ass.Contains(fields, "firstName")
	ass.Equal(graphql.NewNonNull(graphql.String), fields["firstName"].Type)
	ass.Contains(fields, "lastName")
	ass.Equal(graphql.String, fields["lastName"].Type)
	ass.Contains(fields, "age")
	ass.Equal(graphql.Int, fields["age"].Type)

}

func TestError(t *testing.T) {
	t.Parallel()
	ass := assert.New(t)
	_, err := gql_auto.NewEncoder().ArgsOf(reflect.TypeOf("data"))
	ass.Error(err)
}
