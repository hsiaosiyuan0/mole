## Preparations

1. Setup Golang development environment by follow the official [guide](https://go.dev/doc/install)
2. Clone this repo and change working directory to the local directory which handles the cloned stuff
4. Run `setup.sh` under the root of the working directory to do some initializations which currently includes:
   - Create a `.env` file which handle the envs for devtools

3. Run below command to resolve the dependencies:

   ```bash
   make dep
   ```
4. Run below command to run tests

   ```bash
   make test
   ```

## Debugging

[The VS Code Go extension](https://marketplace.visualstudio.com/items?itemName=golang.go) is recommended to easily debug code by:

   1. Choose a test in any `_test.go` file as the entrypoint of the debugging
   2. Create a breakpoint by clicking the [editor margin](https://code.visualstudio.com/docs/editor/debugging#_breakpoints)
   3. Click the `debug test` [CodeLen](https://code.visualstudio.com/blogs/2017/02/12/code-lens-roundup) on the top of the test function name
   4. If the breakpoint in step 2 is on the road of the current test threading then the process will be hung up and the editor will be located to the source location of that breakpoint automatically