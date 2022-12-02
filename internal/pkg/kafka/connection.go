package kafka

import (
	"crypto/tls"
	"fmt"

	"github.com/Shopify/sarama"
	"github.com/electric-saw/kafta/internal/pkg/configuration"
	"github.com/electric-saw/kafta/pkg/cmd/util"
	"github.com/riferrei/srclient"
)

type KafkaConnection struct {
	Client               sarama.Client
	Admin                sarama.ClusterAdmin
	Config               *configuration.Configuration
	Context              *configuration.Context
	SchemaRegistryClient *srclient.SchemaRegistryClient
}

func MakeConnection(config *configuration.Configuration) *KafkaConnection {
	conn, err := MakeConnectionContext(config, config.GetContext())
	util.CheckErr(err)

	return conn
}

func MakeConnectionContext(config *configuration.Configuration, context *configuration.Context) (*KafkaConnection, error) {
	conn := &KafkaConnection{
		Config:  config,
		Context: context,
	}

	err := conn.Connect()
	return conn, err
}

func (k *KafkaConnection) Connect() error {
	clientConfig := sarama.NewConfig()

	clientConfig.Net.DialTimeout = k.Config.ConnectionConfig().DialTimeout
	clientConfig.Net.ReadTimeout = k.Config.ConnectionConfig().ReadTimeout
	clientConfig.Net.WriteTimeout = k.Config.ConnectionConfig().WriteTimeout

	k.initAuth(clientConfig)

	clientConfig.Version = k.Context.GetVersion()
	client, err := sarama.NewClient(
		k.Context.BootstrapServers,
		clientConfig,
	)

	if err != nil {
		return err
	}

	admin, err := sarama.NewClusterAdminFromClient(client)

	if err != nil {
		return err
	}

	k.Client = client
	k.Admin = admin

	return k.connectSr()
}

func (k *KafkaConnection) Close() {
	util.CheckErr(k.Client.Close())
}

func (k *KafkaConnection) initAuth(clientConfig *sarama.Config) {
	if k.Context.UseSASL {
		clientConfig.Net.SASL.Enable = true
		clientConfig.Net.SASL.User = k.Context.SASL.Username
		clientConfig.Net.SASL.Password = k.Context.SASL.Password
		switch k.Context.SASL.Algorithm {
		case "sha256":
			clientConfig.Net.SASL.SCRAMClientGeneratorFunc = func() sarama.SCRAMClient { return &XDGSCRAMClient{HashGeneratorFcn: SHA256} }
			clientConfig.Net.SASL.Mechanism = sarama.SASLTypeSCRAMSHA256
		case "sha512":
			clientConfig.Net.SASL.SCRAMClientGeneratorFunc = func() sarama.SCRAMClient { return &XDGSCRAMClient{HashGeneratorFcn: SHA512} }
			clientConfig.Net.SASL.Mechanism = sarama.SASLTypeSCRAMSHA512
		default:
			clientConfig.Net.SASL.Mechanism = sarama.SASLTypePlaintext
		}
		clientConfig.Net.SASL.Handshake = true

		if k.Context.TLS {
			clientConfig.Net.TLS.Enable = true
			tlsConfig := &tls.Config{
				InsecureSkipVerify: true, // lgtm [go/disabled-certificate-check]
				ClientAuth:         0,
			}

			clientConfig.Net.TLS.Config = tlsConfig
		}
	}
}

func (k *KafkaConnection) connectSr() error {
	if k.SchemaRegistryClient == nil && k.Context.SchemaRegistry != "" {
		srClient := srclient.CreateSchemaRegistryClient(k.Context.SchemaRegistry)
		srClient.SetCredentials(k.Context.SchemaRegistryAuth.Key, k.Context.SchemaRegistryAuth.Secret)
		_, err := srClient.GetGlobalCompatibilityLevel()
		if err != nil {
			return fmt.Errorf("error connecting to schema registry: %w", err)
		}
		k.SchemaRegistryClient = srClient
	}

	return nil
}
