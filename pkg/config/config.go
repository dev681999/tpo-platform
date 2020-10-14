package config

//go:generate stringer -type=Config
import (
	"fmt"
	"io/ioutil"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"
)

func chechCfgFileName(s string) string {
	if s == "" {
		return "config.yaml"
	}

	return s
}

// Config is configuration for the server
type Config struct {
	DBFile     string `yaml:"dbFile" env:"DB_FILE" env-default:"file:tpo.db?cache=shared&_fk=1"`
	ServerAddr string `yaml:"serverAddr" env:"SERVER_ADDR" env-default:":8080"`
	JWTSecret  string `yaml:"jwtSecret" env:"JWT_SECRET" env-default:"MySuperSecretKey"`
}

// String implements the fmt.Stringer.
func (cfg Config) String() string {
	return fmt.Sprintf("Config(db_file=%v, server_addr=%s, jwt_secret=%s)", cfg.DBFile, cfg.ServerAddr, cfg.JWTSecret)
}

// GenerateDefault generates default config
func GenerateDefault() error {
	return GenerateDefaultWithFileName("")
}

// GenerateDefaultWithFileName generates default config in given file
func GenerateDefaultWithFileName(fileName string) error {
	fileName = chechCfgFileName(fileName)
	var cfg Config
	cleanenv.ReadEnv(&cfg)
	b, err := yaml.Marshal(cfg)
	if err != nil {
		log.Debug().Err(err).Msg("")
		return err
	}
	err = ioutil.WriteFile(fileName, b, 0755)
	if err != nil {
		log.Debug().Err(err).Msg("")
		return err
	}

	return nil
}

// ReadConfig reats config with default filename
func ReadConfig() (Config, error) {
	return ReadConfigWithFileName("")
}

// ReadConfigWithFileName reats config with filename
func ReadConfigWithFileName(fileName string) (Config, error) {
	fileName = chechCfgFileName(fileName)
	var cfg Config
	err := cleanenv.ReadConfig(fileName, &cfg)
	if err != nil {
		log.Debug().Err(err).Msg("")
		return Config{}, err
	}

	return cfg, nil
}
