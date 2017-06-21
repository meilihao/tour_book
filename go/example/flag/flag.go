package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	//FlagSet代表一个已注册的flag的集合。FlagSet零值没有名字且采用ContinueOnError错误处理策略
	flagSet = flag.NewFlagSet("test", flag.ExitOnError)

	conf = flagSet.String("conf", "conf.ini", "conf's path")
	ver  = flagSet.Bool("ver", false, "version")
)

func main() {
	fmt.Println(os.Args)
	flagSet.Parse(os.Args[1:])

	if *ver {
		fmt.Println(*ver)
		return
	}
	fmt.Println(*conf)
}

/*
///#go run test.go -ver=true
[./test -ver=true]
true
*/

/*
flag包要求所有的flag参数都必须出现在尾部位置参数的前面,否则这些flag将会被当作位置参数处理.
*/
