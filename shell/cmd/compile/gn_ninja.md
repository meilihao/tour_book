# gn和ninja
参考:
- [Using GN build](https://docs.google.com/presentation/d/15Zwb53JcncHfEwHpnG_PoIbbzQ3GQi_cpujYwbpcbZo/htmlpresent)
- [*google gn构建系统的介绍](https://blog.csdn.net/liweigao01/article/details/96354649)
- [gn入门](https://www.cnblogs.com/xl2432/p/11844943.html)

GN是一种元构建系统，生成Ninja构建文件（Ninja build files）. gn和ninja的关系就与cmake和make的关系差不多(cmake通过编写CMakeLists.txt，可以控制生成的Makefile，从而控制编译过程). gn把.gn文件转换成.ninja文件，然后ninja根据.ninja文件将源码生成目标程序.

关系:
gn = cmake
nijia = make/Makefile

安装:
- gn

    到https://gn.googlesource.com/gn#getting-a-binary下载linux版本, 解压并mv gn到`/usr/bin`
- nijia

    到https://github.com/ninja-build/ninja/releases下载相应版本, 解压并mv nijia到`/usr/bin`

```bash
# gn --version
# ninja --version
```

doc:
- [gn](https://gn.googlesource.com/gn/+/master/docs/)
- [nijia](https://ninja-build.org/manual.html)

gn语法:
- [GN Language and Operation](https://gn.googlesource.com/gn/+/refs/heads/master/docs/language.md)
- [GN Reference](https://gn.googlesource.com/gn/+/refs/heads/master/docs/reference.md)

## [gn example](https://github.com/meilihao/demo/tree/master/gn/hello_world)
参考:
- [minimal-gn-project](https://github.com/skopf/minimal-gn-project)
- [gn官方例子](https://gn.googlesource.com/gn/+/master/examples)

- .gn : .gn文件所在的目录是GN工具认定的工程的source root, 必须

    其内指定了BUILDCONFIG.gn文件的位置, 可通过`gn help dotfile`查看help
- //build/config/BUILDCONFIG.gn :  指定编译时使用的编译工具链. 必须

    由`.gn`文件定义BUILDCONFIG.gn的确切位置.
    用于设置全局变量和默认设置

- BUILD.gn : 编译脚本用于指定编译的目标. 必须

```bash
# tree -a
.
├── build
│   ├── config
│   │   └── BUILDCONFIG.gn
│   └── toolchain
│       └── BUILD.gn
├── BUILD.gn
├── .gn
└── hello_word.cc

3 directories, 5 files
# cat hello_word.cc 
#include <iostream>

int main()
{
        std::cout << "Hello world: gn build example" << std::endl;
        return 0;
}
# cat BUILD.gn 
executable("hello_world") { # executable是指生成exe的名字
  sources = [ # 指源码文件列表，相对目录
    "hello_world.cc",
  ]
}
# cat .gn 
buildconfig = "//build/config/BUILDCONFIG.gn
# cat build/config/BUILDCONFIG.gn # from gn examples#simple_build/build/BUILDCONFIG.gn
# cat build/toolchain/BUILD.gn # from gn examples#simple_build/build/toolchain/BUILD.gn
# gn gen out -v # out为构建目录, `gn gen`会在其中生成ninja脚本
# ninja -C out -v # ninja构建项目
ninja: Entering directory `out'
[1/2] g++ -MMD -MF obj/hello_world.hello_world.o.d     -c ../hello_world.cc -o obj/hello_world.hello_world.o
[2/2] g++  -o hello_world -Wl,--start-group @hello_world.rsp  -Wl,--end-group 
```

gn subcmds:
- gn args <out_dir> : 设置构建参数, 其实就是调用一个文本编辑器, 编辑`gn gen -C out`生成的`<out_dir>/args.gn`, 格式为`k=v`
- gn args <out_dir> --list : 列出可用的构建参数和它们的缺省值
- gn args <out_dir> --list="is_debug" : 仅查看某个参数
- gn clean <out_dir> : 用于对历史编译进行清理. 它会删除输出目录下除了args.gn外的内容，并创建一个可以重新产生构建配置的Ninja构建环境
- gn check : 用来检查头文件依赖关系的有效性
- gn desc out //:hello_world : 查看该target的描述信息, 包括src文件, 依赖的lib, 编译选项等. gn desc也可显示config的描述信息
- gn gen <out_dir> : 创建新的构建目录
- gn gen <out_dir> --args="is_debug=false proprietary_codecs=true" --ide="xxx" : `--args`, gn gen时指定`<out_dir>/args.gn`的配置参数; `--ide`, 生成针对某种ide的工程文件 
- gn ls <out_dir> : 列出所有的target
- gn ls <out_dir> "//:hello_world2" : 列出匹配到的target. 匹配规则不是正则, 可用`gn help label_pattern`查看help
- gn refs out hello_world.cc : 查看依赖该文件的target. gn refs是用来查找反向的依赖(也就是引用了某些东西的targets)
- gn refs <out_dir> //:hello_world : 查看依赖该target的target
- gn help gen : 获取subcmds的help
- gn gen --help : 显示 gn gen 的帮助信息, 与`gn help`类同
- gn help --args : 显示 `--args` 的详细帮助信息

> 上面的`//`代表从项目根目录开始

gn参数:
- -time : 各步骤的耗时
- --tracelog=mylog.trace : 跟踪命令的执行过程. 可导入Chrome的`about:tracing`页面查看

## gn target
target就是gn一个最小的编译单元，可以将它单独传递给ninja进行编译.

从google文档上看有以下几种target:
- action: Declare a target that runs a script a single time.（指定一段指定的脚本）
- action_foreach: Declare a target that runs a script over a set of files.（为一组输入文件分别执行一次脚本）
- bundle_data: [iOS/macOS] Declare a target without output. （声明一个无输出文件的target）
- copy: Declare a target that copies files. （声明一个只是拷贝文件的target）
- create_bundle: [iOS/macOS] Build an iOS or macOS bundle. （编译MACOS/IOS包）
- executable: Declare an executable target. （生成可执行程序）
- generated_file: Declare a generated_file target.
- group: Declare a named group of targets. （执行一组target编译）
- loadable_module: Declare a loadable module target. （创建运行时加载动态连接库，和deps方式有一些区别）
- rust_library: Declare a Rust library target.
- shared_library: Declare a shared library target. （生成动态链接库，.dll or .so）
- source_set: Declare a source set target. (生成静态库，比static_library要快)
- static_library: Declare a static library target. （生成静态链接库，.lib or .a）
- target: Declare an target with the given programmatic type.

因此上面的gn example的hello_world示例其实也只是增加了一个executable target.

## gn 其他
在gn中，使用deps来实现库的依赖关系.

模板，顾名思义，可以用来定义可重用的代码，比如添加新的target类型等.
通常可以将模板单独定义成一个.gni文件，然后其他文件就可以通过import来引入实现共享. 这部分就比较复杂，具体例子可参阅官方文档.

gn编译的toolchain配置非常关键，决定了编译的方式和产物的用途，chromium自带的toolchains也能实现跨平台，但是太过庞大，日常使用的话，可以借鉴：https://github.com/timniederhausen/gn-build

# ninja
参考:
- [Ninja 构建系统](https://blog.csdn.net/yujiawang/article/details/72627121)


ninja 工具介绍

```bash
# ninja -h
options:
  --version  # 打印版本信息（如当前版本是1.10.1）
  -v       # 显示构建中的所有命令行（这个对实际构建的命令核对非常有用）

  -C DIR   # 在执行操作之前，切换到`DIR`目录
  -f FILE  # 制定`FILE`为构建输入文件. 默认文件为当前目录下的`build.ninja`. 如 ninja -f demo.ninja

  -j N     # 并行执行 N 个作业。默认N=3（需要对应的CPU支持）如 ninja -j 2 all
  -l N     # 如果平均负载大于N，不启动新的作业
  -k N     # 持续构建直到N个作业失败为止。默认N=1
  -n       # 排练（dry run）(不执行命令，视其成功执行. 如 ninja -n -t clean)

  -d MODE  # 开启调试模式 (用 -d list 罗列所有的模式)
  -t TOOL  # 执行一个子工具(用 -t list 罗列所有子命令工具)

    - ninja -t clean : 清理构建
    - ninja -t clean pc : 清理某个模块
    - ninja -t browse --port=8000 --no-browser mytarget : 在浏览器中显示编译target的编译依赖图(此命令会启动一个web server服务)
    - ninja -t query mytarget : 查看target的编译过程
    - ninja -t query all : 查看all targets的编译过程
```

ninja还集成了graphviz等一些对开发非常有用的工具, 具体如下：（也就是执行`ninja -t list`的结果）
 
ninja subtools:
    browse  # 在浏览器中浏览依赖关系图（默认会在8080端口启动一个基于python的http服务）
     clean  # 清除构建生成的文件
  commands  # 罗列重新构建制定目标所需的所有命令
      deps  # 显示存储在deps日志中的依赖关系
     graph  # 为指定目标生成 graphviz dot 文件. 如 ninja -t graph all |dot -Tpng -o graph.png
     query  # 显示一个路径的inputs/outputs
   targets  # 通过DAG中rule或depth罗列target
    compdb  # dump JSON兼容的数据库到标准输出
 recompact  # 重新紧凑化ninja内部数据结构

## example
```bash
# ninja -C out :hello_world 仅构建target hello_world
```

运行ninja时, 默认情况下，它会在当前目录中查找名为build.ninja的文件并构建所有过期目标.

Ninja 默认基于系统中可用的 CPU 数量以并发方式执行指令。因为同时运行的命令们的输出可能混淆，Ninja 会在一个命令完成前缓存其输出, 因此从结果看，如同命令是串行的.

## FAQ
### ninja 与 make 比较
Ninja的定位非常清晰，就是达到更快的构建速度.

ninja的设计是对于make的缺陷的考虑，认为make有下面几点造成编译速度过慢：
- 隐式规则，make包含很多默认
- 变量计算，比如编译参与应该如何计算出来
- 依赖对象计算

ninja认为描述文件应该是这样的:
- 依赖必须显式写明(为了方便可以产生依赖描述文件)
- 没有任何变量计算
- 没有默认规则，没有任何默认值. 针对这点所以基本上可以认为ninja就是make的最最精简版

ninja相对于make增加了下面这些功能：
- 如果构建命令发生变化，那么这个构建也会重新执行
- 所依赖的目录在构建之前都已经创建了，如果不是这样的话，我们执行命令之前都要去生成目录
- 每条构建规则，除了执行命令之外，还允许有一个描述，真正执行打印这个描述而不是实际执行命令
- 每条规则的输出都是buffered的，也就是说并行编译，输入内容不会被搅和在一起

在Linux上CMake 2.8.8版本可以生成Ninja文件. 较新版本的CMake支持在Windows和Mac OS X上生成Ninja文件.

## ninja编译
```
apt install re2c
tar -xf ninja-xxx.tar.gz
cd ninja-xxx
./configure.py --bootstrap
./ninja --version
cp ninja /usr/bin/ninja
```
