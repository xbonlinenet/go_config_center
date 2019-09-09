package go_config_center

import (
	"context"
	"fmt"
	"log"
	"sync"
)

type ConfigCenter struct {
	zk            *ZkClient
	CfgModules    map[string]*ConfigModule
	lock          sync.Mutex
	ctx           context.Context
	zkRoot        string
	zkServers     []string
	localCacheDir string
	configType    string
}

// configType 默认json
func NewConfigCenter(zkRoot string, zkServers []string, localCacheDir string, configType string) *ConfigCenter {
	log.SetFlags(log.Llongfile | log.Ltime)
	center := new(ConfigCenter)
	center.initZk()
	center.lock = sync.Mutex{}
	center.CfgModules = make(map[string]*ConfigModule)
	if zkRoot == "" {
		center.zkRoot = "/config_center"
	} else {
		center.zkRoot = zkRoot
	}
	if zkServers == nil || len(zkServers) == 0 {
		center.zkServers = []string{"127.0.0.1:2181"}
	} else {
		center.zkServers = zkServers
	}
	center.localCacheDir = localCacheDir
	center.configType = configType
	return center
}

func (c *ConfigCenter) Close() {
	if c.zk != nil {
		c.zk.Close()
	}
}

func (c *ConfigCenter) initZk() error {
	var err error
	c.zk, err = NewClient(c.zkServers, c.zkRoot, 3)
	return err
}

func (c *ConfigCenter) getModuleZkPath(module_path string) string {
	return fmt.Sprintf("/config_center%s", module_path)
}

func (c *ConfigCenter) GetModule(modulePath string) *ConfigModule {
	c.lock.Lock()
	defer c.lock.Unlock()
	if module, ok := c.CfgModules[modulePath]; ok {
		return module
	}
	if c.zk == nil {
		err := c.initZk()
		if err != nil {
			log.Println(err)
		}
	}
	var err error
	var data []byte
	zkPath := c.getModuleZkPath(modulePath)
	module := NewConfigModule(modulePath, c.localCacheDir, c.configType)
	if c.zk != nil {
		data, err = c.zk.GetData(zkPath)
		if err != nil {
			log.Println(err)
			//alter
		}
		if data != nil && len(data) > 0 {
			err := module.loadFromBuf(data)
			if err != nil {
				log.Println(err)
			}
		}
	}
	if data == nil || len(data) == 0 {
		err := module.loadFromLocalCache()
		if err != nil {
			log.Println(err)
		}
	}
	c.CfgModules[modulePath] = module

	if c.zk != nil {
		go c.zk.ZkWatch(zkPath, module.onModuleChange)
	}
	return module
}
