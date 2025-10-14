# zap
ref: 
- [golang常用库包：log日志记录-uber的Go日志库zap使用详解](https://www.cnblogs.com/jiujuan/p/17304844.html)

## FAQ
### [不打印caller](https://github.com/go-logr/zapr/issues/4)
```go
func init() {
	cfg := zap.NewDevelopmentConfig()
	cfg.DisableCaller = false // Stop annotating logs with the calling function's file name.

	zapLog, err := cfg.Build()
	if err != nil {
		panic(fmt.Sprintf("error building logger: %v", err))
	}
	defer zapLog.Sync()
}
```

### 丢失若干日志
检查zap.Config的Sampling配置