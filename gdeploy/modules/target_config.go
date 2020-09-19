package modules

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type TargetConfig struct {
	Projects []*Project `yaml:"projects"`
}

func LoadTargetConfig(path string) (*TargetConfig, error) {
	var cfg TargetConfig
	dat, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read target file %s error. err: %v", path, err)
	}

	err = yaml.Unmarshal(dat, &cfg.Projects)
	if err != nil {
		return nil, fmt.Errorf("parse target file %s error. err: %v", path, err)
	}

	if cfg.Projects == nil {
		var project *Project
		err = yaml.Unmarshal(dat, &project)
		if err != nil {
			return nil, fmt.Errorf("parse target file %s error. err: %v", path, err)
		}
		if project != nil {
			cfg.Projects = make([]*Project, 0)
			cfg.Projects = append(cfg.Projects, project)
		}
	}
	return &cfg, nil
}
