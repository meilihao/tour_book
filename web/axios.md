# axios

## FAQ
1. post formdata
```js
// create an axios instance
const service = axios.create({
  baseURL: process.env.BASE_API, // api的base_url
  timeout: 5000, // request timeout
  transformRequest: (data, headers) => {
    if (headers['Content-Type'] && headers['Content-Type'] === 'multipart/form-data') {
      const formData = new FormData()
      for (const name in data) {
        formData.set(name, data[name])
      }
      return formData
    }

    return data
  }
})
```

使用
```js
service({
    url: '/login',
    method: 'post',
    data,
    headers: { 'Content-Type': 'multipart/form-data' }
  })
```