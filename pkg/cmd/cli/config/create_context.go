package config

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	b64 "encoding/base64"

	cliflag "github.com/electric-saw/kafta/internal/pkg/flag"
	"github.com/electric-saw/kafta/internal/pkg/kafka"

	"github.com/electric-saw/kafta/internal/pkg/ui"

	"github.com/Shopify/sarama"
	"github.com/electric-saw/kafta/internal/pkg/configuration"
	cmdutil "github.com/electric-saw/kafta/pkg/cmd/util"
	"github.com/spf13/cobra"
)

type createContextOptions struct {
	config               *configuration.Configuration
	name                 string
	currContext          bool
	schemaRegistry       cliflag.StringFlag
	schemaRegistrySecret cliflag.StringFlag
	schemaRegistryKey    cliflag.StringFlag
	ksql                 cliflag.StringFlag
	bootstrapServers     cliflag.StringFlag
	version              cliflag.StringFlag
	user                 cliflag.StringFlag
	password             cliflag.StringFlag
	useSASL              cliflag.StringFlag
	algorithm            cliflag.StringFlag
	useTLS               cliflag.StringFlag
	clientCertFile       cliflag.StringFlag
	clientKeyFile        cliflag.StringFlag
	caCertFile           cliflag.StringFlag
	parsedVersion        string
	quiet                bool
}

var (
	createContextLong = `
		Sets a context entry in config

		Specifying a name that already exists will merge new fields on top of existing values for those fields.`

	createContextExample = `
		# Set the cluster field on the kafka-dev context entry without touching other values
		kafta config set-context kafka-dev --server=b-1.kafka.example.com,b-2.kafka.example.com,b-3.kafka.example.com`
)

func NewCmdConfigSetContext(config *configuration.Configuration) *cobra.Command {
	options := &createContextOptions{config: config}

	cmd := &cobra.Command{
		Use:                   "set-context [NAME | --current] [--server=server] [--cluster=cluster_nickname] [--schema-registry=url] [--ksql=url]",
		DisableFlagsInUseLine: true,
		Short:                 "Sets a context entry in config",
		Long:                  createContextLong,
		Example:               createContextExample,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(options.complete(cmd))
			name, exists, err := options.run()
			cmdutil.CheckErr(err)
			if exists {
				fmt.Printf("Context %q modified.\n", name)
			} else {
				fmt.Printf("Context %q created.\n", name)
			}
		},
	}

	cmd.Flags().BoolVar(&options.currContext, "current", options.currContext, "Modify the current context")
	cmd.Flags().Var(&options.schemaRegistry, "schema-registry", "schema-registry for the context")
	cmd.Flags().Var(&options.schemaRegistrySecret, "schema-registry-secret", "schema-registry secret")
	cmd.Flags().Var(&options.schemaRegistryKey, "schema-registry-key", "schema-registry key")
	cmd.Flags().Var(&options.ksql, "ksql", "ksql for the context")
	cmd.Flags().Var(&options.bootstrapServers, "server", "server for the cluster entry in Kaftaconfig")
	cmd.Flags().Var(&options.version, "version", "kafka vesion for the cluster entry in Kaftaconfig")
	cmd.Flags().Var(&options.useSASL, "sasl", "Use SASL")
	cmd.Flags().VarP(&options.algorithm, "algorithm", "a", "algorithm for SASL")
	cmd.Flags().VarP(&options.user, "username", "u", "Username")
	cmd.Flags().VarP(&options.password, "password", "p", "Password")
	cmd.Flags().Var(&options.useTLS, "tls", "Use TLS")
	cmd.Flags().VarP(&options.clientCertFile, "clientCertFile", "c", "ClientCertFile")
	cmd.Flags().VarP(&options.clientKeyFile, "clientKeyFile", "k", "ClientKeyFile")
	cmd.Flags().VarP(&options.caCertFile, "caCertFile", "f", "CaCertFile")

	return cmd
}

func (o *createContextOptions) run() (string, bool, error) {
	err := o.validate()
	if err != nil {
		return "", false, err
	}

	name := o.name
	if o.currContext {
		if len(o.config.KaftaData.CurrentContext) == 0 {
			return "", false, errors.New("no current context is set")
		}
		name = o.config.KaftaData.CurrentContext
	}

	startingInstance, exists := o.config.KaftaData.Contexts[name]
	if !exists {
		startingInstance = configuration.MakeContext()
	}
	cmdutil.CheckErr(o.promptNeeded(startingInstance))

	context, err := o.modifyContext(*startingInstance)
	if err != nil {
		cmdutil.CheckErr(err)
		return "", false, fmt.Errorf("could not extract TLS configuration")
	}

	fmt.Printf("\n\n%v\n\n", context)

	err = o.checkConnection(context)
	if err != nil {
		cmdutil.CheckErr(err)
		return "", false, fmt.Errorf("could not connect to %s", context.BootstrapServers)
	}

	o.config.KaftaData.Contexts[name] = context

	return name, exists, nil
}

//gocyclo:ignore
func (o *createContextOptions) modifyContext(context configuration.Context) (*configuration.Context, error) {
	modifiedContext := context

	if o.ksql.Provided() {
		modifiedContext.Ksql = o.ksql.Value()
	}

	if o.schemaRegistry.Provided() {
		modifiedContext.SchemaRegistry = o.schemaRegistry.Value()
	}

	if o.schemaRegistrySecret.Provided() {
		modifiedContext.SchemaRegistryAuth.Secret = o.schemaRegistrySecret.Value()
	}

	if o.schemaRegistryKey.Provided() {
		modifiedContext.SchemaRegistryAuth.Key = o.schemaRegistryKey.Value()
	}

	if o.bootstrapServers.Provided() {
		modifiedContext.BootstrapServers = strings.Split(o.bootstrapServers.Value(), ",")
	}

	if len(o.parsedVersion) > 0 {
		version, err := sarama.ParseKafkaVersion(o.parsedVersion)
		cmdutil.CheckErr(err)
		modifiedContext.KafkaVersion = version.String()
	}

	if o.useSASL.Provided() {
		modifiedContext.UseSASL = true
	}

	if o.algorithm.Provided() {
		modifiedContext.SASL.Algorithm = o.algorithm.Value()
	}

	if o.user.Provided() {
		modifiedContext.SASL.Username = o.user.Value()
	}

	if o.password.Provided() {
		modifiedContext.SASL.Password = o.password.Value()
	}

	if o.useTLS.Provided() {
		modifiedContext.UseTLS = true
	}

	if o.clientCertFile.Provided() {
		contentAsB64, err := extractContentToBase64(o.clientCertFile.Value())
		if err != nil {
			return nil, err
		}
		modifiedContext.TLS.ClientCertFile = contentAsB64
	}

	if o.clientKeyFile.Provided() {
		contentAsB64, err := extractContentToBase64(o.clientKeyFile.Value())
		if err != nil {
			return nil, err
		}
		modifiedContext.TLS.ClientKeyFile = contentAsB64
	}

	if o.caCertFile.Provided() {
		contentAsB64, err := extractContentToBase64(o.caCertFile.Value())
		if err != nil {
			return nil, err
		}
		modifiedContext.TLS.CaCertFile = contentAsB64
	}

	return &modifiedContext, nil
}

func extractContentToBase64(filepath string) (string, error) {
	content, err := os.ReadFile(filepath)
	if err != nil {
		return "", err
	}

	enc := b64.URLEncoding.EncodeToString(content)
	return enc, nil
}

func (o *createContextOptions) complete(cmd *cobra.Command) error {
	args := cmd.Flags().Args()
	if len(args) > 1 {
		return cmdutil.HelpErrorf(cmd, "Unexpected args: %v", args)
	}
	if len(args) == 1 {
		o.name = args[0]
	}

	if o.version.Provided() {
		o.parsedVersion = o.version.Value()
	}

	return nil
}

func (o *createContextOptions) validate() error {
	if len(o.name) == 0 && !o.currContext {
		return errors.New("you must specify a non-empty context name or --current")
	}
	if len(o.name) > 0 && o.currContext {
		return errors.New("you cannot specify both a context name and --current")
	}

	if o.ksql.Provided() && !testHost(o.ksql.String()) {
		return errors.New("failed to connect on ksql")
	}

	// if o.schemaRegistry.Provided() && !testHost(o.schemaRegistry.String()) {
	// 	return errors.New("failed to connect on schema-registry")
	// }

	if o.useSASL.Provided() && o.quiet {
		if !o.user.Provided() {
			return errors.New("user flag is required if SASL is provided")
		}

		if !o.password.Provided() {
			return errors.New("user flag is required if SASL is provided")
		}
	}

	if o.useTLS.Provided() && o.quiet {
		if !o.clientCertFile.Provided() {
			return errors.New("clientCertFile is required if TLS is provided")
		}

		if !o.clientKeyFile.Provided() {
			return errors.New("clientKeyFile is required if TLS is provided")
		}
	}

	return nil
}

func testHost(address string) bool {
	if len(strings.Split(address, ":")) <= 1 {
		fmt.Printf("Port is nedeed on %s!\n", address)
		return false
	}
	conn, err := net.DialTimeout("tcp", address, 3*time.Second)
	if err != nil {
		return false
	} else {
		if conn != nil {
			_ = conn.Close()
			return true
		} else {
			return true
		}
	}
}

func (o *createContextOptions) checkConnection(context *configuration.Context) error {
	if o.bootstrapServers.Provided() {
		conn, err := kafka.MakeConnectionContext(o.config, context)
		if err != nil {
			return err
		}
		conn.Close()
	}

	return nil
}

//gocyclo:ignore
func (o *createContextOptions) promptNeeded(context *configuration.Context) error {
	if o.quiet {
		return nil
	}

	if !o.bootstrapServers.Provided() && len(context.BootstrapServers) == 0 {
		servers, err := ui.GetText("Bootstrap servers", true)
		cmdutil.CheckErr(err)
		err = o.bootstrapServers.Set(servers)
		cmdutil.CheckErr(err)
	}

	if !o.version.Provided() && len(context.KafkaVersion) == 0 {
		version, err := ui.GetText("Kafka version", true)
		cmdutil.CheckErr(err)
		err = o.version.Set(version)
		cmdutil.CheckErr(err)
	}

	if !o.useSASL.Provided() {
		sasl, err := ui.GetConfirmation("Use SASL", false)
		cmdutil.CheckErr(err)
		if sasl {
			err := o.useSASL.Set("true")
			cmdutil.CheckErr(err)

			if !o.algorithm.Provided() && len(context.SASL.Algorithm) == 0 {
				algorithm, err := ui.GetText("SASL Algorithm", true)
				cmdutil.CheckErr(err)
				err = o.algorithm.Set(algorithm)
				cmdutil.CheckErr(err)
			}

			if !o.user.Provided() && len(context.SASL.Username) == 0 {
				user, err := ui.GetText("User", true)
				cmdutil.CheckErr(err)
				err = o.user.Set(user)
				cmdutil.CheckErr(err)
			}

			if !o.password.Provided() && len(context.SASL.Password) == 0 {
				password, err := ui.GetPassword("Password", true)
				cmdutil.CheckErr(err)
				err = o.password.Set(password)
				cmdutil.CheckErr(err)
			}

		}
	}

	if !o.useTLS.Provided() {
		tls, err := ui.GetConfirmation("Use TLS", false)
		cmdutil.CheckErr(err)

		if tls {
			err := o.useTLS.Set("true")
			cmdutil.CheckErr(err)

			if !o.clientCertFile.Provided() && len(context.TLS.ClientCertFile) == 0 {
				clientCertFile, err := ui.GetText("TLS ClientCertFile", true)
				cmdutil.CheckErr(err)
				err = o.clientCertFile.Set(clientCertFile)
				cmdutil.CheckErr(err)
			}

			if !o.clientKeyFile.Provided() && len(context.TLS.ClientKeyFile) == 0 {
				clientKeyFile, err := ui.GetText("TLS ClientKeyFile", true)
				cmdutil.CheckErr(err)
				err = o.clientKeyFile.Set(clientKeyFile)
				cmdutil.CheckErr(err)
			}

			if !o.caCertFile.Provided() && len(context.TLS.CaCertFile) == 0 {
				caCertFile, err := ui.GetText("CaCertFile", true)
				cmdutil.CheckErr(err)
				err = o.caCertFile.Set(caCertFile)
				cmdutil.CheckErr(err)
			}
		}
	}

	if !o.schemaRegistry.Provided() {
		schemaRegistry, err := ui.GetText("Schema registry", false)
		cmdutil.CheckErr(err)
		err = o.schemaRegistry.Set(schemaRegistry)
		cmdutil.CheckErr(err)

		if o.schemaRegistry.Provided() && len(context.SchemaRegistryAuth.Secret) == 0 {
			secret, err := ui.GetText("SchemaRegistry Secret", false)
			cmdutil.CheckErr(err)
			err = o.schemaRegistrySecret.Set(secret)
			cmdutil.CheckErr(err)
		}

		if o.schemaRegistry.Provided() && len(context.SchemaRegistryAuth.Key) == 0 {
			key, err := ui.GetPassword("SchemaRegistry Key", false)
			cmdutil.CheckErr(err)
			err = o.schemaRegistryKey.Set(key)
			cmdutil.CheckErr(err)
		}

	}

	return nil
}
