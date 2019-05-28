# docker实现
参考:
- [Docker背后的内核知识：命名空间资源隔离](https://linux.cn/article-5057-1.html)
- [容器核心技术详解](https://blog.fliaping.com/container-core-technical-details/)

Docker 容器实际上是在创建容器进程时，指定了这个进程所需要启用的一组 Namespace和Cgroups 参数. 所以说，容器其实是一种追加了特定参数的进程而已.

对 Docker 来说，它最核心的原理就是为待创建的用户进程：
1. 启用 Linux Namespace 配置
1. 设置指定的 Cgroups 参数
1. 切换进程的根目录（Change Root）: 优先使用 pivot_root 系统调用，如果系统不支持，才会使用 chroot

## [Namespace](https://lwn.net/Articles/531114/)
Namespace包括PID、Mount、UTS、IPC、Network 和 User.

它主要涉及3个系统调用:
- clone() : 实现线程的系统调用，用来创建一个新的进程，并可以通过Namespace flag达到隔离
- unshare() : 使某进程脱离某个namespace
- setns() : 把某进程加入到某个namespace

### PID
让被隔离进程只看到当前 PID Namespace 里的进程

linux创建进程的system call:
```c
int pid = clone(main_function, stack_size, SIGCHLD, NULL); 
```

使用 PID Namespace即指定 CLONE_NEWPID 参数:
```c
int pid = clone(main_function, stack_size, CLONE_NEWPID | SIGCHLD, NULL);
```

### Mount
让被隔离进程只看到当前 Namespace 里的挂载点信息. Mount Namespace 修改的是容器进程对文件系统“挂载点”的认知. 这就意味着: 只有在“挂载”这个操作发生之后，进程的视图才会被改变, 而在此之前，新创建的容器会直接继承宿主机的各个挂载点. 因此Mount Namespace 跟其他 Namespace 的使用略有不同：它对容器进程视图的改变，一定是伴随着挂载操作（mount）才能生效.

flag : CLONE_NEWNS

每当创建一个新容器时，通过结合使用 Mount Namespace 和 rootfs，容器就能够为进程构建出一个完善的文件系统隔离环境，而不是继承自宿主机的文件系统, 这就用到了pivot_root/chroot(切换进程根目录). 而这个挂载在容器根目录上、用来为容器进程提供隔离后执行环境的文件系统，就是所谓的“容器镜像”, 它还有一个更为专业的名字，叫作：rootfs（根文件系统）.

> 实际上，Mount Namespace 正是基于对 chroot 的不断改良才被发明出来的，它也是 Linux 操作系统里的第一个 Namespace.
> rootfs 只是一个操作系统所包含的文件、配置和目录，**并不包括操作系统内核**. 在 Linux 操作系统中，这两部分是分开存放的，操作系统只有在开机启动时才会加载指定版本的内核镜像. 同一台机器上的所有容器都共享宿主机操作系统的内核. 这也是容器相比于虚拟机的主要缺陷之一：毕竟后者不仅有模拟出来的硬件机器充当沙盒，而且每个沙盒里还运行着一个完整的 Guest OS 给应用随便折腾.

### Network
让被隔离进程看到当前 Namespace 里的网络设备和配置

## Linux Cgroups (Linux Control Group)
限制一个进程组能够使用的资源上限，包括 CPU、内存、磁盘、网络带宽等等. 此外，Cgroups 还能够对进程进行优先级设置、审计，以及将进程挂起和恢复等操作.

在 Linux 中，Cgroups 给用户暴露出来的操作接口是文件系统，即它以文件和目录的方式组织在操作系统的 /sys/fs/cgroup 路径下, 因此它就是**一个子系统目录加上一组资源限制文件**的组合.
可通过`mount -t cgroup `查看, 查看输出就会发现在 /sys/fs/cgroup 下面有很多诸如 cpuset、cpu、 memory 这样的子目录，也叫子系统. 再通过`ls /sys/fs/cgroup/${子系统}`即可查看更具体的限制.

在对应的子系统下面创建一个目录, 这个目录就称为一个“控制组, 操作系统会在新创建的目录下，自动生成该子系统对应的资源限制文件.
控制组分类:
- cpu : 为cpu使用限制
- blkio : 为块设备设置I/O限制, 一般用于磁盘等设备
- cpuset : 为进程分配单独的 CPU 核和对应的内存节点
- memory : 为进程设定内存使用的限制

而对于 Docker 等 Linux 容器项目来说，它们只需要在每个子系统下面，为每个容器创建一个控制组（即创建一个新目录），然后在启动容器进程之后，把这个进程的 PID 填写到对应控制组的 tasks 文件中就可以, 比如:
```bash
$ docker run -it --cpu-quota=20000 busybox /bin/sh # 进/sys/fs/cgroup/cpu/docker/fdcef525207f20f7eb4e8d2def11fe5cc6ebe5839ed25959d3dab2a7ddc56bb9查看cpu.cfs_quota_us即可
```

演示:
1. 进入 /sys/fs/cgroup/cpu, 创建控制组container
```bash
# cat container/cpu.cfs_quota_us
-1 # -1即没有限制
# cat container/cpu.cfs_period_us
100000 # 默认的 100 ms（100000 us）
```
1. 创建一个死循环的bash任务`while : ; do : ; done &`, 通过top查看到它的cpu是100%
1. 写入限制, 查看结果
```bash
# echo 20000 > container/cpu.cfs_quota_us # 在每 100 ms 的时间里，被该控制组限制的进程只能使用 20 ms 的 CPU 时间，也就是说这个进程只能使用到 20% 的 CPU
# echo ${死循环job's pid} > container/tasks # 把被限制的进程的 PID 写入 container/tasks, 上面的设置就会对该进程生效
# top # 会发现死循环任务的cpu一直在20%左右波动了
```

> fs_period 和 cfs_quota 需要组合使用就可以用来限制进程在长度为 cfs_period 的一段时间内，只能被分配到总量为 cfs_quota 的 CPU 时间.
> Cgroups 对资源的限制能力也有很多不完善的地方，被提及最多的自然是 /proc 文件系统的问题. 比如应用程序在容器里读取到的 CPU 核数、可用内存等信息都是宿主机上的数据, 目前可通过lxcfs解决.

## 镜像
容器镜像的层（layer）: 在 rootfs 的基础上，Docker 公司创新性地提出了使用多个增量 rootfs 联合挂载一个完整 rootfs 的方案. 即用户制作镜像的每一步操作，都会生成一个层，也就是一个增量 rootfs, 以避免每开发一个应用，或者升级一下现有的应用，都要重复制作一次 rootfs.

镜像的实现使用了Union File System 也叫 UnionFS，其最主要的功能是将多个不同位置的目录联合挂载（union mount）到同一个目录下.
> deepin 15.10 + docker 18.09.6上镜像默认使用的是 [overlay2](https://docs.docker.com/v17.09/engine/userguide/storagedriver/overlayfs-driver/) 这个联合文件系统. 可通过`docker image inspect ${镜像标识}`或`docker inspect ${容器标识}`的`GraphDriver`查看镜像使用的UnionFS. 对于 overlay2 来说，它最关键的目录结构在 `/var/lib/docker/overlay2` 里
> 以前ubuntu+docker默认使用[AUFS](https://docs.docker.com/v17.09/engine/userguide/storagedriver/aufs-driver/), **不推荐**, 因为AUFS未进入kernel, 而overlay2已进入. 现在最新的 docker 版本中默认ubuntu/centos都是使用的 overlayfs.
> 由于容器镜像的操作是增量式的，每次镜像拉取、推送的内容，比原本多个完整的操作系统的大小要小得多；而共享层的存在，可以使得所有这些容器镜像需要的总空间，也比每个镜像的总和要小. 这样就使得基于容器镜像的团队协作更敏捷. 因此容器镜像将会成为未来软件的主流发布方式.

演示overlay2, 它采用了两层结构，lowerdir为镜像层，只读. upperdir为容器层,可读写:
```bash
$ tree
.
├── low
│   ├── a
│   └── c
├── merged
├── upper
│   ├── b
│   └── c
└── work
$ sudo mount -t overlay  -o lowerdir=./low,upperdir=./upper,workdir=./work overlay ./merged
$ tree
.
├── low # image layer
│   ├── a
│   └── c
├── merged # container mount
│   ├── a
│   ├── b
│   └── c
├── upper # container layer
│   ├── b
│   └── c
└── work # overlayfs实现要用的一个目录, 功能未知
    └── work [error opening dir] # 权限问题, 使用sudo tree避免
$ echo 1 > merged/a # 对merged的修改会映射到upper上
$ cat upper/a
1
$ cat low/a # 没有内容输出
$ echo 3 > upper/c
$ cat merged/c
3
$ sudo tree
.
├── low
│   ├── a
│   └── c
├── merged
│   ├── a
│   ├── b
│   └── c
├── upper # echo 1 > merged/a 体现在upper上
│   ├── a
│   ├── b
│   └── c
└── work
    └── work
```

演示aufs:
```bash
$ tree
.
├── dir1
│   ├── a
│   └── c
├── dir2
│   ├── b
│   └── c
└── dir3

3 directories, 4 files
$ sudo mount -t aufs -o dirs=./dir1:./dir2 none ./dir3
[sudo] chen 的密码：
$ tree
.
├── dir1
│   ├── a
│   └── c
├── dir2
│   ├── b
│   └── c
└── dir3
    ├── a
    ├── b
    └── c
$ mount| grep aufs
none on /home/chen/tmpfs/test2/dir3 type aufs (rw,relatime,si=a2434e0fd58d1a90) # 通过si(AUFS内部 ID)可在/sys/fs/aufs下查看被联合挂载在一起的各个层的信息
$ cat /sys/fs/aufs/si_a2434e0fd58d1a90/br[0-9]*
/home/chen/tmpfs/test2/dir1=rw
/home/chen/tmpfs/test2/dir2=ro
```
与overlay2类似, 对dir3的修改会体现在可读写层的dir1上.
aufs实现删除的方式: 在aufs的rw层创建一个 whiteout 文件(.wh.${deleted_filename})，把只读层里的文件“遮挡”起来.

容器的rw layer可通过docker commit 和 push 指令保存并上传到镜像仓库供其他人使用；而与此同时，原先的只读层里的内容则不会有任何变化. 这就是增量 rootfs 的好处.

既然容器的 rootfs 是以只读方式挂载的，那要如何在容器里修改镜像的内容, 答案是Copy-on-Write: 所有的增删查改操作都只会作用在容器层，相同的文件上层会覆盖掉下层. 因此要修改一个文件的时候，首先会从上到下查找有没有这个文件，找到后就复制到容器层中再进行修改，修改的结果就会作用到下层的文件; 没有就创建一个即可.

实际容器演示:
```bash
$ docker inspect fdcef525207f # 获取容器实际rootfs
[
    ...
        "GraphDriver": {
            "Data": {
                "LowerDir": "/var/lib/docker/overlay2/e372efb2eafc1454e515286b3e16dd18a15b5587025c1f61d6fb6fadd15d8e32-init/diff:/var/lib/docker/overlay2/99d1e99f0f49da75a7f795355c92969e90b74f3d00222176bfb999e936b9d5da/diff",
                "MergedDir": "/var/lib/docker/overlay2/e372efb2eafc1454e515286b3e16dd18a15b5587025c1f61d6fb6fadd15d8e32/merged",
                "UpperDir": "/var/lib/docker/overlay2/e372efb2eafc1454e515286b3e16dd18a15b5587025c1f61d6fb6fadd15d8e32/diff",
                "WorkDir": "/var/lib/docker/overlay2/e372efb2eafc1454e515286b3e16dd18a15b5587025c1f61d6fb6fadd15d8e32/work"
            },
            "Name": "overlay2"
        },
        ...
]
```

`/var/lib/docker/overlay2/e372efb2eafc1454e515286b3e16dd18a15b5587025c1f61d6fb6fadd15d8e32-init`即容器的init layer(只读).
init 层是 Docker 项目单独生成的一个内部层，专门用来存放 /etc/hosts、/etc/resolv.conf 等信息,原因:
这些文件本来属于只读镜像的一部分，但是用户往往需要在启动容器时写入一些指定的值比如 hostname，所以就需要在可读写层对它们进行修改. 可是这些修改往往只对当前的容器有效，我们并不希望执行 docker commit 时，把这些信息连同可读写层一起提交掉. 所以，Docker 做法是在修改了这些文件之后，以一个单独的层挂载了出来, 而用户执行 docker commit 只会提交可读写层，所以是不包含init layer内容的.

也可通过`cat /proc/mounts| grep overlay`或`mount| grep overlay`可查看容器在系统中的挂载信息.

参考:
- [overlayfs技术探究以及docker的使用](https://www.jianshu.com/p/959e8e3da4b2)

## docker exec原理
```bash
$ docker inspect --format '{{ .State.Pid }}'  fdcef525207f # 获取容器在宿主机的pid
8432
$ sudo ls -l /proc/8432/ns # 看到这个 8432 进程的所有 Namespace 对应的文件
总用量 0
lrwxrwxrwx 1 root root 0 5月  27 22:45 cgroup -> cgroup:[4026531835]
lrwxrwxrwx 1 root root 0 5月  27 22:45 ipc -> ipc:[4026532547]
lrwxrwxrwx 1 root root 0 5月  27 22:45 mnt -> mnt:[4026532545]
lrwxrwxrwx 1 root root 0 5月  27 22:40 net -> net:[4026532550]
lrwxrwxrwx 1 root root 0 5月  27 22:45 pid -> pid:[4026532548]
lrwxrwxrwx 1 root root 0 5月  27 22:45 pid_for_children -> pid:[4026532548]
lrwxrwxrwx 1 root root 0 5月  27 22:45 user -> user:[4026531837]
lrwxrwxrwx 1 root root 0 5月  27 22:45 uts -> uts:[4026532546]
```

可以看到: 一个进程的每种 Linux Namespace，都在它对应的 /proc/[进程号]/ns 下有一个对应的虚拟文件，并且链接到一个真实的 Namespace 文件上. 有了这些信息我们就可以做一些有趣的事情了，比如加入到一个已经存在的 Namespace 当中.

这也就意味着：一个进程可以选择加入到某个进程已有的 Namespace 当中，从而达到“进入”这个进程所在容器的目的，这正是 docker exec 的实现原理, 而这个操作所依赖系统调用setns().

演示demo:
```c
#define _GNU_SOURCE
#include <fcntl.h>
#include <sched.h>
#include <unistd.h>
#include <stdlib.h>
#include <stdio.h>
#include <string.h>
#include <errno.h>
 
// 它需要两个参数: argv[1]是当前进程要加入的 Namespace 文件的路径, 而argv[1]是将要在这个 Namespace 里运行的进程
// 这段代码的的核心操作是通过 open() 系统调用打开了指定的 Namespace 文件，并把这个文件的描述符 fd 交给 setns() 使用. 
// 在 setns() 执行后，当前进程就加入了这个文件对应的 Linux Namespace 中.
int main(int argc, char *argv[]) {
    int fd;
    
    fprintf(stdout, "argv1: %s, argv2: %s\n", argv[1],argv[2]);

    fd = open(argv[1], O_RDONLY);
    if (setns(fd, 0) == -1) {
        fprintf(stderr, "setns failed: %s\n", strerror(errno));
    	return -1;
    }
    
    if (execvp(argv[2], &argv[2]) != 0 ) {
    	fprintf(stderr, "failed to execvp argments %s\n", strerror(errno));
    	return -1;
  	}

  	printf("all done!\n");
  	return 0;
}
```

运行:
```bash
$ gcc -o setns setns.c 
$ sudo ./setns /proc/8432/ns/net /bin/bash
$ sudo ps -ef|grep /bin/bash
root       441 32530  0 23:32 pts/2    00:00:00 sudo ./setns /proc/8432/ns/net /bin/bash
root       442   441  0 23:32 pts/2    00:00:00 /bin/bash # 找到进程
$ sudo ls -al /proc/442/ns
总用量 0
dr-x--x--x 2 root root 0 5月  27 23:34 .
dr-xr-xr-x 9 root root 0 5月  27 23:34 ..
lrwxrwxrwx 1 root root 0 5月  27 23:34 cgroup -> cgroup:[4026531835]
lrwxrwxrwx 1 root root 0 5月  27 23:34 ipc -> ipc:[4026531839]
lrwxrwxrwx 1 root root 0 5月  27 23:34 mnt -> mnt:[4026531840]
lrwxrwxrwx 1 root root 0 5月  27 23:34 net -> net:[4026532550] # 与上面进程8432的ns比较, 发现一致, 因此它们的ifconfig结果也会一致
lrwxrwxrwx 1 root root 0 5月  27 23:34 pid -> pid:[4026531836]
lrwxrwxrwx 1 root root 0 5月  27 23:34 pid_for_children -> pid:[4026531836]
lrwxrwxrwx 1 root root 0 5月  27 23:34 user -> user:[4026531837]
lrwxrwxrwx 1 root root 0 5月  27 23:34 uts -> uts:[4026531838]
```

> 此外，Docker 还专门提供了`-net`参数，可以让我们启动一个容器时加入到另一个容器的 Network Namespace 里，比如`docker run -it --net container:${容器id} busybox ifconfig`
> 同时如果指定`–net=host`就意味着这个容器不会为进程启用 Network Namespace. 即这个容器拆除了 Network Namespace 的“隔离墙”，它会和宿主机上的其他普通进程一样，直接共享宿主机的网络栈. 这就为容器直接操作和使用宿主机网络提供了一个渠道.

## docker commit
实际上就是在容器运行起来后，把最上层的“可读写层”加上原先容器镜像的只读层，打包组成了一个新的镜像, 下面这些只读层在宿主机上是共享的，不会占用额外的空间.
而由于使用了联合文件系统，你在容器里对镜像 rootfs 所做的任何修改，都会被操作系统先复制到这个可读写层，然后再修改, 这就是所谓的`Copy-on-Write`.
而正如前所说init 层的存在，就是为了避免你执行 docker commit 时，把 Docker 自己对 /etc/hosts 等文件做的修改也一起提交掉.

## 扩展
### clone demo
```c
// from [DOCKER基础技术：LINUX NAMESPACE（上）](https://coolshell.cn/articles/17010.html)
#define _GNU_SOURCE
#include <sys/mount.h> 
#include <sys/types.h>
#include <sys/wait.h>
#include <stdio.h>
#include <sched.h>
#include <signal.h>
#include <unistd.h>
#include <errno.h>
#include <string.h> // for strerror
#include <stdlib.h> // for exit
#define STACK_SIZE (1024 * 1024) /* 定义一个给 clone 用的栈，栈大小1M */
static char container_stack[STACK_SIZE];
char* const container_args[] = {
  "/bin/bash",
  NULL
};

int container_main(void* arg)
{  
  printf("Container - inside the container!\n");
  if (execv(container_args[0], container_args) != 0 ) {
    fprintf(stderr, "failed to execvp argments %s\n", strerror(errno));
    exit(-1);
  }
  printf("Something's wrong!\n");
  return 1;
}

// $ gcc -o ns ns.c
// $ sudo ./ns
int main()
{
  printf("Parent - start a container!\n");
  int container_pid = clone(container_main, container_stack+STACK_SIZE, SIGCHLD , NULL);
  if (container_pid < 0) {
    fprintf(stderr, "clone failed: %s\n", strerror(errno));
    return -1;
  }
  if (waitpid(container_pid, NULL, 0) == -1) {
    fprintf(stderr, "failed to wait pid %d\n", container_pid);
    return -1;
  }
  printf("Parent - container stopped!\n");
  return 0;
}
```