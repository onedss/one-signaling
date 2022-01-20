package main

import (
	"bufio"
	"fmt"
	"net"
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
		reader := bufio.NewReader(conn)
		for {
			slice, err := reader.ReadSlice('\n')
			if err != nil {
				continue
			}
			fmt.Printf("%s", slice)
		}
	}
}
