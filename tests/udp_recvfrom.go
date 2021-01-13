// udp_recvfrom.go
// go run udp_recvfrom.go
// golangによるudpの送信"だけ"するプログラム

package main

import (
	"fmt"
	"net"
)

func main() {
	addr, _ := net.ResolveUDPAddr("udp", ":12345")
	sock, _ := net.ListenUDP("udp", addr)

	i := 0
	for {
		i++
		buf := make([]byte, 1024)
		rlen, _, err := sock.ReadFromUDP(buf)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(string(buf[0:rlen]))
		//fmt.Println(i)
		//go handlePacket(buf, rlen)
	}
}
