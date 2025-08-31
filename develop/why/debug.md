# debug
## 诡异排查
1. 总在输出固定(或模式)内容的位置segmentation

	排查上下文, 可能是:
	1. 使用了已free的memory

## 程序闪退
1. coredump/panic

	go的panic和c的coredump不同, 会输出堆栈信息, c的在coredump里
1. 当前cpu不支持程序使用的指令集
1. 程序panic
1. 程序和环境的arch不匹配
1. OOM
1. 程序依赖库缺失或版本不对
1. 被其他程序kill

	比如cron触发的病毒

> 自身崩溃才可能有coredump, 被kill没有

排查顺序:
1. 终端日志
2. 系统日志
3. 应用日志
4. 审计日志