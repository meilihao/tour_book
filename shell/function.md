# function

定义:
```bash
[function] fname（)
{
    statements;
    [return]
}
```

参数传递方式为："fname"（不需要传递参数）或"fname agr1 arg2"（需要传递两个参数）

> 函数必须先定义后使用
> 如果重新定义了函数，新的函数就会覆盖旧的函数
> 在函数中使用`$<N>`实际是使用了函数的入参而不是调用脚本时的入参, 此时可将脚本的入参赋值给一个变量, 再在函数中使用该变量.

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

init_is_upstart()
{
   if [ -x /sbin/initctl ] && /sbin/initctl version 2>/dev/null | /bin/grep -q upstart; then
       return 0
   fi
   return 1
}

if init_is_upstart; then
     echo "----upstart" # output:"----upstart", 因为函数init_is_upstart返回0, 表示函数执行成功, 因此if的判断条件等价于`if true`
fi
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

## printf
```bash
# printf "%-5s %-10s %-4.2f\n" 3 Jeff 77.564 # 字符串默认是右对齐
```

## bash fork炸弹
```
# :() { :|:& };:
```

等价于:
```
:()
{
    :|:& # 第一个":"是调用":"函数本身; "|"使用管道会fork一个新进程来执行并将左侧的输出作为右侧的输入; 第二个":"+"&"是在后台执行函数":". 因此此时一次调用会有2个进程同时执行.
};
:
```

等价于:
```
bomb()
{
    bomb|bomb&
};
bomb
```

函数的名称为`:`,主要的核心代码是`：|：&`，可以看出这是一个函数本身的递归调用，通过`&`实现在后台开启新进程运行，通过管道实现进程呈几何形式增长.


### c实现
```c
#include <unistd.h>

int main()
{
  while(1)
    fork();
  return 0;
}

```