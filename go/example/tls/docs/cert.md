# cert

参考:
- [Go和HTTPS](http://tonybai.com/2015/04/30/go-and-https/)
- [使用Go实现TLS 服务器和客户端](http://colobu.com/2016/06/07/simple-golang-tls-examples/)
- [golang-tls.md](https://gist.github.com/denji/12b3a568f092ab951456)

## 生成Root证书
```sh
$ openssl genrsa -out rootca.key 2048
$ openssl req -sha256 -new -x509 -days 3650 -subj "/CN=ZXB Global CA/O=ZXB Inc/C=CN" -key rootca.key -out rootca.crt
$ openssl x509 -in rootca.crt -noout -text # 查看证书信息
```

## 查看rootca.crt内容
```
$ openssl x509 -in rootca.crt -noout -text
```

## 生成Server证书
v3_server.ext:
```
# v3.ext(服务端证书必须有subjectAltName,否则chrome报错):
authorityKeyIdentifier=keyid,issuer
basicConstraints=CA:FALSE
extendedKeyUsage = serverAuth
subjectAltName = @alt_names

[alt_names]
DNS.1 = localhost
```

```sh
$ openssl genrsa -out server.key 2048
$ openssl req -sha256 -new -subj "/CN=localhost/O=ZXB Inc/C=CN" -key server.key -out server.csr
$ openssl x509 -req -days 3650 -in server.csr -extfile v3_server.ext -CA rootca.crt -CAkey rootca.key -CAcreateserial -out server.crt
```

## 测试server.crt
```
$ curl -k https://localhost # 不检查服务端证书
$ curl -v --cacert rootca.crt  https://localhost # 检查服务端证书
```

## 生成Client证书
v3_client.ext:
```
authorityKeyIdentifier=keyid,issuer
basicConstraints=CA:FALSE
extendedKeyUsage = clientAuth
```

```sh
$ openssl genrsa -out client.key 2048
$ openssl req -sha256 -new -subj "/CN=AAA/O=AAA Inc/C=CN" -key client.key -out client.csr
$ openssl x509 -req -days 3650 -in client.csr -extfile v3_client.ext -CA rootca.crt -CAkey rootca.key -CAcreateserial -out client.crt
```

## 测试client.crt
```
curl --cacert rootca.crt --cert client.crt --key client.key https://localhost
```