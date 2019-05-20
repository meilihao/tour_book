### JDK
1. 使用ppa/源方式安装
```
sudo add-apt-repository ppa:webupd8team/java
sudo apt-get update
```
1. 安装oracle-java-installer
参考:
  - [Ubuntu 安装 JDK 7 / JDK8 的两种方式](http://www.cnblogs.com/a2211009/p/4265225.html)

```
# jdk7
sudo apt-get install oracle-java7-installer
#　jdk8
sudo apt-get install oracle-java8-installer
```

  如果因为防火墙或者其他原因,导致相应的installer下载速度很慢,可以中断操作，然后先下载好相应jdk的tar.gz包,放在文件夹:

  /var/cache/oracle-jdk7-installer             (jdk7)

  /var/cache/oracle-jdk8-installer             (jdk8)

  下面,再安装一次installer即可．此时installer会默认使用该tar.gz包．

1. 设置系统默认jdk
```
# jdk7
sudo update-java-alternatives -s java-7-oracle
# jdk8
sudo update-java-alternatives -s java-8-oracle
```

1. 测试jdk 是是否安装成功
```
java -version
javac -version
```

#### 安装openjdk
sudo apt install  openjdk-8-jdk

#### 参考
1. 
