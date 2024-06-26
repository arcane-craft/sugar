# Sugar

A toolkit for extending the syntax of Go programming language.  

## Features

- Rust like question mark for types **Result[T]** and **Option[T]**.  


- Try/Catch based error handling for function calls which has return type **error** at last position.  



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
