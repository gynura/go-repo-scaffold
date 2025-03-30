// Package config defines the configuration parameters of the service
package config

import (
	"fmt"
	"userservice/internal/graph"

	"github.com/asaskevich/govalidator"
	"github.com/spf13/viper"
)

const (
	configFileName      = "service-config.yaml"
	configFileExtension = "yaml"
)

// Config defines the static configuration of the service
type Config struct {
	// DbConfig defines the database connection needed for the service to start
	DbConfig DBConfig `mapstructure:"database" valid:"required"`
	// QueueConfig defines the queue connection needed for the service to start
	QueueConfig QueueConfig `mapstructure:"queue" valid:"required"`
}

// DBConfig defines the database configuration
type DBConfig struct {
	PostgreSQL *PostgreSQL `mapstructure:"postgresql"`
}

// QueueConfig defines the queue configuration
type QueueConfig struct {
	RabbitMQ *RabbitMQ `mapstructure:"rabbitmq"`
}

// RabbitMQ defines the connection configuration to a RabbitMQ queue
type RabbitMQ struct {
	// Host address of the queue
	Host string `mapstructure:"host"`
	// Port of queue system
	Port string `mapstructure:"port"`
	// Username to use in order to connect to the queue
	Username string `mapstructure:"username"`
	// Password to use in order to connect to the queue
	Password string `mapstructure:"password"`
	// Vhost name to connect to in the queue
	VHost string `mapstructure:"vhost"`
}

// PostgreSQL defines the connection configuration to a PostgresSQL DB
type PostgreSQL struct {
	// Host of database system
	Host string `mapstructure:"host" valid:"required~hosts is mandatory in SQL DB config"`
	// Port of database system
	Port int `mapstructure:"port" valid:"required~port is mandatory in SQL DB config"`
	// Scheme of database system
	Scheme string `mapstructure:"scheme" valid:"required~scheme is mandatory in SQL DB config"`
	// Username to use in database system
	Username string `mapstructure:"username" json:"-"`
	// Password to use with username in database system
	Password string `mapstructure:"password" json:"-"`
	// SSLMode to use in database system
	SSLMode string `mapstructure:"sslmode" valid:"required~sslmode is mandatory in SQL DB config"`
	// Database to access to in database system
	Database string `mapstructure:"database" valid:"required~database is mandatory in SQL DB config"`
	// SQLClient database client configuration
	SQLClient *PostgresSQLClient `mapstructure:"sqlClient"`
}

// PostgresSQLClient configures the database client
type PostgresSQLClient struct {
	// MaxIdleConnections max idle connections for the database/sql handle
	MaxIdleConnections *int `mapstructure:"maxIdleConnections"`
	// MaxOpenConnections max open connections for the database/sql handle
	MaxOpenConnections *int `mapstructure:"maxOpenConnections"`
	// MaxConnectionLifetime max connection lifetime for the database/sql handle
	MaxConnectionLifetime *int `mapstructure:"maxConnectionLifetime"`
}

func readConfig(path string) (*Config, error) {
	viper.SetConfigName(configFileName)
	viper.SetConfigType(configFileExtension)
	viper.AddConfigPath(path)
	err := viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("error reading config file placed in path [%s] [error:%w]", path, err)
	}

	var staticConfiguration Config
	err = viper.Unmarshal(&staticConfiguration)
	if err != nil {
		panic(err)
	}
	valid, err := govalidator.ValidateStruct(staticConfiguration)
	if !valid || err != nil {
		return nil, fmt.Errorf("error validating static config file placed in path[%s] [error:%w]", path, err)
	}

	return &staticConfiguration, nil
}

func GetGraphConfig(path string) (*graph.Config, error) {
	config, err := readConfig(path)
	if err != nil {
		return nil, err
	}

	graphConfig := graph.Config{
		Libraries: graph.LibrariesConfig{
			PersistenceFw: graph.PersistenceFwConfig{},
		},
		QueueConfig: graph.QueueConfig{},
	}

	if config.QueueConfig.RabbitMQ != nil {
		rabbitConfig := &graph.RabbitMQ{
			Host:     config.QueueConfig.RabbitMQ.Host,
			Port:     config.QueueConfig.RabbitMQ.Port,
			Username: config.QueueConfig.RabbitMQ.Username,
			Password: config.QueueConfig.RabbitMQ.Password,
			VHost:    config.QueueConfig.RabbitMQ.VHost,
		}
		graphConfig.QueueConfig.RabbitMQ = rabbitConfig
	}

	if config.DbConfig.PostgreSQL != nil {
		postgresConfig := &graph.PostgresSQLConfig{
			Host:     config.DbConfig.PostgreSQL.Host,
			Port:     &config.DbConfig.PostgreSQL.Port,
			Scheme:   &config.DbConfig.PostgreSQL.Scheme,
			Username: config.DbConfig.PostgreSQL.Username,
			Password: config.DbConfig.PostgreSQL.Password,
			SSLMode:  config.DbConfig.PostgreSQL.SSLMode,
			Database: config.DbConfig.PostgreSQL.Database,
		}

		graphConfig.Libraries.PersistenceFw.PostgreSQL = postgresConfig

		if config.DbConfig.PostgreSQL.SQLClient != nil {
			graphConfig.Libraries.PersistenceFw.PostgreSQL.SQLClient = &graph.PostgresSQLClientConfig{
				MaxIdleConnections:    config.DbConfig.PostgreSQL.SQLClient.MaxIdleConnections,
				MaxOpenConnections:    config.DbConfig.PostgreSQL.SQLClient.MaxOpenConnections,
				MaxConnectionLifetime: config.DbConfig.PostgreSQL.SQLClient.MaxConnectionLifetime,
			}
		}
	}

	return &graphConfig, nil
}
