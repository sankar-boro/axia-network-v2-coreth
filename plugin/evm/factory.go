// (c) 2019-2020, Axia Systems, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package evm

import (
	"github.com/sankar-boro/axia/ids"
	"github.com/sankar-boro/axia/snow"
	"github.com/sankar-boro/axia/vms"
)

var (
	// ID this VM should be referenced by
	ID = ids.ID{'e', 'v', 'm'}

	_ vms.Factory = &Factory{}
)

type Factory struct{}

func (f *Factory) New(*snow.Context) (interface{}, error) {
	return &VM{}, nil
}
