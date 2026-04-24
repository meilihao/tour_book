# dmesg
显示开机信息

kernel会将开机信息保存在系统缓冲区中(ring buffer)以及`/var/log/dmesg`.

## 选项
- -c : 显示开机信息后,清空ring buffer.
- -n : 设置记录信息的level
- -s : 设置缓冲区的大小, 默认是8192

## FAQ
### dmesg中出现`tpm_tis 00:05: A TPM error (6) occurred attempting to read a pcr value`

可能会导致开机时卡住几秒钟.

tpm资料,[TPM安全芯片](http://baike.baidu.com/view/687208.htm),[Trusted Platform Module](https://wiki.archlinux.org/index.php/Trusted_Platform_Module)

在BIOS里把TPM芯片禁用(推荐,thinkpad t430:F1->Security->Security Chip);或者禁用tpm系统模块(`echo blacklist tpm_tis > /etc/modprobe.d/tpm_tis.conf`),经测试,dmesg中还是有该错误信息.

### `[drm:cpt_set_fifo_underrun_reporting] *ERROR* uncleared pch fifo underrun on pch transcoder A [drm:cpt_serr_int_handler] *ERROR* PCH transcoder A FIFO underrun`

导致开机卡住几秒钟,drm错误,待官方解决.