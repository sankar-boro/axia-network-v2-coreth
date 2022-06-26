// (c) 2019-2020, Axia Systems, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package main

import (
	"fmt"
	"os"

	"github.com/sankar-boro/axia/utils/logging"
	"github.com/sankar-boro/axia/utils/ulimit"
	"github.com/sankar-boro/axia/vms/rpcchainvm"

	"github.com/sankar-boro/coreth/plugin/evm"
)

func main() {
	version, err := PrintVersion()
	if err != nil {
		fmt.Printf("couldn't get config: %s", err)
		os.Exit(1)
	}
	if version {
		fmt.Println(evm.Version)
		os.Exit(0)
	}
	if err := ulimit.Set(ulimit.DefaultFDLimit, logging.NoLog{}); err != nil {
		fmt.Printf("failed to set fd limit correctly due to: %s", err)
		os.Exit(1)
	}

	rpcchainvm.Serve(&evm.VM{IsPlugin: true})
}
