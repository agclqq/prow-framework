package servicediscover

import (
	"context"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/endpoints"
	"go.etcd.io/etcd/client/v3/naming/resolver"
	grpcresolver "google.golang.org/grpc/resolver"
)

type EtcdSD struct {
	etcd        *clientv3.Client
	manager     endpoints.Manager
	ServiceName string
	ttl         int64
	lease       *clientv3.LeaseGrantResponse
}

func NewEtcdDiscover(etcdServer []string, username, password, serviceName string) (*EtcdSD, error) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   etcdServer,
		Username:    username,
		Password:    password,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	manager, err := endpoints.NewManager(client, serviceName)
	if err != nil {
		return nil, err
	}
	dis := &EtcdSD{
		etcd:        client,
		manager:     manager,
		ServiceName: serviceName,
	}
	return dis, err
}
func (d *EtcdSD) Register(ctx context.Context, key string, val string, ttl int64) error {
	lease, err := d.etcd.Grant(context.Background(), ttl)
	if err != nil {
		return nil
	}
	d.ttl = ttl
	d.lease = lease
	return d.manager.AddEndpoint(ctx, d.ServiceName+"/"+key, endpoints.Endpoint{Addr: val, Metadata: nil}, clientv3.WithLease(lease.ID))
}

func (d *EtcdSD) Renew(ctx context.Context) error {
	_, err := d.etcd.KeepAliveOnce(ctx, d.lease.ID)
	if err != nil {
		return err
	}
	return nil
}

func (d *EtcdSD) AutoRenew(ctx context.Context) error {
	renewTime := d.ttl / 2
	if renewTime <= 1 {
		renewTime = 1
	}
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(time.Duration(renewTime) * time.Second):
			err := d.Renew(ctx)
			if err != nil {
				return err
			}
		}
	}
}

func (d *EtcdSD) UnRegister(ctx context.Context, key string) error {
	return d.manager.DeleteEndpoint(ctx, d.ServiceName+"/"+key)
}

func (d *EtcdSD) Watch(ctx context.Context) (endpoints.WatchChannel, error) {
	return d.manager.NewWatchChannel(ctx)
}
func (d *EtcdSD) Resolve() (grpcresolver.Builder, error) {
	return resolver.NewBuilder(d.etcd)
}
func (d *EtcdSD) Close() error {
	return d.etcd.Close()
}
