//go:build !sugar_production

package main

import (
	"encoding/json"
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
	file := From(os.Open("hello.txt")).Q()
	defer file.Close()
	content := From(io.ReadAll(file)).Q()
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
	return Deocde([]byte(`{"hello":"world"}`)).Q().Get("hello").Q().String()
}
