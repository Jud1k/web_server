package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main(){
	addr,err := net.ResolveUDPAddr("udp",":42069")
	if err!=nil{
		fmt.Println(err)
		return
	}
	udpConn,err := net.DialUDP("udp",nil,addr)
	if err!=nil{
		fmt.Println(err)
	}
	defer udpConn.Close()
	for{
		fmt.Println(">")
		reader := bufio.NewReader(os.Stdin)
		userLine,err := reader.ReadString('\n') 
		if err!=nil{
			fmt.Println(err)
			break
		}
		_,err = udpConn.Write([]byte(userLine))
		if err!=nil{
			fmt.Println(err)
			break
		}
	}
}

