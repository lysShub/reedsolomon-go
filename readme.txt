这是专为网络编码设计reedsolomon fec, 与其他reedsolomon包的主要区别是:

1. 无状态的, 在编解码时, 只需传入对应的Para 

2. 不要求block长度一致 

以上特性为网络数据包的编码与恢复提供便利。


reference:

https://github.com/klauspost/reedsolomon

https://github.com/templexxx/reedsolomon

https://vearne.cc/archives/39331
