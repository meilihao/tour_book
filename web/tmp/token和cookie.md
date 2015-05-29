# OAuth Token和cookie

## 区别

token机制是为了防止cookie被清除，另外cookie是会在所有域名请求都携带上，无意中增加了服务端的请求量，token只需要在有必要的时候携带。