package config

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"time"
)

type CfgVolatile struct {
	TimeBankering    time.Duration
	TimeChooseBanker time.Duration
	TimeBet          time.Duration
	TimeCloseBet     time.Duration
	TimeBalance      time.Duration
	TimeReward       time.Duration

	NextResult struct {
		D1 int
		D2 int
		D3 int
	}
}

type Cfg struct {
	LocalAddress  string
	CenterAddress string
	MaxConnection int
	Volatile      CfgVolatile
}

func Load(dir string) (*Cfg, error) {
	if !filepath.IsAbs(dir) {
		absDir, err := filepath.Abs(dir)
		if err != nil {
			return nil, err
		}
		dir = absDir
	}

	bin, err := ioutil.ReadFile(dir)
	if err != nil {
		return nil, err
	}

	cfg := new(Cfg)
	err = json.Unmarshal(bin, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
