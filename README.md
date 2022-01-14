# Mole

Mole is a toolkit written in Golang to facilitate frontend development experience by providing various functionalities like lint, code-format, bundle, .etc.

<details>
  <summary>Why Golang</summary>

A little bit explanation is good for why Golang is preferred in this project. Nowadays, a programming language is not only the grammar things, it's consist of runtime, stdlib, 3rd-party modules and a healthy community, all these are out-of-box by using Golang, more specifically:

- Golang is productive, its simplicity philosophy(something like Grammar and Garbage-collection) saves more time to the functionalities themselves.
- the functionalities like lint and bundle maybe needed to run as web services while Golang has been proved by many impressive projects such k8s that it's good at service things.
</details>

## Features

- [x] JavaScript Parser
  - Supports syntaxes up to [ES2021](https://262.ecma-international.org/12.0/) and [JSX](https://github.com/facebook/jsx)
  - [ESTree](https://github.com/estree/estree) compatible output, [AST explorer on WASM](http://blog.thehardways.me/mole-is-more/#/)
- [ ] TypeScript Parser **WIP**
- [ ] Indenter(like prettier)
- [ ] CSS/Less/Scss

## Preview

It's easy for OSX users to get a preview binary via [brew](https://brew.sh/):

```bash
brew tap hsiaosiyuan0/mole
brew install mole
```

For users on the other platforms could download the executable binary on [Releases](https://github.com/hsiaosiyuan0/mole/releases)

Run below command to test the javascript parser shipped within mole binary:

```bash
mole -ast -file path_to_your_test_file.js
```

## Development

1. Setup Golang environment by follow the official [guide](https://go.dev/doc/install)
2. Clone this repo then change working directory to the local directory which contains the cloned stuff
3. Run below command to resolve the dependencies:

   ```bash
   make dep
   ```
4. Run below command to run tests

   ```bash
   make test
   ```
- [The VS Code Go extension](https://marketplace.visualstudio.com/items?itemName=golang.go) is recommended to easily debug code by:

   1. Choose a test as the entrypoint of debugging
   2. Create a breakpoint by clicking the [editor margin](https://code.visualstudio.com/docs/editor/debugging#_breakpoints)
   3. Click the `debug test` [CodeLen](https://code.visualstudio.com/blogs/2017/02/12/code-lens-roundup) on the top of the test function name
   4. If the breakpoint in step 2 is on the line of the current test threading then the process will be hung up and the editor will be located to the location of that breakpoint automatically