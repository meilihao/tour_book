# chroot
chroot，即 change root directory (更改 root 目录).

在 linux 系统中，系统默认的目录结构都是以 /，即以根 (root) 开始的; 而在使用 chroot 之后，系统的目录结构将以指定的位置作为`/`.

## example
```
$ mkdir rootfs
$ docker export $(docker create alpine) | tar -C rootfs -xvf -
$ sudo chroot rootfs /bin/ls -l # chroot 在 /sbin下
total 68
drwxr-xr-x    2 1000     1000          4096 May 29 14:20 bin
...
drwxr-xr-x   12 1000     1000          4096 May 29 14:20 var
```