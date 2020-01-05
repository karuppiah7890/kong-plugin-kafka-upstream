package main

import (
	"github.com/karuppiah7890/go-pdk"
)

type PluginConfig struct {}

func New() interface{} {
	return &PluginConfig{}
}

func (conf *PluginConfig) Access (kong *pdk.PDK) {
	err := kong.Response.Exit(200, map[string]interface{}{"message": "you are da best! :D"}, nil)
	if err != nil {
		_ = kong.Log.Err(err.Error())
	}
}