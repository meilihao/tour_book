# atomic
atomic包封装了CPU实现的部分原子操作指令,为用户层提供体验良好的原子操作函数,因此atomic包中提供的原语更接近硬件底层,也更为低级,它常被用于实现更为高级的并发同步技术(比如channel和sync包中的同步原语).

以atomic.SwapInt64函数在x86_64平台上的实现为例:
```go
// $GOROOT/src/sync/atomic/doc.go
func SwapInt64(addr *int64, new int64) (old int64)

// $GOROOT/src/sync/atomic/asm.s
TEXT ·SwapInt64(SB),NOSPLIT,$0
    JMP runtime⁄internal⁄atomic·Xchg64(SB)

// $GOROOT/src/runtime/internal/asm_amd64.s
TEXT runtime⁄internal⁄atomic·Xchg64(SB), NOSPLIT, $0-24
    MOVQ ptr+0(FP), BX
    MOVQ new+8(FP), AX
    XCHGQ AX, 0(BX)
    MOVQ AX, ret+16(FP)
    RET
```

从函数SwapInt64的实现中可以看到,它基本就是对x86_64CPU实现的原子操作指令XCHGQ的直接封装.

对共享整型变量的无锁读写性能:
1. 利用原子操作的无锁并发写的性能随着并发量增大几乎保持恒定
1. 利用原子操作的无锁并发读的性能随着并发量增大有持续提升的趋势,并且性能约为读锁的200倍

对共享自定义类型变量的无锁读写:
1. 利用原子操作的无锁并发写的性能随着并发量的增大而小幅下降
1. 利用原子操作的无锁并发读的性能随着并发量增大有持续提升的趋势,并且性能约为读锁的100倍

结论: atomic包更适合一些对性能十分敏感、并发量较大且读多写少的场合