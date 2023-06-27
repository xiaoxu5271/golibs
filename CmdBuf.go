package golibs

import (
	"bytes"
	"sync"
)

type CmdBuf struct {
	buf bytes.Buffer
	lck sync.Mutex
}

func (c *CmdBuf) Write(p []byte) (n int, err error) {
	c.lck.Lock()
	defer c.lck.Unlock()
	return c.buf.Write(p)
}

func (c *CmdBuf) TruncateString() string {
	c.lck.Lock()
	defer c.lck.Unlock()
	s := c.buf.String()
	c.buf.Reset()
	return s
}

func (c *CmdBuf) Reset() {
	c.lck.Lock()
	defer c.lck.Unlock()
	c.buf.Reset()
}
