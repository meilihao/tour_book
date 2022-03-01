livecd
ref:
- [创建一个最小系统的squashfs镜像 by 
livemedia-creator](https://www.codeleading.com/article/48364192708/)
- [使用『grub-install』安裝GRUB](https://hugh712.gitbooks.io/grub/content/installing-grub-using-grub-install.html)

## [基于已有Live CD iso定制自已的Live CD](https://www.codeleading.com/article/88694192825/)
> 前提: squashfs-tools, genisoimage

从CentOS官网下载到 CentOS-7-x86_64-LiveGNOME-2003.iso, 主要工作分为三个部分:

1. 将ISO中的squashfs.img解压出来，以便进行定制修改。这部分可使用下面的脚本来完成（需要root权限）

	```bash
	mkdir mnt
	mount -o loop CentOS-7-x86_64-LiveGNOME-2003.iso mnt/
	mkdir iso
	cp -rfp mnt/* iso/.
	cd iso/LiveOS
	unsquashfs squashfs.img
	cd squashfs-root/LiveOS/
	mkdir ext3fs
	mount -o loop ext3fs.img ext3fs
	chroot ext3fs
	````

	当前目录下需要创建 mnt, iso，用于mount iso和拷贝iso中的内容(推荐使用`rsync -a -v`).

	之后将squashfs.img解压，再将解压出来的ext3fs.img 映像mount到ext3fs目录。

	最后chroot到ext3fs目录，就可以对LiveCD进行需要的定制了。

2. 定制LiveCD

	这里完全是根据自己的需要进行，可以安装，删除软件包. 例如下面将LiveCD启动由原来的进入graphical 界面改为进入console界面.
	
	```bash
	systemctl set-default multi-user.target
	exit
	```

	完成定制后执行exit，退回原来的文件系统.

3. 重新打包iso

	```bash
	umount ext3fs
	rm -rf ext3fs
	cd ../..
	ls
	rm -rf squashfs.img
	mksquashfs squashfs-root/ squashfs.img -noappend -always-use-fragments
	rm -rf squashfs-root/
	ls -l
	cd ..
	pwd
	ls
	mkisofs -b isolinux/isolinux.bin -c isolinux/boot.cat \
	    --no-emul-boot --boot-load-size 4 --boot-info-table \
	    -R -J -v -T -V "CentOS-7-x86_64-GNOME-2003" \
	    -o ../customer_iso.iso .
	cd ..
	ls -l customer_iso.iso
	umount mnt
	echo "Finished"
	```

	生成的customer_iso.iso就是最终定制的LiveCD.


	> xorriso也可构建iso: `xorriso -pathspecs as_mkisofs -as mkisofs -iso-level 3 -full-iso9660-filenames -volid "CentOS 7 x86_64" -appid  -publisher  -preparer "prepared by me" -eltorito-boot isolinux/isolinux.bin -eltorito-catalog isolinux/boot.cat -no-emul-boot -boot-load-size 4 -boot-info-table -isohybrid-mbr ../boot.mbr -eltorito-alt-boot -e images/efiboot.img -no-emul-boot -isohybrid-gpt-basdat -output ../repacked.iso .`

