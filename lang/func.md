# func
[在线调试工具](https://c.runoob.com/compile/11)

形参出现在函数定义中，在整个函数体内都可以使用， 离开该函数则不能使用.
实参出现在主调函数中，进入被调函数后，实参变量也不能使用.

## fgets
```c
char *fgets(char *str, int n, FILE *stream)
```

从指定的流 stream 读取一行，并把它存储在 str 所指向的字符串内. 当读取 (n-1) 个字符时，或者读取到换行符时，或者到达文件末尾时，它会停止，具体视情况而定.

### 参数:
- str : 指向一个字符数组的指针，该数组存储了要读取的字符串
- n : 要读取的最大字符数（**包括最后的空字符**）, 通常是使用以 str 传递的数组长度
- stream  : 指向 FILE 对象的指针，该 FILE 对象标识了要从中读取字符的流

### 返回值
如果成功，该函数返回相同的 str 参数. 如果到达文件末尾或者没有读取到任何字符，str 的内容保持不变，并返回一个空指针.
如果发生错误，返回一个空指针

## strcmp
按ascii逐个比较两个字符串中的各个字符, 直至出现不同的字符或遇到`\0`.

比较结果:
- 相等 : 0
- str1 > str2 : 正整数
- str1 < str2 : 负整数

## strncat
```c
char *strncat(char *dest, const char *src, size_t n) // 追加字符串
```
把 src 所指向的字符串追加到 dest 所指向的字符串的结尾，直到 n 字符长度为止.

> strncat()会将dest字符串最后的'\0'覆盖掉，字符追加完成后，再追加'\0'.

### 参数
- dest : 指向目标数组，该数组包含了一个 C 字符串，且足够容纳追加后的字符串，包括额外的空字符.
- src : 要追加的字符串
- n : 要追加的最大字符数. 如果n大于字符串src的长度，那么仅将src全部追加到dest的尾部.

### 返回值
返回一个指向最终的目标字符串 dest 的指针

## strncpy
```c
char *strncpy(char *dest, const char *src, size_t n)
```
把 src 所指向的字符串复制到 dest，最多复制 n 个字符. 当 src 的长度小于 n 时，dest 的剩余部分将用空字节填充; 如果strlen(src)的值大于或等于len，那么只有len个字符被复制到dst中, 但此时它的结果将不会以`\0`结尾.

> strncpy复制时(dest和n合适的话)会包括src的`\0`.

### 参数
- dest : 指向用于存储复制内容的目标数组
- src : 要复制的字符串
- n : 要从源中复制的字符数

###  返回值
返回最终复制的字符串

## memset
```c
void *memset(void *str, int c, size_t n)
```
复制字符 c（一个无符号字符）到参数 str 所指向的字符串的前 n 个字符.

### 参数
- 要被填充的内存块的首地址
- 要被设置成的值
- 要被设置的内存大小, 单位是字节

### 返回值
返回一个指向存储区 str 的指针

## malloc
```c
void *malloc(size_t size)
```

分配所需的内存空间，并返回一个指向它的指针

### 参数
- size : 内存块的大小，以字节为单位

### 返回值
返回一个指针 ，指向已分配大小的内存. 如果请求失败，则返回 NULL.