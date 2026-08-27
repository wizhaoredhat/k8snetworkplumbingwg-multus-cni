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

package types

import (
	"reflect"
	"testing"
)

func TestPodDeviceAllocationAddAndNext(t *testing.T) {
	alloc := NewPodDeviceAllocation()
	alloc.AddContainerDevices("dpdk1", "intel.com/sriov", []string{"0000:d8:00.5"})
	alloc.AddContainerDevices("dpdk2", "intel.com/sriov", []string{"0000:d8:00.2"})

	id, ok := alloc.NextDeviceID("dpdk1", "intel.com/sriov")
	if !ok {
		t.Fatal("expected device for dpdk1")
	}
	if id != "0000:d8:00.5" {
		t.Fatalf("got %q, want %q", id, "0000:d8:00.5")
	}

	id, ok = alloc.NextDeviceID("dpdk2", "intel.com/sriov")
	if !ok {
		t.Fatal("expected device for dpdk2")
	}
	if id != "0000:d8:00.2" {
		t.Fatalf("got %q, want %q", id, "0000:d8:00.2")
	}
}

func TestPodDeviceAllocationSortPerContainer(t *testing.T) {
	alloc := NewPodDeviceAllocation()
	alloc.AddContainerDevices("app", "intel.com/sriov", []string{"0000:03:02.3", "0000:03:02.0"})
	alloc.SortDeviceIDsPerContainer()

	want := []string{"0000:03:02.0", "0000:03:02.3"}
	got := alloc.ByContainer["app"]["intel.com/sriov"].DeviceIDs
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}
