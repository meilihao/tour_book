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