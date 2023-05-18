package gql_auto

import (
	"fmt"
	"github.com/graphql-go/graphql"
	"reflect"
)

func newErrNotSupported(obj interface{}) error {
	t := reflect.TypeOf(obj)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() == reflect.Map {
		return fmt.Errorf("`%s` is not supported", "map")
	}
	return fmt.Errorf("`%s` is not supported", t.Name())
}

// Option describes how an option will behave when applied to a field.
type Option interface {
	Apply(dst interface{}) error
}

type withName struct {
	message string
}

// WithName creates an `Option` that provides sets the name for
// fields and objects
//
// It can be applied to:
// * Field;
// * Objects;
func WithName(name string) Option {
	return &withName{
		message: name,
	}
}

// Apply sets the description to the object.
func (option *withName) Apply(dst interface{}) error {
	switch t := dst.(type) {
	case *graphql.Field:
		t.Name = option.message
		return nil
	case *graphql.ObjectConfig:
		t.Name = option.message
		return nil
	default:
		return newErrNotSupported(dst)
	}
}

type withDescription struct {
	message string
}

// WithDescription creates an `Option` that provides sets the description for
// fields, objects and arguments.
//
// It can be applied to:
// * Field;
// * Arguments;
// * Objects;
func WithDescription(description string) Option {
	return &withDescription{
		message: description,
	}
}

// Apply sets the description to the object.
func (option *withDescription) Apply(dst interface{}) error {
	switch t := dst.(type) {
	case *graphql.Field:
		t.Description = option.message
		return nil
	case *graphql.ArgumentConfig:
		t.Description = option.message
		return nil
	case *graphql.ObjectConfig:
		t.Description = option.message
		return nil
	default:
		return newErrNotSupported(dst)
	}
}

type withDefaultValue struct {
	defaultValue interface{}
}

// WithDefaultValue creates an `Option` that provides sets the description for
// arguments.
//
// It can be applied to:
// * Arguments;
func WithDefaultValue(defaultValue interface{}) Option {
	return &withDefaultValue{
		defaultValue: defaultValue,
	}
}

// Apply sets the default value to the argument.
func (option *withDefaultValue) Apply(dst interface{}) error {
	switch t := dst.(type) {
	case *graphql.ArgumentConfig:
		t.DefaultValue = option.defaultValue
		return nil
	default:
		return newErrNotSupported(dst)
	}
}

type withDeprecationReason struct {
	message string
}

// WithDeprecationReason creates an `Option` that sets the deprecation reason
// for fields.
//
// It can be applied to:
// * Fields;
func WithDeprecationReason(description string) Option {
	return &withDeprecationReason{
		message: description,
	}
}

// Apply sets the deprecation reason to the object.
func (option *withDeprecationReason) Apply(dst interface{}) error {
	switch t := dst.(type) {
	case *graphql.Field:
		t.DeprecationReason = option.message
		return nil
	default:
		return newErrNotSupported(dst)
	}
}

type withResolver struct {
	resolver graphql.FieldResolveFn
}

// WithResolver creates an `Option` that sets the resolver for fields.
//
// It can be applied to:
// * Fields;
func WithResolver(resolver graphql.FieldResolveFn) Option {
	return &withResolver{
		resolver: resolver,
	}
}

// Apply sets the deprecation reason to the object.
func (option *withResolver) Apply(dst interface{}) error {
	switch t := dst.(type) {
	case *graphql.Field:
		t.Resolve = option.resolver
		return nil
	default:
		return newErrNotSupported(dst)
	}
}

type withArgs struct {
	encoder *encoder
	args    interface{}
}

// WithArgs creates an `Option` that sets the arguments for a field.
//
// It can be applied to:
// * Fields;
func WithArgs(args ...interface{}) Option {
	enc := DefaultEncoder
	var data interface{}
	if len(args) == 2 {
		tmp, ok := args[0].(*encoder)
		if !ok {
			panic("the first parameter of WithArgs must be an encoder")
		}
		enc = tmp
		data = args[1]
	} else if len(args) == 1 {
		data = args[0]
	} else {
		panic("invalid usage")
	}
	return &withArgs{
		encoder: enc,
		args:    data,
	}
}

// Apply sets the arguments of a field.
func (option *withArgs) Apply(dst interface{}) error {
	switch t := dst.(type) {
	case *graphql.Field:
		args, err := option.encoder.Args(option.args)
		if err != nil {
			return err
		}
		t.Args = args
		return nil
	default:
		return newErrNotSupported(dst)
	}
}

type withType struct {
	t graphql.Type
}

// WithType creates an `Option` that sets the type of fields.
//
// It can be applied to:
// * Fields;
func WithType(t graphql.Type) Option {
	return &withType{
		t: t,
	}
}

// Apply sets the type of field.
func (option *withType) Apply(dst interface{}) error {
	switch t := dst.(type) {
	case *graphql.Field:
		t.Type = option.t
		return nil
	default:
		return newErrNotSupported(dst)
	}
}
