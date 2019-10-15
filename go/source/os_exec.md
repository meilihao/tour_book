# [使用os/exec执行命令](https://blog.kowalczyk.info/article/wOYk/advanced-command-execution-in-go-with-osexec.html)

## 1. 执行命令并获得输出结果
```go
func main() {
	cmd := exec.Command("ls", "-lah")
	out, err := cmd.CombinedOutput() // 获得组合在一起的stdout/stderr输出
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	fmt.Printf("combined out:\n%s\n", string(out))
}
```

## 2. 将stdout和stderr分别处理

```go
func main() {
	cmd := exec.Command("ls", "-lah")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
	fmt.Printf("out:\n%s\nerr:\n%s\n", outStr, errStr)
}
```
## 3. 命令执行过程中获得输出
命令执行时间长, 且要获取它的持续输出
```go
func execCmd(shell string, raw []string) (int, error) {
    cmd := exec.Command(shell, raw...)
    stdout, err := cmd.StdoutPipe()
    if err != nil {
        fmt.Println(err)
        return 0, nil
    }
    stderr, err := cmd.StderrPipe()
    if err != nil {
        fmt.Println(err)
        return 0, nil
    }

    if err := cmd.Start(); err != nil {
        fmt.Println(err)
        return 0, nil
    }

    s := bufio.NewScanner(io.MultiReader(stdout, stderr))
    for s.Scan() {
        text := s.Text()
        fmt.Println(text)
    }
    
    if err := cmd.Wait(); err != nil {
        fmt.Println(err)
    }
    return 0, nil
}
```

## 4. 改变执行程序的环境(environment)
```go
cmd := exec.Command("programToExecute")
cmd.Env = append(os.Environ(), "FOO=bar"))
out, err := cmd.CombinedOutput()
if err != nil {
	log.Fatalf("cmd.Run() failed with %s\n", err)
}
fmt.Printf("%s", out)
```

## 5. 预先检查程序是否存在
```go
func checkLsExists() {
	path, err := exec.LookPath("ls")
	if err != nil {
		fmt.Printf("didn't find 'ls' executable\n")
	} else {
		fmt.Printf("'ls' executable is in '%s'\n", path)
	}
}
```

6. 管道

go style:
```go
package main
import (
    "os"
    "os/exec"
)
func main() {
    c1 := exec.Command("ls")
    c2 := exec.Command("wc", "-l")
    c2.Stdin, _ = c1.StdoutPipe()
    c2.Stdout = os.Stdout
    _ = c2.Start()
    _ = c1.Run()
    _ = c2.Wait()
}
```

Trick(伪装)方式:
```go
func main() {
	cmd := "cat /proc/cpuinfo | egrep '^model name' | uniq | awk '{print substr($0, index($0,$4))}'"
	out, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		fmt.Printf("Failed to execute command: %s", cmd)
	}
	fmt.Println(string(out))
}
```