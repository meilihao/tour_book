# map
env: go version go1.12.5 linux/amd64

参考:
- [深度解密Go语言之map](https://www.tuicool.com/articles/ruInEbq)
- [Debugging Go Code with GDB](https://golang.org/doc/gdb)
- [Debugging what you deploy in Go 1.12](https://blog.golang.org/debugging-what-you-deploy)
- [Go里面的map,底层实现](http://www.sreguide.com/2018/05/07/go/go_map_01/)
- [Go Hashmap内存布局和实现](https://studygolang.com/articles/11979)
- [Map 在 Go runtime 中的高效实现（不使用范型）](https://studygolang.com/articles/13226)

Go 语言采用的是哈希查找表，并且使用链表解决哈希冲突.

## 模型
go使用hmap来表示 map.

[hmap](src/runtime/map.go):
```go
// A header for a Go map.
type hmap struct {
	// Note: the format of the hmap is also encoded in cmd/compile/internal/gc/reflect.go.
	// Make sure this stays in sync with the compiler's definition.
	count     int // # live cells == size of map.  Must be first (used by len() builtin) // map的元素个数即len(map)
	flags     uint8
	B         uint8  // log_2 of # of buckets (can hold up to loadFactor * 2^B items) // buckets 的对数, 即buckets 数组的长度是 2^B
	noverflow uint16 // approximate number of overflow buckets; see incrnoverflow for details
	hash0     uint32 // hash seed

	buckets    unsafe.Pointer // array of 2^B Buckets. may be nil if count==0. // 指向 buckets 数组，大小为 2^B. 如果元素个数为0，就为 nil
	oldbuckets unsafe.Pointer // previous bucket array of half the size, non-nil only when growing // 旧buckets 数组, 其长度是新buckets 数组的一半, 仅扩容时有值
	nevacuate  uintptr        // progress counter for evacuation (buckets less than this have been evacuated) // 扩容进度表示小于此地址的 buckets 已迁移完成

	extra *mapextra // optional fields
}
```

[bucket](src/runtime/map.go):
```go
// A bucket for a Go map.
type bmap struct {
	// tophash generally contains the top byte of the hash value
	// for each key in this bucket. If tophash[0] < minTopHash,
	// tophash[0] is a bucket evacuation state instead.
	tophash [bucketCnt]uint8
	// Followed by bucketCnt keys and then bucketCnt values.
	// NOTE: packing all the keys together and then all the values together makes the
	// code a bit more complicated than alternating key/value/key/value/... but it allows
	// us to eliminate padding which would be needed for, e.g., map[int64]int8.
	// Followed by an overflow pointer.
}
```

## 创建 map
通过汇编代码得知make的具体执行函数:
```go
m1:=make(map[string]int,5) // 调用 runtime.makemap_small 因为5<= bucketCnt(8)
m1["A"]=1

m2:=make(map[string]int,10) // 调用 runtime.makemap
m2["A"]=1

fmt.Println(m1,m2)
```

makemap:
```go
// makemap implements Go map creation for make(map[k]v, hint).
// If the compiler has determined that the map or the first bucket
// can be created on the stack, h and/or bucket may be non-nil.
// If h != nil, the map can be created directly in h.
// If h.buckets != nil, bucket pointed to can be used as the first bucket.
func makemap(t *maptype, hint int, h *hmap) *hmap {
	mem, overflow := math.MulUintptr(uintptr(hint), t.bucket.size)
	if overflow || mem > maxAlloc {
		hint = 0
	}

	// initialize Hmap
	if h == nil {
		h = new(hmap)
	}
	h.hash0 = fastrand() // 通过fastrand生成一个随机种子, 作为hash种子

	// Find the size parameter B which will hold the requested # of elements.
	// For hint < 0 overLoadFactor returns false since hint < bucketCnt.
	B := uint8(0)
	for overLoadFactor(hint, B) { // 找到一个 B，使得 map 的装载因子在正常范围内
		B++
	}
	h.B = B

	// allocate initial hash table // 初始化 hash table
	// if B == 0, the buckets field is allocated lazily later (in mapassign) // 如果 B 等于 0，那么 buckets 就会在赋值的时候再分配
	// If hint is large zeroing this memory could take a while. // 如果长度比较大，分配内存会花费长一点
	if h.B != 0 {
		var nextOverflow *bmap
		h.buckets, nextOverflow = makeBucketArray(t, h.B, nil)
		if nextOverflow != nil {
			h.extra = new(mapextra)
			h.extra.nextOverflow = nextOverflow
		}
	}

	return h
}

// overLoadFactor reports whether count items placed in 1<<B buckets is over loadFactor.
func overLoadFactor(count int, B uint8) bool {
	return count > bucketCnt && uintptr(count) > loadFactorNum*(bucketShift(B)/loadFactorDen)
}

// bucketShift returns 1<<b, optimized for code generation.
func bucketShift(b uint8) uintptr {
	if sys.GoarchAmd64|sys.GoarchAmd64p32|sys.Goarch386 != 0 {
		b &= sys.PtrSize*8 - 1 // help x86 archs remove shift overflow checks // 64位操作系统上的指针为8, 所以在64位系统上该值是b &= 63
	}
	return uintptr(1) << b
}
```