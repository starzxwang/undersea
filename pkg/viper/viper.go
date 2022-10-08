package viper

import (
	"strings"
	"sync"
	"time"

	"github.com/spf13/viper"
)

var once = sync.Once{}

func initViper() {
	ins = viper.New()
	ins.SetEnvPrefix("UNDERSEA")
	ins.SetEnvKeyReplacer(EnvKeyReplacer)
	ins.AutomaticEnv()

	initConf()
}

var (
	EnvKeyReplacer = strings.NewReplacer(".", "_")
)

func initConf() {
	ins.SetConfigType("yaml")

	ins.SetConfigFile("./conf/conf.yaml")

	//yamlPath := ins.GetString("conf")
	//if yamlPath != "" {
	//	ins.SetConfigFile(yamlPath)
	//	return
	//}

	//ins.SetConfigName("huskar")
	//curr, err := os.Getwd()
	//if err != nil {
	//	panic(err)
	//}
	//ins.AddConfigPath(curr)
	//for dir := filepath.Dir(curr); dir != curr; dir, curr = filepath.Dir(dir), dir {
	//	ins.AddConfigPath(dir)
	//}

	_ = ins.ReadInConfig()
}

var ins *viper.Viper

func V() *viper.Viper {
	if ins == nil {
		once.Do(initViper)
	}

	return ins
}

func GetInt64(v *viper.Viper, key string, dflt int64) int64 {
	if v.IsSet(key) {
		return v.GetInt64(key)
	}
	return dflt
}

func GetInt(v *viper.Viper, key string, dflt int) int {
	if v.IsSet(key) {
		return v.GetInt(key)
	}
	return dflt
}

func GetFloat64(v *viper.Viper, key string, dflt float64) float64 {
	if v.IsSet(key) {
		return v.GetFloat64(key)
	}
	return dflt
}

func GetDuration(v *viper.Viper, key string, dflt time.Duration) time.Duration {
	if v.IsSet(key) {
		ret, err := time.ParseDuration(v.GetString(key))
		if err != nil {
			return 0 * time.Nanosecond
		}
		return ret
	}
	return dflt
}
