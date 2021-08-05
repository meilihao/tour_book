# meson

## FAQ
### meson build报错`meson_options.txt...\nUnknown type feature`
meson过旧, `pip3 install meson==0.54`, 在`/usr/local/bin/meson`

### `Could not detect Ninja v1.5 or newer`
定义在`mesonbuild/backend/ninjabackend.py`, 引用的是`mesonbuild/environment.py#detect_ninja_command_and_version()`, 它里面规定了默认ninja的版本, 比如meson==0.54是ninja1.7