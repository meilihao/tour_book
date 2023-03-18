# nsis
ref:
- [NSIS脚本详解](https://www.cnblogs.com/visoeclipse/archive/2010/03/29/1699908.html)

## 标签
- SubSection: 选择安装组件的root

	- Section : SubSection的子项
- InstType: 把一个安装类型添加到安装类型列表里

	```conf
	InstType "推荐安装"
	InstType "全部安装"

	Section "区段 1"
	SectionIn 2 
	SectionEnd

	Section "区段 2"
	SectionIn 1 2
	SectionEnd
	```

	当选择"推荐安装"时, 只安装了"区段2"中的内容,当选择"全部安装"时,安装了"区段1"及"区段2"的内容

	`Section /o Sourcecode SEC_SOURCE`解读:
	- Sourcecode: 名称
	- SEC_SOURCE:选中标志
	- `/o`: 该区段默认为不选

	如果Section和名称之间什么也没有表示该section选中

	如果Section全部不选中且没有Section不可操作, 该SubSection不选中.
	如果至少有一个Section不可操作, 那么无论其他Section是否有选中, 该SubSection都是灰色的勾.
	全部Section选中, 该SubSection都是绿色的勾.
	部分Section选中, 该SubSection都是灰色的勾.

	`InstType /NOCUSTOM`: [禁用自定义](https://nsis.sourceforge.io/Reference/InstType)
- SectionIn

	管理到InstType n1[,n2...], n是InstType出现的顺序
- SectionSetFlags: 给Section设置flags

	- `SectionSetFlags ${SEC_SOURCE} 17 # SF_SELECTED & SF_RO`
- `Function .onSelChange`: 设置选中组件时的联动