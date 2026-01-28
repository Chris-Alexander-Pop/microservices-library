// Package cloud provides the core primitives and interfaces for building a Private Cloud (IaaS).
//
// This package encompasses the following domains:
//   - Hypervisor: Virtual Machine management (Libvirt, QEMU, Firecracker)
//   - Provisioning: Bare metal lifecycle (PXE, IPMI)
//   - Scheduler: Intelligent workload placement
//   - Control Plane: Central API and state management
//
// The goal is to provide an AWS-like experience on bare metal infrastructure.
package cloud
