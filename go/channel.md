# channel
参考:
- [总结了才知道，原来channel有这么多用法！](https://segmentfault.com/a/1190000017958702)


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