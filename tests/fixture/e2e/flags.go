// Copyright (C) 2019-2024, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package e2e

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cast"

	"github.com/ava-labs/avalanchego/tests/fixture/tmpnet"
)

const (
	// Ensure that this value takes into account the scrape_interval
	// defined in scripts/run_prometheus.sh.
	networkShutdownDelay = 12 * time.Second

	delayNetworkShutdownEnvName = "TMPNET_DELAY_NETWORK_SHUTDOWN"
)

type FlagVars struct {
	avalancheGoExecPath  string
	pluginDir            string
	networkDir           string
	reuseNetwork         bool
	delayNetworkShutdown bool
	startNetwork         bool
	stopNetwork          bool
	restartNetwork       bool
	nodeCount            int
}

func (v *FlagVars) AvalancheGoExecPath() (string, error) {
	if err := v.validateAvalancheGoExecPath(); err != nil {
		return "", err
	}
	return v.avalancheGoExecPath, nil
}

func (v *FlagVars) validateAvalancheGoExecPath() error {
	if !filepath.IsAbs(v.avalancheGoExecPath) {
		absPath, err := filepath.Abs(v.avalancheGoExecPath)
		if err != nil {
			return fmt.Errorf("avalanchego-path (%s) is a relative path but its absolute path cannot be determined: %w",
				v.avalancheGoExecPath, err)
		}

		// If the absolute path file doesn't exist, it means it won't work out of the box.
		if _, err := os.Stat(absPath); err != nil {
			return fmt.Errorf("avalanchego-path (%s) is a relative path but must be an absolute path", v.avalancheGoExecPath)
		}
	}
	return nil
}

func (v *FlagVars) PluginDir() string {
	return v.pluginDir
}

func (v *FlagVars) NetworkDir() string {
	if !v.reuseNetwork {
		return ""
	}
	if len(v.networkDir) > 0 {
		return v.networkDir
	}
	return os.Getenv(tmpnet.NetworkDirEnvName)
}

func (v *FlagVars) ReuseNetwork() bool {
	return v.reuseNetwork
}

func (v *FlagVars) RestartNetwork() bool {
	return v.restartNetwork
}

func (v *FlagVars) NetworkShutdownDelay() time.Duration {
	if v.delayNetworkShutdown {
		// Only return a non-zero value if the delay is enabled.
		return networkShutdownDelay
	}
	return 0
}

func (v *FlagVars) StartNetwork() bool {
	return v.startNetwork
}

func (v *FlagVars) StopNetwork() bool {
	return v.stopNetwork
}

func (v *FlagVars) NodeCount() int {
	return v.nodeCount
}

func GetEnvWithDefault(envVar, defaultVal string) string {
	val := os.Getenv(envVar)
	if len(val) == 0 {
		return defaultVal
	}
	return val
}

func RegisterFlags() *FlagVars {
	vars := FlagVars{}
	flag.StringVar(
		&vars.avalancheGoExecPath,
		"avalanchego-path",
		os.Getenv(tmpnet.AvalancheGoPathEnvName),
		fmt.Sprintf(
			"[optional] avalanchego executable path if creating a new network. Also possible to configure via the %s env variable.",
			tmpnet.AvalancheGoPathEnvName,
		),
	)
	flag.StringVar(
		&vars.pluginDir,
		"plugin-dir",
		GetEnvWithDefault(tmpnet.AvalancheGoPluginDirEnvName, os.ExpandEnv("$HOME/.avalanchego/plugins")),
		fmt.Sprintf(
			"[optional] the dir containing VM plugins. Also possible to configure via the %s env variable.",
			tmpnet.AvalancheGoPluginDirEnvName,
		),
	)
	flag.StringVar(
		&vars.networkDir,
		"network-dir",
		"",
		fmt.Sprintf("[optional] the dir containing the configuration of an existing network to target for testing. Will only be used if --reuse-network is specified. Also possible to configure via the %s env variable.", tmpnet.NetworkDirEnvName),
	)
	flag.BoolVar(
		&vars.reuseNetwork,
		"reuse-network",
		false,
		"[optional] reuse an existing network previously started with --reuse-network. If a network is not already running, create a new one and leave it running for subsequent usage. Ignored if --stop-network is provided.",
	)
	flag.BoolVar(
		&vars.restartNetwork,
		"restart-network",
		false,
		"[optional] restart an existing network previously started with --reuse-network. Useful for ensuring a network is running with the current state of binaries on disk. Ignored if a network is not already running or --stop-network is provided.",
	)
	flag.BoolVar(
		&vars.delayNetworkShutdown,
		"delay-network-shutdown",
		cast.ToBool(GetEnvWithDefault(delayNetworkShutdownEnvName, "false")),
		"[optional] whether to delay network shutdown to allow a final metrics scrape.",
	)
	flag.BoolVar(
		&vars.startNetwork,
		"start-network",
		false,
		"[optional] start a new network and exit without executing any tests. The new network cannot be reused with --reuse-network. Ignored if either --reuse-network or --stop-network is provided.",
	)
	flag.BoolVar(
		&vars.stopNetwork,
		"stop-network",
		false,
		"[optional] stop an existing network started with --reuse-network and exit without executing any tests.",
	)
	flag.IntVar(
		&vars.nodeCount,
		"node-count",
		tmpnet.DefaultNodeCount,
		"number of nodes the network should initially consist of",
	)

	return &vars
}
