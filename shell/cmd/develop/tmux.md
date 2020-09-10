# tmux
参考:
- [tmux快捷键](http://blog.csdn.net/hcx25909/article/details/7602935)

tmux是一款优秀的终端复用软件，类似 GNU screen，但比screen更出色. 通常当终端关闭后该shell里面运行的任务进程也会随之中断，通过使用tmux就能很容易的解决这个问题.

## example
```bash
$ tmux new -s second-tmux # 创建session, 名称是second-tmux
$ tmux ls # 查看已创建的会话
$ tmux  a  -t SESSION-NAME  # 或`tmux  attach  -t  SESSION-NAME` 进入一个已知会话
$ tmux detach # 暂时离开当前会话
$ tmux  kill-session  -t  SESSION-NAME # 关闭session(在会话内部或外部执行均可)
$ 切换tmux会话终端
```

### 分屏操作
- 水平分屏

    先按 ctrl+b, 放开后再按%
- 垂直分屏
    
    先按 ctrl+b, 放开后再按 "
- 分屏后的窗口中的光标互相切换

    先按ctrl+b, 放开后再按下o

## 快捷键
- 翻页

  有两种方法:
  - ctrl+b进入控制模式, 再用pgUp/pgDn翻页, 按`q`退出翻页
  - ctrl+b进入控制模式, 再按`[`即可用正常的翻页快捷键(pgUp/pgDn/↑/↓/滚轮), 按`q`退出翻页
- 切换tmux会话终端

    先按ctrl+b, 放开后再按s 
- 查看面板编号

    先按ctrl+b, 放开后再按q
- 关闭所有分屏后的窗口，即合并为一个窗口

    先按ctrl+b, 放开后再按！
- 暂时退出当前会话

    先按ctrl+b, 放开后再按 d
- 在当前窗口的基础上再打开一个新的窗口
    先按ctrl+b, 放开后再按c

## FAQ
### `configure: error: "libevent not found"`

    sudo apt-get install libevent-dev

### `configure: error: "curses not found"`

    sudo apt-get install ncurses-dev
    sudo dnf install ncurses-devel

### 快捷键失效

按键顺序是`ctrl+b`松开后再按其他键.例如`ctrl+b ？`，应该先同时按`ctrl+b` 松开后，`shift+/（即输入？）`