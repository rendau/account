package cmd

import (
	"github.com/rendau/dop/dopTools"
	"github.com/spf13/viper"
)

var conf = struct {
	Debug          bool   `mapstructure:"DEBUG"`
	LogLevel       string `mapstructure:"LOG_LEVEL"`
	HttpListen     string `mapstructure:"HTTP_LISTEN"`
	HttpCors       bool   `mapstructure:"HTTP_CORS"`
	SwagHost       string `mapstructure:"SWAG_HOST"`
	SwagBasePath   string `mapstructure:"SWAG_BASE_PATH"`
	SwagSchema     string `mapstructure:"SWAG_SCHEMA"`
	PgDsn          string `mapstructure:"PG_DSN"`
	RedisUrl       string `mapstructure:"REDIS_URL"`
	RedisPsw       string `mapstructure:"REDIS_PSW"`
	RedisDb        int    `mapstructure:"REDIS_DB"`
	RedisKeyPrefix string `mapstructure:"REDIS_KEY_PREFIX"`
	JwtsGrpcUrl    string `mapstructure:"JWTS_GRPC_URL"`
	MsJwtsUrl      string `mapstructure:"MS_JWTS_URL"`
	MsSmsUrl       string `mapstructure:"MS_SMS_URL"`
	NoSmsCheck     bool   `mapstructure:"NO_SMS_CHECK"`
}{}

func confLoad() {
	dopTools.SetViperDefaultsFromObj(conf)

	viper.SetDefault("DEBUG", "false")
	viper.SetDefault("LOG_LEVEL", "info")
	viper.SetDefault("HTTP_LISTEN", ":80")
	viper.SetDefault("SWAG_HOST", "example.com")
	viper.SetDefault("SWAG_BASE_PATH", "/")
	viper.SetDefault("SWAG_SCHEMA", "https")
	viper.SetDefault("REDIS_KEY_PREFIX", "account_")

	viper.SetConfigFile("conf.yml")
	_ = viper.ReadInConfig()

	viper.AutomaticEnv()

	_ = viper.Unmarshal(&conf)
}
