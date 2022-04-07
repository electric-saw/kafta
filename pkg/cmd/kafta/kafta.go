package kafta

import (
	"github.com/electric-saw/kafta/internal/pkg/configuration"
	"github.com/electric-saw/kafta/pkg/cmd/cli/broker"
	"github.com/electric-saw/kafta/pkg/cmd/cli/cluster"
	"github.com/electric-saw/kafta/pkg/cmd/cli/completion"
	configCmd "github.com/electric-saw/kafta/pkg/cmd/cli/config"
	"github.com/electric-saw/kafta/pkg/cmd/cli/consumer"
	"github.com/electric-saw/kafta/pkg/cmd/cli/schema"
	"github.com/electric-saw/kafta/pkg/cmd/cli/topic"
	"github.com/electric-saw/kafta/pkg/cmd/cli/version"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewKaftaCommand(name string) *cobra.Command {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	config := configuration.InitializeConfiguration(name)

	root := &cobra.Command{
		Use:   name,
		Short: "Command line interface for automate process and administration in Kafka clusters",
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			config.UpdateConfig()
		},
	}

	config.BindFlags(root)
	config.EnsureKaftaconfig()

	root.AddCommand(version.NewCmdVersion(config))
	root.AddCommand(topic.NewCmdTopic(config))
	root.AddCommand(configCmd.NewCmdConfig(config))
	root.AddCommand(completion.NewCmdCompletion(config))
	root.AddCommand(broker.NewCmdBroker(config))
	root.AddCommand(consumer.NewCmdConsumer(config))
	root.AddCommand(cluster.NewCmdCluster(config))
	root.AddCommand(schema.NewCmdSchema(config))

	return root
}
