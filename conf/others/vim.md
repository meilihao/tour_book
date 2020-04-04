# vim
## vimrc
位置: `~/.vimrc`

配置:
```conf
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
set statusline=[%n]\ %f%m%r%h\ \|\ \ pwd:\ %{CurDir()}\ \ \|%=\|\ %l/%L,%c\ %p%%\ \|\ ascii=%b,hex=%b%{((&fenc==\"\")?\"\":\"\ \|\ \".&fenc)}\ \|\ %{$USER}\ @\ %{hostname()}\
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
