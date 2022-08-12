# Mole

Mole is a toolkit written in Golang to improve the frontend development by providing various functionalities like lint, code-format, bundle, .etc.

<details>
  <summary>Why Golang</summary>

~~A little bit explanation is good for why Golang is preferred in this project. Nowadays, a programming language is not only the grammar things, it's consist of runtime, stdlib, 3rd-party modules and a healthy community, all these are out-of-box by using Golang, more specifically:~~

- ~~Golang is productive, its simplicity philosophy(something like Grammar and Garbage-collection) saves more time to the functionalities themselves.~~
- ~~the functionalities like lint and bundle maybe needed to run as web services while Golang has been proved by many impressive projects such k8s that it's good at service things.~~

Fine, just all because I'm too fool to use a fancy language
</details>

## Features

- [x] JavaScript Parser
  - ECMAScript up to [ES2021](https://262.ecma-international.org/12.0/)
  - [JSX](https://github.com/facebook/jsx)
  - [ESTree](https://github.com/estree/estree) compatible outputs ([AST explorer on WASM](http://blog.thehardways.me/mole-is-more/#/))
- [x] TypeScript Parser
  - [babel/typescript](https://babeljs.io/docs/en/babel-types#typescript) compatible outputs
- [ ] CSS parser
- [ ] Less parser
- [ ] Scss parser
 
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

See [dev.md](/docs/dev.md) to get more information about how to start development.