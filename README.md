# Mole

Mole is a toolkit written in Golang to facilitate our frontend development experience and will provide various functionalities like lint, code-format, bundle, .etc.

## Features

- [x] JavaScript Parser
  - Supports syntaxes up to [ES2021](https://262.ecma-international.org/12.0/) and [JSX](https://github.com/facebook/jsx)
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
