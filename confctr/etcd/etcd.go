package etcd

import (
	"context"

	clientV3 "go.etcd.io/etcd/client/v3"

	"github.com/agclqq/prow-framework/confctr"
)

type Client struct {
	etcd *clientV3.Client
}

func New(config clientV3.Config) (*Client, error) {
	client, err := clientV3.New(config)
	return &Client{client}, err
}
func (cli *Client) Origin() interface{} {
	return cli.etcd
}
func (cli *Client) Get(key string) ([]confctr.Val, error) {
	return cli.GetWithCtx(context.Background(), key)
}
func (cli *Client) GetWithCtx(ctx context.Context, key string) ([]confctr.Val, error) {
	res, err := cli.etcd.KV.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	var rs = make([]confctr.Val, 0)
	for _, ev := range res.Kvs {
		rs = append(rs, confctr.Val{Value: string(ev.Value)})
	}
	return rs, nil
}

func (cli *Client) Create(key, val string) error {
	return cli.CreateWithCtx(context.Background(), key, val)
}
func (cli *Client) CreateWithCtx(ctx context.Context, key, val string) error {
	_, err := cli.etcd.KV.Put(ctx, key, val)
	if err != nil {
		return err
	}
	return nil
}

func (cli *Client) Update(key, val string) error {
	return cli.CreateWithCtx(context.Background(), key, val)
}
func (cli *Client) UpdateWithCtx(ctx context.Context, key, val string) error {
	_, err := cli.etcd.KV.Put(ctx, key, val)
	if err != nil {
		return err
	}
	return nil
}

func (cli *Client) Delete(key string) error {
	return cli.DeleteWithCtx(context.Background(), key)
}
func (cli *Client) DeleteWithCtx(ctx context.Context, key string) error {
	_, err := cli.etcd.KV.Delete(ctx, key)
	if err != nil {
		return err
	}
	return nil
}

func (cli *Client) Watch(key string, f confctr.CallBack, option ...*confctr.WatchOption) error {
	watchChan := cli.etcd.Watch(context.Background(), key)
	for watchResp := range watchChan {
		for _, event := range watchResp.Events {
			eventName := ""
			if event.Type == clientV3.EventTypePut {
				eventName = "put"
			}
			if event.Type == clientV3.EventTypeDelete {
				eventName = "delete"
			}
			f(eventName, string(event.Kv.Key), string(event.Kv.Value))
		}
	}
	return nil
}
