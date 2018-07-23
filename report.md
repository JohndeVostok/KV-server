#基于Key-Value数据库的文件系统

计53 何琦  计54 马子轩

## 功能实现

我们的文件系统是基于Key-Value数据库的文件系统。实际可以看作是一个类似于网盘的软件。

具体实现的功能就是，将本地文件保存在远端。

其中分为前端与后端两部分

### 前端部分

前端部分使用libfuse实现，主要功能是将文件转化为Key-Value对，过程如下

```
+--------+  Split to blocks  +------------+ Base64 +------------------+
|Raw Data|------------------>|Blocked Data|------->|Base64 String Data|
+--------+  Border handling  +------------+ Encode +------------------+
                                                           |
                                                       As  |  Value
                                                           v
+-----------+   JSON     +---------------+   As    +--------------+
|Stored Keys|----------->|JSON Key String|-------->|Key-Value Pair|
+-----------+ Stringify  +---------------+   Key   +--------------+
                                                           |
                                                     HTTP  |  Pack
                                                           v
                                 +------+   CURL   +-----------------+
                                 |Server|<---------|HTTP POST Request|
                                 +------+ Perform  +-----------------+
```

其中每个文件会在本地存储3个int:q0,q1,len.

Blocksize设定为4KB，数据内容使用base64编码

Key设定为q0,q1,block_id的json

Value为编码后的数据

### 后端部分

后端部分是一个可插拔的Key-Value系统

分为driver和Key-Value数据库两部分.

通过更换不同的driver.可以使用不同的Key-Value数据库使用。

driver通过相应HTTP request提供put和get两个操作

读取前端传来的内容进行解码并连接到Key-Value数据库中

可以使用redis这类成熟的Key-Value数据库。也可以使用自己实现的分布式Key-Value Store.

redis:

就是简单的解码，配合调用redis.

KV-Server:

自己实现的KV-Server，首先实现了写分配，通过在driver进行hash分配到不同的节点.

同时读取通过driver读取较新的内容.

容错机制（依然有bug）：通过log进行checkpoint后的数据恢复.

## 运行环境

机器:Ali-cloud

系统:Linux

依赖:

fuse3 json-c glib-2.0 libcurl

语言:

c/c++ golang python

运行过程:

前端

```
sudo apt-get install libfuse-dev libjson-c-dev libcurl4-gnutls-dev libglib2.0-dev
# You may also need to add several paths to your shell .rc
./compile.sh
./a.out [target directory]
```

后端:redis

```
go run redis-driver.go
```

后端:kv-server

```
go run driver.go
# Other machine
go run client.go <port>
python test.py
# add worker
add name ip:port
```

需要确保代码中的机器设置和添加内容一致.

## 问题与挑战

其中前端遇到的一个问题就是多线程的bug.这个问题通过加print的方式，屏蔽掉了编译优化，当时通过了测试。后来从udp更换为http request的时候使用了外部库，彻底解决了这个问题。

后端遇到的问题是多机之间的同步问题，最后在期限内没解决，回滚了最后一个可用版本。然后实现了一个使用driver进行同步的方法。这样的问题在于driver压力较大，同时容错性降低了，因为本来一个去中心的设计最后变成了中心化设计，可扩展性和容错性都大大降低。

整个实验中，前端问题主要集中在编码解码和C的库使用，后端问题主要集中在同步，而后续KV-Server还会进行更新，主要是要恢复正常的分布式系统同步机制，另外优化接口和配置，使服务器配置更加便捷。