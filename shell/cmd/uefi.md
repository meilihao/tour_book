# uefi
## efibootmgr
UEFI 规范定义了名为 UEFI 启动管理器的一项功能, Linux发行版包含名为efibootmgr 的工具，可用于更改 UEFI 启动管理器的配置, 可用`efibootmgr -v`查看启动项细节.

### example
```bash
# efibootmgr -v
BootCurrent: 0000 # 目前从“启动菜单”的哪个项上启动
Timeout: 0 seconds # 如果固件的 UEFI 启动管理器显示了类似启动菜单的界面，那么这一行表示继续启动默认项之前的超时
BootOrder: 0000,0003,0002,2001,2002,2003 # UEFI 固件将按照BootOrder 中列出的顺序，尝试从“启动菜单”中的每个“项”进行启动, 其余输出显示了实际的启动项
Boot0000* deepin	HD(1,GPT,8c465477-4444-4e2a-9306-6526f24cae36,0x800,0x100000)/File(\EFI\deepin\shimx64.efi)
Boot0002* Linpus lite	HD(1,GPT,8c465477-4444-4e2a-9306-6526f24cae36,0x800,0x100000)/File(\EFI\Boot\grubx64.efi)RC
Boot0003* ubuntu	HD(1,GPT,8c465477-4444-4e2a-9306-6526f24cae36,0x800,0x100000)/File(\EFI\ubuntu\grubx64.efi)RC
Boot2001* EFI USB Device	RC
Boot2002* EFI DVD/CDROM	RC
Boot2003* EFI Network	RC
```