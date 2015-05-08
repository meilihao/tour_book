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