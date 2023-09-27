package gql_auto

import (
	"github.com/graphql-go/graphql"
	"sync"
)

var mutex = &sync.Mutex{}

func AddField(root *graphql.Object, field *graphql.Field) {
	mutex.Lock()
	root.AddFieldConfig(field.Name, field)
	mutex.Unlock()
}
