package gql_auto_test

import (
	"github.com/SbstnErhrdt/gql_auto"
	"github.com/google/uuid"
	"github.com/graphql-go/graphql"
	"github.com/stretchr/testify/assert"
	"testing"
)

type UUIDStruct struct {
	UID  uuid.UUID  `json:"uid" graphql:"UID"`
	UID2 *uuid.UUID `json:"uid2" graphql:"uid2"`
}

func TestEncodeUUID(t *testing.T) {
	t.Parallel()
	ass := assert.New(t)
	obj, err := gql_auto.NewEncoder().Struct(&UUIDStruct{})
	fields := obj.Fields()
	ass.NoError(err)
	ass.NotNil(fields)
	ass.Contains(fields, "UID")
	ass.Equal(fields["UID"].Type.Name(), graphql.String.Name())
	ass.Equal(fields["uid2"].Type.Name(), graphql.String.Name())
}
