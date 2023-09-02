package gql_auto

import (
	"fmt"
	"reflect"
	"strings"
	"unicode"

	"github.com/graphql-go/graphql"
)

type Encoder struct {
	types map[string]graphql.Type
}

func NewEncoder() *Encoder {
	return &Encoder{
		types: make(map[string]graphql.Type),
	}
}

var DefaultEncoder = NewEncoder()

func (enc *Encoder) Struct(obj interface{}, options ...Option) (*graphql.Object, error) {
	t := reflect.TypeOf(obj)
	return enc.StructOf(t, options...)
}

func (enc *Encoder) Args(obj interface{}) (graphql.FieldConfigArgument, error) {
	t := reflect.TypeOf(obj)
	return enc.ArgsOf(t)
}

func toLowerCamelCase(input string) string {
	fieldName := []rune(input)
	lastUppercase := false

	for i := 0; i < len(fieldName); i++ {
		lowerCase := true
		hasAnotherLetter := len(fieldName) > i+1
		// The first letter will always be lower case
		if i > 0 {
			if hasAnotherLetter {
				// Keep only the final capital letter of the start of word
				// i.e. ACTTest -> actTest
				// If the next character is upper
				lowerCase = unicode.IsUpper(fieldName[i+1])
			}
			if (len(fieldName) == (i + 2)) && lastUppercase {
				lowerCase = true
			}
		}
		lastUppercase = unicode.IsUpper(fieldName[i])

		if lowerCase {
			fieldName[i] = unicode.ToLower(fieldName[i])
		}

	}
	return string(fieldName)
}

// StructOf returns a `*graphql.Object` with the description extracted from the
// obj passed.
//
// This method extracts the information needed from the fields of the obj
// informed. All fields tagged with "graphql" are added.
//
// The "graphql" tag can be defined as:
//
// ```
//
//	type T struct {
//	    field string `graphql:"fieldname"`
//	}
//
// ```
//
// * fieldname: The name of the field.
func (enc *Encoder) StructOf(t reflect.Type, options ...Option) (*graphql.Object, error) {
	if r, ok := enc.getType(t); ok {
		if d, ok := r.(*graphql.Object); ok {
			return d, nil
		}
		return nil, fmt.Errorf("%s is not an graphql.Object", r)
	}

	name := t.Name()
	if t.Kind() == reflect.Ptr {
		name = t.Elem().Name()
	}

	objCfg := graphql.ObjectConfig{
		Name:   name,
		Fields: graphql.Fields{},
	}

	// Apply options
	for _, opt := range options {
		err := opt.Apply(&objCfg)
		if err != nil {
			return nil, err
		}
	}

	r := graphql.NewObject(objCfg)
	enc.registerType(t, r)

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	// Goes field by field of the object.

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		// if the field is not exported, ignore it
		if field.PkgPath != "" {
			continue
		}

		tag, ok := field.Tag.Lookup("graphql")
		// if there is no graphql tag look for a json tag
		if !ok {
			tag, ok = field.Tag.Lookup("json")
			if ok {
				tag = strings.Replace(tag, ",omitempty", "", -1)
			}
		}
		// if is tagged with graphql or json, but is not exported, ignore it
		if ok && tag == "-" {
			continue
		}

		objectType, ok := enc.getType(field.Type)
		if !ok {
			ot, err := enc.buildFieldType(field.Type)
			if err != nil {
				return nil, NewErrTypeNotRecognizedWithStruct(err, t, field)
			}
			objectType = ot
			enc.registerType(field.Type, ot)
		}

		// If the tag starts with "!" it is a NonNull type.
		if len(tag) > 0 && tag[0] == '!' {
			objectType = graphql.NewNonNull(objectType)
			tag = tag[1:]
		}

		resolve := fieldResolve(field)

		gqlField := &graphql.Field{
			Type:    objectType,
			Resolve: resolve,
		}
		// check the field name
		fieldName := tag
		// if the field name was set by the tag, use it
		if fieldName == "" {
			// if the field name was not set by the tag
			// build the field name by creating a camel case from the struct field name
			fieldNameRunes := []rune(field.Name)
			if len(tag) > 0 {
				fieldNameRunes = []rune(tag)
			}
			fieldName = toLowerCamelCase(string(fieldNameRunes))
		}
		r.AddFieldConfig(fieldName, gqlField)
	}
	return r, nil
}

func (enc *Encoder) FieldOf(t reflect.Type, options ...Option) (graphql.Field, error) {
	r := graphql.Field{}

	fieldType, err := enc.StructOf(t)
	if err != nil {
		return graphql.Field{}, err
	}
	r.Type = fieldType

	for _, option := range options {
		err = option.Apply(&r)
		if err != nil {
			return graphql.Field{}, err
		}
	}

	return r, nil
}

func (enc *Encoder) Field(t interface{}, options ...Option) (graphql.Field, error) {
	return enc.FieldOf(reflect.TypeOf(t), options...)
}

func (enc *Encoder) ArrayOf(t reflect.Type, options ...Option) (graphql.Type, error) {
	if t.Kind() == reflect.Ptr {
		// If pointer, get the Type of the pointer
		t = t.Elem()
	}
	var typeBuilt graphql.Type
	if cachedType, ok := enc.getType(t); ok {
		return graphql.NewList(cachedType), nil
	}
	if t.Kind() == reflect.Struct {
		bt, err := enc.StructOf(t, options...)
		if err != nil {
			return nil, err
		}
		typeBuilt = bt
	} else {
		ttt, err := enc.buildFieldType(t)
		if err != nil {
			return nil, err
		}
		typeBuilt = ttt
	}
	enc.registerType(t, typeBuilt)
	return graphql.NewList(typeBuilt), nil
}

func (enc *Encoder) InputObjectFieldMap(t reflect.Type) (graphql.InputObjectConfigFieldMap, error) {
	r := graphql.InputObjectConfigFieldMap{}

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return r, fmt.Errorf("cannot build args from a non struct")
	}

	// Goes field by field of the object.
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag, ok := field.Tag.Lookup("graphql")
		// if is tagged with graphql, but is not exported, ignore it
		if ok && tag == "-" {
			continue
		}

		objectType, ok := enc.getType(field.Type)
		if !ok {
			ot, err := enc.buildFieldType(field.Type)
			if err != nil {
				return nil, NewErrTypeNotRecognizedWithStruct(err, t, field)
			}
			objectType = ot
			enc.registerType(field.Type, ot)
		}

		// If the tag starts with "!" it is a NonNull type.
		if len(tag) > 0 && tag[0] == '!' {
			objectType = graphql.NewNonNull(objectType)
			tag = tag[1:]
		}

		inputField := &graphql.InputObjectFieldConfig{
			Type: objectType,
		}

		fieldName := []rune(field.Name)
		if len(tag) > 0 {
			fieldName = []rune(tag)
		}
		fieldNameS := toLowerCamelCase(string(fieldName))

		r[fieldNameS] = inputField
	}

	return r, nil
}

func (enc *Encoder) ArgsOf(t reflect.Type) (graphql.FieldConfigArgument, error) {
	r := graphql.FieldConfigArgument{}

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return r, fmt.Errorf("cannot build args from a non struct")
	}

	// Goes field by field of the object.
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag, ok := field.Tag.Lookup("graphql")
		// if is tagged with graphql, but is not exported, ignore it
		if ok && tag == "-" {
			continue
		}

		objectType, ok := enc.getType(field.Type)
		if !ok {
			ot, err := enc.buildFieldType(field.Type)
			if err != nil {
				return nil, NewErrTypeNotRecognizedWithStruct(err, t, field)
			}
			objectType = ot
			enc.registerType(field.Type, ot)
		}

		// If the tag starts with "!" it is a NonNull type.
		if len(tag) > 0 && tag[0] == '!' {
			objectType = graphql.NewNonNull(objectType)
			tag = tag[1:]
		}

		graphQLArgument := &graphql.ArgumentConfig{
			Type: objectType,
		}

		fieldName := []rune(field.Name)
		if len(tag) > 0 {
			fieldName = []rune(tag)
		}
		fieldNameS := toLowerCamelCase(string(fieldName))

		r[fieldNameS] = graphQLArgument
	}

	return r, nil
}

func (enc *Encoder) getType(t reflect.Type) (graphql.Type, bool) {
	name := t.Name()
	if t.Kind() == reflect.Ptr {
		name = t.Elem().Name()
	}
	if len(name) > 0 {
		gt, ok := enc.types[name]
		return gt, ok
	}
	return nil, false
}

func (enc *Encoder) registerType(t reflect.Type, r graphql.Type) {
	name := t.Name()
	if t.Kind() == reflect.Ptr {
		name = t.Elem().Name()
	}
	if len(name) > 0 {
		enc.types[name] = r
	}
}

func Struct(obj interface{}) *graphql.Object {
	r, err := DefaultEncoder.Struct(obj)
	if err != nil {
		panic(err.Error())
	}
	return r
}

func InputObject(name string, obj interface{}) *graphql.InputObject {
	t := reflect.TypeOf(obj)
	r, err := DefaultEncoder.InputObjectFieldMap(t)
	if err != nil {
		panic(err.Error())
	}
	return graphql.NewInputObject(
		graphql.InputObjectConfig{
			Name:   name,
			Fields: r,
		},
	)
}

// Args Obtain the arguments property of a mutation object
func Args(obj interface{}) graphql.FieldConfigArgument {
	r, err := DefaultEncoder.Args(obj)
	if err != nil {
		panic(err.Error())
	}
	return r
}

func ArgsOf(t reflect.Type) graphql.FieldConfigArgument {
	r, err := DefaultEncoder.ArgsOf(t)
	if err != nil {
		panic(err.Error())
	}
	return r
}

func FieldOf(t reflect.Type, options ...Option) (graphql.Field, error) {
	return DefaultEncoder.FieldOf(t, options...)
}

func Field(t interface{}, options ...Option) graphql.Field {
	r, err := DefaultEncoder.Field(t, options...)
	if err != nil {
		panic(err.Error())
	}
	return r
}
