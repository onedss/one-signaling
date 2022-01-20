package main

import (
	"log"
	"net"
	"strings"
)

func main() {
	listen, err := net.Listen("tcp", "127.0.0.1:8888")
	if err != nil {
		panic(err)
	}
	defer listen.Close()

	for {
		conn, err := listen.Accept()
		if err != nil {
			panic(err)
		}
		for {
			data := make([]byte, 10)

			_, err := conn.Read(data)

			if err != nil {
				log.Printf("%s\n", err.Error())
				break
			}

			receive := string(data)
			log.Printf("receive msg: %s\n", receive)

			send := []byte(strings.ToUpper(receive))
			_, err = conn.Write(send)
			if err != nil {
				log.Printf("send msg failed, error: %s\n", err.Error())
			}

			log.Printf("send msg: %s\n", receive)
		}
	}
}