package gosyncmodules

import (
	"gopkg.in/ini.v1"
)


func GetConfig(configFile string) (*ini.File, error){
	var Cfg *ini.File
	Cfg, err := ini.Load(configFile)
	if err != nil {
		return Cfg, err
	}
	return Cfg, err
}