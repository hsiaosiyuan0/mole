# Mole

Go is the future of Frontend infrastructure.

Mole is a toolkit written in Golang provides various functionalities to process source code of the frontend projects.

## Features

- [x] JavaScript Parser
  - JSX
  - Supports syntaxes up to [ES2021](https://262.ecma-international.org/12.0/)
  - [ESTree](https://github.com/estree/estree) compatible output, [AST explorer on WASM](http://blog.thehardways.me/mole-is-more/#/)
- [ ] Indenter(like prettier)
- [ ] CSS/Less/Scss
- [ ] TypeScript Parser

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
