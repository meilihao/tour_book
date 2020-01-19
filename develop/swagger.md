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

### server env
```yaml
servers:
  - url: http://localhost:9090/v1/api
    description: test
  - url: https://xxx.com/v1/api
    description: prod
```