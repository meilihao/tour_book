# flutter

## FAQ
### 运行时一直卡在`Initializing gradle...`
原因:
- gradle未安装或gradle未安装在默认位置`/home/$USER/.gradle`
- 中国访问不了google, 改为使用镜像访问

解决:
- 安装gradle
```sh
$ cd AndroidStudioProjects/${flutter_app}/android
$ ./gradlew # 安装`AndroidStudioProjects/${flutter_app}/android/gradle/wrapper/gradle-wrapper.properties`中`distributionUrl`指定的gradle版本
```

- 修改`/home/chen/opt/flutter/packages/flutter_tools/gradle/flutter.gradle`:
```
buildscript {
    repositories {
        // google()
        // jcenter()
        maven { url 'https://maven.aliyun.com/repository/google' }
        maven { url 'https://maven.aliyun.com/repository/jcenter' }
    }
    dependencies {
        classpath 'com.android.tools.build:gradle:3.2.1'
    }
}
```

### adb: failed to install ...my_flutter_app/build/app/outputs/apk/app.apk: Failure [INSTALL_FAILED_USER_RESTRICTED: Install canceled by user]
Android开发者选项-打开`USB安装`

### Error connecting to the service protocol: HttpException: Connection closed before full header was received, uri = http://127.0.0.1:34597/ws

[flutter issue](https://github.com/flutter/flutter/issues/14991)

### Android Studio: /dev/kvm device permission denied
```
$ sudo apt install qemu-kvm.
$ ls -al /dev/kvm # 检查该文件的组是否是kvm
$ sudo adduser ${yourname} kvm # 将当前用户加入kvm组
$ grep kvm /etc/group # 检查是否加入成功, 之后注销当前用户或重启即可.
```

### Gradle task assembleDebug failed with exit code 1(待定)
修改`${my_flutter_app}/android/gradle/wrapper/gradle-wrapper.properties`里的`distributionUrl`, 使用gradle的`bin`版本而不是`all`.

### A problem occurred configuring root project 'android'. java.io.FileNotFoundException...
清理gradle缓存`rm -rf /home/chen/.gradle/*`, 再执行`${my_flutter_app}/android/gradlew`.

### 在线gradle下载太慢
flutter gradle通常使用在`${my_flutter_app}/android/gradle/wrapper/gradle-wrapper.properties`中配置的在线gradle, 但也可使用[提前下载](http://services.gradle.org/distributions)好的文件, 比如`distributionUrl=file\:/home/chen/gradle-5.2.1-all.zip`.

> gradle默认安装位置: `/home/${USER}/.gradle/wrapper/dists`
> 也可将`distributionUrl=file...`安装好的gradle拷到`distributionUrl=https...`安装目录下.

### flutter切换版本
- [Upgrading Flutter](https://flutter.io/docs/development/tools/sdk/upgrading)和[How to change channels](https://github.com/flutter/flutter/wiki/Flutter-build-release-channels)
- [flutter版本列表](https://flutter.io/docs/development/tools/sdk/archive?tab=linux)

### flutter命令卡住 & Android Studio 创建flutter(start a new Flutter project)卡住
[国内访问Flutter有时可能会受到限制，Flutter官方为中国开发者搭建了临时镜像](https://flutter.io/community/china)