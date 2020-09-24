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