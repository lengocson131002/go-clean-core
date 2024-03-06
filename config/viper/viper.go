package viper

import (
	"os"
	"time"

	"github.com/lengocson131002/go-clean-core/config"
	"github.com/spf13/viper"
)

type ViperConfig struct {
	cfg *viper.Viper
}

type ConfigFile string

func NewViperConfig(f *ConfigFile) (config.Configure, error) {
	viper := viper.New()
	viper.AutomaticEnv()

	// check if config file exists
	if f != nil {
		cfgFile := string(*f)
		if _, err := os.Stat(cfgFile); err == nil {
			viper.SetConfigFile(cfgFile)
			err := viper.ReadInConfig()
			if err != nil {
				return nil, err
			}
		}
	}

	return &ViperConfig{
		cfg: viper,
	}, nil
}

// Get implements ConfigInterface.
func (c *ViperConfig) Get(key string) interface{} {
	return c.cfg.Get(key)
}

// GetBool implements ConfigInterface.
func (c *ViperConfig) GetBool(key string) bool {
	return c.cfg.GetBool(key)
}

// GetDuration implements ConfigInterface.
func (c *ViperConfig) GetDuration(key string) time.Duration {
	return c.cfg.GetDuration(key)
}

// GetFloat64 implements ConfigInterface.
func (c *ViperConfig) GetFloat64(key string) float64 {
	return c.cfg.GetFloat64(key)
}

// GetInt implements ConfigInterface.
func (c *ViperConfig) GetInt(key string) int {
	return c.cfg.GetInt(key)
}

// GetInt32 implements ConfigInterface.
func (c *ViperConfig) GetInt32(key string) int32 {
	return c.cfg.GetInt32(key)
}

// GetInt64 implements ConfigInterface.
func (c *ViperConfig) GetInt64(key string) int64 {
	return c.cfg.GetInt64(key)
}

// GetIntSlice implements ConfigInterface.
func (c *ViperConfig) GetIntSlice(key string) []int {
	return c.cfg.GetIntSlice(key)
}

// GetString implements ConfigInterface.
func (c *ViperConfig) GetString(key string) string {
	return c.cfg.GetString(key)
}

// GetStringMap implements ConfigInterface.
func (c *ViperConfig) GetStringMap(key string) map[string]any {
	return c.cfg.GetStringMap(key)
}

// GetStringSlice implements ConfigInterface.
func (c *ViperConfig) GetStringSlice(key string) []string {
	return c.cfg.GetStringSlice(key)
}

// GetTime implements ConfigInterface.
func (c *ViperConfig) GetTime(key string) time.Time {
	return c.cfg.GetTime(key)
}

// GetUint implements ConfigInterface.
func (c *ViperConfig) GetUint(key string) uint {
	return c.cfg.GetUint(key)
}

// GetUint16 implements ConfigInterface.
func (c *ViperConfig) GetUint16(key string) uint16 {
	return c.cfg.GetUint16(key)
}

// GetUint32 implements ConfigInterface.
func (c *ViperConfig) GetUint32(key string) uint32 {
	return c.cfg.GetUint32(key)
}

// GetUint64 implements ConfigInterface.
func (c *ViperConfig) GetUint64(key string) uint64 {
	return c.cfg.GetUint64(key)
}
