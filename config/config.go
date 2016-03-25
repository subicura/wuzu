package config

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

var DefaultFilename = ".wuzu.yml"
var defaultContent = `# wuzu config file

build:
  from: 'golang:1.6'
  src: $PWD
  dest: /go/src/github.com/subicura/wuzu
  run: go build -v

`

type WuzuConfig struct {
	Build BuildConfig
}

type BuildConfig struct {
	From string
	Src  string
	Dest string
	Run  string
}

func (c *WuzuConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var conf struct {
		Build map[string]interface{}
	}
	if err := unmarshal(&conf); err != nil {
		return err
	}

	c.Build.From = os.ExpandEnv(conf.Build["from"].(string))
	c.Build.Src = os.ExpandEnv(conf.Build["src"].(string))
	c.Build.Dest = os.ExpandEnv(conf.Build["dest"].(string))
	c.Build.Run = os.ExpandEnv(conf.Build["run"].(string))
	return nil
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func CreateConfigFile() {
	if _, err := os.Stat(DefaultFilename); os.IsNotExist(err) {
		f, err := os.Create(DefaultFilename)
		if err != nil {
			panic(err)
		}
		defer f.Close()

		f.WriteString(defaultContent)
		f.Sync()

		log.Infoln(fmt.Sprintf("%s file is created.", DefaultFilename))
	} else {
		log.Warnln(fmt.Sprintf("%s file is already exist.", DefaultFilename))
	}
}

func Parse(filename string) (WuzuConfig, error) {
	data, err := ioutil.ReadFile(filename)
	checkErr(err)

	wuzuConfig := WuzuConfig{}
	err = yaml.Unmarshal([]byte(data), &wuzuConfig)
	checkErr(err)

	return wuzuConfig, err
}
