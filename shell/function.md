# function

定义:
```bash
[function] fname（）    
{
    statements;    
    [return]    
}
```

参数传递方式为："fname"（不需要传递参数）或"fname agr1 arg2"（需要传递两个参数）

> 函数必须先定义后使用
> 如果重新定义了函数，新的函数就会覆盖旧的函数

```bash
prompt()
{
    eval $3=\"$2\"
    if [ "${OS}" = "Linux" ]
    then
        /bin/echo -e "$1 [$2]: \c"
    else
        /bin/echo "$1 [$2]: \c"
    fi
    read tmp
    if [ ! -t 1 ] ; then echo $tmp ; fi
    if [ -n "$tmp" ] ; then eval $3=\"$tmp\" ; fi
}

ask_str()
{
    default=`eval echo '$'$2`
    prompt "$1" "`echo $default`" answer # answer作为变量, 允许`prompt()`(子函数)修改
    echo "---" $answer
    eval $2=\"$answer\"
    write_str $2
}
```

## 导出函数

    export -f fname # 可使函数的作用域扩展到子进程中.

## read
read: 从标准输入读取单行数据. 当使用重定向的时候，可以读取文件中的一行数据

```bash
while read -r line; do # read从stdin读取一行放入line
    echo "---" $line
done
```
