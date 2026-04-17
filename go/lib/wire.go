# wire
## FAQ
### wire_gen.go中包含`//go:generate go run -mod=mod github.com/google/wire/cmd/wire`导致重新生成wire_gen.go报错
在 go.work 下，Go 强制要求对模块的操作必须是显式的，不允许通过 -mod=mod 这种方式“暗渡陈仓”, 但wire为了确保在 Module 模式下能运行而在wire_gen.go硬编码进去了"//go:generate go run -mod=mod ..."

解决: `sed -i 's/-mod=mod //g' cmd/openhello/wire_gen.go`