# nm
显示关于elf文件中符号的信息

## 选项
- -A: 每个符号前显示文件名
- -D: 显示动态符号
- -g: 仅显示外部符号
- -r: 反序显示符号表

## example
```bash
# file a.out # 检查是否存在符号表
# nm -n a.out |grep main # 查找main的开始地址
# objdump -d a.out --start-address=0x1144 # 通过nm获取的main的开始地址反汇编main()
# nm a.out | grep -Ei 'function|main|globalvar'
```