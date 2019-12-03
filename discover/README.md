运行顺序
先运行watcher再运行register

2.问题

register
(1)为什么要设置租约
server可能异常退出，需要维护一个TTL,当节点宕机的时候无法继续维护(相当于发送心跳),
那么则认为这个节点宕机，这时候会自动取消注册，watcher监听也不在监听到这个节点，不设置租约则会一直存在


Watcher
(2)Watcher作为agent只要启动一个就行了,那经典的LB到底时怎样的呢？


3.
##IP的自动确认

##服务注册
    .客户端注册
    .服务端注册
    
4.
register.go,watcher.go可直接在多台机部署
    