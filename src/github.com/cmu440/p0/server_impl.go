// Implementation of a MultiEchoServer. Students should write their code in this file.

package p0

import(
	"fmt"
	"net"
	"strconv"
	"bufio"
	// "runtime/debug"
)

type clientInfo struct{
	conn *net.Conn
	ch chan []byte
	reader *bufio.Reader
	writer *bufio.Writer
}

type multiEchoServer struct {
	// TODO: implement this!
	listener net.Listener
	port int
	client []clientInfo
	counts int
}

// New creates and returns (but does not start) a new MultiEchoServer.
func New() MultiEchoServer {
	// TODO: implement this!
	s := new(multiEchoServer)

	// s.client.ch = make(chan []byte)

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
			// fmt.Println("Waiting for inbound connection")
			conn, err := mes.listener.Accept()
			if err != nil {
				fmt.Println("Couldn't accept: ", err)
				return err
			}

			mes.clientRegister(conn)

			go mes.readFromCli()
			go mes.broadCastMsg()
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

func (mes *multiEchoServer) clientRegister(c net.Conn) error { //  conn net.Conn

	fmt.Println("new client register")

	mes.counts++

	mes.client = append(mes.client, clientInfo{
		conn: &c,
		reader: bufio.NewReader(c),
		writer: bufio.NewWriter(c),
		ch: make(chan []byte, 20),
	})

	fmt.Println("client is mes.client", mes.client)

	return nil
}

func (mes *multiEchoServer) readFromCli() {
	for {
		for _, c := range mes.client {
			// fmt.Println("+++++++++", c)
			line, err := c.reader.ReadBytes('\n')

			if err != nil {
				// (*(c.conn)).Close()
				fmt.Println("in readFromCli err:", err)
				continue
			}
			// fmt.Println("oops+++:", string(line[:]) )
			for _, c1 := range mes.client {
				c1.ch <- line
			}
		}
	}
}

func (mes *multiEchoServer) broadCastMsg() error {
	for {
		for _, c := range mes.client {
			line := <- c.ch
			// line = append(line, '\n')
			_, err := (*(c.conn)).Write(line)
			// fmt.Println("len of line is", len(line))
			//_, err := fmt.Fprint(c.writer, line)


			if err != nil {
				fmt.Println("error on writing", err)
				break
			}

			// fmt.Println("oops:", string(line[:]) )
			// err = c.writer.Flush()
			// if err != nil {
			// 	return err
			// }

		}

	}
	return nil
}
