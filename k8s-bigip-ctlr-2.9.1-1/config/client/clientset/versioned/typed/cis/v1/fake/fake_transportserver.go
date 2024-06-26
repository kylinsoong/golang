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

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"

	cisv1 "github.com/F5Networks/k8s-bigip-ctlr/config/apis/cis/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeTransportServers implements TransportServerInterface
type FakeTransportServers struct {
	Fake *FakeCisV1
	ns   string
}

var transportserversResource = schema.GroupVersionResource{Group: "cis.f5.com", Version: "v1", Resource: "transportservers"}

var transportserversKind = schema.GroupVersionKind{Group: "cis.f5.com", Version: "v1", Kind: "TransportServer"}

// Get takes name of the transportServer, and returns the corresponding transportServer object, and an error if there is any.
func (c *FakeTransportServers) Get(ctx context.Context, name string, options v1.GetOptions) (result *cisv1.TransportServer, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(transportserversResource, c.ns, name), &cisv1.TransportServer{})

	if obj == nil {
		return nil, err
	}
	return obj.(*cisv1.TransportServer), err
}

// List takes label and field selectors, and returns the list of TransportServers that match those selectors.
func (c *FakeTransportServers) List(ctx context.Context, opts v1.ListOptions) (result *cisv1.TransportServerList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(transportserversResource, transportserversKind, c.ns, opts), &cisv1.TransportServerList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &cisv1.TransportServerList{ListMeta: obj.(*cisv1.TransportServerList).ListMeta}
	for _, item := range obj.(*cisv1.TransportServerList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested transportServers.
func (c *FakeTransportServers) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(transportserversResource, c.ns, opts))

}

// Create takes the representation of a transportServer and creates it.  Returns the server's representation of the transportServer, and an error, if there is any.
func (c *FakeTransportServers) Create(ctx context.Context, transportServer *cisv1.TransportServer, opts v1.CreateOptions) (result *cisv1.TransportServer, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(transportserversResource, c.ns, transportServer), &cisv1.TransportServer{})

	if obj == nil {
		return nil, err
	}
	return obj.(*cisv1.TransportServer), err
}

// Update takes the representation of a transportServer and updates it. Returns the server's representation of the transportServer, and an error, if there is any.
func (c *FakeTransportServers) Update(ctx context.Context, transportServer *cisv1.TransportServer, opts v1.UpdateOptions) (result *cisv1.TransportServer, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(transportserversResource, c.ns, transportServer), &cisv1.TransportServer{})

	if obj == nil {
		return nil, err
	}
	return obj.(*cisv1.TransportServer), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeTransportServers) UpdateStatus(ctx context.Context, transportServer *cisv1.TransportServer, opts v1.UpdateOptions) (*cisv1.TransportServer, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(transportserversResource, "status", c.ns, transportServer), &cisv1.TransportServer{})

	if obj == nil {
		return nil, err
	}
	return obj.(*cisv1.TransportServer), err
}

// Delete takes name of the transportServer and deletes it. Returns an error if one occurs.
func (c *FakeTransportServers) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(transportserversResource, c.ns, name), &cisv1.TransportServer{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeTransportServers) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(transportserversResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &cisv1.TransportServerList{})
	return err
}

// Patch applies the patch and returns the patched transportServer.
func (c *FakeTransportServers) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *cisv1.TransportServer, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(transportserversResource, c.ns, name, pt, data, subresources...), &cisv1.TransportServer{})

	if obj == nil {
		return nil, err
	}
	return obj.(*cisv1.TransportServer), err
}
