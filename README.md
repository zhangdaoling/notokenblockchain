# no token blockchain
改区块链涉及分为两部分或者说分为两种节点。


# key words
chainNode：负责共识和产块，没有经济模型

WorkNode：负责业务功能。例如使用evm实现token功能

block：区块

message：区块中的信息

mainChain：主链，管理子链

childChain：子链


# 计划
1、实现upow共识chainNode，测试共识参数

2、移植vm（暂定evm）实现workNode。实现token功能

3、实现mainChain


# to do
开发no token chainNode

进度：
1、p2pserver：完成

2、db：完成

3、crypto：完成

4、core：
