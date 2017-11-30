# 地址

 参考:
 - [最新县及县以上行政区划代码](http://www.stats.gov.cn/tjsj/tjbz/xzqhdm/)
 - [统计用区划和城乡划分代码](http://www.stats.gov.cn/tjsj/tjbz/tjyqhdmhcxhfdm/),比上面的`最新县及县以上行政区划代码`更细.

## 使用淘宝的收货地址
参考:
- [利用js+php的技术，实现全国地址多级联动的功能](http://blog.csdn.net/qingxinyeren/article/details/51531216)

淘宝地址在[我的淘宝-账户设置-收货地址](https://member1.taobao.com/member/fresh/deliver_address.htm)里,这里还有一个[地址iframe形式的封装](http://member1.taobao.com/member/fresh/deliver_address_frame.htm),它们的功能是一样的.

转换步骤:
1. 下载网页, 用命令`grep -r "A-G"`检索地址入口文件,得到`http://g.alicdn.com//vip/address/6.0.14/index-min.js`(`A-G`是淘宝地址选项卡省份分类的一个标识).
1. 用网上工具格式该js,提取地址信息.
1. 用js解码获取所需地址信息

## 自己处理
```go
package main

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"strings"
	"unicode/utf8"
)

type Place struct {
	Id        string   `json:"id"`
	Name      string   `json:"name"`
	ShortName string   `json:"short_name"`
	ParentId  string   `json:"parent_id"`
	Parent    *Place   `json:"-"`
	Level     int      `json:"level"`
	LevelTrue int      `json:"level_true"`
	Location  string   `json:"location"` // gps
	Children  []*Place `json:"children,omitempty"`
}

func main() {
	root := &Place{
		Id:        "CN",
		Name:      "中国",
		ShortName: "中国",
		ParentId:  "0",
	}

	f, err := os.Open("addr.txt")
	CheckErr(err)
	defer f.Close()

	s := bufio.NewScanner(f)

	var parent *Place
	var tmp string
	for s.Scan() {
		tmp = s.Text()
		if tmp == "" {
			continue
		}

		fds := strings.Fields(tmp)
		if len(fds) != 2 {
			panic(fds)
		}

		p := &Place{
			Id:   fds[0],
			Name: fds[1],
		}

		if strings.HasSuffix(fds[0], "0000") {
			parent = root

			p.ParentId = parent.Id
			p.Parent = parent
			p.Level = parent.Level + 1
			p.LevelTrue = 1
			p.ShortName = CleanName(p)

			root.Children = append(root.Children, p)

			parent = p
		} else if strings.HasSuffix(fds[0], "00") {
			if parent.Level != 1 {
				parent = parent.Parent
			}

			// 挂载到省级节点
			if p.Name == "省直辖县级行政区划" || p.Name == "自治区直辖县级行政区划" || p.Name == "县" {
				continue
			}

			if p.Name == "市辖区" {
				p.Name = parent.Name
			}

			p.ParentId = parent.Id
			p.Parent = parent
			p.Level = parent.Level + 1
			p.LevelTrue = 2
			p.ShortName = CleanName(p)

			parent.Children = append(parent.Children, p)

			parent = p
		} else {
			if p.Name == "市辖区" {
				continue
			}

			p.ParentId = parent.Id
			p.Parent = parent
			p.Level = parent.Level + 1
			p.LevelTrue = 3
			p.ShortName = CleanName(p)

			parent.Children = append(parent.Children, p)
		}
	}
	CheckErr(s.Err())

	// data, err := json.Marshal(root)
	// CheckErr(err)

	// fmt.Println(string(data))
}

func CleanName(p *Place) string {
	switch p.Level {
	case 1:
		r := regexp.MustCompile("(.*?)(市|省|特别行政区|.族.*|维吾尔自治区|自治区)$")
		t := r.FindStringSubmatch(p.Name)
		if len(t) > 0 {
			p.Name = t[1]
		}
	case 2:
		n := utf8.RuneCountInString(p.Name)
		if n > 2 {
			r := regexp.MustCompile("(.*?)(蒙古自治.*|哈萨克.*|布依族.*|.家族.*|.族.*|自治州|市|地区|县|盟|林区)$")
			t := r.FindStringSubmatch(p.Name)
			if len(t) > 0 {
				p.Name = t[1]
			}
        }
        
        // n = utf8.RuneCountInString(p.Name)
		// if n > 2 {
		// 	fmt.Println(p.Name)
		// }
	case 3:
		n := utf8.RuneCountInString(p.Name)
		if n > 2 {
			r := regexp.MustCompile("(.*?)(蒙古族自治.*|蒙古自治.*|哈萨克.*|达斡尔族.*|塔吉克.*|锡伯.*|.族.*|市|新区|区|自治旗|旗|县)$")
			t := r.FindStringSubmatch(p.Name)
			if len(t) > 0 {
				p.Name = t[1]
			}
		}
	}

	return p.Name
}

func CheckErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
```