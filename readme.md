# 基于zookeeper&viper 的 golang 动态配置使用库 (JSON格式) 

### 此库通过zookeeper获取json 配置文件,viper读取配置内容,当zookeeper内容变化时，viper内容自动更新(协程序安全);并且本地缓存配置文件，当zookeeper 不可用时，读取本地缓存 保证服务可用

1. 使用方法
``` golang
center := NewConfigCenter("", nil, "/usr/local/vntop/config_center/local_cache","json")
module := center.GetModule("/test.json")
fmt.Println("--------------------a:", module.GetInt("a"))
```

2. TODO
- [*] 支持多种格式 例如 xml,yaml等
- [] 支持 etcd, consul 配置中心
- [] 支持 自定义watch ,替换viper 自行监控数据内容变化 并且自行解析