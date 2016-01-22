// Implementation of a MultiEchoServer. Students should write their code in this file.
// basicly ref: https://gist.github.com/drewolson/3950226

package p0

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
)

type client struct {
	id       int
	conn     net.Conn
	readMsg  chan string // The message read from the network.
	writeMsg chan string // The message to write to the network.
	reader   *bufio.Reader
	writer   *bufio.Writer
}

type multiEchoServer struct {
	// TODO: implement this!
	port     int
	listener net.Listener

	clients  []client
	join     chan net.Conn
	readMsg  chan string
	writeMsg chan string
	counts   int
}

// New creates and returns (but does not start) a new MultiEchoServer.
func New() MultiEchoServer {
	// TODO: implement this!

	s := new(multiEchoServer)

	s.join = make(chan net.Conn)
	s.readMsg = make(chan string)
	s.writeMsg = make(chan string)

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

	go func() {

		for {
			// fmt.Println("Waiting for inbound connection")
			conn, err := mes.listener.Accept()
			if err != nil {
				fmt.Println("Couldn't accept: ", err)
				continue
			}

			// mes.clientJoin(conn)

			mes.counts++
			mes.join <- conn

		}
	}()

	go mes.run()

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

func (mes *multiEchoServer) clientJoin(c net.Conn) { //  conn net.Conn

	fmt.Println("new client register")
	// mes.counts++

	readChan, writeChan := make(chan string), make(chan string, 100)
	mes.clients = append(mes.clients, client{
		conn:     c,
		reader:   bufio.NewReader(c),
		writer:   bufio.NewWriter(c),
		readMsg:  readChan,
		writeMsg: writeChan,
	})

	cli := mes.clients[len(mes.clients)-1]
	// fmt.Println("client is mes.client", mes.clients)
	go cli.run(mes)
}

func (mes multiEchoServer) run() {
	for {
		select {
		case conn := <-mes.join:
			mes.clientJoin(conn)

		case msg := <-mes.readMsg:

			for _, c := range mes.clients {
				c.writeMsg <- msg

			}

		default:
		}
	}
}

func (c client) run(mes *multiEchoServer) {
	go c.read()
	go c.write()

	for {
		select {
		case msg := <-c.readMsg:

			mes.readMsg <- msg

		default:
		}
	}
}

func (c client) read() {
	for {
		msg, err := c.reader.ReadBytes('\n')

		if err != nil {
			return
		}

		c.readMsg <- string(msg)
	}
}

func (c client) write() {
	for {
		for data := range c.writeMsg {
			_, err := c.writer.WriteString(data)
			// fmt.Println("data is ", data)

			if err != nil {
				return
			}

			c.writer.Flush()
		}
	}
}
