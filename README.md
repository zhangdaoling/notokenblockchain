1、p2p:
 目标：接受发送通用消息，peer管理
 可以借鉴的库：libp2p：https://github.com/libp2p/go-libp2p

2、consensus：什么是一致性或者什么是共识：在信息的内容，信息之间的顺序，可能还要包括信息的执行结果达成一致，
             目标是达成一致，达成的策略可以划分为去中心化，半中心化，中心化。还有别的分法
             不要求性能的场景可以用去中心化，追求性能需要采用半中心化或中心化。
             半中心化一致性bpft，bft，raft。这些算法本身都有实现，这是用的不太一样

3、db：用什么数据库（kv），存什么数据（block，transaction，account等）
   ledb等等

4、合约和vm：
 v8：
 evm：

5、其他：merk，mempool
