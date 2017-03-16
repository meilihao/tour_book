# strings

import "strings"

实现了用于操作字符的简单函数.

### func TrimPrefix(s, prefix string) string

去除s可能含有的prefix(前缀)字符串

参数列表:
- s 原始字符串
- prefix 要去除的前缀字符串


    fmt.Println(strings.TrimPrefix("abc","a")) // "bc"
	fmt.Println(strings.TrimPrefix("abc","b")) // "abc"

### func ToLower(s string) string

将所有字母都转为对应的小写

参数列表:
- s 原始字符串

    fmt.Println(strings.ToLower("aBc")) // "abc"
	fmt.Println(strings.ToLower("abc")) // "abc"

### func ToUpper(s string) string

将所有字母都转为对应的大写

参数列表:
- s 原始字符串

    fmt.Println(strings.ToUpper("aBc")) // "ABC"
	fmt.Println(strings.ToUpper("abc")) // "ABC"