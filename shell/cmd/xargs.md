# xargs

## 描述

参数传递,即将标准输入数据转化成命令行参数(将stdin接收到的数据重新格式化，再将其作为参数传给其他命令),其应该紧跟在管道操作符后面.
**xargs的默认命令是`echo`，空格是默认定界符.** 这意味着通过管道传递给xargs的输入将会包含换行和空白，不过通过xargs的处理，换行和空白将被空格取代.
xargs是构建单行命令的重要组件之一.


## 选项

- -0 : 将`\0`作为定界符
- -d : 自定义一个定界符
- -I : 指定一个占位字符串{}，这个字符串在xargs扩展时会被替换掉，当-I与xargs结合使用，每一个参数命令都会被执行一次
- -n : 设置多行输出,每行最多n个数据(列)

## 例

    # cat  example.txt | xargs  # 多行输入转化成单行输出(用空格替换掉`\n`)
    # cat examplet.txt | xargs -n 3 # 单行输入转换成多行输出，每行n个参数
    # echo "nameXnameXnameXname" | xargs -d X
    # cat args.txt | xargs -n 2 ./cechi.sh # 每次最多传n个参数给脚本
    # cat args.txt | xargs -I {} ./cecho.sh -p {} 1
    # find . -type f -name "*.txt" -print0 | xargs -0 rm -f # find与xargs组合,使用`-0`以避免文件路径包含空格
    # find . -type f -name "*.go" -print0|xargs -0 wc -l # 统计代码行数
