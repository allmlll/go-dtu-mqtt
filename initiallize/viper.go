package initiallize

import (
	"fmt"
	"github.com/spf13/viper"
)

func ViperInit() {
	viper.SetConfigName("conf")
	viper.AddConfigPath("./config")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
		return
	}
}
