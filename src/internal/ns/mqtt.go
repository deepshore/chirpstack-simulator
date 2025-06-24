package ns

import (
	"crypto/tls"
	"crypto/x509"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/deepshore/chirpstack-simulator/internal/config"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var mqttClient mqtt.Client

func Setup(c config.Config) error {
	conf := c.ChirpStack.Gateway.Backend.MQTT

	opts := mqtt.NewClientOptions()
	opts.AddBroker(conf.Server)
	opts.SetUsername(conf.Username)
	opts.SetPassword(conf.Password)
	opts.SetCleanSession(true)

	if conf.TLSConfig.Insecure {
		log.Warn("Using insecure TLS configuration for ChirpStack Application Server MQTT server connection")
		opts.SetTLSConfig(&tls.Config{
			InsecureSkipVerify: true,
		})
	} else {
		caCert, err := os.ReadFile(conf.TLSConfig.CAFile)
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

	mqttClient = mqtt.NewClient(opts)
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		log.WithError(token.Error()).Error("Failed to connect to ChirpStack Application Server MQTT server")
		return token.Error()
	}
	log.Info("Connected to ChirpStack Application Server MQTT server")

	return nil
}

func Client() mqtt.Client {
	return mqttClient
}
