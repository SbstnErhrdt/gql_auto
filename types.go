package gql_auto

import (
	"github.com/google/uuid"
	"github.com/graphql-go/graphql"
	"reflect"
	"time"
)

// GraphqlTyped is the interface implemented by types that will provide a
// special `graphql.Type`.
type GraphqlTyped interface {
	// GraphqlType returns the `graphql.Type` that represents the data type that
	// implements this interface.
	GraphqlType() graphql.Type
}

var (
	graphqlTypedType    = reflect.TypeOf(new(GraphqlTyped)).Elem()
	graphqlResolverType = reflect.TypeOf(new(GraphqlResolver)).Elem()
	timeType            = reflect.TypeOf(time.Time{})
	uuidType            = reflect.TypeOf(uuid.UUID{})
)

func (enc *encoder) buildFieldType(fieldType reflect.Type) (graphql.Type, error) {
	if r, ok := enc.getType(fieldType); ok {
		return r, nil
	}

	if fieldType.Kind() == reflect.Struct && (fieldType != timeType || fieldType != uuidType) {
		// If the type is a struct, we need a pointer to that struct to
		// check if it implements the interface.
		tStruct := reflect.PtrTo(fieldType)
		if tStruct.Implements(graphqlTypedType) {
			vStruct := reflect.New(fieldType)
			return vStruct.Interface().(GraphqlTyped).GraphqlType(), nil
		}
	}

	if fieldType.Implements(graphqlTypedType) {
		vStruct := reflect.New(fieldType.Elem())
		return vStruct.Interface().(GraphqlTyped).GraphqlType(), nil
	}

	// Check if it is a pointer or interface...
	if fieldType.Kind() == reflect.Ptr {
		// Updates the type with the type of the pointer
		fieldType = fieldType.Elem()
	}

	// Special case: If the type is the time.Time type.
	if fieldType == timeType {
		return graphql.DateTime, nil
	}
	// Special case: If the type is the uuid.UUID type.
	if fieldType == uuidType {
		return graphql.String, nil
	}

	switch fieldType.Kind() {
	case reflect.Struct:
		return enc.StructOf(fieldType)
	case reflect.Array, reflect.Slice:
		return enc.ArrayOf(fieldType.Elem())
	case reflect.Bool:
		return graphql.Boolean, nil
	case reflect.String:
		return graphql.String, nil
	case
		reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8,
		reflect.Uint, reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8:
		return graphql.Int, nil
	case
		reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128:
		return graphql.Float, nil
	}
	return nil, NewErrTypeNotRecognized(fieldType)
}
