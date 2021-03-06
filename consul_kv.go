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

// ConsulConfigStorage is the ConfigStorage based on Consul KV
// Configuration format is following:
// [{"Name":"T1","VolumeProportion":30},{"Name":"T2","VolumeProportion":70}]
type ConsulConfigStorage struct {
	kv   *api.KV
	json jsoniter.API
}

// GetConfig is get config from consul kv.
func (cs *ConsulConfigStorage) GetConfig(labName string) (*Config, error) {
	pair, _, err := cs.kv.Get(labName, nil)
	if err != nil {
		return &Config{}, err
	}
	if pair == nil {
		return &Config{}, errors.New("get consul kv is nil")
	}
	config := &Config{}
	if err = cs.json.Unmarshal(pair.Value, &config.treatments); err != nil {
		return config, err
	}
	return config, nil
}

// SetConfig is put config to consul kv
func (cs *ConsulConfigStorage) SetConfig(labName string, config Config) {
	// PUT a new KV pair
	value, _ := cs.json.Marshal(config.treatments)
	p := &api.KVPair{Key: labName, Value: value}
	cs.kv.Put(p, nil)
}
