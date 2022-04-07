# htop

## 输出列
- PROCESSOR : cpu id

## FAQ
### 检查一个进程/线程当前使用的是哪个 CPU
按F2进入设置, 选择Columns, 右移到`Available Columns`下选中 PROCESSOR按回车即可.

> htop/perf cup id从1开始, taskset, ps, top则均是从0开始.