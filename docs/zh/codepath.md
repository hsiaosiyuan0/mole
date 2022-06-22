## AST 与 CFG

在进行静态分支之前，要尽可能的还原和收集程序的信息，用来存放这些信息的两个典型的数据结构就是 AST 和 CFG

AST 是一个树形结构，用来存放解析后的程序语法结构信息。之所以选择树形结构，是因为大部分程序语言的语法都是嵌套结构的，使用树形结构存放可以很好的体现节点之前的语法关系

CFG（Control-flow Graph）则是一个图结构，用来存放程序中指令执行控制流程的流转信息。之所以选择图结构，是因为指令控制流程并不是单向的（考虑跳转指令）

CodePath（ESLint 中的概念），其实就是 CFG

AST 是第一步，这很好理解，因为还原语法结构是起点。有了 AST 之后，就可以构造出 CFG

## Dot notation

在 ESLint 中，为了可视化调试 CFG，会将其中的数据序列化成 [Dot notation](https://graphviz.gitlab.io/doc/info/lang.html) 输出到控制台, 粘贴输出的结果到 [graphviz editor](http://magjac.com/graphviz-visual-editor/) 进行可视化

可以在 [eslint](https://github.com/eslint/eslint.git) 项目根目录中创建测试文件 test.js：

```js
a;
```

在根目录下执行下面的命令即可查看上面代码对应的 dot notation：

```
DEBUG=eslint:code-path ./bin/eslint.js --no-ignore test.js
```

![image](<https://g.gravizo.com/svg?digraph%20G%20%7B%0Agraph%5Bcenter%3Dtrue%20pad%3D.5%5D%0Anode%5Bshape%3Dbox%2Cstyle%3D%22rounded%2Cfilled%22%2Cfillcolor%3Dwhite%5D%3B%0Ainitial%5Blabel%3D%22%22%2Cshape%3Dcircle%2Cstyle%3Dfilled%2Cfillcolor%3Dblack%2Cwidth%3D0.25%2Cheight%3D0.25%5D%3B%0Afinal%5Blabel%3D%22%22%2Cshape%3Ddoublecircle%2Cstyle%3Dfilled%2Cfillcolor%3Dblack%2Cwidth%3D0.25%2Cheight%3D0.25%5D%3B%0As1_1%5Blabel%3D%22Program%3Aenter%5CnExpressionStatement%3Aenter%5CnIdentifier%20(a)%5CnExpressionStatement%3Aexit%5CnProgram%3Aexit%22%5D%3B%0Ainitial-%3Es1_1-%3Efinal%3B%0A%7D>)

> mole 目前也是输出 dot notation，[为什么使用 Dot 而不是 Mermaid 输出调试信息](https://github.com/x-orpheus/mole/issues/4)

## CFG 基本概念

在根据 AST 构造 CFG 之前，需要了解 CFG 的一些基本概念

### 基本块

首先 CFG 中也有节点的概念，但 CFG Node 和 AST Node 不是对等的关系，可以观察下面的代码：

```js
a;
b;
```

如果是 AST，那么其结构应该是：

<img src='https://g.gravizo.com/svg?digraph%20G%20%7B%0A%20%20graph%5Bcenter%3Dtrue%20pad%3D.5%5D%0A%20%20node%5Bshape%3Dbox%2Cstyle%3D%22rounded%2Cfilled%22%2Cfillcolor%3Dwhite%5D%3B%0A%20%20initial%5Blabel%3D%22%22%2Cshape%3Dcircle%2Cstyle%3Dfilled%2Cfillcolor%3Dblack%2Cwidth%3D0.25%2Cheight%3D0.25%5D%3B%0A%20%20final%5Blabel%3D%22%22%2Cshape%3Ddoublecircle%2Cstyle%3Dfilled%2Cfillcolor%3Dblack%2Cwidth%3D0.25%2Cheight%3D0.25%5D%3B%0A%20%20initial%20-%3E%20a%0A%20%20a%20-%3E%20b%0A%20%20b%20-%3E%20final%0A%7D'/>

注意 a 和 b 是通过两个节点分别表示，而如果在 CFG 中，会表示成下面这样：

<img src='https://g.gravizo.com/svg?%0Adigraph%20G%20%7B%0A%20%20graph%5Bcenter%3Dtrue%20pad%3D.5%5D%0A%20%20node%5Bshape%3Dbox%2Cstyle%3D%22rounded%2Cfilled%22%2Cfillcolor%3Dwhite%5D%3B%0A%20%20initial%5Blabel%3D%22%22%2Cshape%3Dcircle%2Cstyle%3Dfilled%2Cfillcolor%3Dblack%2Cwidth%3D0.25%2Cheight%3D0.25%5D%3B%0A%20%20final%5Blabel%3D%22%22%2Cshape%3Ddoublecircle%2Cstyle%3Dfilled%2Cfillcolor%3Dblack%2Cwidth%3D0.25%2Cheight%3D0.25%5D%3B%0A%20%20a%5Blabel%3D%22a%5Cnb%22%5D%0A%20%20initial%20-%3E%20a%0A%20%20a%20-%3E%20final%0A%7D%0A'/>

可以发现 a 与 b 合并到了一起

CFG 中的节点，又被称为 基本块（Basic Block），基本块的定义很简单：

1. 所有进入基本块的线路，必须由基本块头部进入
2. 所有从基本块出去的线路，必须由基本块的尾部输出。换句话说，基本块中间（除去开头和结尾的指令）不能存在进出的线路
3. 当然为了描述的严谨，上述规则不应用于基本块中只包含一条指令的情况
4. 基本块中的指令都是按序执行的

另一个需要注意的概念是，CFG 中，是以指令为最小颗粒度的，比如说：

```js
a + b;
```

需要表示成：

<img src='https://g.gravizo.com/svg?digraph%20G%20%7B%0A%20%20graph%5Bcenter%3Dtrue%20pad%3D.5%5D%0A%20%20node%5Bshape%3Dbox%2Cstyle%3D%22rounded%2Cfilled%22%2Cfillcolor%3Dwhite%5D%3B%0A%20%20initial%5Blabel%3D%22%22%2Cshape%3Dcircle%2Cstyle%3Dfilled%2Cfillcolor%3Dblack%2Cwidth%3D0.25%2Cheight%3D0.25%5D%3B%0A%20%20final%5Blabel%3D%22%22%2Cshape%3Ddoublecircle%2Cstyle%3Dfilled%2Cfillcolor%3Dblack%2Cwidth%3D0.25%2Cheight%3D0.25%5D%3B%0A%20%20a%5Blabel%3D%22binExpr%3Aenter%5Cna%5Cnb%5CnbinExpr%3Aexit%22%5D%0A%20%20initial%20-%3E%20a%0A%20%20a%20-%3E%20final%0A%7D'/>

而不是：

<img src='https://g.gravizo.com/svg?digraph%20G%20%7B%0A%20%20graph%5Bcenter%3Dtrue%20pad%3D.5%5D%0A%20%20node%5Bshape%3Dbox%2Cstyle%3D%22rounded%2Cfilled%22%2Cfillcolor%3Dwhite%5D%3B%0A%20%20initial%5Blabel%3D%22%22%2Cshape%3Dcircle%2Cstyle%3Dfilled%2Cfillcolor%3Dblack%2Cwidth%3D0.25%2Cheight%3D0.25%5D%3B%0A%20%20final%5Blabel%3D%22%22%2Cshape%3Ddoublecircle%2Cstyle%3Dfilled%2Cfillcolor%3Dblack%2Cwidth%3D0.25%2Cheight%3D0.25%5D%3B%0A%20%20a%5Blabel%3D%22a%20%2B%20b%22%5D%0A%20%20initial%20-%3E%20a%0A%20%20a%20-%3E%20final%0A%7D'/>

将 `a + b` 拆开的原因是，这个表达式是无法通过一个机器指令执行的，通常来说需要分为三条指令：

1. 将 a 变量的内容载入某个寄存器
2. 将 b 变量的内容载入某另一个寄存器
3. 执行加法指令

在 ESLint 中的 CodePath 也是以指令的颗粒度进行的处理。那么就可能有疑问，ESLint 中是否有必要以指令这么细的颗粒度来构建 CodePath，对于像 `a + b` 这样没有跳转的语句，是不是视为整体就可以了

采用怎样的颗粒度其实是可选的，指令级别的颗粒度是比较常见的做法，可以适应后续大部分使用情况

### 控制流程

既然 CFG 是将控制流程还原出来，那么就需要知道程序中有哪些控制流程

本质上程序的控制流程只有下面几个，它们和机器指令的语义所对应：

- 按序执行
- 跳转到指定位置继续执行
  - 有条件跳转
  - 无条件跳转

在汇编语言中是存在与机器指令对等的控制流程的

到了高级语言中，提供了更加方便的所谓结构化控制流，就是常见的：

1. logicExpr (andExpr, orExpr)
2. ifStmt
3. repeatStmt（forStmt, whileStmt, etc.）
4. gotoStmt (gotoStmt, contStmt, breakStmt, callExpr etc.)

前三个是常见的高级语言中都提供的，第 4 个放在高级语言中则大概率属于过度设计，比如 JavaScript 中：

```js
outmost: while (1) {
  while (1) {
    break outmost;
  }
}
console.log("ok");
```

上面的代码中做了以下的事情：

- 使用 `outmost` 标记最外层的循环
- 在最内的循环中，使用带参数的 break 语句 `break outmost;` 跳出最外层的循环

像是上面这样的跳转，就是无条件跳转

有条件跳转则是下面的例子：

```js
a || b;
c;
```

在 a 的位置会有一个有条件跳转：

- 如果 a 为 true，则跳转到 c 的开头处继续执行
- 否则跳转到 b 的开头处继续执行

## CFG 的构建方式

要一次为整段代码构建 CFG 是比较难的。因此在 mole 中 CFG 的构建方式，核心原理是：

1. 为每个 AST 节点创建其对应的 CFG subgraph
2. 将连续的 CFG subgraph 连接起来

更具体地说，连接就是将当前产生的 subgraph 和上一个 subgraph 进行连接，因此就涉及到两个问题：

1.  当前 subgraph 的开头，可能可以和上一个 subgraph 的结尾合并 - 因为要生成基本块
2.  上一个 subgraph 可能有多个输出（考虑存在跳转的情况），连接的时候要确保衔接正确

看一个具体的例子：

```js
a || b;
c;
```

为上面的代码生成 CFG 的流程为：

1. 生成 a 的 subgraph

  <img src='https://g.gravizo.com/svg?%20%20digraph%20G%20%7B%0A%20%20%20%20graph%5Bcenter%3Dtrue%20pad%3D.5%5D%0A%20%20%20%20node%5Bshape%3Dbox%2Cstyle%3D%22rounded%2Cfilled%22%2Cfillcolor%3Dwhite%5D%3B%0A%20%20%20%20initial%5Blabel%3D%22%22%2Cshape%3Dcircle%2Cstyle%3Dfilled%2Cfillcolor%3Dblack%2Cwidth%3D0.25%2Cheight%3D0.25%5D%3B%0A%20%20%20%20final%5Blabel%3D%22%22%2Cshape%3Ddoublecircle%2Cstyle%3Dfilled%2Cfillcolor%3Dblack%2Cwidth%3D0.25%2Cheight%3D0.25%5D%3B%0A%20%20%20%20initial%20-%3E%20a%0A%20%20%20%20a%20-%3E%20out1%0A%20%20%20%20a%20-%3E%20out2%5Bcolor%3D%22orange%22%2C%20label%3D%22T%22%5D%0A%20%20%20%20out1%20-%3E%20final%0A%20%20%20%20out2%20-%3E%20final%0A%20%20%7D'/>

> - 黑色表示 Normal 路径
> - 黄色表示 Jump 路径
> - 红色表示 Unreachable 路径
> - out1, out2 表示待定的未知连接点

2. 生成 b 的 subgraph

  <img src='https://g.gravizo.com/svg?%20%20digraph%20G%20%7B%0A%20%20%20%20graph%5Bcenter%3Dtrue%20pad%3D.5%5D%0A%20%20%20%20node%5Bshape%3Dbox%2Cstyle%3D%22rounded%2Cfilled%22%2Cfillcolor%3Dwhite%5D%3B%0A%20%20%20%20initial%5Blabel%3D%22%22%2Cshape%3Dcircle%2Cstyle%3Dfilled%2Cfillcolor%3Dblack%2Cwidth%3D0.25%2Cheight%3D0.25%5D%3B%0A%20%20%20%20final%5Blabel%3D%22%22%2Cshape%3Ddoublecircle%2Cstyle%3Dfilled%2Cfillcolor%3Dblack%2Cwidth%3D0.25%2Cheight%3D0.25%5D%3B%0A%20%20%20%20initial%20-%3E%20b%0A%20%20%20%20b%20-%3E%20out1%5B%5D%0A%20%20%20%20out1%20-%3E%20final%0A%20%20%7D'/>

3. 将 b 的 subgraph 和 a 的 subgraph 连接

  <img src='https://g.gravizo.com/svg?%20%20digraph%20G%20%7B%0A%20%20%20%20graph%5Bcenter%3Dtrue%20pad%3D.5%5D%0A%20%20%20%20node%5Bshape%3Dbox%2Cstyle%3D%22rounded%2Cfilled%22%2Cfillcolor%3Dwhite%5D%3B%0A%20%20%20%20initial%5Blabel%3D%22%22%2Cshape%3Dcircle%2Cstyle%3Dfilled%2Cfillcolor%3Dblack%2Cwidth%3D0.25%2Cheight%3D0.25%5D%3B%0A%20%20%20%20final%5Blabel%3D%22%22%2Cshape%3Ddoublecircle%2Cstyle%3Dfilled%2Cfillcolor%3Dblack%2Cwidth%3D0.25%2Cheight%3D0.25%5D%3B%0A%20%20%20%20initial%20-%3E%20a%0A%20%20%20%20a%20-%3E%20out1%5Bcolor%3D%22orange%22%2Clabel%3D%22T%22%5D%0A%20%20%20%20a%20-%3E%20b%0A%20%20%20%20b%20-%3E%20out2%0A%20%20%20%20out1%20-%3E%20final%0A%20%20%20%20out2%20-%3E%20final%0A%20%20%7D'/>

4. 为 c 生成其 subgraph

  <img src='https://g.gravizo.com/svg?%20%20digraph%20G%20%7B%0A%20%20%20%20graph%5Bcenter%3Dtrue%20pad%3D.5%5D%0A%20%20%20%20node%5Bshape%3Dbox%2Cstyle%3D%22rounded%2Cfilled%22%2Cfillcolor%3Dwhite%5D%3B%0A%20%20%20%20initial%5Blabel%3D%22%22%2Cshape%3Dcircle%2Cstyle%3Dfilled%2Cfillcolor%3Dblack%2Cwidth%3D0.25%2Cheight%3D0.25%5D%3B%0A%20%20%20%20final%5Blabel%3D%22%22%2Cshape%3Ddoublecircle%2Cstyle%3Dfilled%2Cfillcolor%3Dblack%2Cwidth%3D0.25%2Cheight%3D0.25%5D%3B%0A%20%20%20%20initial%20-%3E%20c%0A%20%20%20%20c%20-%3E%20out1%0A%20%20%20%20out1%20-%3E%20final%0A%20%20%7D'/>

5. 将 c 的 subgraph 和 a、b 合并后的 subgraph 连接

  <img src='https://g.gravizo.com/svg?%20%20digraph%20G%20%7B%0A%20%20%20%20graph%5Bcenter%3Dtrue%20pad%3D.5%5D%0A%20%20%20%20node%5Bshape%3Dbox%2Cstyle%3D%22rounded%2Cfilled%22%2Cfillcolor%3Dwhite%5D%3B%0A%20%20%20%20initial%5Blabel%3D%22%22%2Cshape%3Dcircle%2Cstyle%3Dfilled%2Cfillcolor%3Dblack%2Cwidth%3D0.25%2Cheight%3D0.25%5D%3B%0A%20%20%20%20final%5Blabel%3D%22%22%2Cshape%3Ddoublecircle%2Cstyle%3Dfilled%2Cfillcolor%3Dblack%2Cwidth%3D0.25%2Cheight%3D0.25%5D%3B%0A%20%20%20%20initial%20-%3E%20a%0A%20%20%20%20a%20-%3E%20b%0A%20%20%20%20a%20-%3E%20c%5Bcolor%3D%22orange%22%2Clabel%3D%22T%22%5D%0A%20%20%20%20b%20-%3E%20c%0A%20%20%20%20c%20-%3E%20out1%0A%20%20%20%20out1%20-%3E%20final%0A%20%20%7D'/>

由此可见 CFG 的构建过程是先分治在合并的的过程

上面的模式中，还有一个规律点 - 高级语言中结构化后的控制流，既然是结构化后的，其流程是固定可循的，比如：

```js
a || b;
```

固定为：

  <img src='https://g.gravizo.com/svg?digraph%20G%20%7B%0A%20%20graph%5Bcenter%3Dtrue%20pad%3D.5%5D%0A%20%20%20%20node%5Bshape%3Dbox%2Cstyle%3D%22rounded%2Cfilled%22%2Cfillcolor%3Dwhite%5D%3B%0A%20%20%20%20initial%5Blabel%3D%22%22%2Cshape%3Dcircle%2Cstyle%3Dfilled%2Cfillcolor%3Dblack%2Cwidth%3D0.25%2Cheight%3D0.25%5D%3B%0A%20%20%20%20final%5Blabel%3D%22%22%2Cshape%3Ddoublecircle%2Cstyle%3Dfilled%2Cfillcolor%3Dblack%2Cwidth%3D0.25%2Cheight%3D0.25%5D%3B%0A%20%20%20%20initial%20-%3E%20a%0A%20%20%20%20a%20-%3E%20out1%5Bcolor%3D%22orange%22%2Clabel%3D%22T%22%5D%0A%20%20%20%20a%20-%3E%20b%0A%20%20%20%20b%20-%3E%20out2%0A%20%20%20%20out1%20-%3E%20final%0A%20%20%20%20out2%20-%3E%20final%0A%7D'/>

而：

```js
a && b;
```

固定为：

  <img src='https://g.gravizo.com/svg?%20%20digraph%20G%20%7B%0A%20%20%20%20graph%5Bcenter%3Dtrue%20pad%3D.5%5D%0A%20%20%20%20node%5Bshape%3Dbox%2Cstyle%3D%22rounded%2Cfilled%22%2Cfillcolor%3Dwhite%5D%3B%0A%20%20%20%20initial%5Blabel%3D%22%22%2Cshape%3Dcircle%2Cstyle%3Dfilled%2Cfillcolor%3Dblack%2Cwidth%3D0.25%2Cheight%3D0.25%5D%3B%0A%20%20%20%20final%5Blabel%3D%22%22%2Cshape%3Ddoublecircle%2Cstyle%3Dfilled%2Cfillcolor%3Dblack%2Cwidth%3D0.25%2Cheight%3D0.25%5D%3B%0A%20%20%20%20initial%20-%3E%20a%0A%20%20%20%20a%20-%3E%20b%0A%20%20%20%20b%20-%3E%20out2%0A%20%20%20%20a%20-%3E%20out1%5Bcolor%3D%22orange%22%2Clabel%3D%22F%22%5D%0A%20%20%20%20out1%20-%3E%20final%0A%20%20%20%20out2%20-%3E%20final%0A%20%20%7D%27'/>

当然其他结构比如 ifStmt，whileStmt 也是如此有迹可循

## 结构化控制流

这一节会将 JavaScript 中涉及的结构化流程对应的 subgraph 梳理出来

### andExpr

```js
a && b;
```

对应：

<img src='https://g.gravizo.com/svg?digraph%20G%20%7B%0A%20%20graph%5Bcenter%3Dtrue%20pad%3D.5%5D%0A%20%20node%5Bshape%3Dbox%2Cstyle%3D%22rounded%2Cfilled%22%2Cfillcolor%3Dwhite%5D%3B%0A%20%20initial%5Blabel%3D%22%22%2Cshape%3Dcircle%2Cstyle%3Dfilled%2Cfillcolor%3Dblack%2Cwidth%3D0.25%2Cheight%3D0.25%5D%3B%0A%20%20final%5Blabel%3D%22%22%2Cshape%3Ddoublecircle%2Cstyle%3Dfilled%2Cfillcolor%3Dblack%2Cwidth%3D0.25%2Cheight%3D0.25%5D%3B%0A%20%20initial%20-%3E%20a%0A%20%20a%20-%3E%20b%0A%20%20a%20-%3E%20out1%5Bcolor%3D%22orange%22%2Clabel%3D%22F%22%5D%0A%20%20b%20-%3E%20out2%0A%20%20out1%20-%3E%20final%0A%20%20out2%20-%3E%20final%0A%7D'/>

### orExpr

```js
a || b;
```

对应：

<img src='https://g.gravizo.com/svg?digraph%20G%20%7B%0A%20%20graph%5Bcenter%3Dtrue%20pad%3D.5%5D%0A%20%20node%5Bshape%3Dbox%2Cstyle%3D%22rounded%2Cfilled%22%2Cfillcolor%3Dwhite%5D%3B%0A%20%20initial%5Blabel%3D%22%22%2Cshape%3Dcircle%2Cstyle%3Dfilled%2Cfillcolor%3Dblack%2Cwidth%3D0.25%2Cheight%3D0.25%5D%3B%0A%20%20final%5Blabel%3D%22%22%2Cshape%3Ddoublecircle%2Cstyle%3Dfilled%2Cfillcolor%3Dblack%2Cwidth%3D0.25%2Cheight%3D0.25%5D%3B%0A%20%20initial%20-%3E%20a%0A%20%20a%20-%3E%20out1%5Bcolor%3D%22orange%22%2Clabel%3D%22T%22%5D%0A%20%20a%20-%3E%20b%0A%20%20b%20-%3E%20out2%0A%20%20out1%20-%3E%20final%0A%20%20out2%20-%3E%20final%0A%7D'/>

### ifStmt

```js
if (test) cons;
```

对应：

<img src='https://g.gravizo.com/svg?digraph%20G%20%7B%0A%20%20graph%5Bcenter%3Dtrue%20pad%3D.5%5D%0A%20%20node%5Bshape%3Dbox%2Cstyle%3D%22rounded%2Cfilled%22%2Cfillcolor%3Dwhite%5D%3B%0A%20%20initial%5Blabel%3D%22%22%2Cshape%3Dcircle%2Cstyle%3Dfilled%2Cfillcolor%3Dblack%2Cwidth%3D0.25%2Cheight%3D0.25%5D%3B%0A%20%20final%5Blabel%3D%22%22%2Cshape%3Ddoublecircle%2Cstyle%3Dfilled%2Cfillcolor%3Dblack%2Cwidth%3D0.25%2Cheight%3D0.25%5D%3B%0A%20%20initial%20-%3E%20test%0A%20%20test%20-%3E%20cons%0A%20%20test%20-%3E%20out1%5Bcolor%3D%22orange%22%2Clabel%3D%22F%22%5D%0A%20%20cons%20-%3E%20out2%0A%20%20out1%20-%3E%20final%0A%20%20out2%20-%3E%20final%0A%7D'/>

```js
if (test) cons;
else alt;
```

对应：

<img src='https://g.gravizo.com/svg?digraph%20G%20%7B%0A%20%20graph%5Bcenter%3Dtrue%20pad%3D.5%5D%0A%20%20node%5Bshape%3Dbox%2Cstyle%3D%22rounded%2Cfilled%22%2Cfillcolor%3Dwhite%5D%3B%0A%20%20initial%5Blabel%3D%22%22%2Cshape%3Dcircle%2Cstyle%3Dfilled%2Cfillcolor%3Dblack%2Cwidth%3D0.25%2Cheight%3D0.25%5D%3B%0A%20%20final%5Blabel%3D%22%22%2Cshape%3Ddoublecircle%2Cstyle%3Dfilled%2Cfillcolor%3Dblack%2Cwidth%3D0.25%2Cheight%3D0.25%5D%3B%0A%20%20initial%20-%3E%20test%0A%20%20test%20-%3E%20cons%0A%20%20test%20-%3E%20alt%5Bcolor%3D%22orange%22%2Clabel%3D%22F%22%5D%0A%20%20alt%20-%3E%20out1%0A%20%20cons%20-%3E%20out2%0A%20%20out1%20-%3E%20final%0A%20%20out2%20-%3E%20final%0A%7D'/>

### forStmt

```js
for (init; test; update) body;
```

对应：

<img src='https://g.gravizo.com/svg?digraph%20G%20%7B%0A%20%20graph%5Bcenter%3Dtrue%20pad%3D.5%5D%0A%20%20node%5Bshape%3Dbox%2Cstyle%3D%22rounded%2Cfilled%22%2Cfillcolor%3Dwhite%5D%3B%0A%20%20initial%5Blabel%3D%22%22%2Cshape%3Dcircle%2Cstyle%3Dfilled%2Cfillcolor%3Dblack%2Cwidth%3D0.25%2Cheight%3D0.25%5D%3B%0A%20%20final%5Blabel%3D%22%22%2Cshape%3Ddoublecircle%2Cstyle%3Dfilled%2Cfillcolor%3Dblack%2Cwidth%3D0.25%2Cheight%3D0.25%5D%3B%0A%20%20initial%20-%3E%20init%0A%20%20init%20-%3E%20test%0A%20%20test%20-%3E%20body%0A%20%20body%20-%3E%20update%0A%20%20update%20-%3E%20test%5Bcolor%3D%22orange%22%2Ctailport%3Ds%2Cheadport%3Dne%5D%0A%20%20test%20-%3E%20out1%5Bcolor%3D%22orange%22%2Clabel%3D%22F%22%5D%0A%20%20out1%20-%3E%20final%0A%7D'/>

```js
for (init; ; update) body;
```

对应：

<img src='https://g.gravizo.com/svg?digraph%20G%20%7B%0A%20%20graph%5Bcenter%3Dtrue%20pad%3D.5%5D%0A%20%20node%5Bshape%3Dbox%2Cstyle%3D%22rounded%2Cfilled%22%2Cfillcolor%3Dwhite%5D%3B%0A%20%20initial%5Blabel%3D%22%22%2Cshape%3Dcircle%2Cstyle%3Dfilled%2Cfillcolor%3Dblack%2Cwidth%3D0.25%2Cheight%3D0.25%5D%3B%0A%20%20final%5Blabel%3D%22%22%2Cshape%3Ddoublecircle%2Cstyle%3Dfilled%2Cfillcolor%3Dblack%2Cwidth%3D0.25%2Cheight%3D0.25%5D%3B%0A%20%20initial%20-%3E%20init%0A%20%20init%20-%3E%20body%0A%20%20body%20-%3E%20update%0A%20%20update%20-%3E%20body%5Bcolor%3D%22orange%22%2Ctailport%3Ds%2Cheadport%3Dne%5D%0A%20%20final%0A%7D'/>

## whileStmt

```js
while (test) body;
```

对应：

<img src='https://g.gravizo.com/svg?digraph%20G%20%7B%0A%20%20graph%5Bcenter%3Dtrue%20pad%3D.5%5D%0A%20%20node%5Bshape%3Dbox%2Cstyle%3D%22rounded%2Cfilled%22%2Cfillcolor%3Dwhite%5D%3B%0A%20%20initial%5Blabel%3D%22%22%2Cshape%3Dcircle%2Cstyle%3Dfilled%2Cfillcolor%3Dblack%2Cwidth%3D0.25%2Cheight%3D0.25%5D%3B%0A%20%20final%5Blabel%3D%22%22%2Cshape%3Ddoublecircle%2Cstyle%3Dfilled%2Cfillcolor%3Dblack%2Cwidth%3D0.25%2Cheight%3D0.25%5D%3B%0A%20%20initial%20-%3E%20test%0A%20%20test%20-%3E%20body%0A%20%20body%20-%3E%20test%5Bcolor%3D%22orange%22%2Ctailport%3Ds%2Cheadport%3Dne%5D%0A%20%20test%20-%3E%20final%5Bcolor%3D%22orange%22%2Clabel%3D%22F%22%5D%0A%7D'/>

## doWhile

```js
do body while (test)
```

对应：

<img src='https://g.gravizo.com/svg?digraph%20G%20%7B%0A%20%20graph%5Bcenter%3Dtrue%20pad%3D.5%5D%0A%20%20node%5Bshape%3Dbox%2Cstyle%3D%22rounded%2Cfilled%22%2Cfillcolor%3Dwhite%5D%3B%0A%20%20initial%5Blabel%3D%22%22%2Cshape%3Dcircle%2Cstyle%3Dfilled%2Cfillcolor%3Dblack%2Cwidth%3D0.25%2Cheight%3D0.25%5D%3B%0A%20%20final%5Blabel%3D%22%22%2Cshape%3Ddoublecircle%2Cstyle%3Dfilled%2Cfillcolor%3Dblack%2Cwidth%3D0.25%2Cheight%3D0.25%5D%3B%0A%20%20initial%20-%3E%20body%0A%20%20body%20-%3E%20test%0A%20%20test%20-%3E%20body%5Bcolor%3D%22orange%22%2Ctailport%3Ds%2Cheadport%3Dne%2Clabel%3D%22T%22%5D%0A%20%20test%20-%3E%20final%0A%7D'/>

## continue

```js
while (test) {
  body;
  continue;
}
```

对应：

<img src='https://g.gravizo.com/svg?digraph%20G%20%7B%0A%20%20graph%5Bcenter%3Dtrue%20pad%3D.5%5D%0A%20%20node%5Bshape%3Dbox%2Cstyle%3D%22rounded%2Cfilled%22%2Cfillcolor%3Dwhite%5D%3B%0A%20%20initial%5Blabel%3D%22%22%2Cshape%3Dcircle%2Cstyle%3Dfilled%2Cfillcolor%3Dblack%2Cwidth%3D0.25%2Cheight%3D0.25%5D%3B%0A%20%20final%5Blabel%3D%22%22%2Cshape%3Ddoublecircle%2Cstyle%3Dfilled%2Cfillcolor%3Dblack%2Cwidth%3D0.25%2Cheight%3D0.25%5D%3B%0A%20%20initial%20-%3E%20test%0A%20%20test%20-%3E%20body%0A%20%20body%20-%3E%20continue%0A%20%20continue%20-%3E%20test%5Bcolor%3D%22orange%22%2Ctailport%3Ds%2Cheadport%3Dne%5D%0A%20%20test%20-%3E%20final%5Bcolor%3D%22orange%22%2Clabel%3D%22F%22%5D%0A%7D'/>

```js
LabelA: while (test1) {
  while (test2) {
    body;
    continue LabelA;
  }
}
```

对应：

<img src='https://g.gravizo.com/svg?digraph%20G%20%7B%0A%20%20graph%5Bcenter%3Dtrue%20pad%3D.5%5D%0A%20%20node%5Bshape%3Dbox%2Cstyle%3D%22rounded%2Cfilled%22%2Cfillcolor%3Dwhite%5D%3B%0A%20%20initial%5Blabel%3D%22%22%2Cshape%3Dcircle%2Cstyle%3Dfilled%2Cfillcolor%3Dblack%2Cwidth%3D0.25%2Cheight%3D0.25%5D%3B%0A%20%20final%5Blabel%3D%22%22%2Cshape%3Ddoublecircle%2Cstyle%3Dfilled%2Cfillcolor%3Dblack%2Cwidth%3D0.25%2Cheight%3D0.25%5D%3B%0A%20%20initial%20-%3E%20test1%0A%20%20test1%20-%3E%20test2%0A%20%20test2%20-%3E%20body%0A%20%20body%20-%3E%20continue%0A%20%20continue%20-%3E%20test1%5Bcolor%3D%22orange%22%2Ctailport%3Ds%2Cheadport%3Dne%5D%0A%20%20test1%20-%3E%20final%5Bcolor%3D%22orange%22%2Clabel%3D%22F%22%5D%0A%20%20test2%20-%3E%20final%5Bcolor%3D%22orange%22%2Clabel%3D%22F%22%5D%0A%7D'/>

## break

```js
while (test) {
  body;
  break;
}
```

对应：

<img src='https://g.gravizo.com/svg?digraph%20G%20%7B%0A%20%20graph%5Bcenter%3Dtrue%20pad%3D.5%5D%0A%20%20node%5Bshape%3Dbox%2Cstyle%3D%22rounded%2Cfilled%22%2Cfillcolor%3Dwhite%5D%3B%0A%20%20initial%5Blabel%3D%22%22%2Cshape%3Dcircle%2Cstyle%3Dfilled%2Cfillcolor%3Dblack%2Cwidth%3D0.25%2Cheight%3D0.25%5D%3B%0A%20%20final%5Blabel%3D%22%22%2Cshape%3Ddoublecircle%2Cstyle%3Dfilled%2Cfillcolor%3Dblack%2Cwidth%3D0.25%2Cheight%3D0.25%5D%3B%0A%20%20initial%20-%3E%20test%0A%20%20test%20-%3E%20body%0A%20%20body%20-%3E%20break%0A%20%20break%20-%3E%20final%0A%20%20test%20-%3E%20final%5Bcolor%3D%22orange%22%2Clabel%3D%22F%22%5D%0A%7D'/>

```js
LabelA: while (test1) {
  while (test2) {
    body;
    break LabelA;
    unreachable;
  }
}
```

对应：

<img src='https://g.gravizo.com/svg?digraph%20G%20%7B%0A%20%20graph%5Bcenter%3Dtrue%20pad%3D.5%5D%0A%20%20node%5Bshape%3Dbox%2Cstyle%3D%22rounded%2Cfilled%22%2Cfillcolor%3Dwhite%5D%3B%0A%20%20initial%5Blabel%3D%22%22%2Cshape%3Dcircle%2Cstyle%3Dfilled%2Cfillcolor%3Dblack%2Cwidth%3D0.25%2Cheight%3D0.25%5D%3B%0A%20%20final%5Blabel%3D%22%22%2Cshape%3Ddoublecircle%2Cstyle%3Dfilled%2Cfillcolor%3Dblack%2Cwidth%3D0.25%2Cheight%3D0.25%5D%3B%0A%20%20initial%20-%3E%20test1%0A%20%20test1%20-%3E%20test2%0A%20%20test2%20-%3E%20body%0A%20%20body%20-%3E%20break%0A%20%20break%20-%3E%20final%0A%20%20break%20-%3E%20unreachable%5Bcolor%3D%22red%22%5D%0A%20%20unreachable%20-%3E%20test2%5Bcolor%3D%22red%22%2Cheadport%3Dne%5D%0A%20%20test1%20-%3E%20final%5Bcolor%3D%22orange%22%2Clabel%3D%22F%22%5D%0A%20%20test2%20-%3E%20final%5Bcolor%3D%22orange%22%2Clabel%3D%22F%22%5D%0A%7D'/>

## return

```js
function fn() {
  while (test) {
    return;
    unreachable1;
  }
  stmt;
}
```

对应：

<img src='https://g.gravizo.com/svg?digraph%20G%20%7B%0A%20%20graph%5Bcenter%3Dtrue%20pad%3D.5%5D%0A%20%20node%5Bshape%3Dbox%2Cstyle%3D%22rounded%2Cfilled%22%2Cfillcolor%3Dwhite%5D%3B%0A%20%20initial%5Blabel%3D%22%22%2Cshape%3Dcircle%2Cstyle%3Dfilled%2Cfillcolor%3Dblack%2Cwidth%3D0.25%2Cheight%3D0.25%5D%3B%0A%20%20final%5Blabel%3D%22%22%2Cshape%3Ddoublecircle%2Cstyle%3Dfilled%2Cfillcolor%3Dblack%2Cwidth%3D0.25%2Cheight%3D0.25%5D%3B%0A%20%20initial%20-%3E%20test%0A%20%20test%20-%3E%20return%0A%20%20return%20-%3E%20unreachable1%5Bcolor%3D%22red%22%5D%0A%20%20unreachable1%20-%3E%20test%5Bcolor%3D%22red%22%2Cheadport%3Dne%5D%0A%20%20return%20-%3E%20final%0A%20%20test%20-%3E%20stmt%5Bcolor%3D%22orange%22%2Clabel%3D%22F%22%5D%0A%20%20stmt%20-%3E%20implicitReturn%0A%20%20implicitReturn%20-%3E%20final%0A%7D'/>
