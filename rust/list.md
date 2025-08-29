# rust
ref:
- [美国 CISA 建议放弃 C/C++，消除内存安全漏洞](https://www.oschina.net/news/269933/cisa-the-case-for-memory-safe-roadmaps)

## 解题和example
- [github.com/QMHTMY/RustBook](https://github.com/QMHTMY/RustBook)
- [Rust算法题解](https://github.com/rustlang-cn/rust-algos)
- [数据结构和算法](https://www.hackertouch.com/data-structures-and-algorithms.html)
- [aylei / leetcode-rust](https://github.com/aylei/leetcode-rust)
- [Rust算法题解](https://github.com/rustlang-cn/rust-algos)
- [rust-by-practice](https://github.com/sunface/rust-by-practice)

	答案在repo的solutions目录里
- [一个用rust实现 简易应用管理系统](https://github.com/shanliu/lsys)

### 重写
#### 命令行
- [github.com/samuela/rustybox](https://github.com/samuela/rustybox)
- [github.com/uutils/coreutils](https://github.com/uutils/coreutils)

#### other
- [Redis for rust 正式开源，使用 Rust 重塑 Redis 内核](https://www.oschina.net/news/288813/redis-for-rust-open-source)

### blog
- [Actix-Blog](https://github.com/Dengjianping/Actix-Blog)

## 教程
- [easy_rust](https://github.com/Dhghomon/easy_rust)
- [Rust入门秘籍](https://rust-book.junmajinlong.com/about.html)
- [微软开放的Rust视频教程](https://www.youtube.com/playlist?list=PLlrxD0HtieHjbTjrchBwOVks_sr8EVW1x)
- [微软开放的Rust图文教程](https://docs.microsoft.com/zh-cn/learn/paths/rust-first-steps/?WT.mc_id=academic-29077-cxa)
- [张汉东的Rust实战课-课件](https://github.com/ZhangHanDong/inviting-rust)
- [加餐｜这个专栏你可以怎么学，以及Rust是否值得学？](https://time.geekbang.org/column/article/411089)
- [ruiers/os-tutorial-cn](https://github.com/ruiers/os-tutorial-cn)
- [Rust Primer](https://hardocs.com/d/rustprimer)
- [可视化 Rust 各数据类型的内存布局](https://github.com/rustlang-cn/Rustt/blob/main/Articles/%5B2022-05-04%5D%20%E5%8F%AF%E8%A7%86%E5%8C%96%20Rust%20%E5%90%84%E6%95%B0%E6%8D%AE%E7%B1%BB%E5%9E%8B%E7%9A%84%E5%86%85%E5%AD%98%E5%B8%83%E5%B1%80.md)
- [Rust 数据内存布局](https://rustmagazine.github.io/rust_magazine_2021/chapter_6/ant-rust-data-layout.html)
- [Rust 程序设计语言（第二版 & 2018 edition）](https://kaisery.github.io/trpl-zh-cn/)
- [Rust入门第一课](https://rust-book.junmajinlong.com/ch1/00.html)
- [【Rust 新手小册】Day 4. 字节跳动开源的 Volo 框架简介](https://juejin.cn/post/7217644586868031548)
- [Visualizing memory layout of Rust's data types](https://www.youtube.com/watch?v=rDoqT-a6UFg)

## 课件
- [陈天 · Rust 编程第一课](https://github.com/tyrchen/geektime-rust)

## book
- rust权威指南 = [Rust 程序设计语言](https://kaisery.github.io/trpl-zh-cn/title-page.html)
- [Rust 版本指南(中文版)](https://rustwiki.org/zh-CN/edition-guide/)
- [Comprehensive Rust](https://google.github.io/comprehensive-rust/)
- [Hands-On Data Structures and Algorithms with Rust](https://www.amazon.co.jp/-/en/Claus-Matzinger/dp/178899552X)

	[code](https://github.com/PacktPublishing/Hands-On-Data-Structures-and-Algorithms-with-Rust)
- [Rust 秘典（死灵书）](https://nomicon.purewhite.io)
- [Rust By Practice( Rust 练习实践 )](https://practice-zh.course.rs/)

	第18章开始没有更新了 - 2024.3.21
- [Rust 参考手册 - The Rust Language Reference](https://minstrel1977.gitee.io/rust-reference/types.html)
- [Rust语言圣经](https://course.rs)
- [深入RUST标准库](https://github.com/Warrenren/inside-rust-std-library)
- [rust-course](https://github.com/sunface/rust-course)
- [Rust 秘典](https://nomicon.purewhite.io/intro.html)
- [zero-to-production](https://github.com/LukeMathWalker/zero-to-production)

	<<从零构建Rust生产级服务>>
- [rust-library-i18n:Rust 核心库和标准库中文翻译](https://github.com/wtklbm/rust-library-i18n)
- [**Comprehensive Rust**](https://google.github.io/comprehensive-rust/zh-CN/)

## gui
- [Pop!_OS COSMIC 桌面使用 Rust GUI 库 Iced 取代 GTK](https://www.oschina.net/news/212636/cosmic-desktop-iced-toolkit)
- [Slint 1.0 正式发布，Rust 编写的原生 GUI 工具包](https://www.oschina.net/news/235410/slint-1-0-released)

## next
- [用Rust实现用户态高性能存储 - Wang Pu (王璞) from DatenLord](https://weibo.com/1945106210/JAflese1N?type=repost)
- [Rust for Linux](https://rust-for-linux.com/)

## 备份
- [rustic](https://github.com/rustic-rs/rustic)
- [preserve](https://github.com/fpgaminer/preserve)

	deps: `dnf install xz-devel sqlite-devel`

	允许多个备份任务保存在一个备份目标. 存储端当前仅支持file.

	```bash
	# ./preserve -h
	# ./preserve keygen --keyfile test_key
	# ./preserve create my-backup-2016-02-25_11-56-51 /root/test --keyfile test_key --backend file:/root/backups/ # my-backup-2016-02-25_11-56-51是唯一任务标识; /root/test是备份目标(源码里假设是dir, 未测试file的情况); `backend`: 备份存储位置, 当前只支持file
	Gathering list of files...
	Reading files...
	Reading file: src/main.rs
	Progress: 0MB of 0MB
	Reading file: certs/key.pem
	Progress: 0MB of 0MB
	Reading file: certs/cert.pem
	Progress: 0MB of 0MB
	Reading file: Cargo.toml
	Progress: 0MB of 0MB
	Writing archive...
	Backup created successfully
	# ./preserve list --keyfile test_key --backend file:/root/backups # 罗列备份任务
	# mkdir -p /root/tmp
	# ./preserve restore my-backup-2016-02-25_11-56-51 /root/tmp --keyfile test_key --backend file:/root/backups/ # /root/tmp是还原target
	```

	源码:
	- keygen:

		1. `KeyStore::new()`用随机生成的128B作为master_key, 再用pbkdf2导出4把加密各类数据的256B key

			会将256B key分成2个128B bytes, 一个做SivEncryptionKeys.siv_key, 一个做SivEncryptionKeys.cipher_key
		1. 将master_key保存到`--keyfile`指定的位置
	- create:

		1. 构建`ArchiveBuilder::new`, `builder.walk()`遍历目标路径, `builder.read_files()`读取文件内容并存储到`--backend`, 由`builder.create_archive(&backup_name)`组织元数据

			判断文件属性逻辑在`read_file_metadata()`

			存储文件逻辑在`read_file_inner()`

			内容加密逻辑在`keystore.encrypt_block(&buffer)`: 会用`self.block_keys.encrypt(&[], block)`加密内容, 具体步骤:
			1. 构建siv=HMAC-SHA-512-256 (siv_key, aad || plaintext || le64(aad.length) || le64(plaintext.length)), 同时会将siv.0作为blockid, aad是空
			1. HMAC-SHA-512 (cipher_key, nonce) 导出加密密钥, 执行ChaCha20(衍生密钥 from 加密密钥，数据)
		1. `archive.encrypt(&keystore)`->`backend.store_archive()`

			archive内容json序列化后会用`lzma::compress`压缩再`keystore.encrypt_archive_metadata()`加密
	- verify:

		1. 获取archive并解码, 验证archive支付正确
		1. block_list.shuffle()

			打乱block顺序, 使得多次执行(包括中断)尽可能涵盖所有block
		1. verify_blocks()

			解码encrypted_block来验证
	- diff: 比较两次备份文件列表的差异, 比如添加/删除文件等
	- restore: create的逆过程

		`set_file_time`将原文件时间应用到新文件, 见[utimensat](https://man7.org/linux/man-pages/man2/utimensat.2.html)

		`directory_times.reverse()`+`set_file_time`可还原目录时间


- [conserve](https://github.com/sourcefrog/conserve)

	[A comparison to other backup systems](https://github.com/sourcefrog/conserve)

	一个备份任务一个备份目标. 存储端当前仅支持file.

	```bash
	# ./conserve init /root/test/storage/conserve # 初始化备份保存target
	# ./conserve backup /root/test/storage/conserve /root/tmp # 备份/root/tmp
	# ./conserve diff /root/test/storage/conserve /root/tmp # 当前/root/tmp与备份的差异
	# ./conserve versions /root/test/storage/conserve # 显示备份列表
	b0000                2023-08-01T14:14:03+08:00 incomplete # b<xxxx>是备份版本 开始实际 是否完成/完成显示耗时
	b0001                2023-08-01T14:14:26+08:00       0:00
	# ./conserve ls [-b b2] /root/test/storage/conserve/ # 显示指定版本的备份, 默认是最新
	# ./conserve restore -b b2  /root/test/storage/conserve/ /root/tmp/r # 还原指定版本
	# ./conserve validate /root/test/storage/conserve # 验证备份数据完整性
	```

	源码入口在src/bin/conserve.rs:
	- init: `Archive::create`
	- backup: `backup(&Archive::open(open_transport(archive)?)?, source, &options)?;`

## web
- [salvo](https://salvo.rs/zh-hans/)
- [Web Frameworks Benchmark - rust](https://web-frameworks-benchmark.netlify.app/result?l=rust)

## lib
ref:
- [生态系统：有哪些常有的Rust库可以为我所用?](https://time.geekbang.org/column/article/429673)

- clap : 命令行, clap 3 已经整合了 structopt
- dialoguer: 交互式的命令行
- indicati: 在命令行中提供友好的进度条
- crossbeam : 处理并发
- mdbook: 对标 gitbook
- zola: 对标 hugo
- orm: diesel(不支持异步), sea-orm(支持异步), sqlx
- [tklog](https://www.oschina.net/news/294608/tklog-released): log

## example
- 练手[rustlings](https://github.com/rust-lang/rustlings)，小练习 可以让你习惯阅读和编写 Rust

	- `l` + `选题` + `r` : 重置指定练习题
	- [Rust学习 | Rustlings通关记录与题解](https://www.cnblogs.com/Roboduster/p/17781712.html)
- 练手[exercism](https://exercism.org/)，编程语言在线学习网站
- 刷题[codewars](https://www.codewars.com/)，刷题网站，类似 LeetCode
- [talent-plan](https://github.com/pingcap/talent-plan)
- [Writing an OS in Rust](https://os.phil-opp.com/)
- [Learn Rust With Entirely Too Many Linked Lists](https://rust-unofficial.github.io/too-many-lists/index.html)
- [ParvaOS -  用 Rust 语言从头开发的操作系统](https://github.com/gianndev/ParvaOS)
- [用 Rust 重写 SQLite 数据库入门指南](https://avi.im/blag/2025/rickrolling-turso/)
- [RustFS - MinIO 国产化替代方案, 高性能分布式存储](https://github.com/rustfs/rustfs)
- [我用Rust做了一个QQ](https://juejin.cn/post/7262557466172112956)

	https://github.com/SuanCaiYv/prim

### read
- [RobustMQ](https://github.com/robustmq/robustmq-geek)
- [MsgTrans](https://github.com/zoujiaqing/msgtrans/blob/main/README.zh-CN.md)

## ai
- [Dora-rs：下一代机器人开发框架](https://www.oschina.net/news/347272)
	
	- [Dora中文社区](https://doracc.com/)
- [llms-from-scratch-cn](https://github.com/datawhalechina/llms-from-scratch-cn)

## profiling
- [Rust 性能提升 “最后一公里”：详解 Profiling 瓶颈定位与优化｜得物技术](https://my.oschina.net/u/5783135/blog/18687884)