### 传值和传指针

```go
package main

import (
	"fmt"
)

type SocialClient struct {
	Name string
}

type SocialConf []SocialClient

func (s SocialConf) New() {
	s = make([]SocialClient, 0)
	fmt.Printf("%#v\n", s)
}

func (s SocialConf) Add(sc ...SocialClient) { //等价于 func Add(s SocialConf,sc ...SocialClient){}
	fmt.Printf("%p\n", &s)
	s = append(s, sc...)
	fmt.Printf("%#v\n", s)
}

func main() {
	var s SocialConf
	s.New()
	fmt.Printf("%p\n", &s)
	fmt.Printf("%#v\n", s)
	var sc1 = SocialClient{Name: "weibo"}
	var sc2 = SocialClient{Name: "qq"}
	s.Add(sc1, sc2)
	fmt.Printf("%#v\n", s)
}
/*
main.SocialConf{}
0xc20801e020
main.SocialConf(nil)
0xc20801e0a0
main.SocialConf{main.SocialClient{Name:"weibo"}, main.SocialClient{Name:"qq"}}
main.SocialConf(nil)
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
	s = new(SocialConf) //关键,因为s只是一个值为nil的指针
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