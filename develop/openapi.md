# swagger

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