# Duplicated package analysis

Duplicated package analysis can find out the duplicated packages in the project's dependency graph. The definition of the duplicated packages is that there are multiple versions of same package in the project's dependency graph, each version will be considered by the bundle tools as a individual package.

All the duplicated packages will be compiled into the final bundle which leads to the the bundle size problem.

The strategy used for finding out the duplicated package can be briefly described as below:

1. Take some files as the entry points to start the dependency graph analysis.
2. Build the dependency graph, for each new file, treat it as a node in the graph, parse out which their imports to introduce new files to expand the graph until no new file being introduced.
3. For each node in the graph, gather the basic information of it such as its name and version(if it's a umbrella node) as well as the relation information that it was introduced by which nodes and it will introduce which nodes.
4. Traverse the final graph to find out the duplicated packages since each duplicated package is a individual node in it.

## Usage

Use `mole.json` to configure the analyzer:

```json
{
  "target": "react-native", // web, node, react-native
  "entries": ["./src/index.js"] // the entry points
}
```

Get the analysis report by below command:

```
npx molecast -pkg-ana
```
