// Implementation of a MultiEchoServer. Students should write their code in this file.

package p0

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	// "runtime/debug"
)

type nwEvent struct {
	cli      *clientInfo // The client that received a network event.
	readMsg  string      // The message read from the network.
	writeMsg string      // The message to write to the network.
	err      error       // Notifies us that an error has occurred (if non-nil).
}

type clientInfo struct {
	id   int
	conn net.Conn
	//ch     chan []byte
	reader *bufio.Reader
	writer *bufio.Writer
}

type multiEchoServer struct {
	// TODO: implement this!
	listener net.Listener
	port     int
	client   []clientInfo
	counts   int
}

// New creates and returns (but does not start) a new MultiEchoServer.
func New() MultiEchoServer {
	// TODO: implement this!
	s := new(multiEchoServer)

	// s.registerChan = make(chan int)

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

	go func() {

		for {
			// fmt.Println("Waiting for inbound connection")
			conn, err := mes.listener.Accept()
			if err != nil {
				fmt.Println("Couldn't accept: ", err)
				continue
			}

			mes.clientJoin(conn)
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

func (mes *multiEchoServer) clientJoin(c net.Conn) { //  conn net.Conn

	fmt.Println("new client register")

	mes.counts++

	mes.client = append(mes.client, clientInfo{
		conn:   c,
		reader: bufio.NewReader(c),
		writer: bufio.NewWriter(c),
		// ch:     make(chan []byte, 20),
	})

	// cli := mes.client[len(mes.client)-1]
	rc := make(chan *nwEvent)
	go mes.startReadFromCli(rc)
	go mes.startBroadCastMsg(rc)

	fmt.Println("client is mes.client", mes.client)

}

func (mes *multiEchoServer) startReadFromCli(readChan chan<- *nwEvent) {
	for {
		for _, cli := range mes.client {
			cli := cli
			fmt.Println("+++++++++", cli)

			line, err := cli.reader.ReadBytes('\n')

			if err != nil {
				// (*(c.conn)).Close()
				readChan <- &nwEvent{err: err}
				fmt.Println("in readFromCli err:", err)
				return
			}

			readChan <- &nwEvent{
				cli:     &cli,
				readMsg: string(line),
			}

			// fmt.Println("oops+++:", string(line[:]) )
		}
	}
}

func (mes *multiEchoServer) startBroadCastMsg(readChan <-chan *nwEvent) {
	for {
		select {
		case event := <-readChan:
			for _, cli := range mes.client {

				_, msg := event.cli, event.readMsg

				if event.err != nil {
					continue
				}
				// line = append(line, '\n')
				_, err := cli.conn.Write([]byte(msg))
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
		default:
		}
	}
}
