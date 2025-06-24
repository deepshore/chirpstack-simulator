package config

import "time"

var Version string

// Config defines the configuration.
type Config struct {
	ChirpStack struct {
		API struct {
			APIKey   string `mapstructure:"api_key"`
			Server   string `mapstructure:"server"`
			Insecure bool   `mapstructure:"insecure"`
		} `mapstructure:"api"`

		Integration struct {
			MQTT struct {
				Server    string `mapstructure:"server"`
				Username  string `mapstructure:"username"`
				Password  string `mapstructure:"password"`
				TLSConfig struct {
					Insecure bool   `mapstructure:"insecure"`
					CAFile   string `mapstructure:"ca_file"`
				} `mapstructure:"tls_config"`
			} `mapstructure:"mqtt"`
		} `mapstructure:"integration"`

		Gateway struct {
			Backend struct {
				MQTT struct {
					Server    string `mapstructure:"server"`
					Username  string `mapstructure:"username"`
					Password  string `mapstructure:"password"`
					TLSConfig struct {
						Insecure bool   `mapstructure:"insecure"`
						CAFile   string `mapstructure:"ca_file"`
					} `mapstructure:"tls_config"`
				} `mapstructure:"mqtt"`
			} `mapstructure:"backend"`
		} `mapstructure:"gateway"`
	} `mapstructure:"chirpstack"`

	Simulator []struct {
		TenantID       string        `mapstructure:"tenant_id"`
		Duration       time.Duration `mapstructure:"duration"`
		ActivationTime time.Duration `mapstructure:"activation_time"`

		Device struct {
			Count           int           `mapstructure:"count"`
			UplinkInterval  time.Duration `mapstructure:"uplink_interval"`
			FPort           uint8         `mapstructure:"f_port"`
			Payload         string        `mapstructure:"payload"`
			Frequency       int           `mapstructure:"frequency"`
			Bandwidth       int           `mapstructure:"bandwidth"`
			SpreadingFactor int           `mapstructure:"spreading_factor"`
		} `mapstructure:"device"`

		Gateway struct {
			MinCount             int    `mapstructure:"min_count"`
			MaxCount             int    `mapstructure:"max_count"`
			EventTopicTemplate   string `mapstructure:"event_topic_template"`
			CommandTopicTemplate string `mapstructure:"command_topic_template"`
		} `mapstructure:"gateway"`
	} `mapstructure:"simulator"`

	Prometheus struct {
		Bind string `mapstructure:"bind"`
	} `mapstructure:"prometheus"`
}

type DeviceConfig struct {
	DevEUI         string        `mapstructure:"dev_eui"`
	AppKey         string        `mapstructure:"app_key"`
	UplinkInterval time.Duration `mapstructure:"uplink_interval"`
	FPort          uint8         `mapstructure:"f_port"`
	Payload        string        `mapstructure:"payload"`
}

// C holds the global configuration.
var C Config
