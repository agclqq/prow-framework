package consul

import (
	"errors"
	"fmt"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/api/watch"

	"github.com/agclqq/prow-framework/confctr"
)

//Values in the CC store cannot be larger than 512kb.

type Client struct {
	conf   *consulapi.Config
	consul *consulapi.Client
}

func New(config *consulapi.Config) (*Client, error) {
	consul, err := consulapi.NewClient(config)
	return &Client{conf: config, consul: consul}, err
}
func (c *Client) Origin() interface{} {
	return c.consul
}

func (c *Client) Get(key string) ([]confctr.Val, error) {
	res, _, err := c.consul.KV().Get(key, nil)
	if err != nil {
		return nil, err
	}
	val := make([]confctr.Val, 0)
	val = append(val, confctr.Val{Value: string(res.Value)})
	return val, nil
}

func (c *Client) Create(key, value string) error {
	_, err := c.consul.KV().Put(&consulapi.KVPair{Key: key, Value: []byte(value)}, nil)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) Update(key, value string) error {
	_, err := c.consul.KV().Put(&consulapi.KVPair{Key: key, Value: []byte(value)}, nil)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) Delete(key string) error {
	_, err := c.consul.KV().Delete(key, nil)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) Watch(key string, f confctr.CallBack, options ...*confctr.WatchOption) error {
	//https://developer.hashicorp.com/consul/docs/dynamic-app-config/watches#key

	if key == "" {
		return errors.New("key is empty")
	}
	var plan *watch.Plan
	var err error
	var option *confctr.WatchOption
	if len(options) == 0 {
		option = &confctr.WatchOption{}
	} else {
		option = options[0]
	}
	if option.Type == "" {
		option.Type = "key"
	}

	data := make(map[string]interface{})
	//public Watch Option
	if option.Datacenter != "" {
		data["datacenter"] = option.Datacenter
	}
	if option.Token != "" {
		data["token"] = option.Token
	}
	if option.Args != nil {
		data["args"] = option.Args
	}
	if option.Tag != nil {
		data["tag"] = option.Tag
	}

	switch option.Type {
	case "key":
		data["type"] = "key"
		data["key"] = key
	case "keyprefix":
		data["type"] = "keyprefix"
		data["prefix"] = key
	case "services":
		data["type"] = "services"
	case "nodes":
		data["type"] = "nodes"
	case "service":
		data["type"] = "service"
		data["service"] = key
	case "checks":
		data["type"] = "checks"
		data["state"] = key
	case "event":
		data["type"] = "event"
		data["name"] = key
	default:
		return errors.New("type is not supported")
	}
	plan, err = watch.Parse(data)

	plan.Handler = func(idx uint64, val interface{}) {
		if val == nil {
			f("delete", key, "")
		}
		switch v := val.(type) {
		case *consulapi.KVPair:
			t := "create"
			if v.CreateIndex != v.ModifyIndex {
				t = "put"
			}
			f(t, v.Key, string(v.Value))
		case consulapi.KVPairs:
			for _, kv := range v {
				t := "create"
				if kv.CreateIndex != kv.ModifyIndex {
					t = "put"
				}
				f(t, kv.Key, string(kv.Value))
			}
		// services
		case []*consulapi.CatalogService:
			for _, service := range v {
				f(string(option.Type), service.ServiceName, service.ServiceAddress)
			}
		case *consulapi.CatalogService:
			f(string(option.Type), v.ServiceName, v.ServiceAddress)
		case []*consulapi.Node:
			for _, node := range v {
				f(string(option.Type), node.Node, node.Address)
			}
		case []*consulapi.AgentCheck:
			for _, check := range v {
				f(string(option.Type), check.Name, check.Status)
			}
		case []*consulapi.UserEvent:
			for _, event := range v {
				f(string(option.Type), event.Name, string(event.Payload))
			}
		default:
			fmt.Printf("%v\n", v)
		}
	}
	go func() {
		err = plan.RunWithConfig(c.conf.Address, c.conf)
	}()
	return err
}
