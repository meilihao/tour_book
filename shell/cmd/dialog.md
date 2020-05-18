# dialog
shell下的对话框, 有点难用且不好调试.

## example
```bash
 (
            install_platform
        ) | dialog --trace /var/log/my.log --progressbox "Installing ${PKG_NAME}" 20 70 # 将install_platform的内容显示在dialog中, `--trace`将日志记入指定文件

 DIALOG_CANCEL=2 dialog --yes-label "Configure" --no-label "Restore" \ # 部分错误按键可能导致该dialog返回255, 比如`ESC`, 黑屏时通过按键来唤醒屏幕等.
           --yesno \
           "Select 'Configure' start new system configuration\nSelect 'Restore' for old system restoration" \
           10 60 

# cat /var/log/my.log # 记录了一个完整dialog的生命周期
** opened at Sat Apr  4 03:30:48 2020
** dialog 1.2-20130928
# Parameters:
# argv[0] = dialog
# argv[1] = --trace
# argv[2] = /var/log/config.log
# argv[3] = --progressbox
# argv[4] = Installing Federator Server # 标题
# argv[5] = 20
# argv[6] = 70
# discarding 2 parameters starting with argv[1] (--trace)
# init_result
# init_result
# process_common_options, offset 1
#	argv[1] = --progressbox
window 20x70 at 2,4
  0:+--------------------------------------------------------------------+
  1:| Installing Federator Server                                        |
  2:|--------------------------------------------------------------------|
  3:| Installing Federator Server packages ...                           |
  4:|                                                                    |
  5:|                                                                    |
  6:|                                                                    |
  7:|                                                                    |
  8:|                                                                    |
  9:|                                                                    |
 10:|                                                                    |
 11:|                                                                    |
 12:|                                                                    |
 13:|                                                                    |
 14:|                                                                    |
 15:|                                                                    |
 16:|                                                                    |
 17:|                                                                    |
 18:|                                                                    |
 19:+--------------------------------------------------------------------+
# widget returns 0 # 返回的状态
** closed at Sat Apr  4 03:34:16 2020
```

## terminal ui
参考:
- [如何开发富文本的终端UI应用](https://zhuanlan.zhihu.com/p/80356560)

terminal ui是半伪需求, 在有网络的情况下可用`web ui + go`取代(**推荐**); 无网但有RJ45时, 插线互联再参照有网情况; 其他情况[jojomi/gonsole(不完善, 已停更)](https://github.com/jojomi/gonsole)