package rboot

import "sync"

// Node 为记录多机部署时的每个节点的信息
type Node struct {
	// 节点ID
	NodeID string
	// 节点IP地址
	NodeAddr string
	// 节点的Http端口，用于数据传输时使用
	HttpPort string
}

// Bucket 记录了当前网络中存在的所有节点
type Bucket struct {
	Node map[string]Node
	mu   sync.RWMutex
}

func NewBucket() *Bucket {
	return &Bucket{
		Node: make(map[string]Node),
		mu:   sync.RWMutex{},
	}
}

// Len 返回节点长度
func (b *Bucket) Len() int {
	return len(b.Node)
}

// Set 将节点添加到Bucket，如果已经存在则覆盖
func (b *Bucket) Set(n Node) {
	b.mu.Lock()
	b.Node[n.NodeID] = n
	b.mu.Unlock()
}

// Get 根据节点ID获取节点信息
func (b *Bucket) Get(id string) Node {
	b.mu.Lock()
	node := b.Node[id]
	b.mu.Unlock()
	return node
}

// Remove 根据节点ID删除节点
func (b *Bucket) Remove(id string) {
	b.mu.Lock()
	delete(b.Node, id)
	b.mu.Unlock()
}

// Clear 清除所有节点
func (b *Bucket) Clear() {
	b.Node = make(map[string]Node)
}

/*var udpBuffSize = 100

// 响应方法
func (n *Node) pong(conn *net.UDPConn) {
	// 缓冲区
	var buff = make([]byte, udpBuffSize)

	num, addr, err := conn.ReadFromUDP(buff)
	if err != nil {
		log.Error(err)
	}

	if num > 0 {
		pingV := strings.Split(string(buff), ":")
		switch pingV[0] {
		case `ping`:
			// 当机器人第一次启动时发送“ping”命令
			// 数据应该存在三部分，分别是命令“ping”，机器人ID和Http端口，使用“:”隔开
			n.NodeAddr = addr.IP.String()
			if len(pingV) >= 2 {
				n.NodeID = pingV[1]
				n.HttpPort = pingV[2]
				// 成功，返回pong响应
				conn.WriteToUDP([]byte("pong"), addr)
			} else {
				// 失败返回
				conn.WriteToUDP([]byte("failed"), addr)
			}
		case `pong`:
			// 心跳，返回pong
			conn.WriteToUDP([]byte("pong"), addr)
		}

	}
}

// ping 命令
func (n *Node) ping(conn net.Conn, id string) {
	defer conn.Close()

	msg := fmt.Sprintf("ping:%s:%s", id, os.Getenv("WEB_SERVER_PORT"))

	_, err := conn.Write([]byte(msg))
	if err != nil {
		log.Fatal(err)
		return
	}

	resp := make([]byte, 10)

	_, err = conn.Read(resp)
	if err != nil {
		log.Fatal(err)
		return
	}

	if string(resp) == "pong" {
		log.Info("成功连接集群网络...")
	} else if string(resp) == "failed" {
		log.Fatal("连接集群网络失败，请检查数据完整性...")
	}
}

func (n *Node) udpServer() {
	// 获取 web 端口
	port := os.Getenv("UDP_SERVER_PORT")
	if port == "" {
		port = "7866"
	}

	address := ":" + port
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for {
		n.pong(conn)
	}
}*/
