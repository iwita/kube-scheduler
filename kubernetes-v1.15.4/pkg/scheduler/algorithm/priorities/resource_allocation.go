/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package priorities

import (
	"fmt"

	"github.com/iwita/watchapp/pkg/cache"
	v1 "k8s.io/api/core/v1"
	utilfeature "k8s.io/apiserver/pkg/util/feature"
	"k8s.io/klog"
	"k8s.io/kubernetes/pkg/features"
	priorityutil "k8s.io/kubernetes/pkg/scheduler/algorithm/priorities/util"
	schedulerapi "k8s.io/kubernetes/pkg/scheduler/api"
	schedulernodeinfo "k8s.io/kubernetes/pkg/scheduler/nodeinfo"
)

// ResourceAllocationPriority contains information to calculate resource allocation priority.
type ResourceAllocationPriority struct {
	Name   string
	scorer func(requested, allocable *schedulernodeinfo.Resource, includeVolumes bool, requestedVolumes int, allocatableVolumes int) int64
}

type CustomAllocationPriority struct {
	Name   string
	scorer func(nodeName string, pod *v1.Pod) (float64, int, int, *cache.Stress, float64, int32, error)
}

// PriorityMap priorities nodes according to the resource allocations on the node.
// It will use `scorer` function to calculate the score.
func (r *ResourceAllocationPriority) PriorityMap(
	pod *v1.Pod,
	meta interface{},
	nodeInfo *schedulernodeinfo.NodeInfo) (schedulerapi.HostPriority, error) {
	node := nodeInfo.Node()
	if node == nil {
		return schedulerapi.HostPriority{}, fmt.Errorf("node not found")
	}
	allocatable := nodeInfo.AllocatableResource()

	var requested schedulernodeinfo.Resource
	if priorityMeta, ok := meta.(*priorityMetadata); ok {
		requested = *priorityMeta.nonZeroRequest
	} else {
		// We couldn't parse metadata - fallback to computing it.
		requested = *getNonZeroRequests(pod)
	}

	requested.MilliCPU += nodeInfo.NonZeroRequest().MilliCPU
	requested.Memory += nodeInfo.NonZeroRequest().Memory
	var score int64
	// Check if the pod has volumes and this could be added to scorer function for balanced resource allocation.
	if len(pod.Spec.Volumes) >= 0 && utilfeature.DefaultFeatureGate.Enabled(features.BalanceAttachedNodeVolumes) && nodeInfo.TransientInfo != nil {
		score = r.scorer(&requested, &allocatable, true, nodeInfo.TransientInfo.TransNodeInfo.RequestedVolumes, nodeInfo.TransientInfo.TransNodeInfo.AllocatableVolumesCount)
	} else {
		score = r.scorer(&requested, &allocatable, false, 0, 0)
	}

	if klog.V(10) {
		if len(pod.Spec.Volumes) >= 0 && utilfeature.DefaultFeatureGate.Enabled(features.BalanceAttachedNodeVolumes) && nodeInfo.TransientInfo != nil {
			klog.Infof(
				"%v -> %v: %v, capacity %d millicores %d memory bytes, %d volumes, total request %d millicores %d memory bytes %d volumes, score %d",
				pod.Name, node.Name, r.Name,
				allocatable.MilliCPU, allocatable.Memory, nodeInfo.TransientInfo.TransNodeInfo.AllocatableVolumesCount,
				requested.MilliCPU, requested.Memory,
				nodeInfo.TransientInfo.TransNodeInfo.RequestedVolumes,
				score,
			)
		} else {
			klog.Infof(
				"%v -> %v: %v, capacity %d millicores %d memory bytes, total request %d millicores %d memory bytes, score %d",
				pod.Name, node.Name, r.Name,
				allocatable.MilliCPU, allocatable.Memory,
				requested.MilliCPU, requested.Memory,
				score,
			)
		}
	}

	return schedulerapi.HostPriority{
		Host:  node.Name,
		Score: float64(score),
	}, nil
}

func (r *CustomAllocationPriority) PriorityMap(
	pod *v1.Pod,
	meta interface{},
	nodeInfo *schedulernodeinfo.NodeInfo) (schedulerapi.HostPriority, error) {
	node := nodeInfo.Node()
	if node == nil {
		return schedulerapi.HostPriority{}, fmt.Errorf("node not found")
	}
	//allocatable := nodeInfo.AllocatableResource()

	//var requested schedulernodeinfo.Resource
	// if priorityMeta, ok := meta.(*priorityMetadata); ok {
	// 	requested = *priorityMeta.nonZeroRequest
	// } else {
	// 	// We couldn't parse metadata - fallback to computing it.
	// 	requested = *getNonZeroRequests(pod)
	// }

	//requested.MilliCPU += nodeInfo.NonZeroRequest().MilliCPU
	//requested.Memory += nodeInfo.NonZeroRequest().Memory
	var score, finalScore float64
	var socket, socketCoresNum int
	var time int32
	var stress *cache.Stress
	// Check if the pod has volumes and this could be added to scorer function for balanced resource allocation.
	// if len(pod.Spec.Volumes) >= 0 && utilfeature.DefaultFeatureGate.Enabled(features.BalanceAttachedNodeVolumes) && nodeInfo.TransientInfo != nil {
	// 	score = r.scorer(&requested, &allocatable, true, nodeInfo.TransientInfo.TransNodeInfo.RequestedVolumes, nodeInfo.TransientInfo.TransNodeInfo.AllocatableVolumesCount)
	// } else {
	// 	score = r.scorer(&requested, &allocatable, false, 0, 0)
	// }

	score, socket, socketCoresNum, stress, finalScore, time, _ = r.scorer(node.Name, pod)

	// if klog.V(10) {
	// 	if len(pod.Spec.Volumes) >= 0 && utilfeature.DefaultFeatureGate.Enabled(features.BalanceAttachedNodeVolumes) && nodeInfo.TransientInfo != nil {
	// 		klog.Infof(
	// 			"%v -> %v: %v, capacity %d millicores %d memory bytes, %d volumes, total request %d millicores %d memory bytes %d volumes, score %d",
	// 			pod.Name, node.Name, r.Name,
	// 			allocatable.MilliCPU, allocatable.Memory, nodeInfo.TransientInfo.TransNodeInfo.AllocatableVolumesCount,
	// 			requested.MilliCPU, requested.Memory,
	// 			nodeInfo.TransientInfo.TransNodeInfo.RequestedVolumes,
	// 			score,
	// 		)
	// 	} else {
	// 		klog.Infof(
	// 			"%v -> %v: %v, capacity %d millicores %d memory bytes, total request %d millicores %d memory bytes, score %d",
	// 			pod.Name, node.Name, r.Name,
	// 			allocatable.MilliCPU, allocatable.Memory,
	// 			requested.MilliCPU, requested.Memory,
	// 			score,
	// 		)
	// 	}
	// }

	// Added the socket as HostPriority field
	return schedulerapi.HostPriority{
		Host:       node.Name,
		Socket:     socket,
		Score:      float64(score),
		Stress:     stress,
		NumCores:   socketCoresNum,
		Time:       time,
		FinalScore: finalScore,
	}, nil

	// return schedulerapi.HostPriority{
	// 	Host:  node.Name,
	// 	Score: int(score),
	// }, nil
}

func getNonZeroRequests(pod *v1.Pod) *schedulernodeinfo.Resource {
	result := &schedulernodeinfo.Resource{}
	for i := range pod.Spec.Containers {
		container := &pod.Spec.Containers[i]
		cpu, memory := priorityutil.GetNonzeroRequests(&container.Resources.Requests)
		result.MilliCPU += cpu
		result.Memory += memory
	}
	return result
}
