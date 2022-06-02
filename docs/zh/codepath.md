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

## 构建细节

对于下面的语句：

```js
a && b;
c || d;
```

第一条语句的 subgraph（记为 subgraph-stmt1） 其构建方式为：

1. 为 a 生成 subgraph-a
2. 为 b 生成 subgraph-b
3. 连接 subgraph-a 和 subgraph-b，成为 subgraph-stmt1

上面的逻辑是没有问题的，不过有点不方便 subgraph 的消费。考虑下面的例子

```js
a && b;
return;
c || d;
```

1. 从开头要 return 的 subgraph-stmt1 已经生成完毕
2. 为 c 生成 subgraph-c
3. 为 d 生成 subgraph-d
4. 连接 subgraph-c 和 subgraph-d，成为 subgraph-stmt2
5. 连接 subgraph-stmt1 和 subgraph-stmt2 成为整段代码的 graph

问题就出在第2步：

1. 首先我们知道 return 之后的语句都是 unreachable
2. 「为 c 生成 subgraph-c」 的动作发生在 mole 对 ident 节点的 listener 中处理
3. 当 mole 对 ident 的 listener 回调被执行完毕后，随即执行的将是用户的 listeners
4. 而在用户的 listeners 中，由于此时 c 尚未和之前的节点连接，那么 c 将无法感知自己是 unreachable 的
5. 需要等到上述的第4步骤完成后，subgraph-c 和 subgraph-d 内的节点才能感知自身是 unreachable 的

所以要想用户能够知道 c 是不是 unreachable 的，mole 需要在第4步完成后，再遍历一遍 subgraph-c，将 c 的可达性通知用户

为了解决上面的问题，我们需要及时地将正在生成的 subgraph 与之前的节点进行连接

## Visitor or Listenr

遍历 AST 的模式有两种：Visitor 和 Listenr。mole 在生成 CodePath 时选择了 Listener：

- 使用 Listener 可以利用默认的 Visitor 实现，默认 Visitor 会关联 Symtab
- 不过使用 Listener 遍历时的代码相对 Visitor 而言略显不直观

## 测试方式

测试需要覆盖的点：

- Basic Block 的连通性是否正确
- Basic Block 中包含的 ASTNodes 是否正确

所以第一步就是标识 Basic Block。在 ESLint 的实现中，使用自增的 ID 来标识基本块，这个方案在 mole 中有几个问题：

- 自增的 ID 和遍历的顺序存在耦合，实际上只要构建的结果是正确的，节点遍历的顺序应该是无关紧要的。在耦合的情况下，如果未来更改了遍历的顺序，则会导致单测需要重新编写
- mole 中存在 virtual node 的概念，用于表示 a syntax group of basic blocks，比如上文那些固定格式。这就导致自增 ID 也包含一些不确定性

为了解决上面的问题，mole 中使用 Basic Block 中第一个 ASTNode 的  `LineNumber << 32 | ColumnNumber` 作为其 ID