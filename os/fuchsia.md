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
