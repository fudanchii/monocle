package build

import (
	"github.com/moby/moby/api/types"
	"github.com/moby/moby/api/types/container"
)

func (b *Build) ToMobyContainerConfig() *container.Config {
	return &container.Config{}
}

func (b *Build) ToMobyHostConfig() *container.HostConfig {
	return &container.HostConfig{}
}

func (b *Build) ToMobyNetworkConfig() *types.NetworkConfig {
	return &types.NetworkConfig{}
}
