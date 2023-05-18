package gql_auto_test

import (
	"github.com/SbstnErhrdt/gql_auto"
	"github.com/graphql-go/graphql"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type UUIDTyped struct {
	Value string
}

func (t *UUIDTyped) GraphqlType() graphql.Type {
	return graphql.Float
}

type ModelComplete struct {
	ID           UUIDTyped  `graphql:"id"`
	IDPtr        *UUIDTyped `graphql:"idPtr"`
	Name         string     `graphql:"name"`
	NamePtr      *string    `graphql:"namePtr"`
	CreatedAt    time.Time  `graphql:"createdAt"`
	CreatedAtPtr *time.Time `graphql:"createdAtPtr"`
	No           string     `graphql:"-"`
}

func TestEncode(t *testing.T) {
	t.Parallel()
	ass := assert.New(t)
	obj, err := gql_auto.NewEncoder().Struct(&ModelComplete{})
	fields := obj.Fields()
	ass.NoError(err)
	ass.NotNil(fields)
	ass.Len(fields, 6)
	ass.Contains(fields, "id")
	ass.Contains(fields, "idPtr")
	ass.Contains(fields, "name")
	ass.Contains(fields, "namePtr")
	ass.Contains(fields, "createdAt")
	ass.Contains(fields, "createdAtPtr")
}
