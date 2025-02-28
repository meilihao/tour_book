# leak
## kernel
ref:
- [诊断Linux内核内存泄漏等问题](https://zhuanlan.zhihu.com/p/551624495)

1. 看`smem -twk`
1. 看`/proc/meminfo`

    - slab高: 看`slabtop`

### tools

#### systemtap
#### kmemleak
它并不是针对某一个slab，而是针对所有的内核内存