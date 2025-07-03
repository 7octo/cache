package config

type Config struct {
    Server struct {
        Port string `yaml:"port" env:"SERVER_PORT" env-default:":3000"`
    } `yaml:"server"`
}

func Load() (*Config, error) {
    // In a real app, load from file/env
    cfg := &Config{}
    cfg.Server.Port = ":3000"
    return cfg, nil
}