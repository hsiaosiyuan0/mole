# Mole

Mole is a collection of parsers made with 💕 to process the frontend stuffs.

<details>
  <summary>Why Golang</summary>

~~A little bit explanation is good for why Golang is preferred in this project. Nowadays, a programming language is not only the grammar things, it's consist of runtime, stdlib, 3rd-party modules and a healthy community, all these are out-of-box by using Golang, more specifically:~~

- ~~Golang is productive, its simplicity philosophy(something like Grammar and Garbage-collection) saves more time to the functionalities themselves.~~
- ~~the functionalities like lint and bundle maybe needed to run as web services while Golang has been proved by many impressive projects such k8s that it's good at service things.~~

Fine, just all because I'm too fool to use a fancy language

</details>

## Features

- JavaScript Parser

  - ECMAScript up to [ES2021](https://262.ecma-international.org/12.0/)
  - [JSX](https://github.com/facebook/jsx)
  - [ESTree](https://github.com/estree/estree) compatible outputs ([AST explorer on WASM](http://blog.thehardways.me/mole-is-more/#/))

- TypeScript Parser

  - [babel/typescript](https://babeljs.io/docs/en/babel-types#typescript) compatible outputs

### WIP

- [ ] CSS parser
- [ ] Less parser
- [ ] Scss parser

## Using the parser

For using the parser in Mole, the first step is using `go get` to add Mole as your project's dependency

```bash
go get github.com/hsiaosiyuan0/mole

# or using go get via a proxy if you have some network issues
https_proxy=127.0.0.1:1080 go get github.com/hsiaosiyuan0/mole
```

After the installation is succeeded, below code can be used as a demo:

<details>
  <summary>Click to expand the demo</summary>

```go
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"

	"github.com/hsiaosiyuan0/mole/ecma/estree"
	"github.com/hsiaosiyuan0/mole/ecma/parser"
	"github.com/hsiaosiyuan0/mole/span"
)

func main() {
	// imitate the source code you want to parse
	code := `console.log("hello world")`

	// create a Source instance to handle to the source code
	s := span.NewSource("", code)

	// create a parser, here we use the default options
	opts := parser.NewParserOpts()
	p := parser.NewParser(s, opts)

	// inform the parser do its parsing process
	ast, err := p.Prog()
	if err != nil {
		log.Fatal(err)
	}

	// by default the parsed AST is not the ESTree form because the latter has a little redundancy,
	// however Mole supports to convert its AST to ESTree by using the `estree.ConvertProg` function
	b, err := json.Marshal(estree.ConvertProg(ast.(*parser.Prog), estree.NewConvertCtx()))
	if err != nil {
		log.Fatal(err)
	}

	// below is nothing new, we just print the ESTree in JSON form
	var out bytes.Buffer
	json.Indent(&out, b, "", "  ")
	fmt.Println(out.String())
}
```

The produced AST can be consumed by the ast-walker in Mole, more runnable demos see [mole-demo](https://github.com/hsiaosiyuan0/mole-demo)

</details>

## Development

See [dev.md](/docs/dev.md) to get more information about how to start development.
