# Sugar

A toolkit for extending the syntax of Go programming language.  

## Features

- Rust-like question mark for types **Result[T]** and **Option[T]**.  

https://github.com/arcane-craft/sugar/blob/d804e5ad894ad92958b0530e84d0465e3cb26fd2/examples/question/main.go#L19-L24

https://github.com/arcane-craft/sugar/blob/d804e5ad894ad92958b0530e84d0465e3cb26fd2/examples/question/main.go#L56-L58

- Try/Catch based error handling for function call which has return type **error** at last position.  

https://github.com/arcane-craft/sugar/blob/d804e5ad894ad92958b0530e84d0465e3cb26fd2/examples/exception/main.go#L17-L28

- Try function with *defer* error handling inspired on proposal [#32437](https://github.com/golang/go/issues/32437).  

https://github.com/arcane-craft/sugar/blob/d804e5ad894ad92958b0530e84d0465e3cb26fd2/examples/tryfunc/main.go#L18-L37

- C-like macro identifiers.  

https://github.com/arcane-craft/sugar/blob/d804e5ad894ad92958b0530e84d0465e3cb26fd2/examples/predef/main.go#L15-L21

## Usage

1. Write your code.  
2. Transform the code by this command:  
```bash
go run -mod=mod github.com/arcane-craft/sugar/tool/transform@latest [PROJECT_ROOT_DIR]
```
3. Build your project with additional tag *sugar_production*:  
```bash
go build -tags=sugar_production -v [package_name]
```
