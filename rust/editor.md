# editor


## git
- [pre-commit](https://github.com/tyrchen/geektime-rust/blob/master/.pre-commit-config.yaml)

	在根目录生成 .pre-commit-config.yaml 后，运行 pre-commit install，以后 git commit 时就会自动做这一系列的检查，保证提交代码的最基本的正确性

	除此之外，在根目录下最好还声明一个 deny.toml，使用 cargo-deny 来确保使用的第三方依赖没有不该出现的授权（比如不使用任何 GPL/APGL 的代码）、没有可疑的来源（比如不是来自某个 fork 的 GitHub repo 下的 commit），以及没有包含有安全漏洞的版本

## workflow
- [github workflow](https://github.com/tyrchen/geektime-rust/blob/master/.github/workflows/build.yml)

## vscode
### 调试环境配置
ref:
	- [在 VSCode 中配置 Rust 工程](https://www.runoob.com/rust/cargo-tutorial.html)

## vim
参考:
- [将 Vim 设置为 Rust IDE](https://linux.cn/article-12530-1.html)

```bash
$ vim ~/.vimrc
filetype plugin indent on    # 打开检测、插件和缩进配置
syntax on                    # 启用语法高亮
```

## 插件
### sublime3插件
```
Rust Enhanced // 需先通过package control的disable package禁用st3自带的Rust插件
RustAutoComplete
```

### vscode插件
- rust-analyzer: 它会实时编译和分析 Rust 代码，提示代码中的错误，并对类型进行标注; 也可以使用官方的 rust 插件取代
- rust syntax：为代码提供语法高亮
- crates：帮助分析当前项目的依赖是否是最新的版本
- better toml：Rust 使用 toml 做项目的配置管理. better toml 可以帮你语法高亮, 并展示 toml 文件中的错误
- rust test lens：可以帮快速运行某个 Rust 测试
- Tabnine：基于 AI 的自动补全，可以帮助开发者更快地撰写代码