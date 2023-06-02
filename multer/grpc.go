// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package multer

import (
	"context"

	hclog "github.com/hashicorp/go-hclog"
	plugin "github.com/hashicorp/go-plugin"
	proto "github.com/hashicorp/go-plugin/examples/bidirectional/protob"
	"google.golang.org/grpc"
)

// GRPCClient is an implementation of KV that talks over RPC.
type GRPCClient struct {
	broker *plugin.GRPCBroker
	client proto.MulterClient
}

func (m *GRPCClient) Put(key string, value int64, a AddHelper) error {
	addHelperServer := &GRPCAddHelperServer{Impl: a}

	var s *grpc.Server
	serverFunc := func(opts []grpc.ServerOption) *grpc.Server {
		s = grpc.NewServer(opts...)
		proto.RegisterAddHelperServer(s, addHelperServer)

		return s
	}

	brokerID := m.broker.NextId()
	go m.broker.AcceptAndServe(brokerID, serverFunc)

	_, err := m.client.Put(context.Background(), &proto.PutRequest{
		AddServer: brokerID,
		Key:       key,
		Value:     value,
	})

	s.Stop()
	return err
}

func (m *GRPCClient) Get(key string) (int64, error) {
	resp, err := m.client.Get(context.Background(), &proto.GetRequest{
		Key: key,
	})
	if err != nil {
		return 0, err
	}

	return resp.Value, nil
}

// Here is the gRPC server that GRPCClient talks to.
type GRPCServer struct {
	// This is the real implementation
	Impl Multer

	broker *plugin.GRPCBroker
}

func (m *GRPCServer) Put(ctx context.Context, req *proto.PutRequest) (*proto.Empty, error) {
	conn, err := m.broker.Dial(req.AddServer)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	a := &GRPCAddHelperClient{proto.NewAddHelperClient(conn)}
	return &proto.Empty{}, m.Impl.Put(req.Key, req.Value, a)
}

func (m *GRPCServer) Get(ctx context.Context, req *proto.GetRequest) (*proto.GetResponse, error) {
	v, err := m.Impl.Get(req.Key)
	return &proto.GetResponse{Value: v}, err
}

// GRPCClient is an implementation of KV that talks over RPC.
type GRPCAddHelperClient struct{ client proto.AddHelperClient }

func (m *GRPCAddHelperClient) Sum(a, b int64) (int64, error) {
	resp, err := m.client.Sum(context.Background(), &proto.MultRequest{
		A: a,
		B: b,
	})
	if err != nil {
		hclog.Default().Info("add.Sum", "client", "start", "err", err)
		return 0, err
	}
	return resp.R, err
}

// Here is the gRPC server that GRPCClient talks to.
type GRPCAddHelperServer struct {
	// This is the real implementation
	Impl AddHelper
}

// This is the rpc server function which in turn calls the function in the main app.

func (m *GRPCAddHelperServer) Sum(ctx context.Context, req *proto.MultRequest) (resp *proto.MultResponse, err error) {
	r, err := m.Impl.Sum(req.A, req.B)
	if err != nil {
		return nil, err
	}
	return &proto.MultResponse{R: r}, err
}
