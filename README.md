# Etcd的应用场景

## 场景一: 服务的注册和服务发现
* [register](https://github.com/friendlyhank/goetcd/blob/master/discover/example/register.go)
* [watcherr](https://github.com/friendlyhank/goetcd/blob/master/discover/example/watcher.go)

## 场景二: 消息发布和订阅

##场景三: 负载均衡
* [selector|按随机选取](https://github.com/friendlyhank/goetcd/blob/master/balancer/client/selector_test.go)
* [selector|按循环选取](https://github.com/friendlyhank/goetcd/blob/master/balancer/client/selector_test.go)
* [selector|按加权选取](https://github.com/friendlyhank/goetcd/blob/master/balancer/client/selector_test.go)
* [selector|按ping选取](https://github.com/friendlyhank/goetcd/blob/master/balancer/client/selector_test.go)
* [selector|按Hash选取](https://github.com/friendlyhank/goetcd/blob/master/balancer/client/selector_test.go)
* [selector|按地理位置选取](https://github.com/friendlyhank/goetcd/blob/master/balancer/client/selector_test.go)

##场景四: 分布式通知与协调

##场景五: 分布式锁

##场景六:分布式队列

##场景七： 集群监控与Leader竞选
