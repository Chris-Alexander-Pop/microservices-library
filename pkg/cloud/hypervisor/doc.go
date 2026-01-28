// Package hypervisor provides interfaces and adapters for managing Virtual Machines.
//
// It supports creating, starting, stopping, and deleting VMs across different backends.
//
// Basic usage:
//
//	import (
//		"github.com/chris-alexander-pop/system-design-library/pkg/cloud/hypervisor"
//		"github.com/chris-alexander-pop/system-design-library/pkg/cloud/hypervisor/adapters/memory"
//	)
//
//	hyp := memory.New()
//	id, err := hyp.CreateVM(ctx, hypervisor.VMSpec{
//		Name: "test-vm",
//		InstanceType: cloud.InstanceTypeSmall,
//	})
package hypervisor
