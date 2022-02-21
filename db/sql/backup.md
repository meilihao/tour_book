# backup
ref:
- [rman备份与恢复](https://zhuanlan.zhihu.com/p/143866731)

从数据提取的角度来分类:
- 逻辑备份：指通过逻辑导出对数据进行备份，逻辑备份的数据只能基于备份时刻进行数据转储，所以恢复时也只能恢复到备份时保存的数据。对于备份点和故障点之间的数据，逻辑备份也是无能为力的，所以逻辑备份适合那些很少变化的数据表。如果通过逻辑备份进行全库备份，通常需要重建数据库，对于可用性很高的是数据库，这种恢复时间太长，通常不被使用。由于逻辑备份具有平台无关性，所以更为常见的是，逻辑备份被作为一种数据迁移及移动的主要手段。
- 物理备份：是指通过物理文件拷贝的方式对数据库进行备份。

从数据是否正在被使用角度来分类：
- 冷备份: 是指对**数据库进行关闭后的的拷贝备份**，这样的备份具有一致和完整的时间点数据，恢复时只需恢复所有文件就可以启动数据库。
- 热备份: 进行热备份的数据需要在运行归档模式，热备份时不需要关闭数据库。在进行恢复时，通过备份的数据文件及归档日志等文件，数据库就可以完全恢复，恢复到一致进行到最后一个归档模式。当然，如果是为了恢复某些用户错误，热备份的恢复完全可以在某一个时间点上停止恢复，也就是不完全恢复。

从数据库的备份角度分类：
- 完全备份：每次对数据库进行完整备份，当发生数据丢失的灾难时，完全备份无需依赖其他信息即可实现100%的数据恢复，其恢复时间最短且操作最方便。
- 增量备份：只有那些在上次完全备份或增量备份后被修改的文件才会被备份。优点是备份数据量小，需要的时间短，缺点是恢复的时候需要依赖以前备份记录，出问题的风险较大。
- 差异备份：备份那些自从上次完全备份之后被修改过的文件。从差异备份中恢复数据的时间较短，因此只需要两份数据---最后一次完整备份和最后一次差异备份，缺点是每次备份需要的时间较长