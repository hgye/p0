package main

import (
	"fmt"
	"net"
	"os"
	"runtime/debug"
	"bufio"
)

const (
	defaultHost = "localhost"
	defaultPort = 9999
)

// To test your server implementation, you might find it helpful to implement a
// simple 'client runner' program. The program could be very simple, as long as
// it is able to connect with and send messages to your server and is able to
// read and print out the server's echoed response to standard output. Whether or
// not you add any code to this file will not affect your grade.
func main() {
	fmt.Println("Not implemented.")
	conn, err := net.Dial("tcp", defaultHost+":9999")
	//defer conn.Close()

	if err != nil {
		fmt.Println("error ", err)
		os.Exit(-1)
	}

	buf := []byte("here it is\n fuckoff\n")

	_, err = conn.Write(buf)
	if err != nil {
		fmt.Println("error ", err)
		os.Exit(-1)
	}

	for i:=0; i< 50; i++ {

		message, _ := bufio.NewReader(conn).ReadString('\n')
		fmt.Print("Message from server: "+message)
		/*
		var buf1 []byte
		_, err = conn.Read(buf1)

		fmt.Println("echo back: ", string(buf1[:]))

		if err != nil{
			fmt.Println("error ++++", err)
			//debug.PrintStack()
			//os.Exit(-1)
		}*/
	}

	select{

	}
	debug.PrintStack()
	os.Exit(0)
}
