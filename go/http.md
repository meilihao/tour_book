### http例子

```go
package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"io/ioutil"
)

func main() {
	u := url.URL{}
	u.Scheme = "http"
	u.Host = "sendcloud.sohu.com"
	u.Path = "/webapi/mail.send_template.json"
	fmt.Println("url:",u.String())

	q := url.Values{}
	q.Set("entry", "weibo")
	q.Add("gateway", "1")
	fmt.Println("url.Values:",q.Encode())

	u.RawQuery=q.Encode()
	fmt.Println("url+url.Values:",u.String())
	u.RawQuery=""

    //data := `substitution_vars={"to": ["` + to + `"],"sub":{"%name%": ["` + user_name + `"],"%usrl_string%":["` + url_string + `"]}}`
    //将字符串以URL编码,再传入http.Post
	//body := url.QueryEscape(data)

	http.Get("www.baidu.com?" + q.Encode())
	//参数要设置成"application/x-www-form-urlencoded",否则post参数无法传递
	http.Post(u.String(), "application/x-www-form-urlencoded", strings.NewReader(q.Encode()))
	http.PostForm(u.String(), q)

	//复杂请求
	//需要在请求的时候设置头参数、cookie之类的数据
	//post请求时,必须要设定Content-Type为application/x-www-form-urlencoded,post参数才可正常传递
	client := &http.Client{}
	req, _ := http.NewRequest("POST", "http://www.01happy.com/demo/accept.php", strings.NewReader(q.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", "name=anny")

	resp, err := client.Do(req)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {

	}
	fmt.Println("http.Do:",string(body))
}
```

### http post 带文件

client:

```go
    //see "Golang Multipart File Upload Example":http://matt.aimonetti.net/posts/2013/07/01/golang-multipart-file-upload-example/
	file_data, _ := ioutil.ReadFile("main.go")

	body := new(bytes.Buffer)

	//multipart实现了MIME的multipart解析
	w := multipart.NewWriter(body)
	//header
	//方法返回w对应的HTTP multipart请求的Content-Type的值，多以multipart/form-data起始
	content_type := w.FormDataContentType()
	fmt.Println(content_type)

	//form's other param
	w.WriteField("name", "test")
    //upload file
	file, _ := w.CreateFormFile("file", "main.go")
	file.Write(file_data)

	w.Close()//需在http.NewRequest()前关闭,否则服务端FormFile()时会碰到错误"ERROR: unexpected EOF"

	fmt.Println(body)
	//request
	req, _ := http.NewRequest("POST", "http://127.0.0.1:8080/test", body)
	req.Header.Set("Content-Type", content_type)

	resp, _ := http.DefaultClient.Do(req)
	data, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close() // 仅在作为客户端时,http包才需要手动关闭Response.Body;如果是作为服务端,http包会自动处理Request.Body. 目前http包的实现逻辑是只有当应答的Body中的内容被全部读取完毕且调用了Body.Close(),默认的HTTP客户端才会重用带有keep-alive标志的HTTP连接,否则每次HTTP客户端发起请求都会单独向服务端建立一条新的TCP连接,这样做的消耗要比重用连接大得多
	fmt.Println(resp.StatusCode)
	fmt.Println(string(data))
```

server:

```go
//推荐
func HandFile1(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20) //和r.MultipartReader一起用会冲突
	//MultipartForm是解析好的多部件表单，包括上传的文件.该字段只有在调用ParseMultipartForm后才有效
	form := r.MultipartForm
	fmt.Println(form)

	file, _, err := r.FormFile("file")
	if err != nil {
		log.Println(err)
		return
	}
	defer file.Close()

	err = os.MkdirAll("test", 0644)
	//在dir目录下创建一个新的、使用prefix为前缀的临时文件，以读写模式打开该文件并返回os.File指针。
	out, err := ioutil.TempFile("test", "file")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer out.Close()

	if _, err := io.Copy(out, file); err != nil {
		log.Fatal(err)
	}

	fmt.Fprint(w, []byte("123"))
}

func HandFile2(w http.ResponseWriter, r *http.Request) {
	//获取上传文件大小(有较小误差)
	fmt.Println("file's zise ≈", r.ContentLength)
	//如果请求是multipart/form-data POST请求，MultipartReader返回一个multipart.Reader接口，否则返回nil和一个错误。使用本函数代替ParseMultipartForm，可以将r.Body作为流处理。
	rd, err := r.MultipartReader()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for {
		p, err := rd.NextPart()
		if err == io.EOF {
			break
		}

		if p.FormName() == "file" {
			filename := p.FileName()
			file, err := ioutil.TempFile("test", filename)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if _, err := io.Copy(file, p); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}

	fmt.Fprint(w, []byte("123"))
}
func main() {
	http.HandleFunc("/test", HandFile2)
	http.ListenAndServe(":8081", nil)
}
```

### tls

- http://colobu.com/2016/06/07/simple-golang-tls-examples/
- http://tonybai.com/2015/04/30/go-and-https/
- https://github.com/nareix/blog/blob/master/posts/golang-tls-guide.md

## FAQ
### HTTP客户端默认不会及时关闭已经用完的HTTP连接
Go标准库HTTP客户端的默认实现并不会及时关闭已经用完的HTTP连接(仅当服务端主动关闭或要求关闭时才会关闭).

解决方法:
1. 将http.Request中的字段Close设置为true
1. 通过创建一个http.Client新实例来实现的(不使用DefaultClient)

	```go
	tr := &http.Transport{ // http.Transport 默认启用连接复用（Keep-Alive）
		DisableKeepAlives: true,
	}
	cli := &http.Client{
		Transport: tr,
	}
	```
### 常用client
```go
bClient:=&http.Client{
	Timeout: 15*time.Minute, // 全局超时为 15 分钟, 涵盖从发起到完成的总时间（包括连接、传输、重定向等）
	Transport: &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 5*time.Second
		}).Dial,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		}
	}
}
```