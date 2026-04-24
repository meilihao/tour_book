# nvim
## 快捷键
- ctrl+shift+v : 粘贴

## FAQ
### neovim乱码
参考:
 - [Nvim shows weird symbols (�[2 q) when changing modes](https://github.com/neovim/neovim/wiki/FAQ#nvim-shows-weird-symbols-2-q-when-changing-modes)
 - [[RDY] Fix incorrect DECSCUSR fixup codes #6997](https://github.com/neovim/neovim/pull/6997)

 问题出现在xterm及其兼容版本下, 解决方法:
 ```sh
 $ echo 'set guicursor=' > ~/.config/nvim/init.vim # 禁用guicursor
 ```

 > 查看当前使用的term: `echo $TERM`
 > 查看系统支持的term: `tree /usr/share/terminfo`