package golibs

import (
	"encoding/json"
	"strings"
	"sync"

	"github.com/nats-io/nats.go"
)

type TopicFilt struct {
	client *NClient
	sync.RWMutex

	//主题对应消息队列
	chHandl map[string]chan stdHandl
}

type StdHead struct {
	Src     string      `json:"src"`
	Topic   string      `json:"topic"`
	Content interface{} `json:"content"`
}

type stdHandl struct {
	Subj    string
	Head    StdHead
	Content []byte
}

type SubjHandler func(t *TopicFilt, subj *string, head *StdHead, content *[]byte)

func (t *TopicFilt) Init(client *NClient) *TopicFilt {
	t.client = client
	t.chHandl = make(map[string]chan stdHandl)
	return t
}

func (t *TopicFilt) RegistInterestTopic(subj, topic string, capacity int, handler SubjHandler) {
	t.RWMutex.Lock()
	defer t.RWMutex.Unlock()

	t.chHandl[subj+"."+topic] = make(chan stdHandl, capacity)
	go t.doSubject(t.chHandl[subj+"."+topic], handler)
}

func (t *TopicFilt) doSubject(chHandl chan stdHandl, hdl SubjHandler) {
	for {
		msg := <-chHandl
		hdl(t, &msg.Subj, &msg.Head, &msg.Content)
	}
}

func (t *TopicFilt) onSubject(m *nats.Msg) {
	var head StdHead
	subj := m.Subject
	data := &m.Data

	e := json.Unmarshal(*data, &head)
	if e != nil {
		return
	}
	content, e := json.Marshal(head.Content)
	if e != nil {
		return
	}
	head.Content = nil
	var msg stdHandl
	msg.Content = content
	msg.Head = head
	msg.Subj = subj

	t.RWMutex.RLock()
	for k, v := range t.chHandl {
		if isMatchTopic(subj+"."+head.Topic, k) {
			v <- msg
		}
	}
	t.RWMutex.RUnlock()
}

func isMatchTopic(s string, p string) bool {
	if s == p {
		return true
	}
	ss := strings.Split(s, ".")
	sp := strings.Split(p, ".")

	i, j := 0, 0
	for i = 0; i < len(ss); {
		if j < len(sp) && (ss[i] == sp[j] || sp[j] == "*") {
			i++
			j++
		} else if j < len(sp) && sp[j] == ">" {
			return true
		} else {
			return false
		}
	}
	return j == len(sp)
}

func (t *TopicFilt) ReplyTo(head StdHead, content interface{}) {
	t.client.Publish(head.Src, "rsp"+head.Topic[3:], content)
}

func (t *TopicFilt) ReplyAs(role string, head StdHead, content interface{}) {
	t.client.PublishAs(role, head.Src, "rsp"+head.Topic[3:], content)
}
