# openobserve
注意:
- openobserve v0.91.1 + otlp_grpc会报错, 先用otlphttp

```bash
# docker compose -f docker-compose.yaml up -d
```

其他相关:
```bash
# echo 'cm9vdEBleGFtcGxlLmNvbTpvMm9pX09vRHg0STY0VUs0Wlp2NUtUbERKZnF1OU5WTmVtdHAy' | base64 -d # 解码openobserve auth
```

## FAQ
### otlphttp正常但otlp_grpc报"rpc error: code = Unauthenticated desc = No valid auth token[5]"
ref:
- [`opentelemetry gRPC Connection failed#12540`](https://github.com/openobserve/openobserve/issues/12540)
