# nm
显示关于elf文件中符号的信息

## example
```bash
# file a.out # 检查是否存在符号表
# nm -n a.out |grep main # 查找main的开始地址
# objdump -d a.out --start-address=0x1144 # 通过nm获取的main的开始地址反汇编main()
```