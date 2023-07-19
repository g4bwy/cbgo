package cbgo

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"os"
	"strings"
)

type CbIpc struct {
	conn	net.Conn
	reader	*bufio.Reader
}

func (i *CbIpc) Open() error {
	socketPath := os.Getenv("CAGEBREAK_SOCKET")
	if socketPath == "" {
		return fmt.Errorf("CAGEBREAK_SOCKET is not set. run cagebreak with -e option")
	}

	var err error
	i.conn, err = net.Dial("unix", socketPath)

	if err == nil {
		i.reader = bufio.NewReader(i.conn)
	}

	return err
}

func (i *CbIpc) Close() {
	i.conn.Close()
}

func (i *CbIpc) SendCmd(cmd string) {
	if !strings.HasSuffix(cmd, "\n") {
		cmd += "\n"
	}
	if _, err := i.conn.Write([]byte(cmd)); err != nil {
		fmt.Println("can't send IPC command '%s': %s", strings.TrimSuffix(cmd, "\n"), err)
	}
}

func (i *CbIpc) ReadEvent() ([]byte, error) {
	buf, err := i.reader.ReadBytes('\x00')
	if err != nil {
		return []byte{}, err
	}

	if !bytes.HasPrefix(buf, []byte("cg-ipc{")) {
		return []byte{}, fmt.Errorf("invalid IPC line: '%s'", buf)
	}

	d := bytes.TrimPrefix(buf, []byte("cg-ipc"))
	return bytes.TrimSuffix(d, []byte{0x00}), nil
}
