# vim
参考:
- [vim快捷键大全图解](https://segmentfault.com/a/1190000016056004)

文本编辑器

Vim 的三种模式:
- 命令模式：控制光标移动，可对文本进行复制、粘贴、删除和查找等工作
- 输入模式：正常的文本录入
- 末行(控制)模式：保存或退出文档，以及设置编辑环境

![vim小抄](/misc/img/shell/2011061002323916.jpeg)
![](/misc/img/shell/1148777987-547fc1847c4b4.webp)

## example
```
# vim + file # 打开文件, 将光标置于最后一行
# vim +n file # 打开文件, 将光标置于第n行
# vim +/pattern file # 打开文件, 将光标置于第一个与pattern匹配的位置
# vim -R file # 只读打开
# vim -r file # 恢复删除vim打开时崩溃的文件
# vim xxx.txt -c "e ++enc=gbk" # 指定使用gbk编码打开文件
```

## 编码
```
:set fileencoding # 查看编码
:set fileencoding=utf-8 # 编码修改
:set list # 显示出特殊字符, 比如换行
:set ff=unix # 将文件换行转为unix换行. ff=fileencoding
:set ff=dos # 将文件换行转为windows换行
```

## 命令/快捷键

### 其他
- :paste 从剪切板粘贴时保留格式
- :set nu 显示行号
- :set nonu 不显示行号
- :命令 执行该命令
- :整数 跳转到该行
- <百分比>% 按百分比跳转
- v : 进入Visual模式
- u : 撤销上一步的操作
- :<start>,<end>d + 回车
- >> 按2次> 当前行增加缩进 
- << 按2次< 当前行减少缩进
- :10,100> : 第10行至第100行缩进
- :20,80< : 第20行至第80行反缩进
- Visual 模式下缩进: normal 模式下按v 即进入VISUAL模式，可选择多行. 在选择好需要缩进的行后，按一次大于号’>’缩进一次，按’6>’缩进六次，按’<’回缩
- Visual 模式下注释: normal 模式下按Control+v进入VISUAL模式, 再按`I`进入插入模式, 输入注释符“#”或者是"//"(此时仅在光标处显示注释符), 然后立刻按下ESC（两下)后显示全部注释符;Ctrl + v 进入块选择模式，选中要删除的行首的注释符号, 注意 // 要选中两个, 选好之后按 d 即可删除注释，ESC 保存退出.
- 命令模式下注释

    - :10,20s#^#//#g : 在 10 - 20 行添加 // 注释
    - :10,20s#^//##g : 在 10 - 20 行删除 // 注释
### 移动光标
- w : 光标向右移动一个单词
- nw : 光标向右移动n个单词
- b : 光标向左移动一个单词
- nb : 光标向左移动n个单词
- nG : 光标移至第n行的行首
- G : 将光标移动到文件的最后一行的开头(第一个非空字符)
- n+ : 光标下移n行
- n- : 光标上移n行
- n$ : 相对于当前行, 将光标后移至n行的行尾
- H : 将光标移至当前屏幕的顶行
- M : 将光标移至当前屏幕的中间行
- L : 将光标移至当前屏幕的最底行
- ^ : 将光标移至当前行的开头(第一个非空字符)
- 0 : 将光标移至当前行首
- $ : 将光标移至当前行尾
- :$ : 将光标移至文尾
- { : 将光标移动到前面的`{`处, 对于c/go之类的编程语言很有用
- } : 将光标移动到后面的`}`处

### 屏幕翻滚
- ctrl+u : 相对于屏幕, 向文件首翻半屏
- ctrl+d : 相对于屏幕, 向文件尾翻半屏
- ctrl+b : 相对于屏幕, 向文件首翻一屏
- ctrl+f : 相对于屏幕, 向文件尾翻一屏
- ctrl+e : 相对于屏幕, 向下翻一行
- ctrl+y : 相对于屏幕, 向上翻一行
- nz+Enter : 将文件的第n行滚至屏幕顶部, 如果不指定n值, 将当前行滚至屏幕顶部

### 插入/删除
- Esc : 退出命令行模式
- i : 在光标处插入
- I : 在光标所在行的行首插入
- a : 在光标处后插入
- A : 在光标所在行的行尾插入
- o : 在当前行之下新开一行
- O : 在当前行之上新开一行
- r : 替换光标所在的字符, 输入r后,在键盘输入新字符即可完成替换
- yy : 将光标所在行复制到剪贴板
- nyy : 类似yy, 但会复制n行
- yw : 将光标所在的单词复制到剪贴板
- nyw : 类似yw, 但会复制n个单词
- p : 将剪贴板的内容复制到光标后
- P : 将剪贴板的内容复制到光标前, 即将之前删除（dd）或复制（yy）过的数据粘贴到光标后面
- x : 删除光标所在的字符
- nx : 删除光标所在及之后的n-1个字符
- X : 删除光标所在的前一个字符
- nX : 删除光标所在及之前的n-1个字符
- dw : 删除光标所在的单词
- ndw : 删除光标所在及之后的n-1个单词
- d0 : 删除当前行中光标所在前的所有字符
- d$ : 删除当前行中光标所在后的所有字符
- dd : 删除光标所在的行. 删除的内容自动保存到剪贴板
- <n>dd : 删除光标所在的行及其向下的n-1行. 删除的内容自动保存到剪贴板
- <n>d+上方向键 : 删除光标所在的行及其向上的n行
- <n>d+下方向键 : 删除光标所在的行及其向下的n行

- ggvG/ggVG : 全选（高亮显示)
- ggyG : 全部复制
- dG : 全部删除

### 搜索/替换
参考:
- [Vim查找与替换命令大全，功能完爆IDE!](https://segmentfault.com/a/1190000022323247)

- `/abc[\C|\c]` : 从上向下查找字符串abc, `\C`区分大小写, `\c`不区分大小写, 也可用`:set noic|:set ic`实现, ic是ignorecase简写. 带`/`字符串可用`\/`转义
- `?abc` : 从下向上查找字符串abc
- `n` : 在同一方向上重复执行上次的搜索命令
- `N` : 在相反方向上重复执行上次的搜索命令
- `:s/one/two` 将当前光标所在行的第一个 one 替换成 two 
- `:s/a1/a2/g` : 将当前光标所在的行中的所有a1替换为a2
- `:n1,n2s/a1/a2/g` : 将n1~n2行中的所有a1替换为a2
- `:g/a1/a2/g` : 将文中的所有a1替换为a2
- `:%s/a1/a2/g` : 将文中的所有a1替换为a2
- `:set nu/number` : 添加行号

### 保存/退出
- `:wq` : 保存并退出
- `:w [filename]` : 保存, 但不退出, 提供文件名时是另存为.
- `:w! [filename]` : 强制保存, 但不退出, 提供文件名时是另存为.
- `:w file` : 另存为
- `:q` : 不保存就退出
- `:q!` : 不保存, 强制退出

## vimrc
参考:
- [Vim 配置入门](https://www.ruanyifeng.com/blog/2018/09/vimrc.html)

位置: `~/.vimrc`

配置:
```conf
" 注释的文本在左侧使用双引号即可
set hlsearch " 搜索时，高亮显示匹配结果
set showcmd
set showmode
set pastetoggle=<F9>
set showmode " 在底部显示，当前处于命令模式还是插入模式
set ts=4 ":set tabstop=4 设定tab宽度为4个字符
set expandtab " 空格代替tab
syntax on " 语法高亮
" 状态栏
set laststatus=2      " 总是显示状态栏
highlight StatusLine cterm=bold ctermfg=yellow ctermbg=black
" 获取当前路径，将$HOME转化为~
function! CurDir()
        let curdir = substitute(getcwd(), $HOME, "~", "g")
        return curdir
endfunction
set statusline=[%n]\ %f%m%r%h\ \|\ %{CurDir()}\ \ \|%=\|\ %l/%L,%c\ %p%%\ \|\ ascii=%b,hex=%b%{((&fenc==\"\")?\"\":\"\ \|\ \".&fenc)}\ \|\ %{$USER}\ @\ %{hostname()}
set ruler "在编辑过程中，在右下角显示光标位置的状态行
set number
if has('mouse')
    set mouse=a " 使用鼠标光标定位即"鼠标点哪光标跳哪", **不好用, 不推荐**
endif
```

## nvim
Neovim 是能使用 vim 的配置文件的，如果有 vim 的配置，直接软链接就好：
```bash
$ ln -s ~/.vim ~/.config/nvim
$ ln -s ~/.vimrc ~/.config/nvim/init.vim
```

> 如果想 nvim 单独使用一个配置，那就在 `~/.config` 下创建配置文件就行:`~/.config/nvim`和`~.config/nvim/init.vim`

有时 neovim 的某些指令在 vim 中是不能使用的, 所以可使用 has('nvim') 来判断当前使用的版本：
```conf
if has('nvim')
    ...
endif
```

## vim-plug

能在 vim/neovim 中使用的插件管理工具有不少，这里**推荐**的是 [vim-plug](https://github.com/junegunn/vim-plug), 安装见readme.

vim-plug使用:
```bash
$ vim ~/.vimrc # 在~/.vimrc开头添加要安装的plugin
call plug#begin('~/.vim/plugged')
Plug 'neoclide/coc.nvim', {'branch': 'release'}
if has('nvim')
  Plug 'Shougo/defx.nvim', { 'do': ':UpdateRemotePlugins' }
else
  Plug 'Shougo/defx.nvim'
  Plug 'roxma/nvim-yarp'
  Plug 'roxma/vim-hug-neovim-rpc'
endif
call plug#end()
```

然后重启vim, 再执行`:PlugInstall`开始安装plugins即可.

## FAQ
### 支持系统剪切板
```
$ vim --version | grep "clipboard" # 查看vim版本是否支持clipboard, 如果`clipboard`前面有一个减号则表示不支持
```

### 查看keymap冲突
```
:verbose imap <tab> # <tab>为快捷键
```

### vim粘贴代码格式变乱
复制粘贴代码到vim出现格式错乱. 其原因是vim开启了smartindent(智能缩减)或autoindent(自动对齐)模式. 为了保持代码的格式，在粘贴前可以先停止上面的两种模式，命令为：
```conf
set nosmartindent # 不推荐
set noautoindent # 不推荐

set pastetoggle=<F9> # 推荐, Vim的编辑模式(Insert)中，还有一个Paste模式，在该模式下，可将文本原本的粘贴到Vim中，以避免一些格式错误. 通过“:set paste”和“:set nopaste”进入和退出该模式; 也可通过快捷键按 F9 键来打开和关闭paste选项.
```

> 原因是终端把粘贴的文本存入键盘缓存（Keyboard Buffer）中，Vim则把这些内容作为用户的键盘输入来处理. 导致在遇到换行符的时候，如果Vim开启了自动缩进，就会默认的把上一行缩进插入到下一行的开头，最终使代码变乱.

### vim打开中文乱码
```
$ vim ~/.vimrc
...
set fileencodings=utf-8
set termencoding=utf-8
set encoding=utf-8
```

### 保存文件时报"readonly option is set"
解决方法： =
1. 按ESC键
1. 输入:set noreadonly
1. 然后就可以用 :wq 保存并退出了

### 获取vim的选项
`vim -c ':options'`

### 查询某个配置项是打开还是关闭
`:set number?`

### 修改压缩包内容
vim可以直接打开压缩包, 然后对想要的压缩文件进行编辑保存, 一切都是vim的命令完成, 修改后`w`就保存一直使用`q`退出了. 注意: 被编辑的文件会丢失可执行权限, 因此不推荐用来编辑脚本.

### vim中处理windows下的文档换行符
`:set fileformat=unix`再`:wq`即可.

### [如何在基于 Ubuntu 的 Linux 发行版上安装最新的 Vim 9.0](https://linux.cn/article-14899-1.html)
```bash
sudo add-apt-repository ppa:jonathonf/vim
sudo add-apt-repository -r ppa:jonathonf/vim # 移除repo
```

### vim配置文件位置
进入vim,输入命令`:version`,找到vimrc相关内容即是.

### 关闭自动创建备份档(即*~的文件)

> 检查vimrc是否存在(用户配置文件为~/.vimrc，相关的文件位于~/.vim/；全局配置文件为/etc/vimrc，相关的文件位于/usr/share/vim/),不存在时,`cp /usr/share/vim/vim74/vimrc_example.vim ~/.vimrc`
> 打开配置文件，找到这一句：if has("vms"),将这个判断里的if部分保留，else部分注释（Vim的注释符是"）掉`set backup`即可.