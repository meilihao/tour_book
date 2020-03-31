# gtags
安装: `sudo apt install global`

## vim配置
```
$ sudo dpkg -L global |grep vim # 查找global for vim的plugin
$ mkdir -p ~/.vim/plugin/
$ cp /usr/share/vim/addons/plugin/gtags* ~/.vim/plugin/
$ vim ~/.vimrc
"gtags 设置项
set cscopetag " 使用 cscope 作为 tags 命令
set cscopeprg='gtags-cscope' " 使用 gtags-cscope 代替 cscope
let GtagsCscope_Auto_Load = 1
let CtagsCscope_Auto_Map = 1
let GtagsCscope_Quiet = 1
let gtags_file=findfile("GTAGS", ";") "查找 gtags 文件
if !empty(gtags_file)
    exe "cs add" gtags_file
endif
$ cd ~/linux-stable # 进入linux kernel源码目录
$ gtags -v # 生成索引后即可使用vim 查看源码
```

命令行操作:
```
$ global -x sched_entity # 查看定义
```

vim 命令模式:
```conf
:Gtags func   # 查看定义处
:Gtags -r func   # 查看引用处
:Gtags -s text   # 查看未被数据库定义的tags
:copen   # 打开quick fix显示窗口
:cclose   # 关闭quick fix显示窗口
:cn   # 下一项
:cp   # 上一项
:cl   # 列出查询到的相关项
:ccN   # 到列表中第N个符号处
:Gtags -g pattern   # 搜索pattern指定的字符串
:Gtags -gie -pattern   # -e选项可以用于搜索’-‘字符，但是基础搜索，没有元字符，-i选项忽略大小写，类似于grep的选项
:GtagsCuorsor   # 取决于光标位置，要是在定义处，查询其引用，要是在引用处，跳转至其定义处，否则就是Gtags -s命令
:Gtags -P text   # 查询包含text的路径名,Gtags -P后接/dir/为列出叫做dir目录下文件，后接\.h$列出所有的include文件
:Gtags -f file   # 列出file里的符号，Gtags -f %则列出当前文件的符号
```

快捷键:
- Ctrl+} 跳转到函数定义处
- Ctrl+t 跳转回来