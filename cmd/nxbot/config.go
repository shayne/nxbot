package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

var errNoConfig = errors.New("NoConfig")

type config struct {
	NxIPPort           string   `yaml:"nx_ip_port"`
	NxUser             string   `yaml:"nx_user"`
	NxPass             string   `yaml:"nx_pass"`
	HTTPIPPort         string   `yaml:"http_ip_port"`
	TgToken            string   `yaml:"tg_token"`
	TgUserWhitelist    []int    `yaml:"tg_user_whitelist"`
	TgGroupWhitelist   []int64  `yaml:"tg_group_whitelist"`
	TgMotionRecipients []string `yaml:"tg_motion_recipients"`
}

func loadConfig() (*config, error) {
	ec, err := loadEnv()
	if err != nil {
		return nil, err
	}
	yc, err := loadYaml()
	if err != nil && err != errNoConfig {
		return nil, err
	}
	var c *config
	if yc == nil { // no config
		c = ec
	} else {
		c = applyEnvOverwrite(ec, yc)
	}
	err = validateConfig(c)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func loadYaml() (*config, error) {
	if _, err := os.Stat("config.yml"); os.IsNotExist(err) {
		return nil, errNoConfig
	}
	dat, err := ioutil.ReadFile("config.yml")
	if err != nil {
		return nil, err
	}
	var c config
	yaml.Unmarshal(dat, &c)
	return &c, nil
}

func loadEnv() (*config, error) {
	uwl, err := stringListToInts(os.Getenv("TG_USER_WHITELIST"))
	if err != nil {
		return nil, err
	}
	gwl, err := stringListToInt64s(os.Getenv("TG_GROUP_WHITELIST"))
	if err != nil {
		return nil, err
	}
	mrtrimmed := strings.TrimSpace(os.Getenv("TG_MOTION_RECIPIENTS"))
	mr := strings.Split(mrtrimmed, ",")
	if mrtrimmed == "" {
		mr = make([]string, 0, 0)
	}
	return &config{
		NxIPPort:           os.Getenv("NX_IP_PORT"),
		NxUser:             os.Getenv("NX_USER"),
		NxPass:             os.Getenv("NX_PASS"),
		HTTPIPPort:         os.Getenv("HTTP_IP_PORT"),
		TgToken:            os.Getenv("TG_TOKEN"),
		TgUserWhitelist:    uwl,
		TgGroupWhitelist:   gwl,
		TgMotionRecipients: mr,
	}, nil
}

func stringListToInts(strlist string) ([]int, error) {
	trimmed := strings.TrimSpace(strlist)
	if trimmed == "" {
		return make([]int, 0, 0), nil
	}
	sl := strings.Split(trimmed, ",")
	isl := make([]int, 0, len(sl))
	for _, v := range sl {
		i, err := strconv.Atoi(v)
		if err != nil {
			return nil, err
		}
		isl = append(isl, i)
	}
	return isl, nil
}

func stringListToInt64s(strlist string) ([]int64, error) {
	trimmed := strings.TrimSpace(strlist)
	if trimmed == "" {
		return make([]int64, 0, 0), nil
	}
	sl := strings.Split(trimmed, ",")
	isl := make([]int64, 0, len(sl))
	for _, v := range sl {
		i, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return nil, err
		}
		isl = append(isl, i)
	}
	return isl, nil
}

func applyEnvOverwrite(env *config, yaml *config) *config {
	c := yaml

	if env.NxIPPort != "" {
		c.NxIPPort = env.NxIPPort
	}
	if env.NxUser != "" {
		c.NxUser = env.NxUser
	}
	if env.NxPass != "" {
		c.NxPass = env.NxPass
	}
	if env.HTTPIPPort != "" {
		c.HTTPIPPort = env.HTTPIPPort
	}
	if env.TgToken != "" {
		c.TgToken = env.TgToken
	}
	if len(env.TgUserWhitelist) > 0 {
		c.TgUserWhitelist = env.TgUserWhitelist
	}
	if len(env.TgGroupWhitelist) > 0 {
		c.TgGroupWhitelist = env.TgGroupWhitelist
	}
	if len(env.TgMotionRecipients) > 0 {
		c.TgMotionRecipients = env.TgMotionRecipients
	}

	return c
}

func validateConfig(c *config) error {
	if c.NxIPPort == "" {
		return fmt.Errorf("config: nx_ip_port (yaml) / NX_IP_PORT (env) not set")
	}
	if c.NxUser == "" {
		return fmt.Errorf("config: nx_user (yaml) / NX_USER (env) not set")
	}
	if c.NxPass == "" {
		return fmt.Errorf("config: nx_pass (yaml) / NX_PASS (env) not set")
	}
	if c.HTTPIPPort == "" {
		return fmt.Errorf("config: http_ip_port (yaml) / HTTP_IP_PORT (env) not set")
	}
	if c.TgToken == "" {
		return fmt.Errorf("config: tg_token (yaml) / TG_TOKEN (env) not set")
	}
	if len(c.TgMotionRecipients) == 0 {
		return fmt.Errorf("config: tg_motion_recipients (yaml) / TG_MOTION_RECIPIENTS (env) not set")
	}

	return nil
}
