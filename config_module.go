package go_config_center

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/spf13/viper"
)

var DEFAULT_LOCAL_CACHE_DIR = "./config_center"

type ConfigModule struct {
	cfgFile       string
	buf           []byte
	cfg           *viper.Viper
	lock          sync.RWMutex
	localCacheDir string
	configType    string
}

func NewConfigModule(modulePath string, localCacheDir string, configType string) *ConfigModule {
	c := &ConfigModule{}
	if localCacheDir == "" {
		c.localCacheDir = DEFAULT_LOCAL_CACHE_DIR
	} else {
		c.localCacheDir = localCacheDir
	}
	if strings.HasPrefix(modulePath, "/") {
		c.cfgFile = c.localCacheDir + modulePath
	} else {
		c.cfgFile = c.localCacheDir + "/" + modulePath
	}
	dir := filepath.Dir(c.cfgFile)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		log.Println("mkdir:", dir)
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			log.Println("dir:", dir, err)
		}
	}
	log.Println("mkdir:", dir)
	c.cfg = viper.New()
	if configType == "" {
		c.configType = "json"
	} else {
		c.configType = configType
	}
	suffix := path.Ext(modulePath)[1:]
	if suffix != "" {
		c.configType = suffix
	} //直接使用后缀
	c.cfg.SetConfigType(c.configType)
	return c
}

// 存盘
func (c *ConfigModule) loadFromBuf(data []byte) error {
	// log.Println("data:", string(data))
	c.buf = make([]byte, len(data))
	copy(c.buf, data)
	c.lock.Lock()
	cfg := viper.New()
	cfg.SetConfigType(c.configType)
	err := cfg.ReadConfig(bytes.NewReader(c.buf))
	if err != nil {
		c.lock.Unlock()
		return err
	}
	c.cfg = cfg
	c.lock.Unlock()
	log.Println("config from loadFromBuf load ok")
	f, err := ioutil.TempFile(c.localCacheDir, "tmp")
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(data)
	if err != nil {
		return err
	}
	os.Rename(f.Name(), c.cfgFile)
	return nil
}

// 读到buf
func (c *ConfigModule) loadFromLocalCache() error {
	data, err := ioutil.ReadFile(c.cfgFile)
	if err != nil {
		return err
	}
	c.buf = make([]byte, len(data))
	copy(c.buf, data)
	c.lock.Lock()
	cfg := viper.New()
	cfg.SetConfigType(c.configType)
	err = cfg.ReadConfig(bytes.NewReader(c.buf))
	if err != nil {
		c.lock.Unlock()
		return err
	}
	c.cfg = cfg
	c.lock.Unlock()
	log.Println("config from loadFromLocalCache load ok")
	return nil
}

// watch回调
func (c *ConfigModule) onModuleChange(data []byte) {
	err := c.loadFromBuf(data)
	if err != nil {
		log.Println(err)
	}
}

// func (c *ConfigModule) GetCfg() *viper.Viper {
// 	c.lock.RLock()
// 	defer c.lock.RUnlock()
// 	return c.cfg
// }

func (c *ConfigModule) Get(key string) interface{} {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.cfg.Get(key)
}

func (c *ConfigModule) GetInt(key string) int {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.cfg.GetInt(key)
}

func (c *ConfigModule) GetString(key string) string {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.cfg.GetString(key)
}

func (c *ConfigModule) GetStringSlice(key string) []string {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.cfg.GetStringSlice(key)
}

func (c *ConfigModule) GetBool(key string) bool {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.cfg.GetBool(key)
}

func (c *ConfigModule) GetFloat64(key string) float64 {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.cfg.GetFloat64(key)
}

func (c *ConfigModule) GetStringMapString(key string) map[string]string {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.cfg.GetStringMapString(key)
}

func (c *ConfigModule) GetStringMap(key string) map[string]interface{} {
	c.lock.RLock()
	defer c.lock.RUnlock()

	return c.cfg.GetStringMap(key)
}

func (c *ConfigModule) GetAll() map[string]interface{} {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.cfg.AllSettings()
}

func (c *ConfigModule) Raw() []byte {
	c.lock.RLock()
	defer c.lock.RUnlock()
	dest := make([]byte, len(c.buf))
	copy(dest, c.buf)
	return dest
}
