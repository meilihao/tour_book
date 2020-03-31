# vim
## vimrc
位置: `~/.vimrc`

配置:
```conf
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

## FAQ
### 支持系统剪切板
```
$ vim --version | grep "clipboard" # 查看vim版本是否支持clipboard, 如果`clipboard`前面有一个减号则表示不支持
```

