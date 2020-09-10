# crc

[CRC(Cyclic redundancy check, 循环冗余校验)](https://en.wikipedia.org/wiki/Cyclic_redundancy_check)是一种根据网络数据包或电脑文件等数据产生简短固定位数校验码的一种散列函数，主要用来检测或校验数据传输或者保存后可能出现的错误.

在zip和许多其他地方找到的CRC32使用多项式0x04C11DB7(0x04C11DB7表示代号);它的相反形式0xEDB88320可能更为人所知,通常在little-endian实现中找到.

CRC32C使用不同的多项式(0x1EDC6F41,取反0x82F63B78),但计算方式相同.结果自然是不同的.这也称为Castagnoli CRC32,并且在较新的Intel CPU中最为明显,后者可以在3个周期内计算出完整的32位CRC步长.这就是CRC32C变得越来越流行的原因,因为它允许高级实现,尽管有3个周期的延迟,但每个周期有效地处理一个32位字(通过并行处理3个数据流并使用线性代数来组合结果).

> [rocksdb使用了CRC32C, 但飞腾2000+ cpu上在Put 大value时该硬件实现会报错, x64上正常](https://github.com/facebook/rocksdb/issues/7363).