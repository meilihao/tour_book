# tmux
参考:
- [tmux快捷键](http://blog.csdn.net/hcx25909/article/details/7602935)
- [Tmux 使用教程](https://www.ruanyifeng.com/blog/2019/10/tmux.html)

tmux是一款优秀的终端复用软件，类似 GNU screen，但比screen更出色. 通常当终端关闭后该shell里面运行的任务进程也会随之中断，通过使用tmux就能很容易的解决这个问题.

tmux也可用于会话共享.

## example
```bash
$ tmux list-keys # 列出所有快捷键，及其对应的 Tmux 命令
$ tmux list-commands # 列出所有 Tmux 命令及其参数
$ tmux info # 列出当前所有 Tmux 会话的信息
$ tmux source-file ~/.tmux.conf # 重新加载当前的 Tmux 配置
```

### 分屏操作
- 水平分屏

    先按 ctrl+b, 放开后再按%
- 垂直分屏
    
    先按 ctrl+b, 放开后再按 "
- 分屏后的窗口中的光标互相切换

    先按ctrl+b, 放开后再按下o

## 功能
- session

  - tmux new [-s second-tmux [-n <window-name>]] # 创建session, 名称是second-tmux
  - tmux ls # 查看已创建的会话
  - tmux  a  [-t SESSION-NAME/SESSION-ID]  # 或`tmux  attach  -t  SESSION-NAME` 进入一个已知会话
  - tmux detach # 暂时离开当前会话
  - tmux  kill-session  -t  SESSION-NAME # 关闭session(在会话内部或外部执行均可)
  - tmux switch -t <session-name/session-id> # 切换会话
  - tmux rename-session -t 0 <new-name> # rename session

  - Ctrl+b d：分离当前会话
  - Ctrl+b s：列出所有会话. "->"支持展开session列出其中的terminals, 此时可选择切换终端
  - Ctrl+b $：重命名当前会话

- 窗口

  - tmux new-window [-n <window-name>] : 新建一个(指定名称的)窗口
  - tmux select-window -t <window-number/window-id> : 切换窗口
  - mux rename-window <new-name> : 重命名窗口

  - Ctrl+b c：创建一个新窗口，状态栏会显示多个窗口的信息
  - Ctrl+b p：切换到上一个窗口（按照状态栏上的顺序）
  - Ctrl+b n：切换到下一个窗口
  - Ctrl+b <number>：切换到指定编号的窗口，其中的<number>是状态栏上的窗口编号
  - Ctrl+b w：从列表中选择窗口
  - Ctrl+b ,：窗口重命名
  - Ctrl+b f：查找窗口

  窗口排序:

  - swap-window -s 3 -t 1  交换 3 号和 1 号窗口
  - swap-window -t 1       交换当前和 1 号窗口
  - move-window -t 1       移动当前窗口到 1 号
- 窗格

  - tmux split-window [-h] # 划分窗格(一个窗口划分成多个窗格),  `-h`表示左右划分, 默认为上下划分
  - tmux select-pane -<U/D/L/R> # 上下左右移动光标
  - tmux swap-pane -<U/D> # 窗格上下移动, 实现窗格交换

  - Ctrl+b %：划分左右两个窗格
  - Ctrl+b "：划分上下两个窗格
  - Ctrl+b <arrow key>：光标切换到其他窗格. <arrow key>是指向要切换到的窗格的方向键，比如切换到下方窗格，就按方向键↓
  - Ctrl+b ;：光标切换到上一个窗格
  - Ctrl+b o：光标切换到下一个窗格
  - Ctrl+b {：当前窗格与上一个窗格交换位置
  - Ctrl+b }：当前窗格与下一个窗格交换位置
  - Ctrl+b Ctrl+o：所有窗格向前移动一个位置，第一个窗格变成最后一个窗格
  - Ctrl+b Alt+o：所有窗格向后移动一个位置，最后一个窗格变成第一个窗格
  - Ctrl+b x：关闭当前窗格
  - Ctrl+b !：关闭所有分屏后的窗格，并将当前窗格合并为一个窗口
  - Ctrl+b z：当前窗格全屏显示，再使用一次会变回原来大小
  - Ctrl+b Ctrl+<arrow key>：按箭头方向调整窗格大小
  - Ctrl+b q：显示窗格编号

- 翻页

  有两种方法:
  - ctrl+b进入控制模式, 再用pgUp/pgDn翻页, 按`q`/ESC退出翻页
  - ctrl+b进入控制模式, 再按`[`即可用正常的翻页快捷键(pgUp/pgDn/↑/↓/滚轮), 按`q`/ESC退出翻页

## FAQ
### `configure: error: "libevent not found"`

    sudo apt-get install libevent-dev

### `configure: error: "curses not found"`

    sudo apt-get install ncurses-dev
    sudo dnf install ncurses-devel

### 快捷键失效

按键顺序是`ctrl+b`松开后再按其他键.例如`ctrl+b ？`，应该先同时按`ctrl+b` 松开后，`shift+/（即输入？）`