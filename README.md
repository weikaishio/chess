# 简述
![Alt text](https://github.com/weikaishio/chess/blob/master/img.png?raw=true)
分为接入层server_gate、逻辑层server_game、数据层，redis和server_center为协调层。 

一个server_gate可以带一个或者多个server_game，每个server_gate都有一个唯一的gateid。当server_gate连接上server_game时，会把gateid发送给server_game。  

### 客户端接入流程
1. client -> login: 客户端先从server_login登录成功后，得到token和gate地址
2. client -> gate: 
* TCP长链：客户端的通过1拿到的server_gate地址和token，登录gate成功后，server_gate转发到server_game，网络连接由gateid和connid(连接id)标识，由server_game保存到server_center，并有server_center分发给所有的server_game。这样，每个server_game都有所有客户端的连接信息的一个缓存，并根据连接信息来转发消息。如果客户端连接的gateid和server_game的相同，则直接转发给server_gate；如果不一样，则插入到相应的redis的消息队列中，相应的server_gate会取出并转发。  
* HTTP短链：客户端与server_gate的连接是无状态的，每次客户端请求server_gate需要带上token的auth信息，同样由server_gate转发消息给server_game处理

数据层包括server_table和user_db。  
user_db保存玩家信息，可以为任何支持redis协议的数据库。  
server_table处理棋牌桌子逻辑，包括配桌、查询和保存桌子信息、查询玩家的房间位置。配桌成功时，往redis插入消息，server_game获取并处理。牌桌有超时消息时，往redis插入消息，server_game获取并处理。  

server_login处理账号登录，如微信登录等。账号登录成功，向redis插入玩家登录成功的记录，并给客户端返回游戏登录地址和登录token。  
server_login和server_game之间通过redis通信，如验证token、获取玩家信息。可以根据客户端版本选择游戏登录地址，方便更新。

## server_gate
* 客户端网络连接（TCP&HTTP）管理
* 转发客户端消息给逻辑层 
``` golang
client->gate
type ClientGame struct {
	Userid  uint32
	Msgid   uint16
	Token   []byte -- 增加支持http连接所用的每次请求校验
	MsgBody []byte
}
gate->game
type GateBackend struct {
	Msgid  uint16
	Connid uint32
	Token  []byte
	MsgBuf []byte -- ClientGame CBCDecrypt
}
```
* 转发逻辑层消息给客户端
* 发送广播消息

## server_game
* 处理业务逻辑
* 维护客户端连接信息(逻辑层中称为session)
* 转发响应消息到接入层
* 非长链客户端从gate过来的请求增加token校验
* server_table有个上行rpc连接和 下行tpc连接
```golang
game->gate
type BackendGate struct {
	Connid  uint32
	Connids []uint32
	MsgBuf  []byte  -- GameClient CBCEncrypt
} 
type GameClient struct {
	Msgid   uint16
	Result  uint16
	MsgBody []byte
}
```

## server_table
* 进入房间、离开房间、配桌
* 查询和更新桌子信息
* 查询玩家的房间位置

## server_center
* 客户端连接信息管理
* 当增加、删除连接信息时，分发给逻辑层
* 连接信息持久化，避免重启丢失

# 扩展性
* 接入层和逻辑层可以任意扩展
* server_center只管理连接信息，经过测试增删达到每秒3万次左右，不会成为瓶颈
* server_table主要功能在于桌子信息的查询和更新(分别测试，每秒可达8w次左右)，如果达到瓶颈，可以每个房间开一个(这种情况下，查询玩家在哪个房间要查询所有的server_table)
* user_db推荐使用具有持久化功能的ssdb，数据达到一定规模时，仍然具有较高的性能，扩展可以根据userid做哈希。另外也可以使用LedisDB。

### 注：forked from ![chess](https://github.com/gochenzl/chess chess) 做了一定适应性修改