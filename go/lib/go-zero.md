# go-zero

## generate
```bash
# goctl api go --api demo.api --dir demo # 从api定义生成项目
```

## FAQ
### go-zero为什么每个请求重新生成NewXXXLoginc()
`通过在每个请求中重新生成 NewXXXLogic() 对象，可以确保系统的线程安全性和请求独立性，从而提高系统的稳定性和可靠性` from chatgpt