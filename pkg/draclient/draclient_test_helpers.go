// Copyright (c) 2026 Multus Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package draclient

import (
	. "github.com/onsi/gomega"
	omegatypes "github.com/onsi/gomega/types"

	"gopkg.in/k8snetworkplumbingwg/multus-cni.v4/pkg/types"
	v1 "k8s.io/api/core/v1"
)

const draTestContainerName = "test-container"

func podContainerWithClaim(claimRef string) []v1.Container {
	return []v1.Container{{
		Name: draTestContainerName,
		Resources: v1.ResourceRequirements{
			Claims: []v1.ResourceClaim{{Name: claimRef}},
		},
	}}
}

func expectContainerDevices(alloc *types.PodDeviceAllocation, container, resource string, matcher omegatypes.GomegaMatcher) {
	Expect(alloc.ByContainer).To(HaveKey(container))
	Expect(alloc.ByContainer[container]).To(HaveKey(resource))
	Expect(alloc.ByContainer[container][resource].DeviceIDs).To(matcher)
}
