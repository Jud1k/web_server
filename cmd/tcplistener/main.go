package main

import (
	"fmt"
	"io"
	"net"
	"strings"
	"request"
)

func main(){
	listener,err := net.Listen("tcp",":42069")
	if err!=nil{
		fmt.Println(err)
		return
	}
	defer listener.Close()
	for{
		conn,err := listener.Accept()
		if err !=nil{
			fmt.Println(err)
			return 
		}
		fmt.Println("connection accepted")
		Req
		for val := range(ch){
			fmt.Printf("read: %s\n",val)
		}
	}
	
}
