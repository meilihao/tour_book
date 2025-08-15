# lint
## style
- [Go Styleguide](https://github.com/bahlo/go-styleguide/pull/8/commits/d706a55f5538a97b6ed65f490a5dac22efc6f71e)
- [uber-go/guide](https://github.com/xxjwxc/uber_go_guide_cn)
- [golangci/golangci-lint](https://github.com/golangci/golangci-lint)

## 规范代码
- gofmt
- goimports
- go vet

	go vet(Go 1.14)默认已经不再支持变量遮蔽检查了, 可安装`go install golang.org/x/tools/go/analysis/passes/shadow/cmd/shadow`配合`go vet`重新启用它.
- golangci-lint

## [golangci-lint](https://golangci-lint.run/)
ref:
- [golangci-lint的使用](https://blog.huweihuang.com/golang-notes/standard/golangci-lint/)
- [github.com/360EntSecGroup-Skylar/goreporter](https://github.com/qax-os/goreporter)

    已停止维护

```yml
version: "2"

formatters:
  enable:
    - gofumpt # gofmt 的严格超集
    - gci # gci 提供了比 goimports 更细粒度的导入分组控制
    - golines
    - swaggo

linters:
  enable: 
    - govet
    - staticcheck
    - errcheck
    - ineffassign
    - unused
    - misspell
    - gocyclo
    - gosec
  disable:
    - lll  # 行过长检查，不强制
```