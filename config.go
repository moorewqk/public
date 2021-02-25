package public

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type Server struct {
	Redis   *Redis   `mapstructure:"redis" json:"redis" yaml:"redis"`
}

var (
	SV Server
)


func ViperConfig(path string) *viper.Viper {

	v := viper.New()
	v.SetConfigFile(path)
	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	v.WatchConfig()

	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("config file changed:", e.Name)
		if err := v.Unmarshal(&SV); err != nil {
			fmt.Println(err)
		}
	})
	if err := v.Unmarshal(&SV); err != nil {
		fmt.Println(err)
	}
	return v
}
