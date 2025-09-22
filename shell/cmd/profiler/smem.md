# [smem](https://www.selenic.com/smem/)
统计物理内存使用情况, 包括共享内存.

`sudo apt install smem`

smem可以报告比例集大小(PSS)，唯一集大小(USS)和居民集大小(RSS):
- 比例集大小(PSS, proportion set size)：表示USS + 均摊到的共享内存大小.

    **所有使用某共享库的程序均分该共享库占用的内存时,显然所有进程的PSS之和就是系统的内存的使用量**

- 唯一集大小(USS, unique set size)：进程独自占用内存, 不包含任何共享的部分
- 驻留集大小(RSS, resident set size)：进程所使用的非交换区的物理内存. 将各个进程中的RSS值相加后,一般都会超出整个系统的内存消耗,这是因为RSS中包含了各个进程之间的共享内存.
- VSS : virtual set size (total virtual memory mapped)

USS和PSS仅包括物理内存使用情况, 它们不包括已换出到磁盘的内存.

## 选项
- -c 开关指定要显示的列, 比如只对 pss 列感兴趣
- -P 开关过滤进程，比如只包括那些名字里有 firefox 的进程
- -k 开关显示以 MB/GB 为单位的内存使用情况，而不是单纯的字节数
- -m 按映射统计
- -p 以百分比显示
- -P 进程过滤器
- -r 反向排序（从大到小）
- -s 按指定列排序, 比如pss, uss, rss
- -t 显示总数
- -w 按系统地址范围(system wide)统计
- -u 按用户统计
- --- # 输出图形
- --bar=BAR     Show bar graph.
- --pie=PIE     Show pie graph.
- --sort=pss: 按pss排序

## example
```bash
# smem -k # 默认按pss递增
# smem -t -k -s pss -r # 按pss从大到小并汇总大小, 常用
# sudo smem -r -k | sort -k 8 -nr | head -n 20 # 按 PSS 排名
# sudo smem -r -k | sort -k 6 -nr | head -n 20 # 按 USS 排名
# sudo smem -r -k | sort -k 5 -nr | head -n 20 # 按 RSS 排名
# smem -c "name pss uss rss" -s pss -r -k | head -20 # 按进程名显示内存占用排名
# smem -uk # 统计每个用户的内存占用
# smem -mk # 统计映射的内存使用情况
# smem -c "name user pss rss" # 仅显示部分列
# smem -p # 按百分比形式输出内存使用情况
# smem -c pss -P firefox -k -t | tail -n 1 # 获得 Firefox 的总内存使用量
# smem --pie name -c pss # 显示总的内存使用情况并以图形输出
# smem -P root --bar pid -c "pss rss" # 展示每个root用户进程的pss rss
# echo 'smem -c pss -P "$1" -k -t | tail -n 1' > ~/bin/memory-use && chmod +x ~/bin/memory-use # 内存统计脚本, 用法`memory-use firefox`
# smem --bar name -c "pss uss rss" # 显示USS, PSS和RSS的条形图组合
# smem -twk # 按kernel/user space区分显示
# smem --userfilter="root" --bar pid -c "pss rss" #  root 用户的所有进程
```