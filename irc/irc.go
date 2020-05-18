package irc

import (
	"bufio"
	"fmt"
	"net"
	"net/textproto"
	"strings"
	"sync"

	"github.com/maliur/sodaville/logger"
)

type OnMessageFunc = func(rawMessage string)
type Client struct {
	logger        *logger.Slogger
	conn          net.Conn
	closed        bool
	Debug         bool
	OnMessageFunc OnMessageFunc
}

func NewClient(logger *logger.Slogger, conn net.Conn, debug bool, handler OnMessageFunc) *Client {
	return &Client{
		logger,
		conn,
		false,
		debug,
		handler,
	}
}

func (c *Client) Run() {
	var wg sync.WaitGroup
	wg.Add(1)
	go c.startReadLoop(&wg)
	wg.Wait()
}

func (c *Client) Close() {
	c.closed = true
	c.conn.Close()
}

func (c *Client) WriteMessage(message string) {
	if c.closed {
		return
	}
	_, err := c.conn.Write([]byte(message + "\r\n"))
	if err != nil {
		c.logger.Error(fmt.Sprintf("failed to write: %v", err))
	}

	if c.Debug {
		if strings.Contains(message, "PASS") {
			c.logger.Info("< PASS ********")
		} else {
			c.logger.Info("< " + message)
		}
	}
}

// Read packets from the tcp connection
func (c *Client) startReadLoop(wg *sync.WaitGroup) {
	defer wg.Done()
	tp := textproto.NewReader(bufio.NewReader(c.conn))

	for {
		if c.closed {
			fmt.Println("it's closed")
			return
		}

		line, err := tp.ReadLine()
		if err != nil {
			c.logger.Error(fmt.Sprintf("error reading line: %v", err))
			return
		}
		if c.Debug {
			c.logger.Info("> " + line)
		}

		c.OnMessageFunc(line)
	}
}
