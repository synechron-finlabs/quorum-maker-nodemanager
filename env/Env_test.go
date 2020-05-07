package env

import (
	"testing"
)

func TestGetAppConfig(t *testing.T) {
	appConfig := GetAppConfig(true)

	if "" == appConfig.HomeDir {
		t.Errorf("Could not load App Config")
	}
}


func TestGetSetupConfig(t *testing.T) {
	setupConfig := GetSetupConf(true)

	if "" == setupConfig.ContractAdd {
		t.Errorf("Could not load Setup Config")
	}
}