package go_config_center

import (
	"context"
	"errors"
	"log"
	"sync/atomic"
	"time"

	"github.com/samuel/go-zookeeper/zk"
)

type ZkClient struct {
	zkServers []string
	conn      *zk.Conn
	zkRoot    string             // 服务根节点，这里是/
	status    int32              //client 状态 1 runing
	cancel    context.CancelFunc //手动关闭ctx
	cancelCtx context.Context    //手动关闭ctx
	callback  func([]byte)       //TODO
}

func NewClient(zkServers []string, zkRoot string, timeout int) (*ZkClient, error) {
	client := new(ZkClient)
	client.zkServers = zkServers
	client.zkRoot = zkRoot
	client.cancelCtx, client.cancel = context.WithCancel(context.Background())
	client.status = 1
	// 连接服务器
	// option := zk.WithEventCallback(client.EventCallback)
	// conn, _, err := zk.Connect(zkServers, time.Duration(timeout)*time.Second, option)
	conn, _, err := zk.Connect(zkServers, time.Duration(timeout)*time.Second)
	if err != nil {
		return nil, err
	}
	log.Println("zookeeper connetion ok")
	client.conn = conn
	// 创建服务根节点
	if err := client.ensureRoot(); err != nil {
		client.Close()
		return nil, err
	}
	log.Println("ensureRoot ok")
	return client, nil
}
func (s *ZkClient) EventCallback(ch_event zk.Event) {
	log.Println("============================EventCallback", ch_event.Type)
	if ch_event.Type == zk.EventNodeCreated {
		//TODO
	} else if ch_event.Type == zk.EventNodeDeleted {
		//TODO
	} else if ch_event.Type == zk.EventNodeDataChanged {
		data, err := s.GetData(ch_event.Path)
		if err == nil {
			s.callback(data)
		} else {
			//alter
			log.Println(err)
		}
	} else if ch_event.Type == zk.EventNodeChildrenChanged {
		//todo
	}
}
func (s *ZkClient) Close() {
	s.cancel()
	atomic.StoreInt32(&s.status, 0)
	if s.conn != nil {
		s.conn.Close()
	}
}

func (s *ZkClient) ensureRoot() error {
	exists, _, err := s.conn.Exists(s.zkRoot)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("ErrNodeExists")
	}
	return nil
}

func (s *ZkClient) ZkWatch(path string, callback func([]byte)) {
	log.Println("Watch path:", path)
Loop:
	for atomic.LoadInt32(&s.status) == 1 {
		v, _, get_ch, err := s.conn.GetW(path)
		if err != nil {
			log.Println(err)
			time.Sleep(10 * time.Second) //5s 之后
			goto Loop
		}
		callback(v)
		select {
		case ch_event := <-get_ch:
			{
				log.Println("watch callback path:", ch_event.Path, "event_type:", ch_event.Type)
				if ch_event.Type == zk.EventNodeCreated {
					//TODO
				} else if ch_event.Type == zk.EventNodeDeleted {
					//TODO
				} else if ch_event.Type == zk.EventNodeDataChanged {
					data, err := s.GetData(ch_event.Path)
					if err == nil {
						callback(data)
					} else {
						//alter
						log.Println(err)
					}
				} else if ch_event.Type == zk.EventNodeChildrenChanged {
					//todo
				}
			}

		case <-s.cancelCtx.Done():
			break
		}
	}
}

func (s *ZkClient) GetData(path string) ([]byte, error) {
	data, _, err := s.conn.Get(path)
	return data, err
}
