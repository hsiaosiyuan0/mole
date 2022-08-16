# Duplicated package analysis

Duplicated package analysis is used to out the duplicated packages in your project's dependency graph.

The definition of the duplicated packages is that there are multiple versions of the same package in the dependency graph and each version will be treated by the bundle tools as if they are individual package, and be compiled into the final bundle leads to the the bundle size problem.

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

Get the analysis report by using below command:

```
npx molecast -pkg-ana
```
