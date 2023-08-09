package gqlstructs

import (
	"errors"
	"reflect"

	"github.com/graphql-go/graphql"
)

var EmbeddedResolver = func(p graphql.ResolveParams) (res interface{}, err error) {
	res = getField(p.Source, p.Info.FieldName)
	return
}

type Embedded interface {
	GetField(name string) interface{}
}

type FieldsList map[string]interface{}

type FieldTypes struct {
	In  graphql.Input
	Out graphql.Output
}

func initOutputSchema(obj *graphql.Object, fields FieldsList, embedded bool) {
	for name, val := range fields {
		addOutputField(obj, name, val, embedded)
	}
}

func addOutputField(obj *graphql.Object, name string, val interface{}, embedded bool) {
	var resolver graphql.FieldResolveFn
	if embedded {
		resolver = EmbeddedResolver
	}

	switch value := val.(type) {
	case graphql.Output:
		obj.AddFieldConfig(name, &graphql.Field{
			Type:    value,
			Resolve: resolver,
		})

	case FieldTypes:
		obj.AddFieldConfig(name, &graphql.Field{
			Type:    value.Out,
			Resolve: resolver,
		})

	case FieldsList:
		initOutputSchema(obj, value, true)

	default:
		panic(errors.New("invalid type"))

	}
}

func initInputSchema(obj *graphql.InputObject, fields FieldsList) {
	for name, val := range fields {
		addInputField(obj, name, val)
	}
}

func addInputField(obj *graphql.InputObject, name string, val interface{}) {
	switch value := val.(type) {
	case graphql.Output:
		obj.AddFieldConfig(name, &graphql.InputObjectFieldConfig{
			Type: value,
		})

	case FieldTypes:
		obj.AddFieldConfig(name, &graphql.InputObjectFieldConfig{
			Type: value.In,
		})

	case FieldsList:
		initInputSchema(obj, value)

	default:
		panic(errors.New("invalid type"))

	}
}

func getField(obj interface{}, field string) interface{} {
	r := reflect.ValueOf(obj)
	f := reflect.Indirect(r).FieldByName(field)
	return f.Interface()
}
