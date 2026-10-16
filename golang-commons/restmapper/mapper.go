/*
Copyright The Platform Mesh Authors.

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

// Package restmapper provides a static meta.RESTMapper for resource
// kinds not registered in a cluster's live discovery.
package restmapper

import (
	"strings"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type GVRMapping struct {
	GVR   schema.GroupVersionResource
	Scope meta.RESTScopeName
}

// NewStatic builds a DefaultRESTMapper from the given mappings
func NewStatic(mappings map[schema.GroupKind]GVRMapping) meta.RESTMapper {
	if len(mappings) == 0 {
		return nil
	}

	groupVersions := make([]schema.GroupVersion, 0, len(mappings))
	for _, m := range mappings {
		groupVersions = append(groupVersions, m.GVR.GroupVersion())
	}
	mapper := meta.NewDefaultRESTMapper(groupVersions)
	for gk, m := range mappings {
		gvk := m.GVR.GroupVersion().WithKind(gk.Kind)
		plural := m.GVR.GroupVersion().WithResource(m.GVR.Resource)
		singular := m.GVR.GroupVersion().WithResource(strings.ToLower(gk.Kind))
		scope := meta.RESTScopeNamespace
		if m.Scope == meta.RESTScopeNameRoot {
			scope = meta.RESTScopeRoot
		}
		mapper.AddSpecific(gvk, plural, singular, scope)
	}
	return mapper
}
