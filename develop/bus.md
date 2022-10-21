# bus
## list
- [事件驱动的分布式事务架构设计](https://developer.aliyun.com/article/970658)
- [阿里云 EventBridge 事件驱动架构实践](https://www.shouxicto.com/article/2608.html)
- [全面异步化：淘宝反应式架构升级探索](https://www.infoq.cn/article/2uphtmd0poeunmhy5-ay)

	消息驱动强调无阻塞、**无 callback**，所以不会有线程挂在那里，不会有持续的资源消耗

	my ps: 目前发现涉及工作流的, 就需要保存和恢复现场来实现`无 callback`; 此时用callback更简单.