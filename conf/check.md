### 检查虚拟机类型

```
yum install -y gcc gcc-c++ gdb
wget http://people.redhat.com/~rjones/virt-what/files/virt-what-1.15.tar.gz
tar zxvf virt-what-1.15.tar.gz
cd virt-what-1.15/
./configure
make && make install
virt-what
```
