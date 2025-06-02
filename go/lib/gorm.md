# gorm
ref:
- [Go 语言数据库开发利器：Gorm 框架指南](https://zhuanlan.zhihu.com/p/1908102185616646147)

Gorm 是一个用 Go 语言编写的全功能 ORM 库.

与其他 Go 语言的 ORM 框架相比，Gorm 具有以下优势：

1. 与 SQLx 对比

    SQLx 更接近原生 SQL，需要开发者编写更多 SQL 语句
    Gorm 提供更高级的抽象，减少 SQL 编写，更适合快速开发
    Gorm 的自动迁移和关联关系管理是 SQLx 所不具备的

2. 与 XORM 对比

    XORM 功能丰富但 API 相对复杂
    Gorm 的 API 设计更简洁直观，学习曲线更平缓
    Gorm 的社区活跃度和文档质量高且支持中、英、法等11种语言

3. 与 Ent 对比

    Ent 是 Facebook 开发的图数据库 ORM，专注于图数据模型
    Gorm 更通用，适用于传统关系型数据库
    Gorm 的 API 更符合 Go 语言习惯，而 Ent 的 API 风格更接近 GraphQL


## log
ref:
- [Monitoring Gin and GORM with OpenTelemetry](https://dev.to/vmihailenco/monitoring-gin-and-gorm-with-opentelemetry-53o0)