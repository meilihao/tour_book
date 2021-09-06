# containerd
参考:
- [作为k8s容器运行时，containerd跟docker的对比](https://cloud.tencent.com/developer/article/1450788)

crictl 是 CRI 兼容的容器运行时命令行接口, 可以使用它来检查和调试 Kubernetes 节点上的容器运行时和应用程序. crictl 和它的源代码在 cri-tools 代码库.

ctr(**不推荐**))是用于与containerd守护程序进行交互,用于不受支持的调试和管理功能的客户端. 因为不受支持，所以不能保证命令，选项和操作的向后兼容性.

[cri-containerd-cni-<version>-linux-amd64.tar.gz](https://github.com/containerd/containerd/releases)包括:
- crictl
- runc
- cni-bin

## example
```bash
# containerd config default # 获取默认配置
# containerd config dump # 当前配置(仅是从配置文件中读取并展示配置, 让containerd使用它需要restart containerd)
```

## 开发
参考:
- [Getting started with containerd](https://containerd.io/docs/getting-started/)

```go
package main

import (
	"context"
	"fmt"
	"log"
	"syscall"
	"time"

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/cio"
	"github.com/containerd/containerd/oci"
	"github.com/containerd/containerd/namespaces"
)

func main() {
	if err := redisExample(); err != nil {
		log.Fatal(err)
	}
}

func redisExample() error {
	// create a new client connected to the default socket path for containerd
	client, err := containerd.New("/run/containerd/containerd.sock")
	if err != nil {
		return err
	}
	defer client.Close()

	// create a new context with an "example" namespace
	ctx := namespaces.WithNamespace(context.Background(), "example")

	// pull the redis image from DockerHub
	image, err := client.Pull(ctx, "docker.io/library/redis:alpine", containerd.WithPullUnpack)
	if err != nil {
		return err
	}

	// create a container
	container, err := client.NewContainer(
		ctx,
		"redis-server",
		containerd.WithImage(image),
		containerd.WithNewSnapshot("redis-server-snapshot", image),
		containerd.WithNewSpec(oci.WithImageConfig(image)),
	)
	if err != nil {
		return err
	}
	defer container.Delete(ctx, containerd.WithSnapshotCleanup)

	// create a task from the container
	task, err := container.NewTask(ctx, cio.NewCreator(cio.WithStdio))
	if err != nil {
		return err
	}
	defer task.Delete(ctx)

	// make sure we wait before calling start
	exitStatusC, err := task.Wait(ctx)
	if err != nil {
		fmt.Println(err)
	}

	// call start on the task to execute the redis server
	if err := task.Start(ctx); err != nil {
		return err
	}

	// sleep for a lil bit to see the logs
	time.Sleep(3 * time.Second)

	// kill the process and get the exit status
	if err := task.Kill(ctx, syscall.SIGTERM); err != nil {
		return err
	}

	// wait for the process to fully exit and print out the exit status

	status := <-exitStatusC
	code, _, err := status.Result()
	if err != nil {
		return err
	}
	fmt.Printf("redis-server exited with status: %d\n", code)

	return nil
}
```

## FAQ
### enable debug log
```
#  cat /etc/containerd/config.toml
[debug]
  level = "debug"
# systemctl -f -xeu containerd
```
