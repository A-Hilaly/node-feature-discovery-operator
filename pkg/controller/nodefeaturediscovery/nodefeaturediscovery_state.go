/*
Copyright 2020 The Kubernetes Authors.

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

package nodefeaturediscovery

import (
	"errors"

	nfdv1alpha1 "github.com/kubernetes-sigs/node-feature-discovery-operator/pkg/apis/nfd/v1alpha1"
)

type NFD struct {
	resources []Resources
	controls  []controlFunc
	rec       *ReconcileNodeFeatureDiscovery
	ins       *nfdv1alpha1.NodeFeatureDiscovery
	idx       int
}

func addState(n *NFD, path string) error {

	res, ctrl := addResourcesControls(path)

	n.controls = append(n.controls, ctrl)
	n.resources = append(n.resources, res)

	return nil
}

func (n *NFD) init(r *ReconcileNodeFeatureDiscovery,
	i *nfdv1alpha1.NodeFeatureDiscovery) error {
	n.rec = r
	n.ins = i
	n.idx = 0

	addState(n, "/opt/nfd/master")
	addState(n, "/opt/nfd/worker")

	return nil
}

func (n *NFD) step() error {

	for _, fs := range n.controls[n.idx] {

		stat, err := fs(*n)
		if err != nil {
			return err
		}
		if stat != Ready {
			return errors.New("ResourceNotReady")
		}
	}

	n.idx = n.idx + 1

	return nil
}

func (n NFD) validate() {
	// TODO add custom validation functions
}

func (n NFD) last() bool {
	if n.idx == len(n.controls) {
		return true
	}
	return false
}
