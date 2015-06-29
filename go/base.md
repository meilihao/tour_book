### 传值和传指针

**go中所有传递均是值拷贝**

```go
package main

import (
	"fmt"
	"unsafe"
)

type SocialClient struct {
	Name string
}

type SocialConf []SocialClient

func (s SocialConf) New() {
	s = make([]SocialClient, 0)
}

func (s SocialConf) Add(sc ...SocialClient) { //等价于 func Add(s SocialConf,sc ...SocialClient){}
	fmt.Printf("2: %p\n", &s)
	s = append(s, sc...)
	p := (*struct {
		array uintptr
		len   int
		cap   int
	})(unsafe.Pointer(&s)) //获取slice底层结构的内容

	fmt.Printf("2: output: %+v\n", p)
	fmt.Printf("2: %#v\n", s)
}

func main() {
	var s SocialConf
	fmt.Printf("0: %p\n", &s)
	s.New() //与下面的s.Add()的传值情况相同
	var sc1 = SocialClient{Name: "weibo"}
	var sc2 = SocialClient{Name: "qq"}
	s.Add(sc1, sc2) //相当于将s的底层结构拷贝一份传入"func (s SocialConf) Add()",拷贝内容变化,但原值s的内容不变
	p := (*struct {
		array uintptr
		len   int
		cap   int
	})(unsafe.Pointer(&s))
	fmt.Printf("0: output: %+v\n", p)
	fmt.Printf("0: %#v\n", s)
}
/*
0: 0xc20801e020
2: 0xc20801e040
2: output: &{array:833357996128 len:2 cap:2}
2: main.SocialConf{main.SocialClient{Name:"weibo"}, main.SocialClient{Name:"qq"}}
0: output: &{array:0 len:0 cap:0}
0: main.SocialConf(nil)
*/
```

```go
package main

import (
	"fmt"
)

type SocialClient struct {
	Name string
}

type SocialConf []SocialClient

func (s *SocialConf) New() {
	fmt.Printf("%#v\n", s)
	*s = make([]SocialClient, 0)
}

func (s *SocialConf) Add(sc ...SocialClient) {
	*s = append(*s, sc...)
}

func main() {
	var s *SocialConf
	s = new(SocialConf) //关键,因为原先s只是一个值为nil的指针
	s.New()
	var sc1 = SocialClient{Name: "weibo"}
	var sc2 = SocialClient{Name: "qq"}
	s.Add(sc1, sc2)
	fmt.Printf("%#v\n", *s)
}
/*
&main.SocialConf(nil)
main.SocialConf{main.SocialClient{Name:"weibo"}, main.SocialClient{Name:"qq"}}
*/
```

### 初始化

```go
package main

import (
	"fmt"
)

type T struct {
	Name string
}

func (t *T) Print() string {
	t.Name = "test"
	return t.Name
}

var t *T

func main() {
	fmt.Printf("%#v\n", t)
	//t = new(T)<=>t = &T{} //ok
	*t = T{} //error:*t此时为nil,根本没有分配内存过(即没指向变量),谈何赋值
	fmt.Printf("%#v\n", t)
	fmt.Println(t.Print())
}
```

## 指针

unsafe模块的文档中提到几条转换规则，理解了以后就很容易做指针运算了：

A pointer value of any type can be converted to a Pointer.
A Pointer can be converted to a pointer value of any type.
A uintptr can be converted to a Pointer.
A Pointer can be converted to a uintptr.

unsafe.Pointer类似C的void *,即通用指针类型,可表示任意类型值的指针,用于转换不同类型指针，但它不可以参与指针运算;uintptr是golang的内置类型，是能存储指针的uint整型,是变量的首地址,可用于指针运算.