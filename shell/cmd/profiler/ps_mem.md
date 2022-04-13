# ps_mem
一个用于精确报告 Linux 各用户态应用内存用量的简单 Python 脚本.

它会分别计算一个程序私有内存总量和共享内存总量，并以更准确的方式给出了总的内存使用量.

安装方法:
```bash
dnf install ps_mem
pip install ps_mem
# ps_mem
# ps_mem [-s] # 打印出应用的全路径
# ps_mem -p 2886,4386 # 指定pid
# ps_mem w 2 # 每 2 秒报告一次内存使用情况
```