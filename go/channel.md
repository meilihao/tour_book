# channel
参考:
- [总结了才知道，原来channel有这么多用法！](https://segmentfault.com/a/1190000017958702)
- [**今日头条Go建千亿级微服务的实践**](https://mp.weixin.qq.com/s?__biz=MjM5MDE0Mjc4MA==&mid=2650996069&idx=1&sn=63e7f5d5f91f9d84f1c3278426f6edf6)


## 并发控制
### 并发请求
```go
n:=2
wg:=sync.WaitGroup{}
errCh:=make(chan err, n)
fn:=func(){
	defer wg.Done()

	if err:=xxx();err!=nil{
		errCh<-err
		return
	}
}
for i:=0;i<n;i++{
	wg.Add(1)
	go fn()
}
wg.Wait()
close(errCh)

for err=range errCh {
	if err！=nil {
		return err
	}
}

return nil
```

### 并发请求+超时
```go
n:=2
wg:=sync.WaitGroup{}
done := make(chan struct{})
errCh:=make(chan err, n)
fn:=func(){
	defer wg.Done()

	if err:=xxx();err!=nil{
		errCh<-err
		return
	}
}
for i:=0;i<n;i++{
	wg.Add(1)
	go fn()
}

go func(){
    wg.Wait()
    close(errCh)
    close(done)
}

select {
	// 错误快返回,适用于get接口
	case err := <-errChan:
		return err
	case <-done:
	case <-time.After(500 * time.Millisecond):
        return fmt.Errorf()
}

return nil
```

### 控制并发数量
ref:
- [Go 语言高性能编程 - 利用 channel 的缓存区](https://geektutu.com/post/hpg-concurrency-control.html)

```go
func main() {
	var wg sync.WaitGroup
	ch := make(chan struct{}, 3)

	for i := 0; i < 10; i++ {
		ch <- struct{}{}
		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			log.Println(i)
			time.Sleep(time.Second)
			
			<-ch
		}(i)
	}

	wg.Wait()
}
```

### 并发控制+数据/错误处理
```go
package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	var err error
	ls := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13} // 可触发end flag
	//ls := []int{1, 2, 3, 4, 5, 6, 7, 8, 9} // 没足够数量触发end flag
	var flag int32
	errCh := make(chan error)
	dataCh := make(chan int)
	ch := make(chan struct{}, 3) // 并发控制
	doCh := make(chan int, 2)    // 一个处理err, 一个处理data. 可用WaitGroup代替
	wg := sync.WaitGroup{}

	// 先有receiver, 避免all groutines locks
	go func() {
		for e := range errCh {
			if err == nil {
				err = e
			}
			fmt.Println("get err:", e)
		}
		doCh <- 1
	}()

	go func() {
		for i := range dataCh {
			fmt.Println("get data:", i)
		}

		doCh <- 1
	}()

	for i := range ls {
		if atomic.LoadInt32(&flag) > 0 {
			fmt.Println("---end flag")
			break
		}

		//fmt.Println("doing:", i)
		ch <- struct{}{}
		wg.Add(1)

		go func(i int) {
			defer wg.Done()
			fmt.Println("in:", i)

			time.Sleep(time.Second)
			if i < 6 {
				dataCh <- i
			} else {
				atomic.AddInt32(&flag, 1)
				errCh <- fmt.Errorf("send err: %d", i)
			}
			fmt.Println("out:", i)

			<-ch
		}(i)
	}

	wg.Wait()
	close(errCh)
	close(ch)
	close(dataCh)
	<-doCh
	<-doCh
	close(doCh)
	fmt.Println(err)
}
```