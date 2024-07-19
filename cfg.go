package main

import (
	"public_service/config"
	"strings"
)

type ClientCfg struct {
	cfg *config.Config
}

func NewClientCfg(configPath string) (*ClientCfg, error) {
	cfg, err := config.LoadConfigFile(configPath)
	if err != nil {
		return nil, err
	}

	return &ClientCfg{
		cfg: cfg,
	}, err
}

func (cc *ClientCfg) GetMasters() []string {
	mstr := cc.cfg.GetString("masterAddr")
	return strings.Split(mstr, ",")
}

func (cc *ClientCfg) GetPublicServer() string {
	publicServer := cc.cfg.GetString("publicServer")
	return publicServer
}

func (cc *ClientCfg) GetGit() string {
	git := cc.cfg.GetString("git")
	return git
}

func (cc *ClientCfg) GetGitLocal() string {
	git := cc.cfg.GetString("git_local")
	return git
}
