package config

import (
	"fmt"
	"time"

	cleanenvport "github.com/wb-go/wbf/config/cleanenv-port"
	"github.com/wb-go/wbf/logger"
)

type Config struct {
	Server   ServerConfig   `yaml:"server"    validate:"required"`
	Postgres PostgresConfig `yaml:"postgres"  validate:"required"`
	Logger   LoggerConfig   `yaml:"logger"    validate:"required"`
	Gin      GinConfig      `yaml:"gin"       validate:"required"`
	Retry    RetryConfig    `yaml:"retry"`
}

// LogLevel преобразует строковый уровень в logger.Level из wbf.
func (c LoggerConfig) LogLevel() logger.Level {
	switch c.Level {
	case "debug":
		return logger.DebugLevel
	case "warn":
		return logger.WarnLevel
	case "error":
		return logger.ErrorLevel
	default:
		return logger.InfoLevel
	}
}

// LogEngine преобразует строковый движок в logger.Engine из wbf.
func (c LoggerConfig) LogEngine() logger.Engine {
	return logger.Engine(c.Engine)
}

type LoggerConfig struct {
	Engine string `yaml:"engine" env:"LOG_ENGINE" env-default:"slog"  validate:"required,oneof=slog zap zerolog logrus"`
	Level  string `yaml:"level"  env:"LOG_LEVEL"  env-default:"info"  validate:"required,oneof=debug info warn error"`
}

type GinConfig struct {
	Mode string `yaml:"mode" env:"GIN_MODE" env-default:"debug" validate:"required,oneof=debug release test"`
}

type ServerConfig struct {
	Addr            string        `yaml:"addr"          env:"SERVER_ADDR"          env-default:":8080" validate:"required"`
	ReadTimeout     time.Duration `yaml:"read_timeout"  env:"SERVER_READ_TIMEOUT"  env-default:"10s"   validate:"gt=0"`
	WriteTimeout    time.Duration `yaml:"write_timeout" env:"SERVER_WRITE_TIMEOUT" env-default:"10s"   validate:"gt=0"`
	IdleTimeout     time.Duration `yaml:"idle_timeout"  env:"SERVER_IDLE_TIMEOUT"  env-default:"60s"   validate:"gt=0"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout" env:"SERVER_SHUTDOWN_TIMEOUT" env-default:"15s"   validate:"gt=0"`
}

type PostgresConfig struct {
	Host            string        `yaml:"host"              env:"DB_HOST"              env-default:"localhost"    validate:"required"`
	Port            int           `yaml:"port"              env:"DB_PORT"              env-default:"5432"         validate:"required,min=1,max=65535"`
	User            string        `yaml:"user"              env:"DB_USER"              env-default:"postgres"     validate:"required"`
	Password        string        `yaml:"password"          env:"DB_PASSWORD"          env-default:"postgres"     validate:"required"`
	Database        string        `yaml:"database"          env:"DB_NAME"              env-default:"salestracker"  validate:"required"`
	SSLMode         string        `yaml:"sslmode"           env:"DB_SSLMODE"           env-default:"disable"      validate:"required,oneof=disable require verify-ca verify-full"`
	MaxOpenConns    int           `yaml:"max_open_conns"    env:"DB_MAX_OPEN_CONNS"    env-default:"10"           validate:"min=1"`
	MaxIdleConns    int           `yaml:"max_idle_conns"    env:"DB_MAX_IDLE_CONNS"    env-default:"5"            validate:"min=1"`
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime" env:"DB_CONN_MAX_LIFETIME" env-default:"5m"           validate:"gt=0"`
}

func (p *PostgresConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		p.Host, p.Port, p.User, p.Password, p.Database, p.SSLMode,
	)
}

type RetryConfig struct {
	Attempts int           `yaml:"attempts" env:"RETRY_ATTEMPTS" env-default:"3"    validate:"min=1"`
	Delay    time.Duration `yaml:"delay"    env:"RETRY_DELAY"    env-default:"500ms" validate:"gt=0"`
	Backoff  float64       `yaml:"backoff"  env:"RETRY_BACKOFF"  env-default:"2"     validate:"min=1"`
}

func MustLoad() *Config {
	var cfg Config
	if err := cleanenvport.Load(&cfg); err != nil {
		panic(fmt.Sprintf("failed to load config: %v", err))
	}
	return &cfg
}
