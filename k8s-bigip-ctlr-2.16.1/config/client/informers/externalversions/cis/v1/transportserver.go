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

// Code generated by informer-gen. DO NOT EDIT.

package v1

import (
	"context"
	time "time"

	cisv1 "github.com/F5Networks/k8s-bigip-ctlr/v2/config/apis/cis/v1"
	versioned "github.com/F5Networks/k8s-bigip-ctlr/v2/config/client/clientset/versioned"
	internalinterfaces "github.com/F5Networks/k8s-bigip-ctlr/v2/config/client/informers/externalversions/internalinterfaces"
	v1 "github.com/F5Networks/k8s-bigip-ctlr/v2/config/client/listers/cis/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// TransportServerInformer provides access to a shared informer and lister for
// TransportServers.
type TransportServerInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1.TransportServerLister
}

type transportServerInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewTransportServerInformer constructs a new informer for TransportServer type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewTransportServerInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredTransportServerInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredTransportServerInformer constructs a new informer for TransportServer type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredTransportServerInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.CisV1().TransportServers(namespace).List(context.TODO(), options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.CisV1().TransportServers(namespace).Watch(context.TODO(), options)
			},
		},
		&cisv1.TransportServer{},
		resyncPeriod,
		indexers,
	)
}

func (f *transportServerInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredTransportServerInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *transportServerInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&cisv1.TransportServer{}, f.defaultInformer)
}

func (f *transportServerInformer) Lister() v1.TransportServerLister {
	return v1.NewTransportServerLister(f.Informer().GetIndexer())
}
