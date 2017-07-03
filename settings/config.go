package settings

import (
	"github.com/swanwish/fileserver/common"
	"github.com/swanwish/go-common/config"
	"github.com/swanwish/go-common/logs"
)

var (
	ConfigFilePath = "conf/app.ini"
)

func LoadAppSetting() {
	config.Load(ConfigFilePath)

	if logLevel, err := config.Get(common.SETTING_KEY_LOG_LEVEL); err == nil {
		logs.SetLogLevel(logLevel)
	}
}
