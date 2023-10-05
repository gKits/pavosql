package cli

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

type Client struct {
	Addr string
	Port uint16
	conn net.Conn
}

func NewClient() Client {
	return Client{}
}

func (cli *Client) Run() error {
	if err := cli.connect(); err != nil {
		return err
	}

	reader := bufio.NewScanner(os.Stdin)
	for reader.Scan() {
		text := reader.Text()
		cli.conn.Write([]byte(text))
	}

	if err := cli.disconnect(); err != nil {
		return err
	}

	return nil
}

func (cli *Client) connect() error {
	var err error

	cli.conn, err = net.Dial("tcp", fmt.Sprintf("%s:%d", cli.Addr, cli.Port))
	if err != nil {
		return err
	}

	return nil
}

func (cli *Client) disconnect() error {
	return cli.conn.Close()
}
