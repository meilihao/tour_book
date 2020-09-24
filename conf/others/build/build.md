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
```