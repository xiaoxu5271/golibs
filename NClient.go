package golibs

type NClient struct {
	NatsWraper
	TopicFilt
}

// default name=nclient, default url=localhost:6000
func (n *NClient) Init(server, mySubj, myName string, capacity int) *NClient {
	n.TopicFilt.Init(n)
	n.NatsWraper.Init(server, mySubj, myName, capacity, n)
	go n.requestService()
	return n
}

func (n *NClient) Connect() {
	n.connect()
}

func (n *NClient) GetMySubj() string {
	return n.mySubj
}
