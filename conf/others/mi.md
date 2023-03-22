# mi
## FAQ
### 查看电池损耗
1. `*#*#6485#*#*`

	```log
	# --- 小米9+miui12没有MF_02, MF_05, MF_06
	MB_06：电池的健康情况（Good就是正常 ）
	MF_02：电池的循环周期
	MF_05：目前电池实际容量
	MF_06：出厂的电池容量. 也可见`设置-我的设备-电池容量`
	```

	MF_05/MF_06<80%, 就可以考虑更换了, 这样的电池不耐用续航会大打折扣.
1. `*#*#284#*#*`, 让手机会生成Bug检测报告`bugreport-2023-01-04-151244.zip`

	获取步骤, 推荐在pc上解压:
	1. 解压bugreport-2023-01-04-151244.zip
	1. 解压bugreport-cepheus-RKQ1.200826.002-2023-01-04-15-12-44.zip. 这里用miui自带文件管理的解压功能会报错
	1. 查找bugreport-cepheus-RKQ1.200826.002-2023-01-04-15-12-44.txt中的`healthd`, 得到`healthd: battery l=70 v=4041 t=30.0 h=2 st=2 c=-718 fc=3259000 cc=578 tl=0 ct       =USB_HVDCP chg=a`

		字段说明:
		- l=70 标识当前电量剩余58%
		- cc=578 表示手机充电循环578次
		- fc =3259000 表示这块电池剩余容量还有3259mAh