package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Env          string             `yaml:"env"`
	GRPC         GRPCConfig         `yaml:"grpc"`
	HTTP         HTTPConfig         `yaml:"http"`
	Postgres     PostgresConfig     `yaml:"postgres"`
	SearchEngine SearchEngineConfig `yaml:"search_engine"`
	Upload       UploadConfig       `yaml:"upload"`
	Logger       LoggerConfig       `yaml:"logger"`
}

type GRPCConfig struct {
	Port int `yaml:"port"`
}

type HTTPConfig struct {
	Port int `yaml:"port"`
}

type PostgresConfig struct {
	DSN      string `yaml:"dsn"`
	MaxConns int32  `yaml:"max_conns"`
	MinConns int32  `yaml:"min_conns"`
}

type SearchEngineConfig struct {
	GRPCAddr string `yaml:"grpc_addr"`
}

type UploadConfig struct {
	MaxFileSizeBytes int64 `yaml:"max_file_size_bytes"`
	ChunkSize        int   `yaml:"chunk_size"`
	ChunkOverlap     int   `yaml:"chunk_overlap"`
}

type LoggerConfig struct {
	Level string `yaml:"level"`
}

func Load(path string) (*Config, error) {
	_ = godotenv.Overload()

	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config file %q: %w", path, err)
	}

	cfg := &Config{}
	if err := yaml.Unmarshal(raw, cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	applyEnvOverrides(cfg)

	return cfg, nil
}

func applyEnvOverrides(cfg *Config) {
	if v := os.Getenv("DOCUMENT_PROCESSOR_GRPC_PORT"); v != "" {
		cfg.GRPC.Port = atoiOrDefault(v, cfg.GRPC.Port)
	}
	if v := os.Getenv("DOCUMENT_PROCESSOR_HTTP_PORT"); v != "" {
		cfg.HTTP.Port = atoiOrDefault(v, cfg.HTTP.Port)
	}
	if v := os.Getenv("DOCUMENT_PROCESSOR_POSTGRES_DSN"); v != "" {
		cfg.Postgres.DSN = v
	}
	if v := os.Getenv("DOCUMENT_PROCESSOR_SEARCH_ENGINE_GRPC_ADDR"); v != "" {
		cfg.SearchEngine.GRPCAddr = v
	}
	if v := os.Getenv("DOCUMENT_PROCESSOR_LOG_LEVEL"); v != "" {
		cfg.Logger.Level = v
	}
}

func atoiOrDefault(s string, def int) int {
	n := 0
	for _, c := range s {
		if c < '0' || c > '9' {
			return def
		}
		n = n*10 + int(c-'0')
	}
	return n
}
