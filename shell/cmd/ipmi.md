# ipmi
ref:
- [云服务器带外管理-从IPMI到RedFish](https://blog.csdn.net/asmartkiller/article/details/106558952)

IPMI（Intelligent Platform Management Interface）即智能平台管理接口是使硬件管理具备“智能化”的新一代通用接口标准. 用户可以利用 IPMI 监视服务器的物理特征，如温度、电压、电扇工作状态、电源供应以及机箱入侵等.

Ipmi 最大的优势在于它是独立于BIOS 和 OS 的，由BMC(Baseboard Management Controller)芯片管理, 所以用户无论在开机还是关机的状态下，只要接通电源就可以实现对服务器的监控.

**RedFishs是ipmi的取代者, PMI也在2015年公布2.0 v1.1标准后, 不再更新, 被RedFish永久代替**

## ipmitool
ref:
- [ipmitool常用命令总结](https://www.cnblogs.com/HByang/p/16127044.html)

ipmitool 是一种可用在 linux 系统下的命令行方式的 ipmi 平台管理工具，它支持 ipmi 1.5 规范（最新的规范为 ipmi 2.0），通过它可以实现获取传感器的信息、显示系统日志内容、网络远程开关机等功能.

### example
```
# ipmitool sel get <event_id:7640>
# ipmitool sel list
# ipmitool fru print # 查看主板信息, 依赖mod: ipmi_devintf, ipmi_si
# ipmitool sensor list # 查看ipmi sensor, 不带属性名
# ipmitool sensor get "Chassis_Fan1" # 查看背板风扇1的sensor信息, 带属性名
```