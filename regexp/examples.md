# 例子

- `^.*xxx.*$`: 匹配含xxx的行
- `^$`: 空行
- `^\s*$` : 匹配任意空或空白字符,即空行
- `^[0-9]*$` :	只能输入数字
- `^\d{n}$` :	只能输入n位的数字
- `^\d{n,}$` :	只能输入至少n位的数字
- `^\d{m,n}$` :	只能输入m~n位的数字
- `^(0|[1-9][0-9]*)$` :	只能输入零和非零开头的数字
- `^[0-9]+(.[0-9]{2})?$` :	只能输入有两位小数的正实数
- `^[0-9]+(.[0-9]{1,3})?$` :	只能输入有1~3位小数的正实数
- `^\+?[1-9][0-9]*$` :	只能输入非零的正整数
- `^\-[1-9][0-9]*$` :	只能输入非零的负整数
- `^.{3}$` :	只能输入长度为3的字符
- `^[A-Za-z]+$` :	只能输入由26个英文字母组成的字符串
- `^[A-Za-z0-9]+$` :	只能输入由数字和26个英文字母组成的字符串
- `^\w+$` :	只能输入由数字、26个英文字母或者下划线组成的字符串
- `^[a-zA-Z]\w{5,17}$` :	验证用户密码：以字母开头，长度在6~18之间，只能包含字符、数字和下划线。
- `[^%&',;=?$\x22]+` :	验证是否含有^%&',;=?$\"等字符
- `^[\u4e00-\u9fa5]{0,}$` :	只能输入汉字
- `/^[\u4e00-\u9fa5]+$/` : utf8中文
- `^\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*$` :	验证Email地址
- `^http://([\w-]+\.)+[\w-]+(/[\w-./?%&=]*)?$` :	验证InternetURL
- `^\d{15}|\d{18}$` :	验证身份证号（15位或18位数字）
- `^((2[0-4]\d|25[0-5]|[01]?\d\d?)\.){3}(2[0-4]\d|25[0-5]|[01]?\d\d?)$` :	验证IP地址
- `(\w)\1` :	匹配重叠出现的字符
- `<(?<tag>[^\s>]+)[^>]*>.*</\k<tag>>` : 匹配成对的HTML标签
- `^\d{5,12}$` : 	QQ号必须为5位到12位数字
- `\(?0\d{2}[) -]?\d{8}` : 	(010)88886666，或022-22334455，或029 12345678
- `0\d{2}-\d{8}|0\d{3}-\d{7}`	: 8位本地号(如010-12345678)和7位本地号(0376-2233445)
- `(\d{1,3}\.){3}\d{1,3}` :	IP地址匹配
- `\b(\w+)\b\s+\1\b`	 	匹配相连并重复的单词
- `\b\w+(?=ing\b)` :	匹配以ing结尾的单词的前面部分(除了ing以外的部分)，如查找`I’m singing while you’re dancing.`时，会匹配sing和danc
- `^(?![A-Za-z]+$)(?!\d+$)(?![\W_]+$)\S{6,16}$` : 密码(6-16), 字母、数字、符号至少包含2种
- `((?=.*\d)(?=.*\D)|(?=.*[a-zA-Z])(?=.*[^a-zA-Z]))^.{6,16}$` : 密码(6-16), 字母、数字、符号至少包含2种
- ` *$` : 匹配末尾的空格
- `^/([^/]*)/?.*$` : 提取url中的第一个分片, 比如`/a/b`中的`a`, `/?`表示后面的`/`可有可无. 用于匹配`/a`的情况
- `^[[:space:]]*[^#]*/` : 匹配nfs的/etc/exports有有效的导出规则即`/`前有任意个空白字符和任意个非`#`的字符.
- `^(([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\-]*[a-zA-Z0-9])\.)*([A-Za-z0-9]|[A-Za-z0-9][A-Za-z0-9\-]*[A-Za-z0-9])$` : [匹配hostname](https://stackoverflow.com/questions/106179/regular-expression-to-match-dns-hostname-or-ip-address)
- [`(?=^.{3,50}$)(^[a-z0-9][a-z0-9-]+[a-z0-9]$)`](https://stackoverflow.com/questions/50480924/regex-for-s3-bucket-name/50484916) : 匹配minio, aliyun ,huawei, tencent obs bucket name

## replace
- `: .\[1;\d\dm(.+).\[0m` -> `: \1` ; `: [1;37mINFO[0m -`->`: INFO -` # 去掉log中的ANSI escape code
