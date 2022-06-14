# sed

## 描述

sed是一个很好的文件处理工具，本身是一个管道命令，主要是以行为单位进行处理，可以将数据行进行替换、删除、新增、选取等特定工作.

sed是非交互式的编辑器.它不会修改文件，除非使用shell重定向来保存结果.默认情况下，所有的输出行都被打印到屏幕上.

sed编辑器逐行处理文件（或输入），并将结果发送到屏幕.具体过程如下：首先sed把当前正在处理的行保存在一个临时缓存区中（也称为模式空间），然后处理临时缓冲区中的行，完成后把该行发送到屏幕上.sed每处理完一行就将其从临时缓冲区删除，然后将下一行读入，进行处理和显示.处理完输入文件的最后一行后，sed便结束运行.sed把每一行都存在临时缓冲区中，对这个副本进行编辑，所以不会修改原文件.

## 语法格式

```
sed [OPTIONS] [FILE...]
```
将多个命令合并使用,使用分号分隔

## 选项

- -d : 删除，删除选择的行
- -e ∶ 直接在指令列模式上进行 sed 的动作编辑,不修改原文件仅在terminal输出；
- -f ∶ 直接将 sed 的动作写在一个档案内， -f filename 则可以执行 filename 内的sed 动作
- --follow-symlinks : 处理软连接
- -i ∶ 直接修改原文件，而不在屏幕输出. **注意目标文件是软连接时结合`--follow-symlinks`使用, 否则软链会变成常规文件**
- -n ∶ 使用安静(silent)模式。在一般 sed 的用法中，所有来自 STDIN的资料一般都会被列出到萤幕上。但如果加上 -n 参数后，则只有经过sed 特殊处理的那一行(或者动作)才会被列出来。
- -p : 打印指定的数据
- -r ∶ 启用扩展的正规语法.(预设是基础正规表示法语法), 与**perl正则存在差异**, 有且仅有贪婪模式. [此时可用`perl -pe <rule>`代替](https://segmentfault.com/q/1010000005690165)
- -= : 打印当前行号码

sed 's/要替换的字符串/新的字符串/g'（要替换的字符串可以用正则表达式）,末端的g表示在行内进行全局替换.替换时sed默认使用反斜杠（/）分割,紧跟在s命令后的字符就是查找串和替换串之间的分隔符.分隔符默认为正斜杠，但可以改变。无论什么字符（换行符、反斜线除外），只要紧跟s命令，就成了新的串分隔符.如果“原内容”或“新内容”中含有特殊字符（比如”/”或者”#”等），可以使用其它符号把各部分隔开,比如`,`,`#`,`:`,`~`.

function：
- a ：新增， a 的后面可以接字串，而这些字串会在新的一行出现(目前的下一行)～
- c ：取代， c 的后面可以接字串，这些字串可以取代 n1,n2 之间的行！
- d ：删除，因为是删除啊，所以 d 后面通常不接任何咚咚；
- i ：插入， i 的后面可以接字串，而这些字串会在新的一行出现(目前的上一行)；
- p ：列印，亦即将某个选择的数据印出。通常 p 会与参数 sed -n 一起运行～
- s ：取代，可以直接进行取代的工作哩！通常这个 s 的动作可以搭配正规表示法！例如 1,20s/old/new/g 就是啦！

sed内置命令N的作用: 不会清空模式空间(pattern space)内容, 并且从输入文件中读取下一行数据, 追加到模式空间中, 两行数据以`\n`连接. 即将当前读入行和用N命令添加的下一行拼成"一行".

参考:
- [sed命令详解](http://www.cnblogs.com/edwardlost/archive/2010/09/17/1829145.html)

## 例
```
$ sed -i 's/^#\s*StrictHostKeyChecking.*$/    StrictHostKeyChecking no/' /etc/ssh/ssh_config # 整行替换
$ sed -n '3,7 p' data # 打印3~7行
$ sed -n '/linux/I p' data # 查找指定字符串
$ sed -i "s:/static:/blog/static:" `grep /static -rl ./` # 检索当前目录下的文件,将其包含的字符串"/static"替换为"/blog/static".
$ sed -i  "s:"action='"'":"action='"'/blog":" `grep action= -rl ./` 将`action="/comment/{{.Content.Id}}/"`替换为`action="/blog/comment/{{.Content.Id}}/"`
$ sed -i '1,nd' 2016.txt # 删除前n行
$ sed 's/text/replace_text/' file  # 替换每一行的第一处匹配的text
$ sed 's/text/replace_text/g' file  # 全局替换
$ sed '/^$/d' file # 移除空白行
$ sed '/qq/d' file # 移除包含"qq"的行
$ echo this is en example | sed 's/\w+/[&]/g' # 已匹配的字符串通过标记&来引用
$ p=patten
$ r=replaced
$ echo "line con a patten" | sed "s/$p/$r/g" # 双引号会对表达式求值
$ sed -n '2,6p' access.log # 输出指定行(2~6)
$ sed -e 'i\head' access.log # 在每行的前面插入行"head" # 行尾的反斜杠( \ )允许你在该流中插入新行而不会提前执行命令,类似于延迟执行.
$ sed -e 'a\end' access.log # 在每行的后面追加行"end"
$ sed -e '/google/c\hello' access.log # 将包含"google"的行替换成字符串"hello"
$ sed -n '1,5p;1,5=' access.log # 显示1~5行,并在结果的1~5行上添加行号
$ sed -i '/原行内容/a要添加的新一行内容' 文件
$ sed -i '325a要添加的新一行内容' 文件
$ sed -i '6i 要添加的新一行内容' 文件 # 插入内容
$ sed -n "/test_string/p" a.log # **仅检测sed匹配到的内容**
$ sed -i "s@if \[ -f /etc/exports \] && grep -q@if [ -f /etc/exports ] # \&\& grep -q@" b.bak # `&`表示引用及引用被匹配到的字符串, 因此必须转义, 否则此例会导致在被替换字符串中出现两次匹配内容.
$ sed -i "/LFS_Sources_Root=\${LFSRoot}\/sources/d" build/tcl.sh # 删除行
$ echo sources/linux-5.8.1.tar.xz |sed 's/linux\-\(.*\)\.tar\.xz/\1/g' => sources/5.8.1 # 提取变量, 因为sed是匹配替换因此"sources/"被保留了
$ sed -i '/2222222222/a\3333333333' test.txt # 在匹配行前插入新行
$ sed -i '/2222222222/i\3333333333' test.txt # 在匹配行后追加新行
$ sed -i "s/max_dbs_open = 100.*/max_dbs_open = 1000/g" /etc/couchdb/default.ini # 匹配里使用通配符
$ sed -i -r "s/([0-9]{1,3}\.){3}([0-9]{1,3})/127.0.0.1/g" <file> # 替换文件中的ip
$ pattern1=XXX
$ sed -i "s/aaa/$pattern1/g" inputfile # 如果要使用shell变量，就需要使用双引号
$ for file in `ls | grep .txt` # 批量文件重命名
do
 newfile=`echo $file | sed 's/\([a-z]\+\)\([0-9]\+\)/\1-\2/'`
 mv $file $newfile
done
```

ps:
sed一般使用单引号，**sed引用shell变量时使用双引号**即可，因为双引号是弱转义，不会去除$的变量表示功能，而单引号为强转义，会把$作为一般符号表示，所以不会表示为变量