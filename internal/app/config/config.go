package config

type PostgreSQLConfig struct {
	URI string `mapstructure:"uri"`
}

type ServerConfig struct {
	Port int `mapstructure:"port"`
}

type JWTConfig struct {
	Secret string `mapstructure:"secret"`
}

type CryptoConfig struct {
	Key string `mapstructure:"key"`
}

type Config struct {
	Debug      bool             `mapstructure:"debug"`
	Server     ServerConfig     `mapstructure:"server"`
	PostgreSQL PostgreSQLConfig `mapstructure:"postgresql"`
	JWT        JWTConfig        `mapstructure:"jwt"`
	Crypto     CryptoConfig     `mapstructure:"crypto"`
	ApiPrefix  string           `mapstructure:"apiPrefix"`
}
