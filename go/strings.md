## 统计字符串长度

- 使用 bytes.Count() 统计
- 使用 strings.Count() 统计
- 将字符串转换为 []rune 后调用 len 函数进行统计
- 使用 utf8.RuneCountInString() 统计,**推荐**
