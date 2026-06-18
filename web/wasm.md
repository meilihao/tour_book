# wasm
## example
```js
import init, { greet, add } from "./pkg/your_wasm_module.js"; // init 函数是 WebAssembly (WASM) 模块的异步初始化加载器，它的作用是为你的 JavaScript 环境准备好运行 WebAssembly 代码所需的一切. 简单来说，它执行了三个关键步骤：下载 .wasm 文件、编译成机器码、实例化（建立 JS 与 WASM 的内存桥梁）

// 1. 必须先调用 init() 来加载和初始化模块
await init(); 

// 2. 初始化完成后，才能安全地调用其他导出的函数
console.log(greet("World")); 
console.log(add(1, 2)); 
```
