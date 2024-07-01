//go:build sugar_production

package main

import (
	"encoding/json"
	"io"
	"os"

	. "github.com/arcane-craft/sugar/option"
	. "github.com/arcane-craft/sugar/result"

	fmt_O1ASOQ7UHS "fmt"
)

func main() {
	ResultQuestion()
	OptionQuestion()
}

func ResultQuestion() Result[string] {
	var07QIGRQF8K := From(os.Open("hello.txt"))
	if var07QIGRQF8K.IsErr() {
		return Err[string](fmt_O1ASOQ7UHS.Errorf("func ResultQuestion() github.com/arcane-craft/sugar/result.Result[string]: %w", var07QIGRQF8K.UnwrapErr()))
	}
	file := var07QIGRQF8K.Unwrap()

	defer file.Close()
	var60KDCIE178 := From(io.ReadAll(file))
	if var60KDCIE178.IsErr() {
		return Err[string](fmt_O1ASOQ7UHS.Errorf("func ResultQuestion() github.com/arcane-craft/sugar/result.Result[string]: %w", var60KDCIE178.UnwrapErr()))
	}
	content := var60KDCIE178.Unwrap()

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
	varCEV3QKHRVO := Deocde([]byte(`{"hello":"world"}`))
	if varCEV3QKHRVO.IsNone() {
		return None[string]()
	}
	var4OEKS2K7VS := varCEV3QKHRVO.Unwrap().Get("hello")
	if var4OEKS2K7VS.IsNone() {
		return None[string]()
	}
	return var4OEKS2K7VS.Unwrap().String()
}
