package go_config_center

import (
	"fmt"
	"io/ioutil"
	"testing"
	"time"
)

//go test -v -run TestA
func TestA(t *testing.T) {
	center := NewConfigCenter("", nil, "/usr/local/vntop/config_center/local_cache", "json")
	module := center.GetModule("/gateway/test.yaml")
	fmt.Println("============module================")
	fmt.Println("--------------------a:", module.GetInt("a"))
	fmt.Println("--------------------b:", module.GetInt("b"))
	time.Sleep(10 * time.Second)
	fmt.Println("--------------------a:", module.GetInt("a"))
	center.Close()
	time.Sleep(10 * time.Second)
}

func TestTempFile(t *testing.T) {
	f, err := ioutil.TempFile(DEFAULT_LOCAL_CACHE_DIR, "tmp")
	f.Close()
	fmt.Println(err)
}

//go test -v -bench="."  -benchtime=10s -run=BenchmarkGetInt
func BenchmarkGetInt(b *testing.B) {
	center := NewConfigCenter("", nil, "/usr/local/vntop/config_center/local_cache", "json")
	module := center.GetModule("/gateway/test.json")
	fmt.Println("============module================")
	for i := 0; i < b.N; i++ {
		fmt.Println("--------------------a:", module.GetInt("a"))
	}
}
