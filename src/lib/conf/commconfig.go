package conf

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"

	"github.com/go-yaml/yaml"
)

const (
	filePath = "D:/golang/src/github.com/linnull/messaging/conf"
	fileName = "comm.yaml"
)

type CommConfig struct {
	sync.Mutex
	fileName string

	LogBasePath string `yaml:"log_base_path"`
	LogLevel    string `yaml:"log_level"`

	EtcdAddress []string `yaml:"etcd_address"`
}

func NewCommConfig() *CommConfig {
	conf := new(CommConfig)
	conf.fileName = fmt.Sprintf("%s/%s", filePath, fileName)
	return conf
}

func (c *CommConfig) load() error {
	file, err := os.Open(c.fileName)
	if err != nil {
		return err
	}
	defer file.Close()
	data, _ := ioutil.ReadAll(file)
	return yaml.Unmarshal(data, c)
}

func (c *CommConfig) reload() {
	file, err := os.Open(c.fileName)
	if err != nil {
		log.Printf("file open err:%v/n", err)
		return
	}
	defer file.Close()
	data, _ := ioutil.ReadAll(file)
	err = yaml.Unmarshal(data, c)
	if err != nil {
		log.Printf("yaml.Unmarshal err:%v/n", err)
		return
	}
	log.Printf("reload done config:%v/n", c)
}

var commConfig *CommConfig

func init() {
	commConfig = NewCommConfig()
	err := commConfig.load()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		var tick = time.NewTicker(time.Second * 60)
		defer tick.Stop()
		for {
			select {
			case <-tick.C:
				commConfig.reload()
			}
		}
	}()
}

func GetCommConfig() CommConfig {
	commConfig.Lock()
	defer commConfig.Unlock()
	return *commConfig
}
