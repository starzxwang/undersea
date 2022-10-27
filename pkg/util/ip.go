package util

import (
	"fmt"
	"net"
	"os"
)

func GetIpAddr() string {
	addrs, err1 := net.InterfaceAddrs()

	if err1 != nil {
		fmt.Println(err1)
		os.Exit(1)
	}

	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {

			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}
