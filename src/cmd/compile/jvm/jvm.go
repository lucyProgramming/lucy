package jvm

import (
	"fmt"
	"github.com/756445638/lucy/src/cmd/compile/ast"
	"strings"
)

type Jvm struct {
}

func ParseType(name string) (*ast.VariableType, error) {
	if name == "V" {
		return &ast.VariableType{
			Typ: ast.VARIABLE_TYPE_VOID,
		}, nil
	}
	if name == "B" {
		return &ast.VariableType{
			Typ: ast.VARIABLE_TYPE_BYTE,
		}, nil
	}
	if name == "C" {
		return &ast.VariableType{
			Typ: ast.VARIABLE_TYPE_CHAR,
		}, nil
	}
	if name == "D" {
		return &ast.VariableType{
			Typ: ast.VARIABLE_TYPE_DOUBLE,
		}, nil
	}
	if name == "F" {
		return &ast.VariableType{
			Typ: ast.VARIABLE_TYPE_FLOAT,
		}, nil
	}
	if name == "I" {
		return &ast.VariableType{
			Typ: ast.VARIABLE_TYPE_INT,
		}, nil
	}
	if name == "J" {
		return &ast.VariableType{
			Typ: ast.VARIABLE_TYPE_LONG,
		}, nil
	}
	if name == "S" {
		return &ast.VariableType{
			Typ: ast.VARIABLE_TYPE_SHORT,
		}, nil
	}
	if name == "Z" {
		return &ast.VariableType{
			Typ: ast.VARIABLE_TYPE_SHORT,
		}, nil
	}
	if strings.HasPrefix(name, "L") {
		return &ast.VariableType{
			Typ:  ast.VARIABLE_TYPE_CLASS,
			Name: name[1:],
		}, nil
	}
	if strings.HasPrefix(name, "[") {
		t, err := ParseType(name[1:])
		if err != nil {
			return nil, err
		}
		return &ast.VariableType{
			Typ:             ast.VARIABLE_TYPE_ARRAY,
			CombinationType: t,
		}, nil
	}
	panic(fmt.Errorf("unkown type:%v", name))
	return nil, fmt.Errorf("unkown type:%v", name)
}
