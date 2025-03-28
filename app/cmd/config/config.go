// Package config defines the configuration parameters of the service
package config

import (
	"fmt"

	"github.com/asaskevich/govalidator"
	"github.com/spf13/viper"
)

const (
	configFileName      = "service-config.yml"
	configFileExtension = "yaml"
)

// StaticConfiguration defines the static configuration of the service
type StaticConfiguration struct {
	// DbConfig defines the database connection needed for the service to start
	DbConfig DBConfig `mapstructure:"database" valid:"required"`
}

// DBConfig defines the database configuration
type DBConfig struct {
	// PostgreSQL database configuration
	PostgreSQL *PostgreSQLInfo `mapstructure:"postgresql"`
}

// PostgreSQLInfo defines the connection configuration to a PostgresSQL DB
type PostgreSQLInfo struct {
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

func GetStaticConfiguration(path string) (*StaticConfiguration, error) {
	viper.SetConfigName(configFileName)
	viper.SetConfigType(configFileExtension)
	viper.AddConfigPath(path)
	err := viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("error reading static config file placed in path [%s] [error:%w]", path, err)
	}

	var staticConfiguration StaticConfiguration
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
