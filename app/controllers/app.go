package controllers

import (
	"io/ioutil"

	"github.com/revel/revel"
	"gopkg.in/yaml.v2"
)

type Breakpoint struct {
	Enabled    bool     `yaml:"enabled"`
	Return     string   `yaml:"return"`
	ReturnCode string   `yaml:"return_code"`
	ReturnCall []string `yaml:"return_call"`
	PreCall    []string `yaml:"pre_call"`
	PostCall   []string `yaml:"post_call"`
}

func (b Breakpoint) IsReturnCode() bool {
	return len(b.ReturnCode) > 0
}

type Config struct {
	Breakpoints map[string]Breakpoint `yaml:"breakpoints"`
}

type App struct {
	*revel.Controller
}

func (c App) Index() revel.Result {
	yaml_file, err := ioutil.ReadFile("/home/itarato/.dynamic_debugger.config.yml")
	if err != nil {
		c.Log.Error("Cannot read config file.")
	}

	config := Config{}
	err_yaml_load := yaml.Unmarshal(yaml_file, &config)

	if err_yaml_load != nil {
		c.Log.Error("YAML cannot be unamrshalled.")
	}

	c.ViewArgs["config"] = config

	return c.Render()
}

func (c App) Update() revel.Result {
	return c.Redirect(App.Index)
}
