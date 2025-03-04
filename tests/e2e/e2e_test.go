// Copyright (C) 2019-2024, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package e2e_test

import (
	"testing"

	"github.com/onsi/ginkgo/v2"

	// ensure test packages are scanned by ginkgo
	_ "github.com/ava-labs/avalanchego/tests/e2e/banff"
	_ "github.com/ava-labs/avalanchego/tests/e2e/c"
	_ "github.com/ava-labs/avalanchego/tests/e2e/faultinjection"
	_ "github.com/ava-labs/avalanchego/tests/e2e/p"
	_ "github.com/ava-labs/avalanchego/tests/e2e/x"
	_ "github.com/ava-labs/avalanchego/tests/e2e/x/transfer"

	"github.com/ava-labs/avalanchego/tests/e2e/vms"
	"github.com/ava-labs/avalanchego/tests/fixture/e2e"
	"github.com/ava-labs/avalanchego/tests/fixture/tmpnet"
)

func TestE2E(t *testing.T) {
	ginkgo.RunSpecs(t, "e2e test suites")
}

var flagVars *e2e.FlagVars

func init() {
	flagVars = e2e.RegisterFlags()
}

var _ = ginkgo.SynchronizedBeforeSuite(func() []byte {
	// Run only once in the first ginkgo process

	nodes := tmpnet.NewNodesOrPanic(flagVars.NodeCount())
	subnets := vms.XSVMSubnetsOrPanic(nodes...)
	return e2e.NewTestEnvironment(
		e2e.NewEventHandlerTestContext(),
		flagVars,
		&tmpnet.Network{
			Owner:   "avalanchego-e2e",
			Nodes:   nodes,
			Subnets: subnets,
		},
	).Marshal()
}, func(envBytes []byte) {
	// Run in every ginkgo process

	// Initialize the local test environment from the global state
	e2e.InitSharedTestEnvironment(ginkgo.GinkgoT(), envBytes)
})
