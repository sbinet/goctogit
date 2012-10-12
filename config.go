package main

import (
	"fmt"
	"os"
	"path/filepath"

	gocfg "github.com/sbinet/go-config/config"
)

var Cfg = gocfg.NewDefault()
var CfgFname = os.ExpandEnv(filepath.Join("${HOME}", ".config", "go-octogit", "config.ini"))

func init() {
	cfgdir := filepath.Dir(CfgFname)
	if !path_exists(cfgdir) {
		err := os.MkdirAll(cfgdir, 0700)
		if err != nil {
			panic(err.Error())
		}
	}
	
	fname := CfgFname

	if !path_exists(fname) {
		section := "go-octogit"
		if !Cfg.AddSection(section) {
			err := fmt.Errorf("go-octogit: could not create section [%s] in file [%s]", section, fname)
			panic(err.Error())
		}
		for k,v := range map[string]string{
			"username": "",
			"password": "",
			"token": "",
		}{
			if !Cfg.AddOption(section, k, v) {
				err := fmt.Errorf("go-octogit: could not add option [%s] to section [%s]", k, section)
				panic(err.Error())
			}
		}
	} else {
		cfg, err := gocfg.ReadDefault(fname)
		if err != nil {
			panic(err.Error())
		}
		Cfg = cfg
	}
}

// EOF
