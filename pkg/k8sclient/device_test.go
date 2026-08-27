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
	"strings"
	"testing"

	"gopkg.in/k8snetworkplumbingwg/multus-cni.v4/pkg/types"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestResolveContainerNameExplicit(t *testing.T) {
	pod := &v1.Pod{
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{Name: "dpdk1", Resources: v1.ResourceRequirements{Limits: v1.ResourceList{
					"intel.com/sriov": resource.MustParse("1"),
				}}},
				{Name: "dpdk2", Resources: v1.ResourceRequirements{Limits: v1.ResourceList{
					"intel.com/sriov": resource.MustParse("1"),
				}}},
			},
		},
	}
	net := &types.NetworkSelectionElement{ContainerName: "dpdk2"}
	name, err := resolveContainerName(pod, net, "intel.com/sriov", nil, 0, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if name != "dpdk2" {
		t.Fatalf("got %q, want dpdk2", name)
	}
}

func TestResolveContainerNameInferByOrder(t *testing.T) {
	pod := &v1.Pod{
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{Name: "dpdk1", Resources: v1.ResourceRequirements{Limits: v1.ResourceList{
					"intel.com/sriov": resource.MustParse("1"),
				}}},
				{Name: "dpdk2", Resources: v1.ResourceRequirements{Limits: v1.ResourceList{
					"intel.com/sriov": resource.MustParse("1"),
				}}},
			},
		},
	}
	networks := []*types.NetworkSelectionElement{
		{Name: "net-a"},
		{Name: "net-b"},
	}
	nadResourceNames := []string{"intel.com/sriov", "intel.com/sriov"}

	name, err := resolveContainerName(pod, networks[1], "intel.com/sriov", networks, 1, nadResourceNames)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if name != "dpdk2" {
		t.Fatalf("got %q, want dpdk2", name)
	}
}

func TestResolveContainerNameAmbiguous(t *testing.T) {
	pod := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{Name: "p"},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{Name: "dpdk1", Resources: v1.ResourceRequirements{Limits: v1.ResourceList{
					"intel.com/sriov": resource.MustParse("1"),
				}}},
				{Name: "dpdk2", Resources: v1.ResourceRequirements{Limits: v1.ResourceList{
					"intel.com/sriov": resource.MustParse("1"),
				}}},
			},
		},
	}
	networks := []*types.NetworkSelectionElement{
		{Name: "net-a"},
		{Name: "net-b"},
		{Name: "net-c"},
	}
	nadResourceNames := []string{"intel.com/sriov", "intel.com/sriov", "intel.com/sriov"}

	_, err := resolveContainerName(pod, networks[2], "intel.com/sriov", networks, 2, nadResourceNames)
	if err == nil {
		t.Fatal("expected ambiguity error")
	}
	if !strings.Contains(err.Error(), "ambiguous device assignment") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestResolveDeviceIDPerContainer(t *testing.T) {
	pod := &v1.Pod{
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{Name: "dpdk1", Resources: v1.ResourceRequirements{Limits: v1.ResourceList{
					"intel.com/sriov": resource.MustParse("1"),
				}}},
				{Name: "dpdk2", Resources: v1.ResourceRequirements{Limits: v1.ResourceList{
					"intel.com/sriov": resource.MustParse("1"),
				}}},
			},
		},
	}
	alloc := types.NewPodDeviceAllocation()
	alloc.AddContainerDevices("dpdk1", "intel.com/sriov", []string{"0000:d8:00.5"})
	alloc.AddContainerDevices("dpdk2", "intel.com/sriov", []string{"0000:d8:00.2"})

	networks := []*types.NetworkSelectionElement{{Name: "net-a"}, {Name: "net-b"}}
	nadResourceNames := []string{"intel.com/sriov", "intel.com/sriov"}

	id, err := resolveDeviceID(pod, networks[0], "intel.com/sriov", alloc, networks, 0, nadResourceNames)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id != "0000:d8:00.5" {
		t.Fatalf("got %q, want 0000:d8:00.5", id)
	}

	id, err = resolveDeviceID(pod, networks[1], "intel.com/sriov", alloc, networks, 1, nadResourceNames)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id != "0000:d8:00.2" {
		t.Fatalf("got %q, want 0000:d8:00.2", id)
	}
}
