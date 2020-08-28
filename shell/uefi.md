# uefi
## efibootmgr
1. 创建一个新的boot option

    efibootmgr -c

1. 修改boot 顺序

    efibootmgr -o X,Y,... : 指定标号为X的启动项顺序在Y之前

1. 启用/禁用boot option

    - efibootmgr -a -b X : 启用标号为X的启动项
    - efibootmgr -A -b X : 禁用标号为X的启动项
