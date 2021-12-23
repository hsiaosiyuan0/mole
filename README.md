# Mole

Mole is a toolkit written in Golang to facilitate frontend development experience by providing various functionalities like lint, code-format, bundle, .etc.

A little bit explanation is good for why Golang is preferred in this project. Nowadays, a programming language is not only the grammar things, it's consist of runtime, stdlib, 3rd-party modules and a healthy community, all these are out-of-box by using Golang, more specifically:

- Golang is productive, its simplicity philosophy(something like Grammar and Garbage-collection) saves more time to the functionalities themselves.
- the functionalities like lint and bundle maybe needed to run as web services while Golang has been proved by many impressive projects such k8s that it's good at service things.

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

### Project structure

Project structure is complied with [Standard Go Project Layout](https://github.com/golang-standards/project-layout)
