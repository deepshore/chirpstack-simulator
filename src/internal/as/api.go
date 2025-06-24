package as

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"os"

	"github.com/chirpstack/chirpstack/api/go/v4/api"
	"github.com/deepshore/chirpstack-simulator/internal/config"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type jwtCredentials struct {
	token string
}

var clientConn *grpc.ClientConn
var mqttClient mqtt.Client

func (j *jwtCredentials) GetRequestMetadata(ctx context.Context, url ...string) (map[string]string, error) {
	return map[string]string{
		"authorization": "Bearer " + j.token,
	}, nil
}

func (j *jwtCredentials) RequireTransportSecurity() bool {
	return false
}

func Setup(cfg config.Config) error {
	conf := cfg.ChirpStack

	// conntet gRPC server
	log.Info("Connecting to ChirpStack Application Server gRPC server at", conf.API.Server)

	options := []grpc.DialOption{
		grpc.WithPerRPCCredentials(&jwtCredentials{token: conf.API.APIKey}),
	}

	if conf.API.Insecure {
		options = append(options, grpc.WithTransportCredentials(insecure.NewCredentials()))
	} else {
		log.Warn("Using TLS for ChirpStack Application Server gRPC server connection is not supported yet")
		options = append(options, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	conn, err := grpc.NewClient(
		conf.API.Server,
		grpc.WithPerRPCCredentials(&jwtCredentials{token: conf.API.APIKey}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.WithError(err).Error("Failed to connect to ChirpStack Application Server gRPC server")
		return err
	}

	clientConn = conn

	// connect MQTT server
	opts := mqtt.NewClientOptions()
	opts.AddBroker(conf.Integration.MQTT.Server)
	opts.SetUsername(conf.Integration.MQTT.Username)
	opts.SetPassword(conf.Integration.MQTT.Password)
	opts.SetCleanSession(true)
	opts.SetAutoReconnect(true)
	if conf.Integration.MQTT.TLSConfig.Insecure {
		log.Warn("Using insecure TLS configuration for ChirpStack Application Server MQTT server connection")
		opts.SetTLSConfig(&tls.Config{InsecureSkipVerify: true})
	} else {
		caCert, err := os.ReadFile(conf.Integration.MQTT.TLSConfig.CAFile)
		if err != nil {
			log.WithError(err).Error("Failed to read CA file for ChirpStack Application Server MQTT server connection")
			return err
		}

		caCertPool := x509.NewCertPool()
		if !caCertPool.AppendCertsFromPEM(caCert) {
			log.Error("Failed to append CA certificate for ChirpStack Application Server MQTT server connection")
			return err
		}
		tlsConfig := &tls.Config{
			InsecureSkipVerify: false,
			RootCAs:            caCertPool,
		}

		opts.SetTLSConfig(tlsConfig)
	}

	mqttClient := mqtt.NewClient(opts)
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		log.WithError(token.Error()).Error("Failed to connect to ChirpStack Application Server MQTT server")
		return token.Error()
	}
	log.Info("Connected to ChirpStack Application Server MQTT server")

	return nil
}

func Tenant() api.TenantServiceClient {
	return api.NewTenantServiceClient(clientConn)
}

func Gateway() api.GatewayServiceClient {
	return api.NewGatewayServiceClient(clientConn)
}

func DeviceProfile() api.DeviceProfileServiceClient {
	return api.NewDeviceProfileServiceClient(clientConn)
}

func Application() api.ApplicationServiceClient {
	return api.NewApplicationServiceClient(clientConn)
}

func Device() api.DeviceServiceClient {
	return api.NewDeviceServiceClient(clientConn)
}

// MQTTClient returns the MQTT client for the Application Server MQTT integration.
func MQTTClient() mqtt.Client {
	return mqttClient
}
