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

	w.Close()

	fmt.Println(body)
	//request
	req, _ := http.NewRequest("POST", "http://127.0.0.1:8080/test", body)
	req.Header.Set("Content-Type", content_type)

	resp, _ := http.DefaultClient.Do(req)
	data, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
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