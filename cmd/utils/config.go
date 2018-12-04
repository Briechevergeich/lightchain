package utils

import (
	ethUtils "github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/consensus/ethash"
	"github.com/ethereum/go-ethereum/eth"
	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/params"
	"gopkg.in/urfave/cli.v1"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

const (
	// Client identifier to advertise over the network
	clientIdentifier = "lightchain"
	// Environment variable for home dir
	emHome = "EMHOME"
)

var (
	// GenesisTargetGasLimit is the target gas limit of the Genesis block.
	// #unstable
	GenesisTargetGasLimit = uint64(100000000)
)

// General settings
var GenesisPathFlag = ethUtils.DirectoryFlag{
	Name:  "genesis",
	Usage: "Genesis path",
}

var projectRootPath = filepath.Join(os.Getenv("GOPATH"), "src/github.com/lightstreams-network", "lightchain")

type ethstatsConfig struct {
	URL string `toml:",omitempty"`
}

type gethConfig struct {
	Eth      eth.Config
	Node     node.Config
	Ethstats ethstatsConfig
}

// makeDataDir retrieves the currently requested data directory
// #unstable
func MakeDataDir(ctx *cli.Context) string {
	path := node.DefaultDataDir()

	emHome := os.Getenv(emHome)
	if emHome != "" {
		path = emHome
	}

	if ctx.GlobalIsSet(ethUtils.DataDirFlag.Name) {
		path = ctx.GlobalString(ethUtils.DataDirFlag.Name)
	}

	if path == "" {
		ethUtils.Fatalf("Cannot determine default data directory, please set manually (--datadir)")
	}

	return path
}

// DefaultNodeConfig returns the default configuration for a go-ethereum node
// #unstable
func DefaultNodeConfig() node.Config {
	cfg := node.DefaultConfig
	cfg.Name = clientIdentifier
	cfg.Version = params.Version
	cfg.HTTPModules = append(cfg.HTTPModules, "eth")
	cfg.WSModules = append(cfg.WSModules, "eth")
	cfg.IPCPath = "geth.ipc"

	emHome := os.Getenv(emHome)
	if emHome != "" {
		cfg.DataDir = emHome
	}

	return cfg
}

// SetLightchainNodeConfig takes a node configuration and applies lightchain specific configuration
// #unstable
func SetLightchainNodeConfig(cfg *node.Config) {
	cfg.P2P.MaxPeers = 0
	cfg.P2P.NoDiscovery = true
}

// SetLightchainEthConfig takes a ethereum configuration and applies lightchain specific configuration
// #unstable
func SetLightchainEthConfig(cfg *eth.Config) {
	// @TODO(ggarri): Replace hardcoded PoW mode by cmd argument
	cfg.Ethash.PowMode = ethash.ModeFake
	// (ggarri): Properties had been removed and seem to not have a correspondence on the newer version
	//cfg.MaxPeers = 0
}

func MakeGenesisPath(ctx *cli.Context) string {
	genesisPath := ctx.Args().First()
	if genesisPath != "" {
		return genesisPath
	} else if ctx.GlobalIsSet(GenesisPathFlag.Name) {
		genesisPath = ctx.GlobalString(GenesisPathFlag.Name)
	} else {
		lightchainDataDir := MakeDataDir(ctx)
		genesisPath = path.Join(lightchainDataDir, "genesis.json")
	}

	return genesisPath
}

func ReadGenesisPath(genesisPath string) ([]byte, error) {
	genesisBlob, err := ioutil.ReadFile(genesisPath)
	if err != nil {
		return nil, err
	}

	return genesisBlob, nil
}

func ReadDefaultGenesis() ([]byte, error) {
	path, err := filepath.Abs(filepath.Join(projectRootPath, "setup/genesis.json"))
	if err != nil {
		return nil, err
	}
	return ReadGenesisPath(path)
}


func ReadDefaultKeystore() (map[string][]byte, error) {
	path, err := filepath.Abs(filepath.Join(projectRootPath, "setup/keystore"))
	if err != nil {
		return nil, err
	}

	var files = make(map[string][]byte)
	err = filepath.Walk(path, func(file string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		
		if info.IsDir() {
			return nil
		}
		content, err := ioutil.ReadFile(file)
		if err != nil {
			return err
		}
		files[info.Name()] = content
		return nil
	})

	if err != nil {
		panic(err)
	}

	return files, nil
}