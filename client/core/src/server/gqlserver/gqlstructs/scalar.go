package gqlstructs

import (
	"strconv"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
)

var StatusEnum = graphql.NewEnum(graphql.EnumConfig{
	Name: "StatusEnum",
	Values: graphql.EnumValueConfigMap{
		"Stopped": &graphql.EnumValueConfig{
			Value: 10,
		},
		"Working": &graphql.EnumValueConfig{
			Value: 20,
		},
		"Stopping": &graphql.EnumValueConfig{
			Value: 30,
		},
		"Ready": &graphql.EnumValueConfig{
			Value: 40,
		},
		"Error": &graphql.EnumValueConfig{
			Value: 50,
		},
	},
})

var Uint = graphql.NewScalar(graphql.ScalarConfig{
	Name:        "Uint",
	Description: "The scalar type represents an Uint",
	Serialize: func(value interface{}) interface{} {
		return value
	},
	ParseValue: func(value interface{}) interface{} {
		switch value := value.(type) {
		case string:
			res, _ := strconv.ParseUint(value, 10, 64)
			return uint(res)
		case *string:
			res, _ := strconv.ParseUint(*value, 10, 64)
			return uint(res)
		case float64:
			return uint(value)
		case *float64:
			return uint(*value)
		default:
			return value
		}
	},
	ParseLiteral: func(valueAST ast.Value) interface{} {
		switch valueAST := valueAST.(type) {
		case *ast.StringValue:
			res, _ := strconv.ParseInt(valueAST.Value, 10, 64)
			return res
		}
		return nil
	},
})
