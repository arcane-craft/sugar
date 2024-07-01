# Sugar

A toolkit for extending the syntax of Go programming language.  

## Features

- Rust like question mark for types **Result[T]** and **Option[T]**.  

https://github.com/arcane-craft/sugar/blob/2d7d6cf2196bdb4e9d767f1a07cd3168d01ccea5/examples/question/main.go#L19-L24

https://github.com/arcane-craft/sugar/blob/2d7d6cf2196bdb4e9d767f1a07cd3168d01ccea5/examples/question/main.go#L56-L58

- Try/Catch based error handling for function call which has return type **error** at last position.  

https://github.com/arcane-craft/sugar/blob/2d7d6cf2196bdb4e9d767f1a07cd3168d01ccea5/examples/exception/main.go#L17-L28

- Try function with *defer* error handling inspired on proposal [#32437](https://github.com/golang/go/issues/32437).  

https://github.com/arcane-craft/sugar/blob/2d7d6cf2196bdb4e9d767f1a07cd3168d01ccea5/examples/tryfunc/main.go#L18-L37

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
