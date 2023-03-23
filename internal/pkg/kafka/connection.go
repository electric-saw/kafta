package kafka

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
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

func (k *KafkaConnection) newTLSConfig() (*tls.Config, error) {
	// #nosec
	tlsConfig := tls.Config{
		InsecureSkipVerify: true, // lgtm [go/disabled-certificate-check]
	}

	// Load CA cert
	if len(k.Context.TLS.CaCertFile) > 0 {
		caCertBlock, _ := base64.URLEncoding.DecodeString(k.Context.TLS.CaCertFile)

		if caCertBlock != nil {
			caCertPool := x509.NewCertPool()
			caCertPool.AppendCertsFromPEM(caCertBlock)
			tlsConfig.RootCAs = caCertPool
		}
	}

	// Load client cert
	if len(k.Context.TLS.ClientCertFile) > 0 || len(k.Context.TLS.ClientKeyFile) > 0 {
		clientCert, _ := base64.URLEncoding.DecodeString(k.Context.TLS.ClientCertFile)
		clientKey, _ := base64.URLEncoding.DecodeString(k.Context.TLS.ClientKeyFile)
		cert, err := tls.X509KeyPair(clientCert, clientKey)
		if err != nil {
			return &tlsConfig, err
		}
		tlsConfig.Certificates = []tls.Certificate{cert}
	}

	return &tlsConfig, nil
}

func (k *KafkaConnection) Connect() error {
	clientConfig := sarama.NewConfig()

	clientConfig.Net.DialTimeout = k.Config.ConnectionConfig().DialTimeout
	clientConfig.Net.ReadTimeout = k.Config.ConnectionConfig().ReadTimeout
	clientConfig.Net.WriteTimeout = k.Config.ConnectionConfig().WriteTimeout

	err := k.initAuth(clientConfig)
	if err != nil {
		return err
	}

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

func (k *KafkaConnection) initAuth(clientConfig *sarama.Config) error {
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
	}

	if k.Context.UseTLS {

		tlsConfig, err := k.newTLSConfig()
		if err != nil {
			return err
		}

		clientConfig.Net.TLS.Enable = true
		clientConfig.Net.TLS.Config = tlsConfig
	}

	return nil
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
