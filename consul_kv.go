package onlinelab

import (
	"errors"
	"github.com/hashicorp/consul/api"
	"github.com/json-iterator/go"
)

// NewConsulConfigStorage is ConsulConfigStorage constructor
func NewConsulConfigStorage(consulAPIConfig *api.Config) (*ConsulConfigStorage, error) {
	client, err := api.NewClient(consulAPIConfig)
	if err != nil {
		return nil, err
	}
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	return &ConsulConfigStorage{kv: client.KV(), json: json}, nil
}

// ConsulConfigStorage
type ConsulConfigStorage struct {
	config Config
	kv     *api.KV
	json   jsoniter.API
}

// GetConfig is get config from consul kv. If an error, it will return to the original config
func (cs *ConsulConfigStorage) GetConfig(labName string) (Config, error) {
	pair, _, err := cs.kv.Get(labName, nil)
	if err != nil {
		return cs.config, err
	}
	if pair == nil {
		return cs.config, errors.New("get consul kv is nil")
	}
	if err = cs.json.Unmarshal(pair.Value, &cs.config.treatments); err != nil {
		return cs.config, err
	}
	cs.config.Name = pair.Key
	return cs.config, nil
}

// SetConfig is put config to consul kv
func (cs *ConsulConfigStorage) SetConfig(config Config) {
	// PUT a new KV pair
	value, _ := cs.json.Marshal(config.treatments)
	p := &api.KVPair{Key: config.Name, Value: value}
	cs.kv.Put(p, nil)
	cs.config = config
}
