# function

定义:

    [function] fname（）    
    {
    	statements;    
    	[return]    
    }

参数传递方式为："fname"（不需要传递参数）或"fname agr1 arg2"（需要传递两个参数）；

## 导出函数

    export -f fname # 可使函数的作用域扩展到子进程中.
