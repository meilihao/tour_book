# swagger
ref:
- [swagger和openAPI: 数据类型](https://www.breakyizhan.com/swagger/2969.html)
- [Swagger与其他API文档编写工具对比](https://haofly.net/swagger/)

## api转换
to OpenAPI 3.0:
- [postman-to-swagger, **推荐**](https://www.npmjs.com/package/postman-to-swagger)
- [postman-to-swagger(online)](https://metamug.com/util/postman-to-swagger/)
- [postman-to-swagger : lib](https://github.com/tecfu/postman-to-swagger)

## repo
- [go-swagger3](https://github.com/parvez3019/go-swagger3)
- [openapi](https://github.com/sv-tools/openapi)
  
  OpenAPI v3.1 Spec implementation in Go with generics
- [libopenapi](https://github.com/pb33f/libopenapi)

  libopenapi is a fully featured, high performance OpenAPI 3.1, 3.0 and Swagger parser, library, validator and toolkit for golang applications.

## 快捷键
- 缩进 : table
- 反向缩进 : shit + table

## FAQ
### response file
```yaml
responses:
        200:
          description:  下载文件, 文件名是xxx
          content:
            application/octet-stream:
              schema:
                type: string
                format: binary
```

### server env和multi hosts
```yaml
servers:
  - url: http://localhost:9090/v1/api
    description: test
  - url: https://{customerId}.saas-app.com:{port}/v2
    variables:
      customerId:
        default: demo
        description: Customer ID assigned by the service provider
      port:
        enum:
          - '443'
          - '8443'
        default: '443'
paths:
  /files:
    description: File upload and download operations
    servers:
      - url: https://files.example.com # 覆盖默认的server
        description: Override base path for all operations with the /files path
```