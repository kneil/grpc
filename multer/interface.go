// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// Package shared contains shared data between the host and plugins.
package multer

import (
	"context"

	"google.golang.org/grpc"

	"github.com/hashicorp/go-plugin"
	proto "github.com/hashicorp/go-plugin/examples/bidirectional/protob"
)

// Handshake is a common handshake that is shared by plugin and host.
var Handshake = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "BASIC_PLUGIN",
	MagicCookieValue: "hello",
}

// PluginMap is the map of plugins we can dispense.
var PluginMap = map[string]plugin.Plugin{
	"multer": &MulterPlugin{},
}

type AddHelper interface {
	Sum(int64, int64) (int64, error)
}

// KV is the interface that we're exposing as a plugin.
type Multer interface {
	Put(key string, value int64, a AddHelper) error
	Get(key string) (int64, error)
}

// This is the implementation of plugin.Plugin so we can serve/consume this.
// We also implement GRPCPlugin so that this plugin can be served over
// gRPC.
type MulterPlugin struct {
	plugin.NetRPCUnsupportedPlugin
	// Concrete implementation, written in Go. This is only used for plugins
	// that are written in Go.
	Impl Multer
}

func (p *MulterPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	proto.RegisterMulterServer(s, &GRPCServer{
		Impl:   p.Impl,
		broker: broker,
	})
	return nil
}

func (p *MulterPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &GRPCClient{
		client: proto.NewMulterClient(c),
		broker: broker,
	}, nil
}

var _ plugin.GRPCPlugin = &MulterPlugin{}