# build

## jemalloc
see [INSTALL.md](https://github.com/jemalloc/jemalloc/blob/dev/INSTALL.md)

```bash
tar xjf jemalloc-*.tar.bz2
cd jemalloc-*
./configure
make && make install
echo '/usr/local/lib' > /etc/ld.so.conf.d/local.conf
ldconfig
cat <<EOF > jemalloc.c
#include <stdio.h>
#include <jemalloc/jemalloc.h> // 如果没有include jemalloc的头文件，编译的时候也不需要链接jemalloc库, 但启动的时候需通过LD_PRELOAD指定jemalloc库的路径(比如`/usr/local/lib/libjemalloc.so`)就可以了
  
void jemalloc_test(int i)
{
        malloc(i*100);
}
 
int main(int argc, char **argv)
{
        int i;
        for(i=0;i<1000;i++)
        {
                jemalloc_test(i);
        }
        malloc_stats_print(NULL,NULL,NULL);
        return 0;
}
EOF
gcc jemalloc.c -o jmtest -ljemalloc
```

### FAQ
#### [jemalloc原理](http://tinylab.org/the-builtin-heap-profiling-of-jemalloc/)
#### <jemalloc>: Invalid conf pair
原因: [编译时未使用`--enable-prof`](https://github.com/jemalloc/jemalloc/issues/175)

解决方法: `./configure --enable-prof`

[MALLOC_CONF](http://jemalloc.net/jemalloc.3.html):
- prof:true 启用 profiling
- prof_active:false 一开始不激活(前提prof:true), 之后代码中通过 mallctl 调用来激活/关闭

    ```c
    bool active = true;
    mallctl("prof.active", NULL, NULL, &active, sizeof(active))
    ```
- 采样统计

    - lg_prof_sample:N 分配过程中，每流转 1 « N 个字节，将采样统计数据转储到文件
    - lg_prof_interval:N 每N秒采样一次
    - prof_gdump:true 当总分配量创新高时，将采样统计数据转储到文件
    - 在程序内主动触发转储`mallctl("prof.dump", NULL, NULL, NULL, 0)`

> 除了 jemalloc 的 jeprof 外，还有如下一些堆占用剖析工具gperftools 与 jeprof 类似，不过是基于 tcmalloc 的.

#### MALLOC_CONF运行时没有输出
原因应是程序没有引入`include <jemalloc/jemalloc.h>`, 且又没有设置LD_PRELOAD环境变量导致.

解决方法: `MALLOC_CONF=prof_leak:true,lg_prof_sample:0,prof_final:true LD_PRELOAD=/usr/local/lib/libjemalloc.so.2 ./a.out`

#### 查看`jeprof.*.f.heap`的分析结果
参考:
- [jemalloc的heap profiling](https://www.yuanguohuo.com/2019/01/02/jemalloc-heap-profiling/)

`jeprof a.out jeprof.34447.0.f.heap`

基于base看结果:`eprof multi_demo --base=jeprof.39496.0.i0.heap jeprof.39496.1.i1.heap`

#### 生成heap文件的调用路径图(pdf格式)
`jeprof --show_bytes --pdf a.out jeprof.34447.0.f.heap > a.pdf`

> 依赖: `apt install ghostscript graphviz`