# meson
meson 是一个开源的构建系统，速度快，易使用.

meson类似于 CMake ，Meson 并不直接构建软件，而是使用合适的后端，在 GNU/Linux 使用 ninja，在Windows 上使用 Visual Studio，在 MacOS 上使用 Xcode.

当前gnome, systemd, qemu使用了meson.

## 选择
### [bazel meson cmake autoreconf 优缺点对比](https://www.cnblogs.com/luoyinjie/p/17122127.html)
Bazel：Bazel 是一个由 Google 开发的构建工具，主要用于构建大型、复杂的软件项目。Bazel 的主要优点包括高速、可扩展性好、构建结果可重复、支持多语言等等。Bazel 适用于大型项目和工程师数量庞大的组织。

Meson：Meson 是一个轻量级的构建工具，主要用于构建 C++、C 和其他编程语言的项目。Meson 的主要优点包括易学易用、快速、高效、构建结果可重复等等。Meson 适用于小型到中型的项目，或者需要快速构建和测试原型的项目。

CMake：CMake 是一个跨平台的构建工具，主要用于构建 C++、C 和其他编程语言的项目。CMake 的主要优点包括可扩展性好、支持多个编译器和构建系统、易于学习等等。CMake 适用于小型到大型的项目，并且有广泛的社区支持。

Autoreconf：Autoreconf 是一个由 GNU 开发的自动化工具，用于为 GNU autoconf 生成 configure 脚本。Autoreconf 主要用于 UNIX 和类 UNIX 系统上的软件项目，它的主要优点是易于使用，可以自动生成 configure 脚本。Autoreconf 适用于需要使用 GNU autoconf 的项目。

因此，选择最好的工具需要根据具体需求来做出决定. 如果需要构建大型、复杂的项目，Bazel 是一个不错的选择; 如果需要快速构建和测试原型，Meson 可能更适合; 如果需要跨平台的支持和广泛的社区支持，CMake 是一个不错的选择; 如果需要使用 GNU autoconf，则需要使用 Autoreconf.

## FAQ
### meson build报错`meson_options.txt...\nUnknown type feature`
meson过旧, `pip3 install meson==0.54`, 在`/usr/local/bin/meson`

### `Could not detect Ninja v1.5 or newer`
定义在`mesonbuild/backend/ninjabackend.py`, 引用的是`mesonbuild/environment.py#detect_ninja_command_and_version()`, 它里面规定了默认ninja的版本, 比如meson==0.54是ninja1.7