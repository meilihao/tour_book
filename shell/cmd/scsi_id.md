# scsi_id
其包含在udev程序包中(`/lib/udev/scsi_id`)，可以在配置multipath.conf时通过该程序来获取scsi设备的序号, 通过序号，便可以判断多个路径对应了同一设备.这个是多路径实现的关键.
scsi_id是通过sg驱动，向设备发送EVPD page80或page83 的inquery命令来查询scsi设备的标识, 但一些设备并不支持EVPD 的inquery命令，所以它们无法被用来生成multipath设备; 但可以改写scsi_id，为不能提供scsi设备标识的设备虚拟一个标识符，并输出到标准输出.

> VMware Virtual disk没有wwid
