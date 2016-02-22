# script

## 描述

录制终端会话,和scriptreplay连用.

## 选项

- -a ∶ 将录制时运行的命令及其输出内容追加到指定文件
- -f : 每次输入后刷新输出
- -t ∶ 将时序信息导入stderr或文件

## 例

    # script -t 2> timing.log -a output.session # `2>`将stderr重定向到timing.log
    cmd...
    exit # 表示结束录制

### 在多个用户间广播terminal session

假设有term1(主播),term2(听众,可以有多个)

1. term1

        # mkfifo scriptfifo

2. term2

        # cat scriptfifo

3. term1

        # script -f scriptfifo
        cmds...
        exit # 退出录制

此时,term1的操作将通过管道广播到term2.步骤3可以和步骤1合并,只是term2此时cat可能得到很多内容(因为term1可能有操作),容易混乱,不推荐.
