# source
k8s 1.20.5

# 构建
see [README](https://github.com/kubernetes/kubernetes), 并**按照其规划路径**.

> 因为k8s repo包含软连接, 因此需选择releases中的`.tar.gz`压缩包

> 单独构建某一个组件，如kubectl组件，则需要指定WHAT参数: `make WHAT=cmd/kubectl`, 具体可参考[Makefile](https://github.com/kubernetes/kubernetes/blob/master/build/root/Makefile#L79)

当前k8s支持三种如下三种构建方式.

### 1. 本地构建
k8s构建一切都始于Makefile, 二进制输出路径是`_output/bin/`:
- Makefile ：顶层Makefile文件，描述了整个项目所有代码文件的编译顺序、编译规则及编译后的二进制输出等

    ```makefile
    all: generated_files # generated_files用于代码生成（Code Generation）
    hack/make-rules/build.sh $(WHAT) # $(WHAT)参数表示要指定构建的Kubernetes组件名称，不指定该参数则默认构建Kubernetes的所有组件
    ```
- Makefile.generated_files ：描述了代码生成的逻辑

通过make help命令，可以展示出所有可用的构建选项，从构建到测试的选项都有.

### 2. 容器构建
Kubernetes提供了两种容器环境下的构建方式：make release和make quick-release，它们之间的区别如下:
- make release ：构建所有的目标平台（Darwin、Linux、Windows），构建过程会比较久，并同时执行单元测试过程。
- make quick-release ：快速构建，只构建当前平台，并略过单元测试过程

make quick-release与make release相比多了两个变量，即KUBE_RELEASE_RUN_TESTS和KUBE_FASTBUILD. KUBE_RELEASE_RUN_TESTS变量，将其设为n则跳过运行单元测试；KUBE_FASTBUILD变量，将其设为true则跳过跨平台交叉编译. 通过这两个变量可以实现快速构建，最终执行build/release.sh脚本，运行容器环境构建.

在容器环境构建过程中，有多个容器镜像参与其中，分别如下:
- build容器（kube-cross）： 即构建容器，在该容器中会对代码文件执行构建操作，完成后其会被删除
- data容器： 即存储容器，用于存放构建过程中所需的所有文件
- rsync容器： 即同步容器，用于在容器和主机之间传输数据，完成后其会被删除

[具体构建过程](https://github.com/kubernetes/kubernetes/blob/master/build/release.sh):
1. kube::build::verify_prereqs

    进行构建环境的配置及验证: 该过程会检查本机是否安装了Docker容器环境，而对于Darwin平台，该过程会检查本机是否安装了docker-machine环境
1. kube::build::build_image

    根据Dockerfile文件构建容器镜像. Dockerfile文件来源于build/build-image/Dockerfile.
1. kube::build::run_build_command make cross

    运行构建容器并在构建容器内部执行构建Kubernetes源码的操作

    其中，构建的平台由KUBE_SUPPORTED_SERVER_PLATFORMS变量(hack/lib/golang.sh)控制; 
    构建的组件由KUBE_SERVER_TARGETS变量(hack/lib/golang.sh)控制
1. kube::build::copy_output

    使用同步容器，将编译后的代码文件复制到主机上。
1. kube::release::package_tarballs

    进行打包，将二进制文件打包到_output目录中

最终，代码文件以tar.gz压缩包的形式输出至_output/release-tars文件夹

### 3. bazel
目前Kubernetes已经支持使用Bazel进行构建和测试了，但尚未将Bazel作为默认的构建工具.

make bazel常用操作, 二进制输出路径是`bazel-bin`:
- make bazel-build ：构建所有二进制文件
- make bazel-test ：运行所有单元测试
- make bazel-test-integration ：运行所有集成测试
- make bazel-release ：在容器中进行构建

单独构建: `bazel build //cmd/kubectl/...`, `//cmd/kubectl/...`是bazel中的标记, 用于指定需要构建的包名.

每当开发者对Kubernetes代码进行更新迭代、添加或删除Go语言文件代码，以及更改Go import时，都必须更新各个包下的BUILD和BUILD.bazel文件，更新操作可通过运行hack/update-bazel.sh脚本自动完成.

Bazel第一次构建时须生成Bazel Cache，时间较长，再次构建时无须再生成Bazel Cache，Bazel Cache有利于大大提高构建速度.

## 代码生成
参考:
- [使用 code generator 生成 kubernetes 的 crd 代码](https://www.g5niusx.com/2020/03/kubernetes-1.html)
- [code-generator使用](https://tangxusc.github.io/2019/05/code-generator%E4%BD%BF%E7%94%A8/)

更新code gen方法: `hack/update-codegen.sh`

### Tag
代码生成器通过Tags（标签）来识别一个包是否需要生成代码及确定生成代码的方式，Kubernetes提供的Tags可以分为如下两种:
- 全局Tags ：定义在每个包的doc.go文件中，对整个包中的类型自动生成代码

    比如`pkg/apis/node/doc.go`:
    ```go
    // +k8s:deepcopy-gen=package
    // +groupName=node.k8s.io
    ```

    全局Tags告诉deepcopy-gen代码生成器为该包中的每个类型自动生成DeepCopy函数。其中的//+groupName定义了资源组名称，资源组名称一般使用域名形式命名。
- 局部Tags ：定义在Go语言的类型声明上方，只对指定的类型自动生成代码

    比如`pkg/apis/core/types.go`：
    ```go
    // +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
    ```

    局部Tags定义在Pod资源类型的上方，它定义了该类型有一个代码生成器deepcopy-gen: deepcopy-gen代码生成器为这个资源类型自动生成DeepCopy函数.

    > 关于Tags的位置，局部Tags一般定义在类型声明的上方，但如果该类型有注释信息，则局部Tags的定义需要与类型声明的注释信息之间至少有一个空行. 这是因为Kubernetes的API文档生成器会根据类型声明的注释信息（comment-block）生成文档, 为了避免Tags信息出现在文档中，故将Tags定义在注释的上方并空一行.

Tags的定义规则通常为//+tag-name或//+tag-name=value，它们被定义在注释中.

### 代码生成器
顶层Makefile中定义了generated_files命令，该命令用于构建代码生成器. `make generated_files`生成的二进制工具会被输出到`_output/bin`. 它实际调用的是`hack/make-rules/build.sh`, 会根据传入的代码生成器的main入口文件路径，构建二进制文件.

- conversion-gen

    自动生成Convert函数的代码生成器，用于资源对象的版本转换函数.

    给定一个包的目录路径作为输入源，conversion-gen会遍历包中的所有类型，若类型为types.Struct且过滤掉了私有的Struct类型，则为该类型生成Convert函数，并为该类型同时生成RegisterConversions注册函数，这些函数可以为对象在内部和外部类型之间提供转换函数.

    conversion-gen的生成规则: `vendor/k8s.io/code-generator/cmd/conversion-gen/generators/conversion.go#(*genConversion) convertibleOnlyWithinPackage()`

    为整个包生成Convert相关函数时，其Tags形式如下：`// +k8s:conversion-gen=<peer-pkg>`, 其中的`<peer-pkg>`用于定义包的导入路径

    为整个包生成Convert相关函数且依赖其他包时，其Tags形式如下：`// +k8s:conversion-gen-external-types=<type-pkg>`, 其中的`<type-pkg>`用于定义其他包的路径.

    在排除某个属性后生成Convert相关函数时，其Tags形式如下：`// +k8s:conversion-gen=false`

    构建conversion-gen二进制文件，并执行conversion-gen代码生成器，为k8s.io/kubernetes/pkg/apis/abac/v1beta1包生成zz_generated.conversion.go代码文件:
    ```bash
    $ cd $GOPATH/src/k8s.io/kubernetes
    $ hack/make-rules/build.sh ./vendor/k8s.io/code-generator/cmd/conversion-gen
    $ ./hack/run-in-gopath.sh ./_output/bin/conversion-gen --v 1 --logtostderr -i k8s.io/kubernetes/pkg/apis/abac/v1beta1 --extra-peer-dirs k8s.io/kubernetes/pkg/apis/core,k8s.io/kubernetes/pkg/apis/core/v1,k8s.io/api/core/v1 -O zz_generated.conversion # run-in-gopath.sh不能省略
    ```

- deepcopy-gen

    自动生成DeepCopy函数的代码生成器，用于资源对象的深复制函数.

    给定一个包的目录路径作为输入源, deepcopy-gen会遍历包中的所有类型，若类型为types.Struct，则会为该类型生成深复制函数, 这些函数可以有效地执行每种类型的深复制操作, 避免性能开销.

    deepcopy-gen的生成规则: `vendor/k8s.io/gengo/examples/deepcopy-gen/generators/deepcopy.go#copyableType()`

    为整个包生成DeepCopy相关函数时，其Tags形式：`+k8s:deepcopy-gen=package`

    为单个类型生成DeepCopy相关函数时，其Tags形式：`+k8s:deepcopy-gen=true`

    为整个包生成DeepCopy相关函数时，可以忽略单个类型，其Tags形式：`+k8s:deepcopy-gen=false`

    有时在Kubernetes源码里会看到deepcopy-gen的Tags被定义成runtime.Object，这时deepcopy-gen会为该类型生成返回值为runtime.Obejct类型的DeepCopyObject函数, 原始代码(`pkg/apis/abac/types.go`)如下:
    ```go
    // +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

    // Policy contains a single ABAC policy rule
    type Policy struct {
        metav1.TypeMeta

        // Spec describes the policy rule
        Spec PolicySpec
    }
    ```

    生成代码(`pkg/apis/abac/zz_generated.deepcopy.go`)如下:
    ```go
    // DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
    func (in *Policy) DeepCopyObject() runtime.Object {
        if c := in.DeepCopy(); c != nil {
            return c
        }
        return nil
    }
    ```

    构建deepcopy-gen二进制文件，并执行deepcopy-gen代码生成器，为`k8s.io/kubernetes/pkg/apis/abac/v1beta1`包生成`zz_generated.deepcopy.go`代码文件:
    ```bash
    $ cd $GOPATH/src/k8s.io/kubernetes
    $ hack/make-rules/build.sh ./vendor/k8s.io/code-generator/cmd/deepcopy-gen
    $ ./_output/bin/deepcopy-gen --v 1 --logtostderr -i k8s.io/kubernetes/pkg/apis/abac --bounding-dirs k8s.io/kubernetes -o $GOPATH/src -O zz_generated.deepcopy # `-i`从GOPATH查找package
    ```

    deepcopy-gen参数说明如下:
    - --v ：指定日志级别
    - --logtostderr ：日志输出到“标准错误输出”
    - -i，--input-dirs ：输入源，即.todo文件中的目录列表，以逗号分隔
    - --bounding-dirs ：依赖的包并为其生成深复制的类型
    - -O，--output-file-base ：输出文件的名字

- defaulter-gen

    自动生成Defaulter函数的代码生成器，用于资源对象的默认值函数.

    给定一个包的目录路径作为输入源，defaulter-gen会遍历包中的所有类型，若类型属性拥有特定类型（如TypeMeta、ListMeta、ObjectMeta），则为该类型[生成Defaulter函数](https://github.com/kubernetes/kubernetes/blob/master/pkg/apis/rbac/v1/zz_generated.defaults.go)，并为其生成RegisterDefaults注册函数.

    defaulter-gen的生成规则: `k8s.io/gengo/examples/defaulter-gen/generators/defaulter.go`

    为拥有TypeMeta属性的类型生成Defaulter相关函数时，其Tags形式如下：`// +k8s:defaulter-gen=TypeMeta`
    为拥有ListMeta属性的类型生成Defaulter相关函数时，其Tags形式如下：`// +k8s:defaulter-gen=ListMeta`
    为拥有ObjectMeta属性的类型生成Defaulter相关函数时，其Tags形式如下：`// +k8s:defaulter-gen=ObjectMeta`

    defaulter-gen的Tags都属于全局Tags，没有局部Tags. 其值可以为TypeMeta、ListMeta、ObjectMeta，最常用的是TypeMeta.有时在Kubernetes源码里会看到[defaulter-gen-input](https://github.com/kubernetes/kubernetes/blob/master/pkg/apis/storage/v1/doc.go)，这说明当前包会依赖于指定的路径包，代码示例如下：`// +k8s:defaulter-gen-input=../../../../vendor/k8s.io/api/storage/v1`.

    构建defaulter-gen二进制文件，并执行defaulter-gen代码生成器，为k8s.io/kubernetes/pkg/apis/rbac/v1包生成zz_generated.defaults.go代码文件:
    ```bash
    $ cd $GOPATH/src/k8s.io/kubernetes
    $ hack/make-rules/build.sh ./vendor/k8s.io/code-generator/cmd/defaulter-gen
    $ ./hack/run-in-gopath.sh ./_output/bin/defaulter-gen --v 1 --logtostderr -i k8s.io/kubernetes/pkg/apis/rbac/v1 --extra-peer-dirs k8s.io/kubernetes/pkg/apis/rbac/v1 -o $GOPATH/src -O zz_generated.defaults # run-in-gopath.sh不能省略
    ```

- go-bindata

    是一个第三方工具。它能够将静态资源文件嵌入Go语言中，例如在Web开发中将静态的HTML、JavaScript等静态资源文件嵌入Go语言代码文件中并提供一些操作方法. 给定一个静态资源目录路径作为输入源，go-bindata可以为其生成go文件.

    ```bash
    $ cd $GOPATH/src/k8s.io/kubernetes
    $ hack/make-rules/build.sh ./vendor/github.com/jteeuwen/go-bindata/go-bindata
    $ ./hack/run-in-gopath.sh hack/generate-bindata.sh
    ```

    generate-bindata.sh脚本重点执行如下代码:
    ```bash
    # 为translations静态资源目录生成pkg/kubectl/generated/bindata.go.tmp文件. translations目录存放的是与i18n（国际化）语言包相关的文件，在不修改内部代码的情况下支持不同语言及地区
    go-bindata -nometadata -nocompress -o "${BINDATA_OUTPUT}.tmp" -pkg generated \
    -ignore .jpg -ignore .png -ignore .md -ignore 'BUILD(\.bazel)?' \
    "translations/..."
    ```

- openapi-gen

    自动生成OpenAPI定义文件（OpenAPI Definition File）的代码生成器.

    给定一个包的目录路径作为输入源，openapi-gen会遍历包中的所有类型，若类型为types.Struct并忽略其他类型，则为types.Struct类型生成OpenAPI定义文件，该文件用于kube-apiserver服务上的OpenAPI规范的生成.

    openapi-gen的生成规则: `vendor/k8s.io/kube-openapi/pkg/generators/openapi.go#(openAPITypeWriter) generate()`

    为特定类型或包生成OpenAPI定义文件时，其Tags形式如下：`// +k8s:openapi-gen=true`

    排除为特定类型或包生成OpenAPI定义文件时，其Tags形式如下：`// +k8s:openapi-gen=false`

    构建openapi-gen二进制文件，并执行openapi-gen代码生成器，为k8s.io/kubernetes/vendor/k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1包生成zz_generated.openapi.go代码文件，该代码文件存放在k8s.io/kubernetes/pkg/generated/openapi目录下:
    ```bash
    $ cd $GOPATH/src/k8s.io/kubernetes
    $ hack/make-rules/build.sh ./vendor/k8s.io/code-generator/cmd/openapi-gen
    $ ./hack/run-in-gopath.sh ./_output/bin/openapi-gen --v 1 --logtostderr -i k8s.io/kubernetes/vendor/k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1 -p k8s.io/kubernetes/pkg/generated/openapi -O zz_generated.openapi2 -h vendor/k8s.io/code-generator/hack/boilerplate.go.txt -r _output/violattions.report  # run-in-gopath.sh不能省略
    ```

代码生成器生成代码的流程基本相同，以deepcopy-gen代码生成器为例，生成过程可分为如下3步:
1. 构建deepcopy-gen二进制文件
1. 生成.todo文件(在`.make/_output/bin/*.todo`)

    .todo文件相当于临时文件，用来存放被Tags标记过的包. 通过shell的grep命令可以将所有代码包中被Tags标记过的包目录记录在.todo文件中，这样可方便记录哪些包需要使用代码生成功能.

    ```bash
    // from Makefile.generated_files
    ALL_K8S_TAG_FILES := $(shell                             \
    find $(ALL_GO_DIRS) -maxdepth 1 -type f -name \*.go  \
        | xargs grep --color=never -l '^// *+k8s:'       \
    )


    DEEPCOPY_DIRS := $(shell                                             \
    grep --color=never -l '+k8s:deepcopy-gen=' $(ALL_K8S_TAG_FILES)  \
        | xargs -n1 dirname                                          \
        | LC_ALL=C sort -u                                           \
    )
    ```

    Makefile.generated_files中定义了ALL_K8S_TAG_FILES变量，其用于获取Kubernetes代码中被“/+k8s：”标签标记过的包；也定义了DEEPCOPY_DIRS变量，其用于从ALL_K8S_TAG_FILES中筛选出被`+k8s：deepcopy-gen`标签标记过的包. 最终将筛选出的包目录路径输出到.todo文件中.

1. 生成DeepCopy（深复制）相关函数


    ```bash
    // from Makefile.generated_files
    gen_deepcopy: $(DEEPCOPY_GEN) $(META_DIR)/$(DEEPCOPY_GEN).todo
    if [[ -s $(META_DIR)/$(DEEPCOPY_GEN).todo ]]; then                 \
        pkgs=$$(cat $(META_DIR)/$(DEEPCOPY_GEN).todo | paste -sd, -);  \
        if [[ "$(DBG_CODEGEN)" == 1 ]]; then                           \
            echo "DBG: running $(DEEPCOPY_GEN) for $$pkgs";            \
        fi;                                                            \
        ./hack/run-in-gopath.sh $(DEEPCOPY_GEN)                        \
            --v $(KUBE_VERBOSE)                                        \
            --logtostderr                                              \
            -i "$$pkgs"                                                \
            --bounding-dirs $(PRJ_SRC_PATH),"k8s.io/api"               \
            -O $(DEEPCOPY_BASENAME)                                    \
            "$$@";                                                     \
    fi                                                                 \
    ```

### gengo代码生成核心实现

Kubernetes的代码生成器都是在k8s.io/gengo包的基础上实现的. 代码生成器都会通过一个输入包路径（--input-dirs）参数，根据gengo的词法分析、抽象语法树等操作，最终生成代码并输出（--output-file-base）.

gengo代码目录结构说明如下:
- args ：代码生成器的通用flags参数.
- examples ：包含deepcopy-gen、defaulter-gen、import-boss、set-gen等代码生成器的生成逻辑.
- generator ：代码生成器通用接口Generator.
- namer ：命名管理，支持创建不同类型的名称. 例如，根据类型生成名称，定义type foo string，能够生成func FooPrinter（f*foo）{Print（string（*f））}.
- parser ：代码解析器，用来构造抽象语法树.
- types ：类型系统，用于数据类型的定义及类型检查算法的实现.

gengo的代码生成逻辑与编译器原理非常类似，大致可分为如下几个过程:
1. Gather The Info ：收集Go语言源码文件信息及内容

    gengo收集Go包信息可分为两步:
    1. 为生成的目标代码文件设置构建标签

        在[Default函数](https://github.com/kubernetes/kubernetes/blob/master/vendor/k8s.io/gengo/args/args.go#L42)中定义了默认的GeneratedBuildTag字符串，在每次构建时，代码生成器会将GeneratedBuildTag作为构建标签打入生成的代码文件中. 每个代码生成器都会通过Packages函数执行该操作，以deepcopy-gen代码生成器为例: deepcopy-gen代码生成器中的[Packages函数](https://github.com/kubernetes/kubernetes/blob/master/vendor/k8s.io/gengo/examples/deepcopy-gen/generators/deepcopy.go#L128)将GeneratedBuildTag字段进行拼接，每一个通过deepcopy-gen代码生成器生成的代码文件（如zz_generated.deepcopy.go），第1行总是构建标签, 最后生成代码的构建标签如下: `// +build !ignore_autogenerated`, 表示该文件是由代码生成器自动生成的，不需要人工干预或人工编辑该文件.


    2. 收集Go包信息并读取源码内容.

        代码生成器通过[--input-dirs参数](https://github.com/kubernetes/kubernetes/blob/master/vendor/k8s.io/gengo/args/args.go#L131)指定传入的Go包路径，通过build.Import方法收集Go包的信息，build.Import支持多种模式，其中build.ImportComment用于解析import语句后的注释信息；build.FindOnly用于查找包所在的目录，不读取其中的源码内容。代码函数层级为`b.AddDir→[b.importPackage](https://github.com/kubernetes/kubernetes/blob/master/vendor/k8s.io/gengo/parser/parse.go#L124)→[b.addDir](https://github.com/kubernetes/kubernetes/blob/master/vendor/k8s.io/gengo/parser/parse.go#L295)`.

        通过build.Import方法获得Go包信息以后，就可以得到包下面的所有源码文件的路径了，将所有Go源码内容读入内存中，等待Lexer词法解析器的下一步处理.


1. Lexer/Parser ：通过Lexer词法分析器进行一系列词法分析 -> AST Generator ：生成抽象语法树 -> Type Checker ：对抽象语法树进行类型检查

    Kubernetes gengo是在Go语言标准库支持代码解析功能的基础上进行的功能封装.

    gengo的代码解析:

    首先，通过[token.NewFileSet](https://github.com/kubernetes/kubernetes/blob/master/vendor/k8s.io/gengo/parser/parse.go#L107)实例化得到token.FileSet对象，该对象用于记录文件中的偏移量、类型、原始字面量及词法分析的数据结构和方法等.

    得到Tokens后，在[addFile](https://github.com/kubernetes/kubernetes/blob/master/vendor/k8s.io/gengo/parser/parse.go#L180)函数中，使用parser.ParseFile解析器对Tokens数据进行处理，Parser解析器将传入两种标识，其中parser.DeclarationErrors表示报告声明错误，parser.ParseComments表示解析代码中的注释并将它们添加到抽象语法树中, 最终得到抽象语法树结构.

    得到抽象语法树结构后，就可以[对其进行类型检查](https://github.com/kubernetes/kubernetes/blob/master/vendor/k8s.io/gengo/parser/parse.go#L418)了，通过Go语言标准库go/types下的Check方法进行检查，会对检查过程进行一些优化，使程序执行得更快.

    [gengo的类型系统（Type System）](https://github.com/kubernetes/kubernetes/blob/master/vendor/k8s.io/gengo/types/types.go#L65)在Go语言本身的类型系统之上归类并添加了几种类型. gengo的类型系统在Go语言标准库go/types的基础上进行了封装.

    所有的类型都通过`vendor/k8s.io/gengo/parser/parse.go`的walkType方法进行识别. gengo类型系统中的Struct、Map、Pointer、Interface等，与Go语言提供的类型并无差别. gengo与Go语言不同的类型，例如Builtin、Alias、DeclarationOf、Unknown、Unsupported及Protobuf. 另外，Signature并非是一个类型，它依赖于Func函数类型，用来描述Func函数的接收参数信息和返回值信息等.

    1. Builtin（内置类型）
    Builtin将多种Base类型归类成一种类型，以下几种类型在gengo中统称为Builtin类型:
    - 内置字符串类型——string
    - 内置布尔类型——bool
    - 内置数字类型——int、float、complex64等

    2. Alias（别名类型）
    Alias类型是Go 1.9版本中支持的特性.

    ```go
    type T1 struct{}
    type T2 = T1
    ```

    在Go语言标准库的reflect（反射）包识别T2的原始类型时，会将它识别为Struct类型，而无法将它识别为Alias类型. 原因在于，Alias类型在运行时是不可见的.

    如何让Alias类型在运行时可被识别呢？答案是因为gengo依赖于go/types的Named类型，所以要让Alias类型在运行时可被识别，在声明时将TypeName对象绑定到Named类型即可.

    3. DeclarationOf（声明类型）
    
    DeclarationOf并不是严格意义上的类型，它表示声明过的函数、全局变量或常量，但并未被引用过.

    在pkg/apis/abac/v1beta1/register.go中，AddToScheme变量在声明后未被其他对象引用过，则可以认为它是DeclarationOf类型的.
    
    4. Unknown（未知类型）
    当对象匹配不到以上所有类型的时候，它就是Unknown类型的

    5.Unsupported（未支持类型）

    当对象属于Unknown类型时，则会设置该对象为Unsupported类型，并在其使用过程中报错

    6. Protobuf（Protobuf类型）

    由go-to-protobuf代码生成器单独处理的类型

 
1. Code Generation ：生成代码，将抽象语法树转换为机器代码

    编译器生成的代码一般是二进制代码，而Kubernetes的代码生成器生成的是Go语言代码. gengo的[Generator接口](https://github.com/kubernetes/kubernetes/blob/master/vendor/k8s.io/gengo/generator/generator.go#L90)字段说明如下:

    - Name ：代码生成器的名称，返回值为生成的目标代码文件名的前缀，例如deepcopy-gen代码生成器的目标代码文件名的前缀为zz_generated.deepcopy
    - Filter ：类型过滤器，过滤掉不符合当前代码生成器所需的类型
    - Namers ：命名管理器，支持创建不同类型的名称. 例如，根据类型生成名称
    - Init ：代码生成器生成代码之前的初始化操作
    - Finalize ：代码生成器生成代码之后的收尾操作
    - PackageVars ：生成全局变量代码块，例如var（…）
    - PackageConsts ：生成常量代码块，例如consts（…）
    - GenerateType ：生成代码. 块根据传入的类型生成代码
    - Imports ：获得需要生成的import代码块. 通过该方法生成Go语言的import代码块，例如import（…）
    - Filename ：生成的目标代码文件的全名，例如deepcopy-gen代码生成器的目标代码文件名为zz_generated.deepcopy.go
    - FileType ：生成代码文件的类型，一般为golang，也有protoidl、api-violation等代码文件类型

Kubernetes目前提供的每个代码生成器都可以实现以上方法. 如果代码生成器没有实现某些方法，则继承默认代码生成器（DefaultGen）的方法，DefaultGen定义于vendor/k8s.io/gengo/generator/default_generator.go中.

`./_output/bin/deepcopy-gen --v 1 --logtostderr -i k8s.io/kubernetes/pkg/apis/abac/v1beta1 --bounding-dirs k8s.io/kubernetes,"k8s.io/api" -O zz_generated.deepcopy`生成流程详解:
1. 实例化generator.Packages对象

    deepcopy-gen代码生成器根据输入的包的目录路径（即输入源），实例化[generator.Packages](https://github.com/kubernetes/kubernetes/blob/master/vendor/k8s.io/gengo/examples/deepcopy-gen/generators/deepcopy.go#L208)对象，根据generator.Packages结构生成代码.

    在deepcopy-gen代码生成器的Packages函数中，实例化generator.Packages对象并返回该对象。根据输入源信息，实例化当前Packages对象的结构：PackageName字段为v1beta1，PackagePath字段为k8s.io/kubernetes/pkg/apis/abac/v1beta1. 其中，最主要的是GeneratorFunc定义了Generator接口的实现（即NewGenDeepCopy实现了Generator接口方法）.

2.执行代码生成

    在gengo中，generator定义代码生成器通用接口Generator. 通过[ExecutePackage](https://github.com/kubernetes/kubernetes/blob/master/vendor/k8s.io/gengo/generator/execute.go#L215)函数，调用不同代码生成器（如deepcopy-gen）的Generator接口方法，并生成代码.

    ExecutePackage代码生成执行流程：生成Header代码块→生成Imports代码块→生成Vars全局变量代码块→生成Consts常量代码块→生成Body代码块. 最后，调用assembler.AssembleFile函数，将生成的代码块信息写入zz_generated.deepcopy.go文件，生成pkg/apis/abac/v1beta1/zz_generated.deepcopy.go代码结构.

    deepcopy-gen代码生成器最终生成了代码文件zz_generated.deepcopy.go，该文件的整体结构可分为如下部分:
    1. Header代码块信息，包括build tag和license boilerplate文件（存放开源软件作者及开源协议等信息），其中license boilerplate文件可以从hack/boilerplate/boilerplate.go.txt中获取
    1. Imports代码块信息，引入外部包
    1. Vars全局变量代码块信息，当前代码文件未使用Vars
    1. Consts常量代码块信息，当前代码文件未使用Consts
    1. Body代码块信息，生成DeepCopy深复制函数

    在生成代码的过程中，Filter函数和GenerateType函数非常重要. deepcopy-gen代码生成器根据[Filter类型过滤器](https://github.com/kubernetes/kubernetes/blob/master/vendor/k8s.io/gengo/examples/deepcopy-gen/generators/deepcopy.go#L397)筛选需要生成哪些结构.

    通过Filter→copyableType的实现，deepcopy-gen代码生成器只筛选出了类型为Struct结构的数据（即只为Struct结构的数据生成DeepCopy函数）. 

    [GenerateType函数](https://github.com/kubernetes/kubernetes/blob/master/vendor/k8s.io/gengo/examples/deepcopy-gen/generators/deepcopy.go#L582)会其根据传入的类型生成Body代码块信息. 内部通过Go语言标准库text/template模板语言渲染出生成的Body代码块信息.

    generator.NewSnippetWriter内部封装了text/template模板语言，通过将模板应用于数据结构来执行模板. SnippetWriter对象在实例化时传入模板指令的标识符（即指令开始为$，指令结束为$，有时候也会使用{{}}作为模板指令的标识符）.

    SnippetWriter通过Do函数加载模板字符串，并执行渲染模板. 模板指令中的点（`.`）表示引用args参数传递到模板指令中. 模板指令中的（`|`）表示管道符，即把左边的值传递给右边.