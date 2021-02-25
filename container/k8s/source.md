# source
k8s 1.20.5

源码目录:
```bash
├── api // 存放api规范相关的文档
│   ├── api-rules // 已经存在的违反Api规范的api
│   ├── openapi-spec // OpenApi规范
│   └── OWNERS
├── build // 构建和测试脚本
│   ├── bindata.bzl
│   ├── BUILD
│   ├── build-image
│   ├── code_generation.bzl
│   ├── code_generation_test.bzl
│   ├── common.sh
│   ├── container.bzl
│   ├── copy-output.sh // 在容器中运行该脚本，后面可接多个命令：make, make cross 等
│   ├── dependencies.yaml
│   ├── go.bzl
│   ├── kazel_generated.bzl
│   ├── lib
│   ├── make-build-image.sh
│   ├── make-clean.sh  //清理容器中和本地的_output目录
│   ├── nsswitch.conf
│   ├── openapi.bzl
│   ├── OWNERS
│   ├── package-tarballs.sh
│   ├── pause
│   ├── platforms.bzl
│   ├── README.md
│   ├── release-images.sh
│   ├── release-in-a-container.sh
│   ├── release.sh
│   ├── release-tars
│   ├── root
│   ├── run.sh // 在容器中运行该脚本，后面可接多个命令：make, make cross 等
│   ├── shell.sh // 容器中启动一个shell终端
│   ├── tools.go
│   ├── util.sh
│   ├── visible_to
│   ├── workspace.bzl
│   └── workspace_mirror.bzl
├── BUILD.bazel -> build/root/BUILD.root
├── CHANGELOG
│   ├── CHANGELOG-1.20.md
│   ├── OWNERS
│   └── README.md
├── CHANGELOG.md -> CHANGELOG/README.md
├── cluster // 自动创建和配置kubernetes集群的脚本，包括networking, DNS, nodes等
│   ├── addons
│   ├── BUILD
│   ├── common.sh
│   ├── gce
│   ├── get-kube-binaries.sh
│   ├── get-kube.sh
│   ├── images
│   ├── kubectl.sh
│   ├── kube-down.sh
│   ├── kubemark
│   ├── kube-up.sh
│   ├── kube-util.sh
│   ├── log-dump
│   ├── OWNERS
│   ├── pre-existing
│   ├── README.md
│   ├── skeleton
│   └── validate-cluster.sh
├── cmd # // 内部包含各个组件的入口，具体核心的实现部分在pkg目录下
│   ├── BUILD
│   ├── clicheck
│   ├── cloud-controller-manager
│   ├── dependencycheck
│   ├── gendocs
│   ├── genkubedocs
│   ├── genman
│   ├── genswaggertypedocs
│   ├── genutils
│   ├── genyaml
│   ├── importverifier
│   ├── kubeadm
│   ├── kube-apiserver
│   ├── kube-controller-manager
│   ├── kubectl
│   ├── kubectl-convert
│   ├── kubelet
│   ├── kubemark
│   ├── kube-proxy
│   ├── kube-scheduler
│   ├── linkcheck
│   ├── OWNERS
│   ├── preferredimports
│   └── verifydependencies
├── code-of-conduct.md
├── CONTRIBUTING.md
├── docs
│   ├── BUILD
│   └── OWNERS
├── go.mod
├── go.sum
├── hack
│   ├── benchmark-go.sh
│   ├── boilerplate
│   ├── BUILD
│   ├── build-cross.sh
│   ├── build-go.sh
│   ├── cherry_pick_pull.sh
│   ├── conformance
│   ├── dev-build-and-push.sh
│   ├── dev-build-and-up.sh
│   ├── dev-push-conformance.sh
│   ├── e2e-internal
│   ├── e2e-node-test.sh
│   ├── generate-bindata.sh
│   ├── generate-docs.sh
│   ├── gen-swagger-doc
│   ├── get-build.sh
│   ├── ginkgo-e2e.sh
│   ├── grab-profiles.sh
│   ├── install-etcd.sh
│   ├── jenkins
│   ├── lib
│   ├── lint-dependencies.sh
│   ├── list-feature-tests.sh
│   ├── local-up-cluster.sh
│   ├── make-rules
│   ├── module-graph.sh
│   ├── OWNERS
│   ├── pin-dependency.sh
│   ├── print-workspace-status.sh
│   ├── README.md
│   ├── run-in-gopath.sh
│   ├── testdata
│   ├── test-go.sh
│   ├── test-integration.sh
│   ├── tools
│   ├── update-all.sh
│   ├── update-bazel.sh
│   ├── update-codegen.sh
│   ├── update-generated-api-compatibility-data.sh
│   ├── update-generated-device-plugin-dockerized.sh
│   ├── update-generated-device-plugin.sh
│   ├── update-generated-docs.sh
│   ├── update-generated-kms-dockerized.sh
│   ├── update-generated-kms.sh
│   ├── update-generated-kubelet-plugin-registration-dockerized.sh
│   ├── update-generated-kubelet-plugin-registration.sh
│   ├── update-generated-pod-resources-dockerized.sh
│   ├── update-generated-pod-resources.sh
│   ├── update-generated-protobuf-dockerized.sh
│   ├── update-generated-protobuf.sh
│   ├── update-generated-runtime-dockerized.sh
│   ├── update-generated-runtime.sh
│   ├── update-generated-swagger-docs.sh
│   ├── update-gofmt.sh
│   ├── update-hack-tools.sh
│   ├── update-import-aliases.sh
│   ├── update-openapi-spec.sh
│   ├── update-translations.sh
│   ├── update-vendor-licenses.sh
│   ├── update-vendor.sh
│   ├── update-workspace-mirror.sh
│   ├── verify-all.sh
│   ├── verify-api-groups.sh
│   ├── verify-bazel.sh
│   ├── verify-boilerplate.sh
│   ├── verify-cli-conventions.sh
│   ├── verify-codegen.sh
│   ├── verify-conformance-requirements.sh
│   ├── verify-description.sh
│   ├── verify-external-dependencies-version.sh
│   ├── verify-flags
│   ├── verify-flags-underscore.py
│   ├── verify-generated-device-plugin.sh
│   ├── verify-generated-docs.sh
│   ├── verify-generated-files-remake.sh
│   ├── verify-generated-files.sh
│   ├── verify-generated-kms.sh
│   ├── verify-generated-kubelet-plugin-registration.sh
│   ├── verify-generated-pod-resources.sh
│   ├── verify-generated-protobuf.sh
│   ├── verify-generated-runtime.sh
│   ├── verify-generated-swagger-docs.sh
│   ├── verify-gofmt.sh
│   ├── verify-golint.sh
│   ├── verify-govet-levee.sh
│   ├── verify-govet.sh
│   ├── verify-hack-tools.sh
│   ├── verify-import-aliases.sh
│   ├── verify-import-boss.sh
│   ├── verify-imports.sh
│   ├── verify-linkcheck.sh
│   ├── verify-no-vendor-cycles.sh
│   ├── verify-openapi-spec.sh
│   ├── verify-pkg-names.sh
│   ├── verify-prerelease-lifecycle-tags.sh
│   ├── verify-publishing-bot.py
│   ├── verify-readonly-packages.sh
│   ├── verify-shellcheck.sh
│   ├── verify-spelling.sh
│   ├── verify-staging-meta-files.sh
│   ├── verify-staticcheck.sh
│   ├── verify-test-code.sh
│   ├── verify-test-featuregates.sh
│   ├── verify-test-images.sh
│   ├── verify-typecheck-dockerless.sh
│   ├── verify-typecheck-providerless.sh
│   ├── verify-typecheck.sh
│   ├── verify-vendor-licenses.sh
│   └── verify-vendor.sh
├── LICENSE
├── LICENSES
│   ├── LICENSE
│   ├── OWNERS
│   └── vendor
├── logo // kubernetes的logo
│   ├── LICENSE
│   ├── logo.pdf
│   ├── logo.png
│   ├── logo.svg
│   ├── logo_with_border.pdf
│   ├── logo_with_border.png
│   ├── logo_with_border.svg
│   ├── name_blue.pdf
│   ├── name_blue.png
│   ├── name_blue.svg
│   ├── name_white.pdf
│   ├── name_white.png
│   ├── name_white.svg
│   ├── OWNERS
│   └── usage_guidelines.md
├── Makefile -> build/root/Makefile
├── Makefile.generated_files -> build/root/Makefile.generated_files
├── _output
│   ├── AGGREGATOR_violations.report
│   ├── APIEXTENSIONS_violations.report
│   ├── bin -> /opt/mark/go/src/k8s.io/kubernetes/_output/local/bin/linux/amd64
│   ├── CODEGEN_violations.report
│   ├── KUBE_violations.report
│   ├── local
│   ├── SAMPLEAPISERVER_violations.report
│   └── violattions.report
├── OWNERS
├── OWNERS_ALIASES
├── pkg // 主要代码存放类
│   ├── api
│   ├── apis
│   ├── auth
│   ├── BUILD
│   ├── capabilities
│   ├── client
│   ├── cloudprovider
│   ├── cluster
│   ├── controller
│   ├── controlplane
│   ├── credentialprovider
│   ├── features
│   ├── fieldpath
│   ├── generated
│   ├── kubeapiserver
│   ├── kubectl
│   ├── kubelet
│   ├── kubemark
│   ├── OWNERS
│   ├── printers
│   ├── probe
│   ├── proxy
│   ├── quota
│   ├── registry
│   ├── routes
│   ├── scheduler
│   ├── security
│   ├── securitycontext
│   ├── serviceaccount
│   ├── ssh
│   ├── util
│   ├── volume
│   └── windows
├── plugin
│   ├── BUILD
│   ├── OWNERS
│   └── pkg
├── README.md
├── SECURITY_CONTACTS
├── staging // 列出了独立发布的代码. kubernetes 的一些代码以独立项目的方式发布的, 这些代码虽然以独立项目发布，但是都在 kubernetes 主项目中维护，位于目录 kubernetes/staging/下，且这些代码会被定期同步到各个独立项目中
│   ├── BUILD
│   ├── OWNERS
│   ├── publishing
│   ├── README.md
│   ├── repos_generated.bzl
│   └── src
├── SUPPORT.md
├── test // 测试代码
│   ├── BUILD
│   ├── cmd
│   ├── conformance
│   ├── e2e
│   ├── e2e_kubeadm
│   ├── e2e_node
│   ├── fixtures
│   ├── fuzz
│   ├── images
│   ├── instrumentation
│   ├── integration
│   ├── kubemark
│   ├── list
│   ├── OWNERS
│   ├── soak
│   ├── typecheck
│   └── utils
├── third_party
│   ├── BUILD
│   ├── etcd.BUILD
│   ├── forked
│   ├── intemp
│   ├── multiarch
│   ├── OWNERS
│   └── protobuf
├── translations
│   ├── BUILD
│   ├── extract.py
│   ├── kubectl
│   ├── OWNERS
│   ├── README.md
│   └── test
├── vendor
│   ├── bitbucket.org
│   ├── BUILD
│   ├── cloud.google.com
│   ├── github.com
│   ├── go.etcd.io
│   ├── golang.org
│   ├── go.mongodb.org
│   ├── gonum.org
│   ├── google.golang.org
│   ├── go.opencensus.io
│   ├── gopkg.in
│   ├── go.uber.org
│   ├── k8s.io
│   ├── modules.txt
│   ├── OWNERS
│   └── sigs.k8s.io // 不同国家的语言包，使用poedit查看及编辑
└── WORKSPACE -> build/root/WORKSPACE
```

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
- make release ：构建所有的目标平台（Darwin、Linux、Windows），构建过程会比较久，并同时执行单元测试过程
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

# 核心数据结构
参考:
- [Using The Kubernetes API](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.20/)

在整个Kubernetes体系架构中，资源是Kubernetes最重要的概念，可以说Kubernetes的生态系统都围绕着资源运作. Kubernetes系统虽然有相当复杂和众多的功能，但它本质上是一个资源控制系统——注册、管理、调度资源并维护资源的状态.

在Kubernetes庞大而复杂的系统中，只有资源是远远不够的，Kubernetes将资源再次分组和版本化，形成Group（资源组）、Version（资源版本）、Resource（资源）.
- Group ：被称为资源组，在Kubernetes API Server中也可称其为APIGroup
- Version ：被称为资源版本，在Kubernetes API Server中也可称其为APIVersions
- Resource ：被称为资源，在Kubernetes API Server中也可称其为APIResource
- Kind ：资源种类，描述Resource的种类，与Resource为同一级别

![Group、Version、Resource核心数据结构](/misc/img/container/k8s/group-version-resource.jpg)

Kubernetes系统支持多个Group，每个Group支持多个Version，每个Version支持多个Resource，其中部分资源同时会拥有自己的子资源（即SubResource）. 例如，Deployment资源拥有Status子资源.

资源组、资源版本、资源、子资源的完整表现形式为`<group>/<version>/<resource>/<subresource>`. 以常用的Deployment资源为例，其完整表现形式为`apps/v1/deployments/status`
另外，资源对象（Resource Object）由“资源组+资源版本+资源种类”组成，并在实例化后表达一个资源对象，例如Deployment资源实例化后拥有资源组、资源版本及资源种类，其表现形式为`<group>/<version>，Kind=<kind>`，例如`apps/v1，Kind=Deployment`.

每一个资源都拥有一定数量的资源操作方法（即Verbs），资源操作方法用于Etcd集群存储中对资源对象的增、删、改、查操作. 目前Kubernetes系统支持8种资源操作方法，分别是create、delete、deletecollection、get、list、patch、update、watch操作方法.

每一个资源都至少有两个版本，分别是外部版本（External Version）和内部版本（Internal Version）. 外部版本用于对外暴露给用户请求的接口所使用的资源对象. 内部版本不对外暴露，仅在Kubernetes API Server内部使用.

Kubernetes资源也可分为两种，分别是Kubernetes Resource（Kubernetes内置资源）和Custom Resource（自定义资源）. 开发者通过CRD（即Custom Resource Definitions）可实现自定义资源，它允许用户将自己定义的资源添加到Kubernetes系统中，并像使用Kubernetes内置资源一样使用它们.

## ResourceList
Kubernetes Group、Version、Resource等核心数据结构存放在[vendor/k8s.io/apimachinery/pkg/apis/meta/v1](https://github.com/kubernetes/kubernetes/tree/master/staging/src/k8s.io/apimachinery/pkg/apis/meta/v1)目录中. 它包含了Kubernetes集群中所有组件使用的通用核心数据结构，例如APIGroup、APIVersions、APIResource等. 其中，可以通过APIResourceList数据结构描述所有Group、Version、Resource的结构.

Kubernetes的每个资源可使用metav1.APIResource结构进行描述，它描述资源的基本信息，例如资源名称（即Name字段）、资源所属的命名空间（即Namespaced字段）、资源种类（即Kind字段）、资源可操作的方法列表（即Verbs字段）.

每一个资源都属于一个或多个资源版本，资源所属的版本通过metav1.APIVersions结构描述，一个或多个资源版本通过Versions []string字符串数组进行存储.

在APIResourceList中，通过GroupVersion字段来描述资源组和资源版本，它是一个字符串，当资源同时存在资源组和资源版本时，它被设置为`<group>/<version>`；当资源不存在资源组（Core Group）时，它被设置为`/<version>`. 可以看到Pod、Service资源属于v1版本，而Deployment资源属于apps资源组下的v1版本.

另外，可以通过Group、Version、Resource结构来明确标识一个资源的资源组名称、资源版本及资源名称. Group、Version、Resource简称GVR，在Kubernetes源码中该数据结构被大量使用，它被定义在[vendor/k8s.io/apimachinery/pkg/runtime/schema](https://github.com/kubernetes/kubernetes/blob/master/staging/src/k8s.io/apimachinery/pkg/runtime/schema/group_version.go#L96)中. 在`vendor/k8s.io/apimachinery/pkg/runtime/schema`包中还定义了常用的资源数据结构:
- GroupResource   GR  资源组+资源
- GroupKind   GK  资源组+资源种类
- GroupVersion    GV  资源组+资源版本
- GroupVersionKind    GVK 资源自+资源版本+资源种类

    当资源对象的GVK输出为“/，Kind=”时，我们同样认为它是内部版本的资源对象.
- GroupVersions   GVS 资源组内多个资源版本

## Group（资源组）
Kubernetes系统中定义了许多资源组，这些资源组按照不同功能将资源进行了划分，资源组特点如下:
- 将众多资源按照功能划分成不同的资源组，并允许单独启用/禁用资源组. 当然也可以单独启用/禁用资源组中的资源
- 支持不同资源组中拥有不同的资源版本. 这方便组内的资源根据版本进行迭代升级
- 支持同名的资源种类（即Kind）存在于不同的资源组内
- 资源组与资源版本通过Kubernetes API Server对外暴露，允许开发者通过HTTP协议进行交互并通过动态客户端（即DynamicClient）进行资源发现
- 支持CRD自定义资源扩展
- 用户交互简单，例如在使用kubectl命令行工具时，可以不填写资源组名称

group用[APIGroup](https://github.com/kubernetes/kubernetes/blob/master/staging/src/k8s.io/apimachinery/pkg/apis/meta/v1/types.go#L943)表示:
- Name ：资源组名称
- Versions ：资源组下所支持的资源版本
- PreferredVersion ：首选版本. 当一个资源组内存在多个资源版本时，Kubernetes API Server在使用资源时会选择一个首选版本作为当前版本
在当前的Kubernetes系统中，支持两类资源组，分别是拥有组名的资源组和没有组名的资源组
- 拥有组名的资源组 ：其表现形式为<group>/<version>/<resource>，例如apps/v1/deployments

    拥有组名的资源组的HTTP PATH以/apis为前缀，其表现形式为/apis/<group>/<version>/<resource>，例如http://localhost：8080/apis/apps/v1/deployments.
- 没有组名的资源组 ：被称为Core Groups（即核心资源组）或Legacy Groups，也可被称为GroupLess（即无组）。其表现形式为/<version>/<resource>，例如/v1/pods

    没有组名的资源组的HTTP PATH以/api为前缀，其表现形式为/api/<version>/<resource>，例如http://localhost：8080/api/v1/pods

    > 没有组名的资源组，表示资源组名称为空。在后面会经常出现类似于/v1的表达，用来表示核心资源组下的v1资源版本

## version
Kubernetes的资源版本控制类似于语义版本控制（Semantic Versioning），在该基础上的资源版本定义允许版本号以v开头，例如v1beta1. 每当发布新的资源时，都需要对其设置版本号，这是为了在兼容旧版本的同时不断升级新版本，这有助于帮助用户了解应用程序处于什么阶段，以及实现当前程序的迭代. 语义版本控制应用得非常广泛，目前也是开源界常用的一种版本控制规范.

Kubernetes的资源版本控制可分为3种，分别是Alpha、Beta、Stable，它们之间的迭代顺序为Alpha→Beta→Stable，其通常用来表示软件测试过程中的3个阶段:
1. Alpha是第1个阶段，一般用于Kubernetes开发者内部测试

    该版本是不稳定的，可能存在很多缺陷和漏洞，官方随时可能会放弃支持该版本。在默认的情况下，处于Alpha版本的功能会被禁用. Alpha版本名称一般为v1alpha1、v1alpha2、v2alpha1等.
1. Beta是第2个阶段，该版本已经修复了大部分不完善之处，但仍有可能存在缺陷和漏洞，一般由特定的用户群来进行测试

    相对稳定的版本，Beta版本经过官方和社区很多次测试，当功能迭代时，该版本会有较小的改变，但不会被删除. 在默认的情况下，处于Beta版本的功能是开启状态的, Beta版本命名一般为v1beta1、v1beta2、v2beta1.
1. Stable是第3个阶段，此时基本形成了产品并达到了一定的成熟度，可稳定运行

    为正式发布的版本，Stable版本基本形成了产品，该版本不会被删除。在默认的情况下，处于Stable版本的功能全部处于开启状态. Stable版本命名一般为v1、v2、v3.

version用[APIVersions](https://github.com/kubernetes/kubernetes/blob/master/staging/src/k8s.io/apimachinery/pkg/apis/meta/v1/types.go#L915)表示:
- Versions ：所支持的资源版本列表

## Resource
一个资源被实例化后会表达为一个资源对象（即Resource Object）. 在Kubernetes系统中定义并运行着各式各样的资源对象. 所有资源对象都是Entity, Kubernetes使用这些Entity来表示当前状态. 可以通过Kubernetes API Server进行查询和更新每一个资源对象. Kubernetes目前支持两种Entity:
- 持久性实体（Persistent Entity） ：在资源对象被创建后，Kubernetes会持久确保该资源对象存在. 大部分资源对象属于持久性实体，例如Deployment资源对象
- 短暂性实体（Ephemeral Entity） ：也可称其为非持久性实体（Non-Persistent Entity）. 在资源对象被创建后，如果出现故障或调度失败，不会重新创建该资源对象，例如Pod资源对象.

Resource用[APIResource](https://github.com/kubernetes/kubernetes/blob/master/staging/src/k8s.io/apimachinery/pkg/apis/meta/v1/types.go#L984)表示:
- Name ：资源名称。
- SingularName ：资源的单数名称，它必须由小写字母组成，默认使用资源种类（Kind）的小写形式进行命名。例如，Pod资源的单数名称为pod，复数名称为pods
- Namespaced ：资源是否拥有所属命名空间
- Group ：资源所在的资源组名称
- Version ：资源所在的资源版本
- Kind ：资源种类
- Verbs ：资源可操作的方法列表，例如get、list、delete、create、update等
- ShortNames ：资源的简称，例如Pod资源的简称为po

### 资源外部版本与内部版本
Kubernetes资源代码定义在pkg/apis目录下，在Kubernetes系统中，同一资源对应着两个版本: 外部版本（External Version）与内部版本（Internal Version）. 例如，Deployment资源，它所属的外部版本表现形式为apps/v1，内部版本表现形式为apps/__internal.

- External Object ：外部版本资源对象，也称为Versioned Object（即拥有资源版本的资源对象）. 外部版本用于对外暴露给用户请求的接口所使用的资源对象，例如，用户在通过YAML或JSON格式的描述文件创建资源对象时，所使用的是外部版本的资源对象. 外部版本的资源对象通过资源版本（Alpha、Beta、Stable）进行标识
- Internal Object ：内部版本资源对象。内部版本不对外暴露，仅在Kubernetes API Server内部使用. 内部版本用于多资源版本的转换，例如将v1beta1版本转换为v1版本，其过程为v1beta1→internal→v1，即先将v1beta1转换为内部版本（internal），再由内部版本（internal）转换为v1版本. 内部版本资源对象通过runtime.APIVersionInternal（即__internal）进行标识.

资源版本（如v1beta1、v1等）与外部版本/内部版本概念不同: **拥有资源版本的资源属于外部版本，拥有runtime.APIVersionInternal标识的资源属于内部版本**

资源的外部版本代码定义在`pkg/apis/<group>/<version>/`目录下，资源的内部版本代码定义在`pkg/apis/<group>/`目录下. 例如，Deployment资源，它的外部版本定义在`pkg/apis/apps/{v1，v1beta1，v1beta2}/`目录下，它的内部版本定义在`pkg/apis/apps/`目录下（内部版本一般与资源组在同一级目录下）.

资源的外部版本和内部版本是需要相互转换的，而用于转换的函数需要事先初始化到资源注册表（Scheme）中. 多个外部版本（External Version）之间的资源进行相互转换，都需要通过内部版本（Internal Version）进行中转. 这也是Kubernetes能实现多资源版本转换的关键.

在Kubernetes源码中，外部版本的资源类型定义在vendor/k8s.io/api目录下，其完整描述路径为`vendor/k8s.io/api/<group>/<version>/<resource file>`.例如，Pod资源的外部版本，定义在vendor/k8s.io/api/core/v1/目录下.

资源的外部版本与内部版本的代码定义也不太一样，外部版本的资源需要对外暴露给用户请求的接口，所以资源代码定义了JSON Tags和Proto Tags，用于请求的序列化和反序列化操作。内部版本的资源不对外暴露，所以没有任何的JSON Tags和Proto Tags定义. 可通过`vendor/k8s.io/api/core/v1/types.go`和`pkg/apis/core/types.go`中的`type Pod struct`比较.

### 资源代码定义
Kubernetes资源代码定义在pkg/apis目录下，同一资源对应着内部版本和外部版本，内部版本和外部版本的资源代码结构并不相同.

资源的内部版本定义了所支持的资源类型（types.go）、资源验证方法（validation.go）、资源注册至资源注册表的方法（install/install.go）等. 而资源的外部版本定义了资源的转换方法（conversion.go）、资源的默认值（defaults.go）等.

以Deployment资源为例，它的内部版本定义在pkg/apis/apps/目录下，其资源代码结构如下：
- doc.go ：GoDoc文件，定义了当前包的注释信息。在Kubernetes资源包中，它还担当了代码生成器的全局Tags描述文件
- register.go ：定义了资源组、资源版本及资源的注册信息
- types.go ：定义了在当前资源组、资源版本下所支持的资源类型
- v1 、v1beta1 、v1beta2 ：定义了资源组下拥有的资源版本的资源（即外部版本）
- install ：把当前资源组下的所有资源注册到资源注册表中
- validation ：定义了资源的验证方法
- zz_generated.deepcopy.go ：定义了资源的深复制操作，该文件由代码生成器自动生成

每一个Kubernetes资源目录，都通过register.go代码文件定义所属的资源组和资源版本，内部版本资源对象通过runtime.APIVersionInternal（即__internal）标识.
每一个Kubernetes资源目录，都通过type.go代码文件定义当前资源组/资源版本下所支持的资源类型.

以Deployment资源为例，它的外部版本定义在pkg/apis/apps/{v1，v1beta1，v1beta2}目录下，其资源代码结构如下:
- 其中doc.go和register.go的功能与内部版本资源代码结构中的相似
- conversion.go ：定义了资源的转换函数（默认转换函数），并将默认转换函数注册到资源注册表中
- zz_generated.conversion.go ：定义了资源的转换函数（自动生成的转换函数），并将生成的转换函数注册到资源注册表中。该文件由代码生成器自动生成。
- defaults.go ：定义了资源的默认值函数，并将默认值函数注册到资源注册表中
- zz_generated.defaults.go ：定义了资源的默认值函数（自动生成的默认值函数），并将生成的默认值函数注册到资源注册表中。该文件由代码生成器自动生成。
外部版本与内部版本资源类型相同，都通过register.go代码文件定义所属的资源组和资源版本，外部版本资源对象通过资源版本（Alpha、Beta、Stable）标识

### 将资源注册到资源注册表中
在每一个Kubernetes资源组目录中，都拥有一个install/install.go代码文件，它负责将资源信息注册到资源注册表（Scheme）中, 以apps的`pkg/apis/apps/install/install.go`为例:
```go
func init() {
    Install(legacyscheme.Scheme)
}

// Install registers the API group and adds types to a scheme
func Install(scheme *runtime.Scheme) {
    utilruntime.Must(apps.AddToScheme(scheme))
    utilruntime.Must(v1beta1.AddToScheme(scheme))
    utilruntime.Must(v1beta2.AddToScheme(scheme))
    utilruntime.Must(v1.AddToScheme(scheme))
    utilruntime.Must(scheme.SetVersionPriority(v1.SchemeGroupVersion, v1beta2.SchemeGroupVersion, v1beta1.SchemeGroupVersion))
}
```

legacyscheme.Scheme是kube-apiserver组件的全局资源注册表，Kubernetes的所有资源信息都交给资源注册表统一管理. apps.AddToScheme函数注册apps资源组内部版本的资源. v1.AddToScheme函数注册apps资源组外部版本的资源. scheme.SetVersionPriority函数注册资源组的版本顺序，如有多个资源版本，排在最前面的为资源首选版本, 因此scheme.SetVersionPriority注册版本顺序很重要.

#### 资源首选版本
首选版本（Preferred Version），也称优选版本（Priority Version），一个资源组下拥有多个资源版本，例如，apps资源组拥有v1、v1beta1、v1beta2等资源版本. 当使用apps资源组下的Deployment资源时，在一些场景下，如不指定资源版本，则使用该资源的首选版本
以apps资源组为例，注册资源时会注册多个资源版本，分别是v1、v1beta2、v1beta1.

当通过资源注册表[scheme.PreferredVersionAllGroups](https://github.com/kubernetes/kubernetes/blob/master/staging/src/k8s.io/apimachinery/pkg/runtime/scheme.go#L638)函数获取所有资源组下的首选版本时，将位于最前面的资源版本作为首选版本.

> 在`(*Scheme)versionPriority`结构中并不存储资源对象的内部版本

除了scheme.PreferredVersionAllGroups函数外，还有另两个函数用于获取资源版本顺序相关的操作:
- scheme.PrioritizedVersionsForGroup ：获取指定资源组的资源版本，按照优先顺序返回
- scheme.PrioritizedVersionsAllGroups ：获取所有资源组的资源版本，按照优先顺序返回

### 资源操作方法
在Kubernetes系统中，针对每一个资源都有一定的操作方法（即Verbs），例如，对于Pod资源对象，可以通过kubectl命令行工具对其执行create、delete、get等操作. Kubernetes系统所支持的操作方法目前有8种操作，分别是create、delete、deletecollection、get、list、patch、update、watch. 这些操作方法可分为四大类，分别属于增、删、改、查，对资源进行创建、删除、更新和查询.

资源操作方法可以通过[metav1.Verbs](https://github.com/kubernetes/kubernetes/blob/master/staging/src/k8s.io/apimachinery/pkg/apis/meta/v1/types.go#L1023)进行描述.

不同资源拥有不同的操作方法，例如，针对Pod资源对象与pod/logs子资源对象，Pod资源对象拥有create、delete、deletecollection、get、list、patch、update、watch等操作方法，pod/logs子资源对象只拥有get操作方法，因为日志只需要执行查看操作。Pod资源对象与pod/logs子资源对象的操作方法分别通过metav1.Verbs进行描述.

资源对象的操作方法与存储（Storage）相关联，增、删、改、查实际上都是针对存储的操作。如何了解一个资源对象拥有哪些可操作的方法呢？需要查看与存储相关联的源码包registry，其定义在vendor/k8s.io/apiserver/pkg/registry/目录下. 每种操作方法对应一个操作方法接口（Interface）:

![资源对象操作方法接口说明](/misc/img/container/k8s/resource-verbs.jpg)

以Pod资源对象为例，Pod资源对象的存储（Storage）实现了以上接口的方法，Pod资源对象继承了genericregistry.Store，该对象可以管理存储（Storage）的增、删、改、查操作; 以pod/logs子资源对象为例，该资源对象只实现了get操作方法:
```go
// https://github.com/kubernetes/kubernetes/blob/master/pkg/registry/core/pod/storage/storage.go#L50

// PodStorage includes storage for pods and all sub resources
type PodStorage struct {
    Pod                 *REST
    ...
    Log                 *podrest.LogREST
    ...
}

// REST implements a RESTStorage for pods
type REST struct {
    *genericregistry.Store
    ...
}

// https://github.com/kubernetes/kubernetes/blob/master/staging/src/k8s.io/apiserver/pkg/registry/generic/registry/store.go#L93

func (e *Store) Create(ctx context.Context, obj runtime.Object, createValidation rest.ValidateObjectFunc, options *metav1.CreateOptions) (runtime.Object, error) {}

// https://github.com/kubernetes/kubernetes/blob/master/pkg/registry/core/pod/rest/log.go#L74
func (r *LogREST) Get(ctx context.Context, name string, opts runtime.Object) (runtime.Object, error) {}
```

### 资源与命名空间
Kubernetes系统支持命名空间（Namespace），其用来解决Kubernetes集群中资源对象过多导致管理复杂的问题. 每个命名空间相当于一个“虚拟集群”，不同命名空间之间可以进行隔离，当然也可以通过某种方式跨命名空间通信.

在一些使用场景中，命名空间常用于划分不同的环境，比如将Kubernetes系统划分为3个环境，分别是pro生产环境、test测试环境及dev开发环境，它们之间相互隔离，admin管理员用户对3个环境都拥有权限，而derek作为开发者只对dev开发环境拥有权限.

Kubernetes系统中默认内置了4个命名空间，分别介绍如下:
- default ：所有未指定命名空间的资源对象都会被分配给该命名空间
- kube-system ：所有由Kubernetes系统创建的资源对象都会被分配给该命名空间
- kube-public ：此命名空间下的资源对象可以被所有人访问（包括未认证用户）
- kube-node-lease ：此命名空间下存放来自节点的心跳记录（节点租约信息）

通过运行kubectl get namespace命令查看Kubernetes系统上所有的命名空间信息. 另外，在Kubernetes系统中，大部分资源对象都存在于某些命名空间中（例如Pod资源对象）. 但并不是所有的资源对象都存在于某个命名空间中（例如Node资源对象）. 决定资源对象属于哪个命名空间，可通过资源对象的ObjectMeta.Namespace描述.

### 自定义资源
Kubernetes系统拥有强大的高扩展功能，其中自定义资源（Custom Resource）就是一种常见的扩展方式，即可将自己定义的资源添加到Kubernetes系统中. Kubernetes系统附带了许多内置资源，但是仍有些需求需要使用自定义资源来扩展Kubernetes的功能.

自Kubernetes 1.7开始支持CustomResourceDefinitions（自定义资源定义，简称CRD）.

开发者通过CRD可以实现自定义资源，它允许用户将自己定义的资源添加到Kubernetes系统中，并像使用Kubernetes内置资源一样使用这些资源，例如，在YAML/JSON文件中带有Spec的资源定义都是对Kubernetes中的资源对象的定义，所有的自定义资源都可以与Kubernetes系统中的内置资源一样使用kubectl或client-go进行操作.

### 资源对象描述文件定义
Kubernetes资源可分为内置资源（Kubernetes Resources）和自定义资源（Custom Resources），它们都通过资源对象描述文件（Manifest File）进行定义.

一个资源对象需要用5个字段来描述它，分别是Group/Version、Kind、MetaData、Spec、Status, 这些字段定义在YAML或JSON文件中. Kubernetes系统中的所有的资源对象都可以采用YAML或JSON格式的描述文件来定义，下面是某个Pod文件的资源对象描述文件为例:
```yaml
apiVersion: v1
kind: Pod
metadata:
  name: static-web
  labels:
    role: myrole
spec:
  containers:
    - name: web
      image: nginx
      ports:
        - name: web
          containerPort: 80
          protocol: TCP
```

- apiVersion ：指定创建资源对象的资源组和资源版本，其表现形式为<group>/<version>，若是core资源组（即核心资源组）下的资源对象，其表现形式为<version>
- kind ：指定创建资源对象的种类
- metadata ：描述创建资源对象的元数据信息，例如名称、命名空间等
- spec ：包含有关Deployment资源对象的核心信息，告诉Kubernetes期望的资源状态、副本数量、环境变量、卷等信息。
- status ：包含有关正在运行的Deployment资源对象的信息

每一个Kubernetes资源对象都包含两个嵌套字段，即spec字段和status字段. 其中spec字段是必需的，它描述了资源对象的“期望状态”（Desired State），而status字段用于描述资源对象的“实际状态”（Actual State），它是由Kubernetes系统提供和更新的. 在任何时刻，Kubernetes控制器一直努力地管理着对象的实际状态以与期望状态相匹配.

## Kubernetes内置资源
Kubernetes系统内置了众多“资源组、资源版本、资源”，这才有了现在功能强大的资源管理系统, 可通过如下方式获得当前Kubernetes系统所支持的内置资源:
- kubectl api-versions ：列出当前Kubernetes系统支持的资源组和资源版本，其表现形式为<group>/<version>
- kubectl api-resources ：列出当前Kubernetes系统支持的Resource资源列表

```bash
kubectl api-resources -o wide [--namespaced=true] # 查看所有api resource, `--namespaced`表示该资源是否属于namespace
kubectl explain configmap [--api-version apps/v1] # 查看resource信息
kubectl api-versions # 查看所有api version
kubectl get deployments.v1.apps -n kube-system # 按照resource获取deployment
```

## runtime.Object类型基石
Runtime被称为“运行时”，在很多其他程序或语言中见过它，它一般指程序或语言核心库的实现. Kubernetes Runtime在`vendor/k8s.io/apimachinery/pkg/runtime`中实现，它提供了通用资源类型runtime.Object.

[runtime.Object](https://github.com/kubernetes/kubernetes/blob/master/staging/src/k8s.io/apimachinery/pkg/runtime/interfaces.go#L299)是Kubernetes类型系统的基石. Kubernetes上的所有资源对象（Resource Object）实际上就是一种Go语言的Struct类型，相当于一种数据结构，它们都有一个共同的结构叫runtime.Object. runtime.Object被设计为Interface接口类型，作为资源对象的通用资源对象.

以资源对象Pod为例，该资源对象可以转换成runtime.Object通用资源对象，也可以从runtime.Object通用资源对象转换成Pod资源对象.

runtime.Object提供了两个方法，分别是GetObjectKind和DeepCopyObject:
- GetObjectKind ：用于设置并返回GroupVersionKind
- DeepCopyObject ：用于深复制当前资源对象并返回

深复制相当于将数据结构克隆一份，因此它不与原始对象共享任何内容. 它使代码在不修改原始对象的情况下可以改变克隆对象的任何属性.

那么，如何确认一个资源对象是否可以转换成runtime.Object通用资源对象呢？这时需要确认该资源对象是否拥有GetObjectKind和DeepCopyObject方法. Kubernetes的每一个资源对象都嵌入了metav1.TypeMeta类型，metav1.TypeMeta类型实现了GetObjectKind方法，所以资源对象拥有该方法. 另外，Kubernetes的每一个资源对象都实现了DeepCopyObject方法，该方法一般被定义在zz_generated.deepcopy.go文件中. 因此，可以认为该资源对象能够转换成runtime.Object通用资源对象. 
所以，Kubernetes的任意资源对象都可以通过runtime.Object存储它的类型并允许深复制操作. 

比如实例化Pod资源，得到Pod资源对象，通过runtime.Object将Pod资源对象转换成通用资源对象（得到obj）, 然后通过断言的方式，将obj通用资源对象转换成Pod资源对象（得到pod2）. 最终可通过reflect（反射）来验证转换之前和转换之后的资源对象是否相等.

## Unstructured数据
数据可以分为结构化数据（Structured Data）和非结构化数据（Unstructured Data）. Kubernetes内部会经常处理这两种数据.

预先知道数据结构的数据类型是结构化数据.

无法预知数据结构的数据类型或属性名称不确定的数据类型是非结构化数据，其无法通过构建预定的struct数据结构来序列化或反序列化数据.

Kubernetes非结构化数据通过map[string]interface{}表达，并提供[Unstructured接口](https://github.com/kubernetes/kubernetes/blob/master/staging/src/k8s.io/apimachinery/pkg/runtime/interfaces.go). 在client-go编程式交互的DynamicClient内部，实现了Unstructured类型，用于处理非结构化数据.

## Scheme资源注册表
Kubernetes Scheme资源注册表类似于Windows操作系统上的注册表，只不过注册的是资源类型.

Kubernetes系统拥有众多资源，每一种资源就是一个资源类型，这些资源类型需要有统一的注册、存储、查询、管理等机制. 目前Kubernetes系统中的所有资源类型都已注册到Scheme资源注册表中，其是一个内存型的资源注册表，拥有如下特点:
- 支持注册多种资源类型，包括内部版本和外部版本
- 支持多种版本转换机制
- 支持不同资源的序列化/反序列化机制

Scheme资源注册表支持两种资源类型（Type）的注册，分别是UnversionedType和KnownType资源类型，分别介绍如下:
- UnversionedType ：无版本资源类型，这是一个早期Kubernetes系统中的概念，它主要应用于某些没有版本的资源类型，该类型的资源对象并不需要进行转换。在目前的Kubernetes发行版本中，**无版本类型已被弱化**，几乎所有的资源对象都拥有版本，但在metav1元数据中还有部分类型，它们既属于meta.k8s.io/v1又属于UnversionedType无版本资源类型，例如metav1.Status、metav1.APIVersions、metav1.APIGroupList、metav1.APIGroup、metav1.APIResourceList
- KnownType ：是目前Kubernetes最常用的资源类型，也可称其为“拥有版本的资源类型”

通过runtime.NewScheme可实例化一个新的Scheme资源注册表. 注册资源类型到Scheme资源注册表有两种方式: UnversionedType资源类型的对象通过[scheme.AddUnversionedTypes](https://github.com/kubernetes/kubernetes/blob/master/staging/src/k8s.io/apimachinery/pkg/runtime/scheme.go#L144)方法进行注册，KnownType资源类型的对象通过[scheme.AddKnownTypes](https://github.com/kubernetes/kubernetes/blob/master/staging/src/k8s.io/apimachinery/pkg/runtime/scheme.go#L162)方法进行注册.

在Scheme资源注册表中，不同的资源类型使用的注册方法不同，分别介绍如下:
- scheme.AddUnversionedTypes ：注册UnversionedType资源类型
- scheme.AddKnownTypes ：注册KnownType资源类型
    
    与scheme.AddKnownTypeWithName区别: 在注册资源类型时，无须指定Kind名称，而是通过reflect机制获取资源类型的名称作为资源种类名称
- scheme.AddKnownTypeWithName ：注册KnownType资源类型，须指定资源的Kind资源种类名称

### Scheme资源注册表数据结构
[Scheme](https://github.com/kubernetes/kubernetes/blob/master/staging/src/k8s.io/apimachinery/pkg/runtime/scheme.go#L46)资源注册表数据结构主要由4个map结构组成，它们分别是gvkToType、typeToGVK、unversionedTypes、unversionedKinds.

Scheme资源注册表结构字段说明如下:
- gvkToType ：存储GVK与Type的映射关系
- typeToGVK ：存储Type与GVK的映射关系，一个Type会对应一个或多个GVK
- unversionedTypes ：存储UnversionedType与GVK的映射关系
- unversionedKinds ：存储Kind（资源种类）名称与UnversionedType的映射关系

Scheme资源注册表通过Go语言的map结构实现映射关系，这些映射关系可以实现高效的正向和反向检索，从Scheme资源注册表中检索某个GVK的Type，它的时间复杂度为O(1).

Scheme资源注册表在Kubernetes系统体系中属于非常核心的数据结构.

参考[TestScheme](https://github.com/kubernetes/kubernetes/blob/master/staging/src/k8s.io/apimachinery/pkg/runtime/scheme_test.go#L67)和[TestUnversionedTypes](https://github.com/kubernetes/kubernetes/blob/master/staging/src/k8s.io/apimachinery/pkg/runtime/scheme_test.go#L399)代码.

![资源注册表映射关系](/misc/img/container/k8s/scheme.jpg)

GVK（资源组、资源版本、资源种类）在Scheme资源注册表中以<group>/<version>，Kind=<kind>的形式存在，其中对于Kind（资源种类）字段，在注册时如果不指定该字段的名称，那么默认使用类型的名称，例如corev1.Pod类型，通过reflect机制获取资源类型的名称，那么它的资源种类Kind=Pod.

资源类型在Scheme资源注册表中以Go Type（通过reflect机制获取）形式存在.

另外，需要注意的是，UnversionedType类型的对象在通过scheme.AddUnversionedTypes方法注册时，会同时存在于4个map结构中.

在运行过程中，kube-apiserver组件常对Scheme资源注册表进行查询，它提供了如下方法:
- scheme.KnownTypes ：查询注册表中指定GV下的资源类型
- scheme.AllKnownTypes ：查询注册表中所有GVK下的资源类型
- scheme.ObjectKinds ：查询资源对象所对应的GVK，一个资源对象可能存在多个GVK
- scheme.New ：查询GVK所对应的资源对象
- scheme.IsGroupRegistered ：判断指定的资源组是否已经注册
- scheme.IsVersionRegistered ：判断指定的GV是否已经注册
- scheme.Recognizes ：判断指定的GVK是否已经注册
- scheme.IsUnversioned ：判断指定的资源对象是否属于UnversionedType类型

## Codec编解码器
Codec编解码器与Serializer序列化器之间的差异:
- Serializer ：序列化器，包含序列化操作与反序列化操作。序列化操作是将数据（例如数组、对象或结构体等）转换为字符串的过程，反序列化操作是将字符串转换为数据的过程，因此可以轻松地维护数据结构并存储或传输数据
- Codec ：编解码器，包含编码器与解码器。编解码器是一个通用术语，指的是可以表示数据的任何格式，或者将数据转换为特定格式的过程。所以，可以将Serializer序列化器也理解为Codec编解码器的一种

k8s Codec编解码器通用接口定义在[`vendor/k8s.io/apimachinery/pkg/runtime/interfaces.go`](https://github.com/kubernetes/kubernetes/blob/master/staging/src/k8s.io/apimachinery/pkg/runtime/interfaces.go)

从Codec编解码器通用接口的定义可以看出，Serializer序列化器属于Codec编解码器的一种，这是因为每种序列化器都实现了Encoder与Decoder方法. 只要实现了Encoder与Decoder方法的数据结构，就是序列化器. Kubernetes目前支持3种主要的序列化器: json, yaml, protobuf.

在进行编解码操作时，每一种序列化器都对资源对象的metav1.TypeMeta（即APIVersion和Kind字段）进行验证，如果资源对象未提供这些字段，就会返回错误. 每种序列化器分别实现了Encode序列化方法与Decode反序列化方法，分别介绍如下:
- jsonSerializer ：JSON格式序列化/反序列化器
    
    通过json.NewSerializerWithOptions函数进行实例化, 使用application/json的ContentType作为标识, 文件扩展名为json
- yamlSerializer ：YAML格式序列化/反序列化器

    通过json.NewSerializerWithOptions函数进行实例化, 使用application/yaml的ContentType作为标识, 文件扩展名为yaml
- protobufSerializer ：Protobuf格式序列化/反序列化器

    protobufSerializer通过protobuf.NewSerializer函数进行实例化，使用application/vnd.kubernetes.protobuf的ContentType标识，文件扩展名为pb.

    [Encode](https://github.com/kubernetes/kubernetes/blob/master/staging/src/k8s.io/apimachinery/pkg/runtime/serializer/protobuf/protobuf.go#L163)函数首先验证资源对象是否为proto.Marshaler类型，proto.Marshaler是一个interface接口类型，该接口专门留给对象自定义实现的序列化操作。如果资源对象为proto.Marshaler类型，则通过t.Marshal序列化函数进行编码.

    而且，通过unk.MarshalTo函数在编码后的数据前加上protoEncodingPrefix前缀，前缀为magic-number特殊标识，其用于标识一个包的完整性。所有通过protobufSerializer序列化器编码的数据都会有前缀。前缀数据共4字节，分别是0x6b、0x38、0x73、0x00，其中第4个字节是为编码样式保留的

    [Decode](https://github.com/kubernetes/kubernetes/blob/master/staging/src/k8s.io/apimachinery/pkg/runtime/serializer/protobuf/protobuf.go#L97)函数首先验证protoEncodingPrefix前缀，前缀为magic-number特殊标识，其用于标识一个包的完整性，然后验证资源对象是否为proto.Message类型，最后通过proto.Unmarshal反序列化函数进行解码.

Codec编解码器将Etcd集群中的数据进行编解码操作.

Codec编解码器通过[NewCodecFactory→newSerializersForScheme](https://github.com/kubernetes/kubernetes/blob/master/staging/src/k8s.io/apimachinery/pkg/runtime/serializer/codec_factory.go#L51)函数实例化，在实例化的过程中会将jsonSerializer、yamlSerializer、protobufSerializer序列化器全部实例化.

Kubernetes在jsonSerializer序列化器上做了一些优化，caseSensitiveJsonIterator函数实际封装了github.com/json-iterator/go第三方库，json-iterator有如下几个好处:
- json-iterator支持区分大小写。Go语言标准库encoding/json在默认情况下不区分大小写
- json-iterator性能更优
- json-iterator 100%兼容Go语言标准库encoding/json，可随时切换两种编解码方式

## Converter资源版本转换器
在Kubernetes系统中，同一资源拥有多个资源版本，Kubernetes系统允许同一资源的不同资源版本进行转换，例如Deployment资源对象，当前运行的是v1beta1资源版本，但v1beta1资源版本的某些功能或字段不如v1资源版本完善，则可以将Deployment资源对象的v1beta1资源版本转换为v1版本. 可通过kubectl convert命令进行资源版本转换: `kubectl convert -f v1beta1Deployment.yaml --output-version=apps/v1`

首先，定义一个YAML Manifest File资源描述文件，该文件中定义Deployment资源版本为v1beta1. 通过执行kubect convert命令，--output-version将资源版本转换为指定的资源版本v1. 如果指定的资源版本不在Scheme资源注册表中，则会报错. 如果不指定资源版本，则默认转换为资源的首选版本.

Converter资源版本转换器主要**用于解决多资源版本转换问题**，Kubernetes系统中的一个资源支持多个资源版本，如果要在每个资源版本之间转换，最直接的方式是，每个资源版本都支持其他资源版本的转换，但这样处理起来非常麻烦. 例如，某个资源对象支持3个资源版本，那么就需要提前定义一个资源版本转换到其他两个资源版本（v1→v1alpha1，v1→v1beta1）、（v1alpha1→v1，v1alpha1→v1beta1）及（v1beta1→v1，v1beta1→v1alpha1）.

随着资源版本的增加，资源版本转换的定义会越来越多. 为了解决这个问题，Kubernetes通过内部版本（Internal Version）机制实现资源版本转换. 当需要在两个资源版本之间转换时，例如v1alpha1→v1beta1或v1alpha1→v1。Converter资源版本转换器先将第一个资源版本转换为__internal内部版本，再转换为相应的资源版本。每个资源只要能支持内部版本，就能与其他任何资源版本进行间接的资源版本转换.

### Converter转换器数据结构
Converter转换器数据结构主要存放转换函数（即Conversion Funcs）. [Converter转换器](https://github.com/kubernetes/kubernetes/blob/master/staging/src/k8s.io/apimachinery/pkg/conversion/converter.go#L40)数据结构如下:
- conversionFuncs ：默认转换函数。这些转换函数一般定义在资源目录下的conversion.go代码文件中
- generatedConversionFuncs ：自动生成的转换函数。这些转换函数一般定义在资源目录下的zz_generated.conversion.go代码文件中，是由代码生成器自动生成的转换函数。
- ignoredConversions ：若资源对象注册到此字段，则忽略此资源对象的转换操作。
- nameFunc ：在转换过程中其用于获取资源种类的名称，该函数被定义在vendor/k8s.io/apimachinery/pkg/runtime/scheme.go代码文件中。

Converter转换器数据结构中存放的转换函数（即Conversion Funcs）可以分为两类，分别为默认的转换函数（即conversionFuncs字段）和自动生成的转换函数（即generatedConversionFuncs字段）. 它们都通过[ConversionFuncs](https://github.com/kubernetes/kubernetes/blob/master/staging/src/k8s.io/apimachinery/pkg/conversion/converter.go#L40)来管理转换函数.

ConversionFunc类型函数（即Type Function）定义了转换函数实现的结构，将资源对象a转换为资源对象b. a参数定义了转换源（即source）的资源类型，b参数定义了转换目标（即dest）的资源类型。scope定义了多次转换机制（即递归调用转换函数）.

### Converter注册转换函数
Converter转换函数需要通过注册才能在Kubernetes内部使用，目前Kubernetes支持5个注册转换函数，分别介绍如下:
- scheme.AddIgnoredConversionType ：注册忽略的资源类型，不会执行转换操作，忽略资源对象的转换操作
- [scheme.AddConversionFuncs](https://github.com/kubernetes/kubernetes/blob/master/vendor/k8s.io/apimachinery/pkg/runtime/scheme.go) ：注册多个Conversion Func转换函数
- scheme.AddConversionFunc ：注册单个Conversion Func转换函数
- scheme.AddGeneratedConversionFunc ：注册自动生成的转换函数
- scheme.AddFieldLabelConversionFunc ：注册字段标签（FieldLabel）的转换函数

以[apps/v1](https://github.com/kubernetes/kubernetes/blob/master/pkg/apis/core/v1/conversion.go)资源组、资源版本为例，通过scheme.AddConversionFunc函数注册所有资源的转换函数.

### Converter资源版本转换原理
原理: 通过"_internal"桥梁的转换

例如Deployment资源对象，起初使用v1beta1资源版本，而v1资源版本更稳定，则会将v1beta1资源版本转换为v1资源版本. 在将Deployment v1beta1资源版本转换为内部版本（即__internal版本），得到转换后资源对象的GVK为“/，Kind=”。在这里，会产生疑问，为什么v1beta1资源版本转换为内部版本以后得到的GVK为“/，Kind=”而不是“apps/__internal，Kind=Deployment”。这就需要看Kubernetes源码实现了.

Scheme资源注册表可以通过两种方式进行版本转换:
- scheme.ConvertToVersion ：将传入的（in）资源对象转换成目标（target）资源版本，在版本转换之前，会将资源对象深复制一份后再执行转换操作，相当于安全的内存对象转换操作
- scheme.UnsafeConvertToVersion ：与scheme.ConvertToVersion功能相同，但在转换过程中不会深复制资源对象，而是直接对原资源对象进行转换操作，尽可能高效地实现转换。但该操作是非安全的内存对象转换操作

scheme.ConvertToVersion与scheme.UnsafeConvertToVersion资源版本转换功能都依赖于s.convertToVersion函数来实现.

![](/misc/img/container/k8s/convertion.jpg)

Converter转换器转换流程:
1. 获取传入的资源对象的反射类型
    
    资源版本转换的类型可以是runtime.Object或runtime.Unstructured，它们都属于Go语言里的Struct数据结构，通过Go语言标准库reflect机制获取该资源类型的反射类型，因为在Scheme资源注册表中是以反射类型注册资源的。获取传入的资源对象的反射类型
1. 从资源注册表中查找到传入的资源对象的GVK

    从Scheme资源注册表中查找到传入的资源对象的所有GVK，验证传入的资源对象是否已经注册，如果未曾注册，则返回错误

1. 从多个GVK中选出与目标资源对象相匹配的GVK

    target.KindForGroupVersionKinds函数从多个可转换的GVK中选出与目标资源对象相匹配的GVK. 这里有一个优化点，转换过程是相对耗时的，大量的相同资源之间进行版本转换的耗时会比较长。在Kubernetes源码中判断，如果目标资源对象的GVK在可转换的GVK列表中，则直接将传入的资源对象的GVK设置为目标资源对象的GVK，而无须执行转换操作，缩短部分耗时
1. 判断传入的资源对象是否属于Unversioned类型

    对于UnversionedType, 属于该类型的资源对象并不需要进行转换操作，而是直接将传入的资源对象的GVK设置为目标资源对象的GVK即可
1. 执行转换操作

    在执行转换操作之前，先判断是否需要对传入的资源对象执行深复制操作，然后通过s.converter.Convert转换函数执行转换操作.

    实际的转换函数是通过doConversion函数执行的，执行过程如下:
    - 从默认转换函数列表（即c.conversionFuncs）中查找出pair对应的转换函数，如果存在则执行该转换函数（即fn）并返回
    - 从自动生成的转换函数列表（即generatedConversionFuncs）中查找出pair对应的转换函数，如果存在则执行该转换函数（即fn）并返回
    - 如果默认转换函数列表和自动生成的转换函数列表中都不存在当前资源对象的转换函数，则使用[`(*Converter) Convert()`](https://github.com/kubernetes/kubernetes/blob/master/staging/src/k8s.io/apimachinery/pkg/conversion/converter.go#L281)函数传入的转换函数（即f）。调用f之前，需要将src与dest资源对象通过EnforcePtr函数取指针的值，因为`(*Converter) Convert()`函数传入的转换函数接收的是非指针资源对象

1. 设置转换后资源对象的GVK

    将v1beta1资源版本转换为内部版本（即__internal版本），得到转换后资源对象的GVK为“/，Kind=”。原因在于setTargetKind函数，转换操作执行完成以后，通过setTargetKind函数设置转换后资源对象的GVK，判断当前资源对象是否为内部版本（即APIVersionInternal），是内部版本则设置GVK为[`schema.GroupVersionKind{}`](https://github.com/kubernetes/kubernetes/blob/master/staging/src/k8s.io/apimachinery/pkg/runtime/scheme.go#L577)