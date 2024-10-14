package consul

import (
	"fmt"
	"testing"
	"time"

	consulapi "github.com/hashicorp/consul/api"

	"github.com/agclqq/prow-framework/confctr"
)

func TestNew(t *testing.T) {
	t.Skip()
	client, err := New(&consulapi.Config{
		Address:    "127.0.0.1:8500",
		Datacenter: "dc1",
	})
	if err != nil {
		return
	}
	go func() {
		client.Watch("aa", func(eventName string, key string, val string) {
			fmt.Printf("[%s] %s %s\n", eventName, key, val)
		})
	}()
	time.Sleep(1 * time.Second)
	fmt.Println("-------test create-------")
	err = client.Create("aa", "this is aa")
	if err != nil {
		t.Error(err)
	}

	fmt.Println("-------test get-------")
	res, err := client.Get("aa")
	if err != nil {
		t.Error(err)
	}
	for _, v := range res {
		fmt.Println(v.Value)
	}

	fmt.Println("-------test update-------")
	err = client.Update("aa", "this is bb")
	if err != nil {
		t.Error(err)
	}
	time.Sleep(1 * time.Second)
	fmt.Println("-------test delete-------")
	err = client.Delete("aa")
	if err != nil {
		t.Error(err)
	}
	time.Sleep(10 * time.Second)
}

func TestClient_Watch(t *testing.T) {
	t.Skip()
	type args struct {
		key    string
		f      confctr.CallBack
		option *confctr.WatchOption
	}
	config := &consulapi.Config{
		Address:    "127.0.0.1:8500",
		Datacenter: "dc1",
	}
	nc, err := consulapi.NewClient(config)
	if err != nil {
		t.Error(err)
		return
	}
	cli := &Client{
		conf:   config,
		consul: nc,
	}
	if err != nil {
		t.Error(err)
		return
	}
	err = cli.Create("aa", "create aa")
	if err != nil {
		t.Error(err)
	}

	tests := []struct {
		name    string
		client  *Client
		args    args
		wantErr bool
	}{
		{name: "test1", client: cli, args: args{key: "aa", f: func(eventName string, key string, val string) { t.Log(val, eventName, key) }, option: &confctr.WatchOption{Type: confctr.WatchTypeKey}}, wantErr: false},
		{name: "test2", client: cli, args: args{key: "aa/bb", f: func(eventName string, key string, val string) {
			t.Log(val, eventName, key)
		}, option: &confctr.WatchOption{Type: confctr.WatchTypeKeyPrefix}}, wantErr: false},
		{name: "test3", client: cli, args: args{key: "aa", f: func(eventName string, key string, val string) {
			t.Log(val, eventName, key)
		}, option: &confctr.WatchOption{Type: confctr.WatchTypeEvent}}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.client.Watch(tt.args.key, tt.args.f, tt.args.option); (err != nil) != tt.wantErr {
				t.Errorf("Watch() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	go func() {
		err = cli.Create("bb", "create bb")
		if err != nil {
			t.Error(err)
		}
		err = cli.Create("aa/bb", "create aa/bb")
		if err != nil {
			t.Error(err)
		}
		err = cli.Create("aa/bb/cc", "create aa/bb/cc")
		if err != nil {
			t.Error(err)
		}
		err = cli.Update("aa", "update aa to aa1")
		if err != nil {
			t.Error(err)
		}
		err = cli.Update("bb", "update bb to bb1")
		if err != nil {
			t.Error(err)
		}
		err = cli.Update("aa/bb", "update aa/bb to aa/bb1")
		if err != nil {
			t.Error(err)
		}
		err = cli.Delete("aa")
		if err != nil {
			t.Error(err)
		}
		err = cli.Delete("bb")
		if err != nil {
			t.Error(err)
		}
		err = cli.Delete("aa/bb")
		if err != nil {
			t.Error(err)
		}
		err = cli.Delete("aa/bb/cc")
		if err != nil {
			t.Error(err)
		}
	}()
	time.Sleep(time.Second * 1)

	time.Sleep(time.Second * 1)
}
