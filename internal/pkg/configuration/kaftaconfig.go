package configuration

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/Shopify/sarama"
	"github.com/electric-saw/kafta/pkg/cmd/util"
	"gopkg.in/yaml.v3"
)

type KaftaConfig struct {
	Contexts       map[string]*Context `yaml:"contexts"`
	CurrentContext string              `yaml:"current-context"`
	Connection     ConnectionConfig    `yaml:"config"`
	path           string              `yaml:"-"`
}

type ConnectionConfig struct {
	DialTimeout  time.Duration
	WriteTimeout time.Duration
	ReadTimeout  time.Duration
}

type Context struct {
	SchemaRegistry     string `yaml:"schema-registry"`
	SchemaRegistryAuth struct {
		Key    string
		Secret string
	}
	Ksql             string   `yaml:"ksql"`
	BootstrapServers []string `yaml:"bootstrap-servers"`
	KafkaVersion     string   `yaml:"kafka-version"`
	UseSASL          bool
	UseTLS           bool
	TLS              struct {
		ClientCertFile string
		ClientKeyFile  string
		CaCertFile     string
	}
	SASL struct {
		Algorithm string
		Username  string
		Password  string
	}
}

func MakeContext() *Context {
	return &Context{
		UseSASL:      false,
		KafkaVersion: sarama.V3_3_0_0.String(),
	}
}

func LoadKaftaconfigOrDefault(configPath string) (*KaftaConfig, bool) {
	config := &KaftaConfig{
		path:     configPath,
		Contexts: make(map[string]*Context),
		Connection: ConnectionConfig{
			DialTimeout:  15 * time.Second,
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 30 * time.Second,
		},
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return config, true
	} else {

		rawYaml, err := os.ReadFile(configPath)
		if err != nil {
			return config, true
		}

		err = yaml.Unmarshal(rawYaml, &config)

		util.CheckErr(err)
	}

	return config, false
}

func (k *KaftaConfig) Write() {
	err := os.MkdirAll(path.Dir(k.path), 0760)
	util.CheckErr(err)

	out, err := os.Create(k.path)
	util.CheckErr(err)
	defer out.Close()

	err = yaml.NewEncoder(out).Encode(k)
	util.CheckErr(err)
}

func (k *KaftaConfig) ConfigPath() string {
	return k.path
}

func (k *KaftaConfig) ConfigFileName() string {
	return filepath.Base(k.ConfigPath())
}

func (c *Context) GetVersion() sarama.KafkaVersion {
	if len(c.KafkaVersion) == 0 {
		msg := "No version found, please input with 'kafta config set-context NAME --version XXXXXX"
		fmt.Printf("WARN %s\n", msg)
		c.KafkaVersion = "3.3.0"
	}
	version, err := sarama.ParseKafkaVersion(c.KafkaVersion)
	util.CheckErr(err)
	return version
}
