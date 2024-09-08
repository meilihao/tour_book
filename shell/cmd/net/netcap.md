# [netcap](https://github.com/bytedance/netcap)
ref:
- [Installing BCC](https://github.com/iovisor/bcc/blob/master/INSTALL.md#install-and-compile-bcc-1)
- [字节跳动开源 Linux 内核网络抓包工具 netcap](https://my.oschina.net/u/6150560/blog/15154212)

netcap是字节跳动开源的Linux内核网络抓包工具.

> [官方: 建议kernel 5.15以上](https://github.com/bytedance/netcap/issues/5#issuecomment-2296100614)/[网友: kernel 5.16及以上 + bcc v0.25.0](https://github.com/bytedance/netcap/issues/3)

## FAQ
### netcap执行报`/virtual/main.c:155:21: warning: implicit declaration of function 'bpf_probe_read_kernel' is invalid in C99 [-Wimplicit-function-declaration]`
ref:
- [Program used external function 'bpf_probe_read_kernel' which could not be resolved! ](https://github.com/iovisor/bcc/issues/3395)
- [bpf-helpers(7) — Linux manual page](https://man7.org/linux/man-pages/man7/bpf-helpers.7.html)

```bash
$ vim netcap/pkg/dump/code/c/skb.h
..
#include <uapi/linux/in.h>
#define bpf_probe_read_kernel bpf_probe_read
#define _(P) ({typeof(P) val; bpf_probe_read_kernel(&val, sizeof(val), &P); val;})
```

在bpf_probe_read_kernel使用前插入`#define bpf_probe_read_kernel bpf_probe_read`

kernel 5.4.17里没有bpf_probe_read_kernel, 但有bpf_probe_read.

### netcap执行报`error loading BPF program: permission denied`
ref:
- [bpf: Failed to load program: Permission denied](https://github.com/bytedance/netcap/issues/2)
- [unbounded memory access](https://github.com/bpftrace/bpftrace/blob/master/docs/internals_development.md#unbounded-memory-access)
- [BPF的API](https://woodpenker.github.io/2021/12/05/BPF%E4%B9%8B%E5%B7%85%E7%9A%84%E5%AD%A6%E4%B9%A0--%E8%BF%BD%E8%B8%AA%E7%B3%BB%E7%BB%9F%E5%8E%86%E5%8F%B2%E4%B8%8E%E7%9B%B8%E5%85%B3%E6%8A%80%E6%9C%AF/)

```bash
# netcap skb -f tracepoint:skb:kfree_skb -e "icmp " -S
bpf: Failed to load program: Permission denied
0: (7b) *(u64 *)(r10 -56) = r1
1: (79) r6 = *(u64 *)(r1 +8)
2: (b7) r2 = 0
3: (b7) r1 = 0
4: (7b) *(u64 *)(r10 -32) = r1
last_idx 4 first_idx 0
regs=2 stack=0 before 3: (b7) r1 = 0
5: (63) *(u32 *)(r10 -4) = r2
last_idx 5 first_idx 0
regs=4 stack=0 before 4: (7b) *(u64 *)(r10 -32) = r1
regs=4 stack=0 before 3: (b7) r1 = 0
regs=4 stack=0 before 2: (b7) r2 = 0
6: (18) r1 = 0xffff942d9f375200
8: (bf) r2 = r10
9: (07) r2 += -4
10: (85) call bpf_map_lookup_elem#1
11: (7b) *(u64 *)(r10 -24) = r0
12: (15) if r0 == 0x0 goto pc+108
 R0_w=map_value(id=0,off=0,ks=4,vs=296,imm=0) R6_w=inv(id=0) R10=fp0 fp-8=mmmm???? fp-24_w=map_value fp-32_w=00000000 fp-56_w=ctx
13: (bf) r3 = r6
14: (07) r3 += 192
15: (bf) r1 = r10
16: (07) r1 += -16
17: (b7) r2 = 8
18: (85) call bpf_probe_read#4
last_idx 18 first_idx 0
regs=4 stack=0 before 17: (b7) r2 = 8
...
from 70 to 72: R0=inv(id=0) R1_w=inv(id=0,umax_value=255,var_off=(0x0; 0xff)) R2_w=inv256 R6=inv(id=0) R7_w=inv(id=0) R8=inv(id=0) R9=inv(id=0) R10=fp0 fp-8=mmmm???? fp-16=mmmmmmmm fp-24=map_value fp-32=mmmmmmmm fp-40=mmmmmmmm fp-48=mmmmmmmm fp-56=ctx
72: (bf) r9 = r7
73: (67) r9 <<= 32
74: (77) r9 >>= 32
75: (79) r1 = *(u64 *)(r10 -24)
76: (07) r1 += 40
77: (bf) r2 = r9
78: (bf) r3 = r8
79: (85) call bpf_probe_read#4
 R0=inv(id=0) R1_w=map_value(id=0,off=40,ks=4,vs=296,imm=0) R2_w=inv(id=0,umax_value=4294967295,var_off=(0x0; 0xffffffff)) R3_w=inv(id=0) R6=inv(id=0) R7_w=inv(id=0) R8=inv(id=0) R9_w=inv(id=0,umax_value=4294967295,var_off=(0x0; 0xffffffff)) R10=fp0 fp-8=mmmm???? fp-16=mmmmmmmm fp-24=map_value fp-32=mmmmmmmm fp-40=mmmmmmmm fp-48=mmmmmmmm fp-56=ctx
R2 unbounded memory access, use 'var &= const' or 'if (var < const)'
processed 130 insns (limit 1000000) max_states_per_insn 0 total_states 8 peak_states 8 mark_read 7
```

应该是当前kernel没有bpf_probe_read_kernel, 而用了`#define bpf_probe_read_kernel bpf_probe_read`(from [Program used external function 'bpf_probe_read_kernel' which could not be resolved! ](https://github.com/iovisor/bcc/issues/3395))导致的.

[bpf验证器有限制](https://elixir.bootlin.com/linux/v5.4.17/source/kernel/bpf/verifier.c#L3419), bpf_probe_read的第二个参数需要做长度限制, 但不知道如何修改.