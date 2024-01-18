package etcd

import (
	"fmt"
	"testing"

	clientV3 "go.etcd.io/etcd/client/v3"
)

func TestClient_Get(t *testing.T) {
	client, err := New(clientV3.Config{
		Endpoints: []string{"127.0.0.1:2379"},
		Username:  "",
		Password:  "",
	})
	if err != nil {
		t.Error(err)
		return
	}

	key := "aa"
	val := "this is aa"
	// create
	fmt.Println("-------test create-------")
	err = client.Create(key, val)
	if err != nil {
		t.Error(err)
		return
	}

	// get
	fmt.Println("-------test get-------")
	response, err := client.Get("aa")
	if err != nil {
		return
	}
	if response == nil {
		t.Errorf("get want not nil,got nil")
	}
	for _, kv := range response {
		fmt.Println(kv.Value)
	}

	// update
	fmt.Println("-------test update-------")
	err = client.Update(key, "this is bb")
	if err != nil {
		t.Error(err)
		return
	}
	response, err = client.Get("aa")
	if err != nil {
		return
	}
	if response == nil {
		t.Errorf("get want not nil,got nil")
	}
	for _, kv := range response {
		fmt.Println(kv.Value)
	}

	// delete
	fmt.Println("-------test delete-------")
	err = client.Delete(key)
	if err != nil {
		t.Error(err)
		return
	}
}
