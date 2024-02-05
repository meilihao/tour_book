# pprof
## memory
`go tool pprof -alloc_space/-inuse_space http://ip:8899/debug/pprof/heap`
优先使用-inuse_space来分析，因为直接分析导致问题的现场比分析历史数据肯定要直观的多，一个函数alloc_space多不一定就代表它会导致进程的RSS高.