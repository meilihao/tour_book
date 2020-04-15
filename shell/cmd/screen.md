# screen
它是一款能够实现多窗口远程控制的开源服务程序. 但推荐使用tmux.

下功能如下:
- 会话恢复：即便网络中断，也可让会话随时恢复，确保用户不会失去对远程会话的控制
- 多窗口：每个会话都是独立运行的，拥有各自独立的输入输出终端窗口，终端窗口内显示过的信息也将被分开隔离保存，以便下次使用时依然能看到之前的操作记录
- 会话共享：当多个用户同时登录到远程服务器时，便可以使用会话共享功能让用户之间的输入输出信息共享

## 选项
- -S : 创建会话窗口
- -d : 将指定会话进行离线处理
- -r : 回复指定会话
- -x : 一次性恢复所有的会话
- -ls : 显示当前已有的会话
- -wipe : 把目前无法使用的会话删除

## example
```bash
# screen -S backup # 创建一个名称为 backup 的会话窗口
# screen -ls # 看到当前的会话
# exit # 退出当前会话
# screen -r linux # 恢复会话, 即关闭窗口不会关闭会话. 因此如果一段时间内不再使用某个会话窗口，可以把它设置为临时断开（detach）模式，随后在需要时再重新连接（attach）回来即可. 这段时间内，在会话窗口内运行的程序会继续执行
# ### user1和user2共享会话
# ssh 192.168.10.10 # user1
# screen -S linux
# ssh 192.168.10.10 # user2
# screen -x
```