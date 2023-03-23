package configuration

import (
	"errors"
	"os"
	"path"

	"github.com/electric-saw/kafta/pkg/cmd/util"
	"github.com/spf13/cobra"
)

const (
	DefaultFileName   = "config"
	DefaultFolderName = ".kafta"
)

type Configuration struct {
	DebugMode       bool
	AppName         string
	ActiveContext   string
	KaftaconfigFile string
	KaftaData       *KaftaConfig
}

func homeDir() string {
	home, err := os.UserHomeDir()
	util.CheckErr(err)
	return home
}

func InitializeConfiguration(appName string) *Configuration {
	configName := path.Join(homeDir(), DefaultFolderName, DefaultFileName)

	config := &Configuration{
		DebugMode:       false,
		AppName:         appName,
		KaftaconfigFile: configName,
	}

	config.EnsureKaftaconfig()

	return config
}

func (c *Configuration) BindFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolVarP(&c.DebugMode, "debug", "d", false, "Debug mode")
	cmd.PersistentFlags().StringVarP(&c.ActiveContext, "context", "", c.ActiveContext, "The name of the kafkaconfig context to use")
	cmd.PersistentFlags().StringVarP(&c.KaftaconfigFile, "kafkaconfig", "", c.KaftaconfigFile, "Path to the kafkaconfig file to use for CLI requests.")
}

func (c *Configuration) EnsureKaftaconfig() {
	config, isNew := LoadKaftaconfigOrDefault(c.KaftaconfigFile)
	if isNew {
		err := config.Write()
		util.CheckErr(err)
	}
	c.KaftaData = config
}

func (c *Configuration) UpdateConfig() {
	err := c.KaftaData.Write()
	util.CheckErr(err)
}

func (c *Configuration) GetContext() *Context {
	if len(c.ActiveContext) == 0 {
		c.ActiveContext = c.KaftaData.CurrentContext
	}

	if len(c.ActiveContext) == 0 || c.KaftaData.Contexts[c.ActiveContext] == nil {
		util.CheckErr(errors.New("no context found"))
	}
	return c.KaftaData.Contexts[c.ActiveContext]
}

func (c *Configuration) ConnectionConfig() *ConnectionConfig {
	return &c.KaftaData.Connection
}
