// Implementation of a MultiEchoServer. Students should write their code in this file.

package p0

import(
	"fmt"
	"net"
	"strconv"
	"bufio"
	// "runtime/debug"
)

type connectClient struct{
	conn net.Conn
	ch chan []byte
	reader *bufio.Reader
	writer *bufio.Writer
}

type multiEchoServer struct {
	// TODO: implement this!
	listener net.Listener
	port int
	client connectClient
	counts int
}

// New creates and returns (but does not start) a new MultiEchoServer.
func New() MultiEchoServer {
	// TODO: implement this!
	s := new(multiEchoServer)

	s.client.ch = make(chan []byte)

	return s
}

func (mes *multiEchoServer) Start(port int) error {
	// TODO: implement this!
	// defer debug.PrintStack()
	mes.port = port

	ln, err := net.Listen("tcp", ":"+strconv.FormatInt(int64(port), 10))
	if err != nil {
		fmt.Println("Couldn't listen:", err)
		return err
	}

	mes.listener = ln

	go func() error {
		for {
			fmt.Println("Waiting for inbound connection")
			conn, err := mes.listener.Accept()
			if err != nil {
				fmt.Println("Couldn't accept: ", err)
				return err
			}

			mes.client.conn = conn

			mes.counts++
			mes.client.reader = bufio.NewReader(conn)
			mes.client.writer = bufio.NewWriter(mes.client.conn)

			go mes.handleConn()
			go mes.sendBackMsg(conn)
		}
	}()

	return nil
}

func (mes *multiEchoServer) Close() {
	// TODO: implement this!
	mes.listener.Close()
}

func (mes *multiEchoServer) Count() int {
	// TODO: implement this!
	return mes.counts
}

// TODO: add additional methods/functions below!

func (mes *multiEchoServer) handleConn() error { //  conn net.Conn
	fmt.Println("Reading from connection")

	for {
		line, err := mes.client.reader.ReadBytes('\n')

		if err != nil {
			return err
		}

		mes.client.ch <- line

	}

	// return nil
}


func (mes *multiEchoServer) sendBackMsg(c net.Conn) error {
	//var line []byte

	for {
		// fmt.Println("sending out message")

		line := <- mes.client.ch
		// line = append(line, '\n')
		_, err := c.Write(line)
		fmt.Println("len of line is", len(line))
		// _, err := fmt.Fprint(mes.client.writer, line)

		fmt.Println("oops:", string(line[:]) )
		if err != nil {
			fmt.Println("error on writing", err)
			return err
		}

		//io.Copy(mes.client.writer, mes.client.conn)

		err = mes.client.writer.Flush()
		if err != nil {
			return err
	}

	}
	return nil
}
