package env

import (
	"flag"
	"fmt"
	"log"
	"log/slog"
	"strings"

	"github.com/spf13/viper"
)

type Env struct{}

func New() *Env {
	viper.AddConfigPath(".")
	viper.SetConfigFile(".env")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		slog.Warn(fmt.Sprintf("reading config: %s", err.Error()))
	}
	flag.String("app", "web", "functionality to run")
	flag.Parse()
	return &Env{}
}

func (*Env) MustPresent(key ...string) {
	for _, k := range key {
		if !viper.IsSet(k) {
			log.Fatalf("Environment variable %s not set", k)
		}
	}
}

func (d *Env) MustGet(key string) string {
	d.MustPresent(key)
	return d.Get(key)
}

func (*Env) Get(key string) string {
	function := viper.GetString(key)
	flag.Visit(func(f *flag.Flag) {
		if strings.EqualFold(f.Name, key) {
			function = f.Value.String()
		}
	})
	return function
}

func (d *Env) GetDefault(key string, defaultValue string) string {
	val := d.Get(key)
	if val == "" {
		val = defaultValue
	}
	return val
}
