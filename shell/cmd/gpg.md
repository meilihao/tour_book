# gpg
```
# gpg --verify qemu-5.0.0-rc3.tar.xz.sig qemu-5.0.0-rc3.tar.xz
```

## FAQ
### Can't check signature: No public key
```
gpg --verify qemu-5.0.0-rc3.tar.xz.sig qemu-5.0.0-rc3.tar.xz
gpg: keybox '/home/chen/.gnupg/pubring.kbx' created
gpg: Signature made Thu 16 Apr 2020 07:38:50 AM CST
gpg:                using RSA key CEACC9E15534EBABB82D3FA03353C9CEF108B584
gpg: Can't check signature: No public key
```

解决:
```
# gpg --search-keys CEACC9E15534EBABB82D3FA03353C9CEF108B584
# gpg --recv-key CEACC9E15534EBABB82D3FA03353C9CEF108B584
```

### new key but contains no user ID - skipped
```
gpg [--verbose] --recv-key CEACC9E15534EBABB82D3FA03353C9CEF108B584
gpg: key 3353C9CEF108B584: new key but contains no user ID - skipped
gpg: Total number processed: 1
gpg:           w/o user IDs: 1
```

默认keyserver不受信任, 更换keyserver.

解决:
```
# gpg --keyserver keyserver.ubuntu.com --recv CEACC9E15534EBABB82D3FA03353C9CEF108B584
```