# ansible
Ansible 目前是运维自动化工具中最简单、容易上手的一款优秀软件， 能够用来管理各种资源。用户可以使用 Ansible 自动部署应用程序， 以此实现 IT 基础架构的全面部署.

相较于 Chef、 Puppet、 SaltStack 等 C/S（客户端/服务器） 架构的自动化工具来讲，尽管 Ansible 的性能并不是最好的，但由于它基于 SSH 远程会话协议，不需要客户端程
序，只要知道受管主机的账号密码，就能直接用 SSH 协议进行远程控制，因此使用起来优势明显.

Ansible 服务本身并没有批量部署的功能，它仅仅是一个框架，真正具有批量部署能力的是其所运行的模块。 Ansible 内置了上千个模块，会在安装 Ansible 时一并安装，通过调用指定的
模块，就能实现特定的功能。 Ansible 内置的模块非常丰富，几乎可以满足一切需求， 使用起来也非常简单，一条命令甚至影响上千台主机。 如果需要更高级的功能，也可以运用 Python
语言对 Ansible 进行二次开发.

## 术语
- control node 控制节点 : 安装了 Ansible 服务的主机，也称为 Ansible 控制端，主要是用来发布运行任务、调用功能模块，以及对其他主机进行批量控制
- managed node 受控节点 : 被 Ansible 服务所管理的主机，也被称为受控主机或客户端，是模块命令的被执行对象
- inventory 主机清单 : 受控节点的列表，可以是 IP 地址、主机名或者域名
- module 模块 : 用于实现特定功能的代码； Ansiblie 默认带有上千款模块； 可以在 Ansible Galaxy 中选择更多的模块
- task 任务 : 要在 Ansible 客户端上执行的操作
- playbook 剧本 : 通过 YAML 语言编写的可重复执行的任务列表； 把重复性的操作写入到剧本文件中后，下次可直接调用剧本文件来执行这些操作

	剧本文件的结构由 4 部分组成， 分别是 target、 variable、 task、 handler， 其各自的作用如下:
	- target： 用于定义要执行剧本的主机范围
	- variable： 用于定义剧本执行时要用到的变量
	- task： 用于定义将在远程主机上执行的任务列表
	- handler： 用于定义执行完成后需要调用的后续任务

	YAML 语言编写的 Ansible 剧本文件会按照从上到下的顺序自动运行，其形式类似于Shell 脚本，但格式有严格的要求.

- role 角色 : 从 Ansible 1.2 版本开始引入的新特性，用于结构化地组织剧本, 类似于编程中的封装技术；通过调用角色可实现一连串的功能

	角色的好处就在于将剧本组织成了一个简洁的、可重复调用的抽象对象，使得用户把注意力放到剧本的宏观大局上，统筹各个关键性任务，只有在需要时才去深入了解细节

	角色的获取有 3 种方法:
	1. 加载系统内置角色
	1. 从外部环境获取角色
	1. 自行创建角色

	RHEL 系统自带的角色, 剧本模板文件存放在`/usr/share/doc/rhel-system-roles/`:
	- rhel-system-roles.kdump 配置 kdump 崩溃恢复服务
	- rhel-system-roles.network 配置网络接口
	- rhel-system-roles.selinux 配置 SELinux 策略及模式
	- rhel-system-roles.timesync 配置网络时间协议
	- rhel-system-roles.postfix 配置邮件传输服务
	- rhel-system-roles.firewall 配置防火墙服务
	- rhel-system-roles.tuned 配置系统调优选项

## 配置
Ansible 服务的主配置文件存在优先级的顺序关系，默认存放在/etc/ansible 目录中的主配置文件优先级最低, 优先级是:
- 高 ./ansible.cfg
- 中 ~/ansible.cfg
- 低 /etc/ansible/ansible.cfg

其他:
- /etc/ansible/hosts : 主机清单（ inventory）

Ansible 常用变量:
- ansible_ssh_host 受管主机名
- ansible_ssh_port 端口号
- ansible_ssh_user 默认账号
- ansible_ssh_pass 默认密码
- ansible_shell_type Shell 终端类型

## 模块
- ping 检查受管主机的网络是否能够连通(类似于常用的 ping 命令)
- yum 安装、更新及卸载软件包
- yum_repository 管理主机的软件仓库配置文件, 能够添加、修改及删除软件仓库的配置信息，参数相对比较复杂

	- ansible all -m yum_repository -a 'name="EX294_BASE" description="EX294 base software" baseurl="file:///media/cdrom/BaseOS" gpgcheck=yes enabled=1 gpgkey="file:///media/cdrom/RPM-GPG-KEY-redhat-release"' # add repo
- template 复制模板文件到受管主机
- copy 新建、修改及复制文件
- user 创建、修改及删除用户
- group 创建、修改及删除用户组
- service 启动、关闭及查看服务状态
- get_url 从网络中下载文件
- file 设置文件权限及创建快捷方式
- cron 添加、修改及删除计划任务
- command 直接执行用户指定的命令
- shell 直接执行用户指定的命令（支持特殊字符）
- debug 输出调试或报错信息
- mount 挂载硬盘设备文件
- filesystem 格式化硬盘设备文件
- lineinfile 通过正则表达式修改文件内容
- setup 收集受管主机上的系统及变量信息

	- `ansible all -m setup -a 'filter="*ip*"'` : filter指定过滤条件
- firewalld 添加、修改及删除防火墙策略
- lvg 管理主机的物理卷及卷组设备
- lvol 管理主机的逻辑卷设备

## role
Ansible Galaxy 是 Ansible 的一个官方社区，用于共享角色和功能代码，用户可以在网站自由地共享和下载 Ansible 角色.

在 Ansible 的主配置`/etc/ansible/ansible.cfg`中， 定义有角色保存路径。 如果用户新建的角色信息不在规定的目录内，则无法使用 ansible-galaxy list 命令找到。因此需要手动填写新角色的目录路径，或是进入/etc/ansible/roles 目录内再进行创建.

Ansible 角色的目录结构及含义
- defaults 包含角色变量的默认值（优先级低）
- files 包含角色执行任务时所引用的静态文件
- handlers 包含角色的处理程序定义
- meta 包含角色的作者、许可证、 平台和依赖关系等信息
- tasks 包含角色所执行的任务
- templates 包含角色任务所使用的 Jinja2 模板
- tests 包含用于测试角色的剧本文件
- vars 包含角色变量的默认值（优先级高）

```bash
# ansible-galaxy install 角色名称
# ansible-galaxy install -r nginx.yml # 创建role?
# ansible-galaxy list
# --- 部署apache
# cd /etc/ansible/roles
# ansible-galaxy init apache # 创建一个新的角色信息，且建立成功后便会在当前目录下生成出一个新的同名目录
# cd apache
# vim tasks/main.yml
---
- name: one
  yum:
	name: httpd
	state: latest
- name: two
  service:
	name: httpd
	state: started
	enabled: yes
- name: three
  firewalld:
	service: http
	permanent: yes
	state: enabled
	immediate: yes
- name: four
  template:
	src: index.html.j2
	dest: /var/www/html/index.html
# vim templates/index.html.j2
Welcome to {{ ansible_fqdn }} on {{ ansible_all_ipv4_addresses }}
# cd ~
# vim roles.yml # 用于调用 apache 角色的 yml文件
---
- name: 调用自建角色
  hosts: all
	roles:
		- apache
# ansible-playbook roles.yml # 执行完毕后，在浏览器中随机输入几台主机的 IP 地址，即可访问到包含主机 FQDN 和IP 地址的网页了
# --- 部署lvm
# vim lv.yml
---
- name: 创建和使用逻辑卷
  hosts: all
  tasks:
  	- block: # 将上述的 3 个模块命令作为一个整体（ 相当于对这 3 个模块的执行结果作为一个整体进行判断）
		- name: one
		  lvg:
			vg: research
			pvs: /dev/sdb
			pesize: 150M
		- name: two
		  lvol:
			vg: research
			lv: data
			size: 150M
		- name: three
		  filesystem:
			fstype: ext4
			dev: /dev/research/data
	  rescue: # 使用 rescue 操作符进行救援，且只有 block 块中的模块执行失败后才会调用 rescue 中的救援模块
		- debug:
			msg: "Could not create logical volume of that size"
# ansible-playbook lv.yml
# --- 判断主机组名
# vim issue.yml
---
- name: 修改文件内容
	hosts: all
	tasks:
		- name: one
		  copy:
			content: 'Development'
			dest: /etc/issue
		  when: "inventory_hostname in groups.dev" #  inventory_hostname 的变量，用于定义每台主机所对应的 Ansible 服务的主机组名称
		- name: two
		  copy:
			content: 'Test'
			dest: /etc/issue
		  when: "inventory_hostname in groups.test"
		- name: three
		  copy:
			content: 'Production'
			dest: /etc/issue
		  when: "inventory_hostname in groups.prod"
# ansible-playbook issue.yml
# --- 管理文件属性
# vim chmod.yml
---
- name: 管理文件属性
  hosts: dev
  tasks:
	- name: one # 创建dir
	  file:
		path: /linuxprobe
		state: directory
		owner: root
		group: root
		mode: '2775'
	- name: two # 创建link
	  file:
		src: /linuxprobe
		dest: /linuxcool
		state: link
# ansible-playbook chmod.yml
# --- 管理密码库文件
# vim locker.yml
---
pw_developer: Imadev
pw_manager: Imamgr
# vim /root/secret.txt # 将用于加密locker.yml
whenyouwishuponastar
# chmod 600 /root/secret.txt
# vim /etc/ansible/ansible.cfg
...
vault_password_file = /root/secret.txt # 指定密码所在的文件路径
...
# ansible-vault encrypt locker.yml # 使用 ansible-vault 命令对文件进行加密
# cat locker.yml # 使用 AES 256 加密方式进行加密
$ANSIBLE_VAULT;1.1;AES256
38653234313839336138383931...
# ansible-vault rekey --ask-vault-pass locker.yml # rekey 参数手动对文件进行改密操作，同时应结合--ask-vault-pass 参数进行修改，否则 Ansible 服务会因接收不到用户输入的旧密码值而拒绝新的密码变更请求
# ansible-vault edit locker.yml # 编辑locker.yml
# ansible-vault view locker.yml
```

## 格式
格式: ansible 受管主机节点 -m 模块名称[-a 模块参数]. -a 是要传递给模块的参数，只有功能极其简单的模块才不需要
额外参数，所以大多情况下-m 与-a 参数都会同时出现

## 选项:
- -k 手动输入 SSH 协议的密码
- -i 指定主机清单文件
- -m 指定要使用的模块名
- -M 指定要使用的模块路径
- -S 使用 su 命令
- -T 设置 SSH 协议的连接超时时间
- -a 设置传递给模块的参数
- --version 查看版本信息

## example
在 Ansible 服务中， ansible 是用于执行临时任务的命令，也就在是执行后即结束（ 与剧本文件的可重复执行不同）。在使用 ansible 命令时， 必须指明受管主机的信息，如果已经设
置过主机清单文件（ /etc/ansible/hosts）， 则可以使用 all 参数来指代全体受管主机，或是用 dev、test 等主机组名称来指代某一组的主机.

```bash
# dnf install -y ansible # RHEL 8 系统的镜像文件默认不带有 Ansible 服务程序，需要从 Extra Packages for Enterprise Linux（ EPEL） 扩展软件包仓库获取
# ansible --version # 获取Ansible 的版本及配置信息
# vim /etc/ansible/hosts
[dev]
192.168.10.20
[test]
192.168.10.21
[prod]
192.168.10.22
192.168.10.23
[balancers]
192.168.10.24
[all:vars]
ansible_user=root
ansible_password=redhat
# ansible all -m ping # 对所有主机调用 ping 模块
# ansible-inventory --graph 以结构化的方式显示出受管主机的信息
# vim /etc/ansible/ansible.cfg
...
host_key_checking = False # 不需要 SSH 协议的指纹验证
...
# ansible-doc -l # 列举出当前 Ansible 服务所支持的所有模块信息
# ansible-doc a10_server # 显示出这个模块的作用、可用参数及实例等信息
# --- 角色
# dnf install -y rhel-system-roles # 安装包含系统角色的软件包
# ansible-galaxy list # 查看 RHEL 8 系统中有哪些自带的角色可用
```