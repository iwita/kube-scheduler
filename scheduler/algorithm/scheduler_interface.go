/*
Copyright 2014 The Kubernetes Authors.

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

package algorithm

import (
	v1 "k8s.io/api/core/v1"
	schedulerapi "k8s.io/kubernetes/pkg/scheduler/api"
	schedulernodeinfo "k8s.io/kubernetes/pkg/scheduler/nodeinfo"
)

// SchedulerExtender is an interface for external processes to influence scheduling
// decisions made by Kubernetes. This is typically needed for resources not directly
// managed by Kubernetes.
type SchedulerExtender interface {
	// Name returns a unique name that identifies the extender.
	Name() string

	// Filter based on extender-implemented predicate functions. The filtered list is
	// expected to be a subset of the supplied list. failedNodesMap optionally contains
	// the list of failed nodes and failure reasons.
	Filter(pod *v1.Pod,
		nodes []*v1.Node, nodeNameToInfo map[string]*schedulernodeinfo.NodeInfo,
	) (filteredNodes []*v1.Node, failedNodesMap schedulerapi.FailedNodesMap, err error)

	// Prioritize based on extender-implemented priority functions. The returned scores & weight
	// are used to compute the weighted score for an extender. The weighted scores are added to
	// the scores computed  by Kubernetes scheduler. The total scores are used to do the host selection.
	Prioritize(pod *v1.Pod, nodes []*v1.Node) (hostPriorities *schedulerapi.HostPriorityList, weight int, err error)
	// Bind delegates the action of binding a pod to a node to the extender.
	Bind(binding *v1.Binding) error

	// IsBinder returns whether this extender is configured for the Bind method.
	IsBinder() bool

	// IsInterested returns true if at least one extended resource requested by
	// this pod is managed by this extender.
	IsInterested(pod *v1.Pod) bool

	// ProcessPreemption returns nodes with their victim pods processed by extender based on
	// given:
	//   1. Pod to schedule
	//   2. Candidate nodes and victim pods (nodeToVictims) generated by previous scheduling process.
	//   3. nodeNameToInfo to restore v1.Node from node name if extender cache is enabled.
	// The possible changes made by extender may include:
	//   1. Subset of given candidate nodes after preemption phase of extender.
	//   2. A different set of victim pod for every given candidate node after preemption phase of extender.
	ProcessPreemption(
		pod *v1.Pod,
		nodeToVictims map[*v1.Node]*schedulerapi.Victims,
		nodeNameToInfo map[string]*schedulernodeinfo.NodeInfo,
	) (map[*v1.Node]*schedulerapi.Victims, error)

	// SupportsPreemption returns if the scheduler extender support preemption or not.
	SupportsPreemption() bool

	// IsIgnorable returns true indicates scheduling should not fail when this extender
	// is unavailable. This gives scheduler ability to fail fast and tolerate non-critical extenders as well.
	IsIgnorable() bool
}
