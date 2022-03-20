# cobra

## FAQ
### [`cmd --reboot true`报`Error: unknown command "true" fro "cmd"`](https://blog.csdn.net/beijihukk/article/details/119753501)
这是因为cobra在创建不同类型flag的时候做了区别处理.

bool类型参数不需要跟一个值: 它本身就代表true, 不显示添加参数则表示false; 而string类型参数则需要跟一个值.
