# base64

## 选项
- -d/--decode : 解码

## example
```bash
# echo  'linux' | base64
# echo 'bGludXgK' | base64 -d
```

# basenc
ref:
- [man basenc](https://man.archlinux.org/man/community/man-pages-zh_cn/basenc.1.zh_CN)

> from coreutils
支持base64/base64url/base32/base32hex/base16/base2msbf/base2lsbf/z85

## example
```bash
# echo  'linux' | basenc --base64
# echo  'bGludXgK' | basenc --base64 -d
```