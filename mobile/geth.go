// Copyright 2016 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

// Contains all the wrappers from the node package to support client side node
// management on mobile platforms.

package geth

import (
	"fmt"
	"path"
	"path/filepath"

	"github.com/ethereum/go-ethereum/whisper/shhclient"

	"github.com/ethereum/go-ethereum/p2p/discover"
	"github.com/ethereum/go-ethereum/p2p/discv5"

	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/eth"
	"github.com/ethereum/go-ethereum/eth/downloader"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/ethstats"
	"github.com/ethereum/go-ethereum/les"
	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/p2p/nat"
	"github.com/ethereum/go-ethereum/params"
	whisper "github.com/ethereum/go-ethereum/whisper/whisperv5"
)

// NodeConfig represents the collection of configuration values to fine tune the Geth
// node embedded into a mobile process. The available values are a subset of the
// entire API provided by go-ethereum to reduce the maintenance surface and dev
// complexity.
type NodeConfig struct {
	// MaxPeers is the maximum number of peers that can be connected. If this is
	// set to zero, then only the configured static and trusted peers can connect.
	MaxPeers int

	// EthereumEnabled specifies whether the node should run the Ethereum protocol.
	EthereumEnabled bool

	// EthereumNetworkID is the network identifier used by the Ethereum protocol to
	// decide if remote peers should be accepted or not.
	EthereumNetworkID int64 // uint64 in truth, but Java can't handle that...

	// EthereumDatabaseCache is the system memory in MB to allocate for database caching.
	// A minimum of 16MB is always reserved.
	EthereumDatabaseCache int

	// EthereumNetStats is a netstats connection string to use to report various
	// chain, transaction and node stats to a monitoring server.
	//
	// It has the form "nodename:secret@host:port"
	EthereumNetStats string

	// WhisperEnabled specifies whether the node should run the Whisper protocol.
	WhisperEnabled bool
}

// defaultNodeConfig contains the default node configuration values to use if all
// or some fields are missing from the user's specified list.
var defaultNodeConfig = &NodeConfig{
	MaxPeers:              25,
	EthereumEnabled:       true,
	EthereumNetworkID:     1,
	EthereumDatabaseCache: 128,
}

// NewNodeConfig creates a new node option set, initialized to the default values.
func NewNodeConfig() *NodeConfig {
	config := *defaultNodeConfig
	return &config
}

// Node represents a Geth Ethereum node instance.
type Node struct {
	node *node.Node
}

func getBootstrapNodes(config *NodeConfig) (v5nodes []*discv5.Node, nodes []*discover.Node) {
	switch config.EthereumNetworkID {
	case 3:
		v5nodes = []*discv5.Node{}
		nodes = make([]*discover.Node, len(params.TestnetBootnodes))
		for i, url := range params.TestnetBootnodes {
			nodes[i] = discover.MustParseNode(url)
		}
		break
	case 4:
		v5nodes = make([]*discv5.Node, len(params.RinkebyV5Bootnodes))
		for i, url := range params.RinkebyV5Bootnodes {
			v5nodes[i] = discv5.MustParseNode(url)
		}
		nodes = make([]*discover.Node, len(params.RinkebyBootnodes))
		for i, url := range params.RinkebyBootnodes {
			nodes[i] = discover.MustParseNode(url)
		}
		break
	case 1:
	default:
		v5nodes = make([]*discv5.Node, len(params.DiscoveryV5Bootnodes))
		for i, url := range params.DiscoveryV5Bootnodes {
			v5nodes[i] = discv5.MustParseNode(url)
		}
		nodes = make([]*discover.Node, len(params.MainnetBootnodes))
		for i, url := range params.MainnetBootnodes {
			nodes[i] = discover.MustParseNode(url)
		}
		break
	}

	return
}

func getGenesis(config *NodeConfig) (genesis *core.Genesis) {
	switch config.EthereumNetworkID {
	case 3:
		genesis = core.DefaultTestnetGenesisBlock()
		break
	case 4:
		genesis = core.DefaultRinkebyGenesisBlock()
		break
	case 1:
	default:
		genesis = nil
		break
	}
	return
}

// NewNode creates and configures a new Geth node.
func NewNode(datadir string, config *NodeConfig) (stack *Node, _ error) {
	// If no or partial configurations were specified, use defaults
	if config == nil {
		config = NewNodeConfig()
	}
	if config.MaxPeers == 0 {
		config.MaxPeers = defaultNodeConfig.MaxPeers
	}
	// Create the empty networking stack
	v5BootstrapNodes, bootstrapNodes := getBootstrapNodes(config)
	nodeConf := &node.Config{
		Name:        clientIdentifier,
		Version:     params.Version,
		DataDir:     datadir,
		KeyStoreDir: filepath.Join(datadir, "keystore"), // Mobile should never use internal keystores!
		P2P: p2p.Config{
			NoDiscovery:      false,
			DiscoveryV5:      true,
			ListenAddr:       ":30303",
			DiscoveryV5Addr:  ":30304",
			BootstrapNodes:   bootstrapNodes,
			BootstrapNodesV5: v5BootstrapNodes,
			MaxPeers:         config.MaxPeers,
			NAT:              nat.Any(),
		},
	}
	rawStack, err := node.New(nodeConf)
	if err != nil {
		return nil, err
	}

	// Register the Ethereum protocol if requested
	genesis := getGenesis(config)
	if config.EthereumEnabled {
		ethConf := eth.DefaultConfig
		ethConf.Genesis = genesis
		ethConf.SyncMode = downloader.FastSync
		ethConf.NetworkId = uint64(config.EthereumNetworkID)
		ethConf.DatabaseCache = config.EthereumDatabaseCache
		ethConf.MaxPeers = config.MaxPeers
		ethConf.DatabaseHandles = 1024
		ethConf.EthashCacheDir = path.Join(datadir, ".ethash")
		ethConf.EthashDatasetDir = path.Join(datadir, ".ethash")
		if err := rawStack.Register(func(ctx *node.ServiceContext) (node.Service, error) {
			fullNode, err := eth.New(ctx, &ethConf)
			if fullNode != nil && ethConf.LightServ > 0 {
				ls, _ := les.NewLesServer(fullNode, &ethConf)
				fullNode.AddLesServer(ls)
			}
			return fullNode, err
		}); err != nil {
			return nil, fmt.Errorf("ethereum init: %v", err)
		}

		// If netstats reporting is requested, do it
		if config.EthereumNetStats != "" {
			if err := rawStack.Register(func(ctx *node.ServiceContext) (node.Service, error) {

				var ethServ *eth.Ethereum
				ctx.Service(&ethServ)

				var lesServ *les.LightEthereum
				ctx.Service(&lesServ)

				return ethstats.New(config.EthereumNetStats, ethServ, lesServ)
			}); err != nil {
				return nil, fmt.Errorf("netstats init: %v", err)
			}
		}
	}
	// Register the Whisper protocol if requested
	if config.WhisperEnabled {
		if err := rawStack.Register(func(*node.ServiceContext) (node.Service, error) {
			return whisper.New(&whisper.DefaultConfig), nil
		}); err != nil {
			return nil, fmt.Errorf("whisper init: %v", err)
		}
	}

	return &Node{rawStack}, nil
}

// Start creates a live P2P node and starts running it.
func (n *Node) Start() error {
	return n.node.Start()
}

// Stop terminates a running node along with all it's services. In the node was
// not started, an error is returned.
func (n *Node) Stop() error {
	return n.node.Stop()
}

// GetEthereumClient retrieves a client to access the Ethereum subsystem.
func (n *Node) GetEthereumClient() (client *EthereumClient, _ error) {
	rpc, err := n.node.Attach()
	if err != nil {
		return nil, err
	}
	return &EthereumClient{ethclient.NewClient(rpc)}, nil
}

func (n *Node) GetShhClient() (client *ShhClient, _ error) {
	rpc, err := n.node.Attach()
	if err != nil {
		return nil, err
	}
	return &ShhClient{shhclient.NewClient(rpc)}, nil
}

// GetNodeInfo gathers and returns a collection of metadata known about the host.
func (n *Node) GetNodeInfo() *NodeInfo {
	return &NodeInfo{n.node.Server().NodeInfo()}
}

// GetPeersInfo returns an array of metadata objects describing connected peers.
func (n *Node) GetPeersInfo() *PeerInfos {
	return &PeerInfos{n.node.Server().PeersInfo()}
}
