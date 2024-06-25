//go:build sugar_production

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	. "github.com/arcane-craft/sugar/option"
	. "github.com/arcane-craft/sugar/result"
)

func main() {
	ResultQuestion()
	OptionQuestion()
}

func ResultQuestion() Result[string] {
	varRB3TCFLMUS := FromR1Err(os.Open("hello.txt"))
	if varRB3TCFLMUS.IsErr() {
		return Err[string](fmt.Errorf("func ResultQuestion() github.com/arcane-craft/sugar/result.Result[string]: %w", varRB3TCFLMUS.UnwrapErr()))
	}
	file := varRB3TCFLMUS.Unwrap()

	defer file.Close()
	varO5PQ7UU4G4 := FromR1Err(io.ReadAll(file))
	if varO5PQ7UU4G4.IsErr() {
		return Err[string](fmt.Errorf("func ResultQuestion() github.com/arcane-craft/sugar/result.Result[string]: %w", varO5PQ7UU4G4.UnwrapErr()))
	}
	content := varO5PQ7UU4G4.Unwrap()

	return Ok(string(content))
}

type JSONValue struct {
	any
}

func (v JSONValue) String() Option[string] {
	s, ok := v.any.(string)
	if !ok {
		return None[string]()
	}
	return Some(s)
}

type JSONObject map[string]any

func Deocde(content []byte) Option[JSONObject] {
	var obj JSONObject
	if err := json.Unmarshal(content, &obj); err != nil {
		return None[JSONObject]()
	}
	return Some(obj)
}

func (o JSONObject) Get(key string) Option[JSONValue] {
	value, ok := o[key]
	if !ok {
		return None[JSONValue]()
	}
	return Some(JSONValue{value})
}

func OptionQuestion() Option[string] {
	varGRIG5FKNGC := Deocde([]byte(`{"hello":"world"}`))
	if varGRIG5FKNGC.IsNone() {
		return None[string]()
	}
	var4OEKS2K7VS := varGRIG5FKNGC.Unwrap().Get("hello")
	if var4OEKS2K7VS.IsNone() {
		return None[string]()
	}
	return var4OEKS2K7VS.Unwrap().String()
}
