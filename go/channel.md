# channel
参考:
- [总结了才知道，原来channel有这么多用法！](https://segmentfault.com/a/1190000017958702)
- [今日头条Go建千亿级微服务的实践](http://blog.itpub.net/69946034/viewspace-2670129/)


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