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

package k8sclient

import (
	"fmt"

	"gopkg.in/k8snetworkplumbingwg/multus-cni.v4/pkg/types"
	v1 "k8s.io/api/core/v1"
)

func resolveDeviceID(
	pod *v1.Pod,
	net *types.NetworkSelectionElement,
	resourceName string,
	alloc *types.PodDeviceAllocation,
	networks []*types.NetworkSelectionElement,
	networkIndex int,
	nadResourceNames []string,
) (string, error) {
	if net != nil && net.DeviceID != "" {
		return net.DeviceID, nil
	}

	containerName, err := resolveContainerName(pod, net, resourceName, networks, networkIndex, nadResourceNames)
	if err != nil {
		return "", err
	}
	if containerName == "" {
		return "", nil
	}

	deviceID, ok := alloc.NextDeviceID(containerName, resourceName)
	if !ok {
		return "", fmt.Errorf("no device allocated for container %q resource %q", containerName, resourceName)
	}
	return deviceID, nil
}

func resolveContainerName(
	pod *v1.Pod,
	net *types.NetworkSelectionElement,
	resourceName string,
	networks []*types.NetworkSelectionElement,
	networkIndex int,
	nadResourceNames []string,
) (string, error) {
	if net != nil && net.ContainerName != "" {
		return net.ContainerName, nil
	}

	containers := containersRequestingResource(pod, resourceName)
	if len(containers) == 0 {
		return "", nil
	}
	if len(containers) == 1 {
		return containers[0], nil
	}

	occurrence := resourceNameOccurrenceIndex(networks, networkIndex, resourceName, nadResourceNames)
	if occurrence < len(containers) {
		return containers[occurrence], nil
	}

	return "", fmt.Errorf("ambiguous device assignment for resource %q: specify containerName in networks annotation", resourceName)
}

func containersRequestingResource(pod *v1.Pod, resourceName string) []string {
	if pod == nil {
		return nil
	}
	var names []string
	for _, container := range pod.Spec.Containers {
		if containerRequestsResource(container, resourceName) {
			names = append(names, container.Name)
		}
	}
	return names
}

func containerRequestsResource(container v1.Container, resourceName string) bool {
	if container.Resources.Limits != nil {
		if _, ok := container.Resources.Limits[v1.ResourceName(resourceName)]; ok {
			return true
		}
	}
	if container.Resources.Requests != nil {
		if _, ok := container.Resources.Requests[v1.ResourceName(resourceName)]; ok {
			return true
		}
	}
	return false
}

func resourceNameOccurrenceIndex(
	networks []*types.NetworkSelectionElement,
	networkIndex int,
	resourceName string,
	nadResourceNames []string,
) int {
	if networkIndex < 0 || networkIndex >= len(networks) || networkIndex >= len(nadResourceNames) {
		return -1
	}
	if nadResourceNames[networkIndex] != resourceName {
		return -1
	}
	occurrence := 0
	for i := 0; i < networkIndex; i++ {
		if i < len(nadResourceNames) && nadResourceNames[i] == resourceName {
			occurrence++
		}
	}
	return occurrence
}
