// Copyright (c) 2021 Multus Authors
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

package types

import "sort"

// PodDeviceAllocation holds per-container device allocations from kubelet, checkpoint, or DRA.
type PodDeviceAllocation struct {
	// ByContainer maps containerName -> resourceName -> devices (+ per-pair Index).
	ByContainer map[string]map[string]*ResourceInfo
}

// NewPodDeviceAllocation returns an empty PodDeviceAllocation.
func NewPodDeviceAllocation() *PodDeviceAllocation {
	return &PodDeviceAllocation{
		ByContainer: make(map[string]map[string]*ResourceInfo),
	}
}

// AddContainerDevices appends device IDs for the given container and resource.
func (a *PodDeviceAllocation) AddContainerDevices(containerName, resourceName string, deviceIDs []string) {
	if len(deviceIDs) == 0 {
		return
	}
	byResource, ok := a.ByContainer[containerName]
	if !ok {
		byResource = make(map[string]*ResourceInfo)
		a.ByContainer[containerName] = byResource
	}
	entry, ok := byResource[resourceName]
	if !ok {
		entry = &ResourceInfo{}
		byResource[resourceName] = entry
	}
	entry.DeviceIDs = append(entry.DeviceIDs, deviceIDs...)
}

// NextDeviceID returns the next unassigned device ID for container/resource and advances the index.
func (a *PodDeviceAllocation) NextDeviceID(containerName, resourceName string) (string, bool) {
	byResource, ok := a.ByContainer[containerName]
	if !ok {
		return "", false
	}
	entry, ok := byResource[resourceName]
	if !ok || entry.Index >= len(entry.DeviceIDs) {
		return "", false
	}
	id := entry.DeviceIDs[entry.Index]
	entry.Index++
	return id, true
}

// SortDeviceIDsPerContainer sorts DeviceIDs within each container/resource pair for deterministic ordering.
func (a *PodDeviceAllocation) SortDeviceIDsPerContainer() {
	for _, byResource := range a.ByContainer {
		for _, info := range byResource {
			if info.DeviceIDs != nil {
				sort.Strings(info.DeviceIDs)
			}
		}
	}
}
