package util

import "net"

func GetGlobalUnicastIp() string{
	addrs,_ := net.InterfaceAddrs()
	for _,address := range addrs{
		ipnet, ok := address.(*net.IPNet)
		if !ok {
			continue
		}
		isGlobalUnicast := ipnet.IP.IsGlobalUnicast()
		if !isGlobalUnicast{
			continue
		}
		return ipnet.IP.String()
	}
	return ""
}
