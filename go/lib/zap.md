# zap
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