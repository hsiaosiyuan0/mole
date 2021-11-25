# mole

Tools to deeply process the frontend projects.

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

## Features

- [x] JavaScript Parser, [demo](http://blog.thehardways.me/mole-is-more/#/)
- [ ] JSX

## Development


### Project structure 

Project structure is complied with [Standard Go Project Layout](https://github.com/golang-standards/project-layout)
