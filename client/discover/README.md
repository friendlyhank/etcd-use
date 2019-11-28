运行顺序
先运行watcher再运行register

register
(1)为什么要设置租约
server可能异常退出，需要维护一个TTL,类似于心跳,master可以监听到


Watcher
(2)Watcher作为agent只要启动一个就行了,那经典的LB到底时怎样的呢？