package appctx

import (
	"bytes"
	"github.com/zibilal/datacabinet"
	"github.com/zibilal/datacabinet/persistence"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"sync"
)

const (
	DefaultConfigFlagVal = "config.yml"
)

type Config struct {
	Address  string `yaml:"address"`
	Mode     string `yaml:"mode"`
	Database struct {
		ConnectionString string `yaml:"connection_string"`
		Timeout          int    `yaml:"timeout"`
		Name             string `yaml:"name"`
	} `yaml:"database"`
}

type AppContext struct {
	Config        *Config
	DataConnector connector.Connector
	Persistence   persistence.Persistence
}

func NewAppContext() *AppContext {
	ctx := new(AppContext)
	ctx.Config = new(Config)
	return ctx
}

func (c *AppContext) LoadAppContext(readers ...io.Reader) error {
	buff := bytes.NewBuffer([]byte{})
	for _, reader := range readers {
		dat, err := ioutil.ReadAll(reader)
		if err != nil {
			return err
		}
		buff.Write(dat)
	}

	if err := yaml.Unmarshal(buff.Bytes(), c.Config); err != nil {
		return err
	}

	return nil
}

var (
	Instance *AppContext
	once     sync.Once
)

func GetAppContext() *AppContext {
	if Instance == nil {
		panic("Application Context has not been initialized yet")
	}
	return Instance
}
