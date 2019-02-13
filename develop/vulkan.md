# vulkan

## FAQ
### 设置环境变量
参考`${vulkansdk}/setup-env.sh`

### 如何允许sdk里的examples, samples
直接运行`${vulkansdk}/build_xxx.sh`即可

> 运行`build_examples.sh`时, 没有构建成功且没有报错, 按照脚本自己构建即可.

### VK_ERROR_INCOMPATIBLE_DRIVER
```
$ sudo apt install mesa-vulkan-drivers
```

> mesa是软件实现的OpenGL. 安装了NVidia或AMD的专有驱动程序(会自带vulkan驱动)，则不需要Mesa; 如果想使用开源驱动（nouveau，radeon，radeonhd，intel）就需要 Mesa.

### Could NOT find xcb (missing: XCB_INCLUDE_DIR XCB_LIBRARY)
```
$ sudo apt-get install libxcb1-dev
```