/*
Copyright The Kubernetes Authors.

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

// Code generated by lister-gen. DO NOT EDIT.

package v1

import (
	v1 "github.com/F5Networks/k8s-bigip-ctlr/v2/config/apis/cis/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// ExternalDNSLister helps list ExternalDNSes.
// All objects returned here must be treated as read-only.
type ExternalDNSLister interface {
	// List lists all ExternalDNSes in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1.ExternalDNS, err error)
	// ExternalDNSes returns an object that can list and get ExternalDNSes.
	ExternalDNSes(namespace string) ExternalDNSNamespaceLister
	ExternalDNSListerExpansion
}

// externalDNSLister implements the ExternalDNSLister interface.
type externalDNSLister struct {
	indexer cache.Indexer
}

// NewExternalDNSLister returns a new ExternalDNSLister.
func NewExternalDNSLister(indexer cache.Indexer) ExternalDNSLister {
	return &externalDNSLister{indexer: indexer}
}

// List lists all ExternalDNSes in the indexer.
func (s *externalDNSLister) List(selector labels.Selector) (ret []*v1.ExternalDNS, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.ExternalDNS))
	})
	return ret, err
}

// ExternalDNSes returns an object that can list and get ExternalDNSes.
func (s *externalDNSLister) ExternalDNSes(namespace string) ExternalDNSNamespaceLister {
	return externalDNSNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// ExternalDNSNamespaceLister helps list and get ExternalDNSes.
// All objects returned here must be treated as read-only.
type ExternalDNSNamespaceLister interface {
	// List lists all ExternalDNSes in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1.ExternalDNS, err error)
	// Get retrieves the ExternalDNS from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1.ExternalDNS, error)
	ExternalDNSNamespaceListerExpansion
}

// externalDNSNamespaceLister implements the ExternalDNSNamespaceLister
// interface.
type externalDNSNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all ExternalDNSes in the indexer for a given namespace.
func (s externalDNSNamespaceLister) List(selector labels.Selector) (ret []*v1.ExternalDNS, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.ExternalDNS))
	})
	return ret, err
}

// Get retrieves the ExternalDNS from the indexer for a given namespace and name.
func (s externalDNSNamespaceLister) Get(name string) (*v1.ExternalDNS, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1.Resource("externaldns"), name)
	}
	return obj.(*v1.ExternalDNS), nil
}
