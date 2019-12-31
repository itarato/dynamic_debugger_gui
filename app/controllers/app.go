package controllers

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/revel/revel"
	"gopkg.in/yaml.v2"
)

const ConfigFileName = ".dynamic_debugger.config.yml"

func ConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/%s", homeDir, ConfigFileName), nil
}

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

type ParamBreakpoint struct {
	Enabled      string
	ReturnOption string
	Return       string
	ReturnCode   string
	ReturnCall   []string
	PreCall      []string
	PostCall     []string
}

func LoadConfig() (Config, error) {
	config := Config{}

	config_file_path, file_name_err := ConfigFilePath()
	if file_name_err != nil {
		return config, file_name_err
	}

	yaml_file, err := ioutil.ReadFile(config_file_path)
	if err != nil {
		return config, err
	}

	err_yaml_load := yaml.Unmarshal(yaml_file, &config)
	return config, err_yaml_load
}

func (c App) Index() revel.Result {
	config, err := LoadConfig()
	if err != nil {
		c.Log.Error("Config cannot be loaded")
	}

	c.ViewArgs["config"] = config

	return c.Render()
}

func UpdateConfigFromParam(breakpoints_raw map[string]ParamBreakpoint) error {
	config := Config{make(map[string]Breakpoint)}

	for name, param_breakpoint := range breakpoints_raw {
		breakpoint := Breakpoint{}

		breakpoint.Enabled = param_breakpoint.Enabled == "on"

		if param_breakpoint.ReturnOption == "return" {
			breakpoint.Return = param_breakpoint.Return
		} else {
			breakpoint.ReturnCode = param_breakpoint.ReturnCode
		}

		breakpoint.ReturnCall = param_breakpoint.ReturnCall
		breakpoint.PreCall = param_breakpoint.PreCall
		breakpoint.PostCall = param_breakpoint.PostCall

		config.Breakpoints[name] = breakpoint
	}

	return UpdateConfig(config)
}

func UpdateConfig(config Config) error {
	yaml_out, yaml_err := yaml.Marshal(&config)
	if yaml_err != nil {
		return yaml_err
	}

	config_file_path, file_name_err := ConfigFilePath()
	if file_name_err != nil {
		panic("Cannot generate config file path")
	}
	file_err := ioutil.WriteFile(config_file_path, yaml_out, 0644)
	if file_err != nil {
		return file_err
	}

	return nil
}

func (c App) Update() revel.Result {
	var breakpoints map[string]ParamBreakpoint
	c.Params.Bind(&breakpoints, "breakpoints")

	update_err := UpdateConfigFromParam(breakpoints)
	if update_err != nil {
		c.Log.Error("Cannot update config file.")
	}

	return c.Redirect(App.Index)
}

func (c App) AddBreakpoint() revel.Result {
	name := c.Params.Form.Get("name")

	config, err := LoadConfig()
	if err != nil {
		c.Log.Error("Config cannot be loaded.")
	} else {
		new_breakpoint := Breakpoint{}
		new_breakpoint.Enabled = true
		config.Breakpoints[name] = new_breakpoint
		update_err := UpdateConfig(config)
		if update_err != nil {
			c.Log.Error("Config cannot be updated")
		}
	}

	return c.Redirect(App.Index)
}
