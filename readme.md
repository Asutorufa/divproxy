# DivProxy

![image](https://raw.githubusercontent.com/Asutorufa/divproxy/master/img/view.png?token=AD5OHSUFJ6XK3YAWSZGN7H26ARVHY)

## server

```shell
+---+                +----------------------------+
|i  |                |                            |
|n  |--------------->|        divproxy            |
|b  |                |   +-------------------+    |
|o  |--------------->|   |match domain/cidr  |------------------>
|u  |                |   |select one outbound|    |
|n  |--------------->|   +-------------------+    |            
|d  |                |                            |
+---+                +----------------------------+
```

### socks5

```shell
第一次握手
+-------------------------+
|    socks5 forst verify  |
+-------------------------+
| VER | NMETHODS | METHOD |
+-------------------------+
VER是SOCKS版本，这里应该是0x05；
NMETHODS是METHODS部分的长度；
METHOD是客户端支持的认证方式列表，每个方法占1字节。当前的定义是：
0x00 不需要认证
0x01 GSSAPI
0x02 用户名、密码认证
0x03 - 0x7F由IANA分配（保留）
0x80 - 0xFE为私人方法保留
0xFF 无可接受的方法

服务器应答
+----+--------+
|VER | METHOD |
+----+--------+
| 1  |   1    |
+----+--------+
VER 和 METHOD 与上面的相同

SOCKS5 用户名密码认证方式
在客户端、服务端协商使用用户名密码认证后，客户端发出用户名密码，格式为（以字节为单位）：
+------------+--------------+--------+---------+------+
| 鉴定协议版本 | 用户名长度	| 用户名 | 密码长度| 密码 |
+------------+--------------+--------+---------+------+
|      1	 |       1      |   动态 |    1    | 动态 |
+------------+--------------+--------+---------+------+
鉴定协议版本当前为 0x01.

服务器鉴定后发出如下回应：
+------------+--------+
|鉴定协议版本|鉴定状态|
+------------+--------+
|      1     |    1   |
+------------+--------+
其中鉴定状态 0x00 表示成功，0x01 表示失败。


socks5第二次握手
+-----------------------------------------------------------------------+
|					       socks5 protocol                              |
+-----------------------------------------------------------------------+
| socks_version | link_style | none | ipv4/ipv6/domain | address | port |
+-----------------------------------------------------------------------+
+----------------------------------------------------------------+
|					      socks5协议                             |
+-----------+----------+----------+----------------+------+------+
| socks版本 | 连接方式 | 保留字节 | 域名/ipv4/ipv6 | 域名 | 端口 |
+-----------+----------+----------+----------------+------+------+
VER是SOCKS版本，这里应该是0x05；
CMD是SOCK的命令码
	0x01表示CONNECT请求
	0x02表示BIND请求
	0x03表示UDP转发
RSV 0x00，保留
ATYP DST.ADDR类型
0x01 IPv4地址，DST.ADDR部分4字节长度
0x03 域名，DST.ADDR部分第一个字节为域名长度，DST.ADDR剩余的内容为域名，没有\0结尾。
0x04 IPv6地址，16个字节长度。
DST.ADDR 目的地址
DST.PORT 网络字节序表示的目的端口

第二次服务器应答
+---+---+----+----+--------+--------+
|VER|REP| RSV|ATYP|BND.ADDR|BND.PORT|
+---+---+----+----+--------+--------+
| 1 | 1 |0x00|  1 |  动态  |    2   |
+---+---+----+----+--------+--------+

VER是SOCKS版本，这里应该是0x05；
REP应答字段
    0x00表示成功
    0x01普通SOCKS服务器连接失败
    0x02现有规则不允许连接
    0x03网络不可达
    0x04主机不可达
    0x05连接被拒
    0x06 TTL超时
    0x07不支持的命令
    0x08不支持的地址类型
    0x09 - 0xFF未定义
RSV 0x00，保留
ATYP BND.ADDR类型
    0x01 IPv4地址，DST.ADDR部分4字节长度
    0x03域名，DST.ADDR部分第一个字节为域名长度，DST.ADDR剩余的内容为域名，没有\0结尾。
    0x04 IPv6地址，16个字节长度。
BND.ADDR 服务器绑定的地址
BND.PORT 网络字节序表示的服务器绑定的端口
```

### DNS

```shell
+------------------------------+
|             id               |  16bit
+------------------------------+
|qr|opcpde|aa|tc|rd|ra|z|rcode |
+------------------------------+
|          QDCOUNT             |
+------------------------------+
|          ancount             |
+------------------------------+
|          nscount             |
+------------------------------+
|          arcount             |
+------------------------------+

• ID：这是由生成DNS查询的程序指定的16位的标志符。
该标志符也被随后的应答报文所用，
 申请者利用这个标志将应答和原来的请求对应起来。

• QR：该字段占1位，用以指明DNS报文是请求（0）还是应答（1）。
• OPCODE：该字段占4位，用于指定查询的类型。
 值为0表示标准查询，值为1表示逆向查询，值为2表示查询服务器状态，
 值为3保留，值为4表示通知，值为5表示更新报文，值6～15的留为新增操作用。

• AA：该字段占1位，仅当应答时才设置。
 值为1，即意味着正应答的域名服务器是所查询域名的
 管理机构或者说是被授权的域名服务器。

• TC：该字段占1位，代表截断标志。
 如果报文长度比传输通道所允许的长而被分段，该位被设为1。

• RD：该字段占1位，是可选项，表示要求递归与否。
 如果为1，即意味 DNS解释器要求DNS服务器使用递归查询。

• RA：该字段占1位，代表正在应答的域名服务器可以执行递归查询，
 该字段与查询段无关。
• Z：该字段占3位，保留字段，其值在查询和应答时必须为0。
• RCODE：该字段占4位，该字段仅在DNS应答时才设置。用以指明是否发生了错误。
允许取值范围及意义如下：
0：无错误情况，DNS应答表现为无错误。
1：格式错误，DNS服务器不能解释应答。
2：严重失败，因为名字服务器上发生了一个错误，DNS服务器不能处理查询。
3：名字错误，如果DNS应答来自于授权的域名服务器，
 意味着DNS请求中提到的名字不存在。
4：没有实现。DNS服务器不支持这种DNS请求报文。
5：拒绝，由于安全或策略上的设置问题，DNS名字服务器拒绝处理请求。
6 ～15 ：留为后用。

• QDCOUNT：该字段占16位，指明DNS查询段中的查询问题的数量。
• ANCOUNT：该字段占16位，指明DNS应答段中返回的资源记录的数量，在查询段中该值为0。
• NSCOUNT：该字段占16位，指明DNS应答段中所包括的授权域名服务器的资源记录的数量，在查询段中该值为0。
• ARCOUNT：该字段占16位，指明附加段里所含资源记录的数量，在查询段中该值为0。
(2）DNS正文段
在DNS报文中，其正文段封装在图7-42所示的DNS报文头内。DNS有四类正文段：查询段、应答段、授权段和附加段。
```

## matcher

```shell
            root
           /    \
          0      1
         / \    / \
        0   1  0   1
       /
    当ip匹配到此处,此处已无任何子树,且是某一cidr的末尾时则匹配成功
    若此处节点为null(golang为nil)且不是某一cidr的末尾时则匹配失败
```
  
域名的前缀树相同,只不过域名不再是只有0和1,而且在匹配的时候还需要跳过前面的那些前缀.

```shell
        +---------+
        |  root   |
        +---------+
       /     /      \
facebook   google   twitter  ...
  /        /   \       \
com      com   mail    com   ...
                \
               com

在对域名匹配时
如对 www.play.google.com匹配:
    没有www  跳过
    没有play 跳过
    有google 继续
    有com 且 域名已为最后一个节点 -> 判断trie中是否为最后的一个子树 -> 是 -> 匹配成功
```
