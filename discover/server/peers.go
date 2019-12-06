package server

import(
	"time"
)

type Node struct{
	Key string
	ServiceMeta *ServiceMeta
}

//ServiceMeta-
type ServiceMeta struct {
	IP  string
	Endpoint string
	Weight int32  //权重
	CreateTime time.Time //创建时间
}

