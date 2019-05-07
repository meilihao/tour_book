# fuchsia
## 编译
按照官方文档即可

## FAQ
### Script returned non-zero exit code
```log
fx set core.x64
ERROR at //build/config/fuchsia/zircon.gni:90:3: Script returned non-zero exit code.
  exec_script("//build/zircon/populate_zircon_public.py",
  ^----------
Current dir: /home/chen/git/fuchsia/out/default/
Command: /usr/bin/env /home/chen/git/fuchsia/build/zircon/populate_zircon_public.py /home/chen/git/fuchsia/out/default.zircon/legacy_dirs.json
Returned 1.
stderr:

Traceback (most recent call last):
  File "/home/chen/git/fuchsia/build/zircon/populate_zircon_public.py", line 72, in <module>
    sys.exit(main())
  File "/home/chen/git/fuchsia/build/zircon/populate_zircon_public.py", line 36, in main
    for top_dir, subdirs in dirs.iteritems():
AttributeError: 'dict' object has no attribute 'iteritems'

See //BUILD.gn:5:1: whence it was imported.
import("//build/config/fuchsia/zircon.gni")
```
fx依赖python2.7,使用`sudo update-alternatives --config python`修改即可

### could not extend fvm, unable to stat fvm image???
```bash
$ fx run  -m 3072 -g                                                                                                                        
ERROR: could not extend fvm, unable to stat fvm image
CMDLINE: TERM=xterm-256color kernel.serial=legacy kernel.entropy-mixin=523d486471f16e635e887030974ac81c5476824864082c30819b240a85155da2 kernel.halt-on-panic=true
+ exec /home/chen/git/fuchsia/buildtools/linux-x64/qemu/bin/qemu-system-x86_64 -kernel /home/chen/git/fuchsia/out/default/../default.zircon/multiboot.bin -initrd /tmp/tmp.fYJ0VoVdNY/fuchsia-ssh.zbi -m 3072 -serial stdio -vga std -drive file=/tmp/tmp.fYJ0VoVdNY/fvm.blk,format=raw,if=none,id=mydisk -device ich9-ahci,id=ahci -device ide-drive,drive=mydisk,bus=ahci.0 -net none -smp 4,threads=2 -machine q35 -device isa-debug-exit,iobase=0xf4,iosize=0x04 -cpu Haswell,+smap,-check,-fsgsbase -append 'TERM=xterm-256color kernel.serial=legacy kernel.entropy-mixin=523d486471f16e635e887030974ac81c5476824864082c30819b240a85155da2 kernel.halt-on-panic=true '
qemu-system-x86_64: -drive file=/tmp/tmp.fYJ0VoVdNY/fvm.blk,format=raw,if=none,id=mydisk: Could not open '/tmp/tmp.fYJ0VoVdNY/fvm.blk': No such file or directory
rm: 无法删除'/tmp/tmp.fYJ0VoVdNY/fvm.blk': 没有那个文件或目录
```
