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

	cisv1 "github.com/F5Networks/k8s-bigip-ctlr/v2/config/apis/cis/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeVirtualServers implements VirtualServerInterface
type FakeVirtualServers struct {
	Fake *FakeCisV1
	ns   string
}

var virtualserversResource = schema.GroupVersionResource{Group: "cis.f5.com", Version: "v1", Resource: "virtualservers"}

var virtualserversKind = schema.GroupVersionKind{Group: "cis.f5.com", Version: "v1", Kind: "VirtualServer"}

// Get takes name of the virtualServer, and returns the corresponding virtualServer object, and an error if there is any.
func (c *FakeVirtualServers) Get(ctx context.Context, name string, options v1.GetOptions) (result *cisv1.VirtualServer, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(virtualserversResource, c.ns, name), &cisv1.VirtualServer{})

	if obj == nil {
		return nil, err
	}
	return obj.(*cisv1.VirtualServer), err
}

// List takes label and field selectors, and returns the list of VirtualServers that match those selectors.
func (c *FakeVirtualServers) List(ctx context.Context, opts v1.ListOptions) (result *cisv1.VirtualServerList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(virtualserversResource, virtualserversKind, c.ns, opts), &cisv1.VirtualServerList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &cisv1.VirtualServerList{ListMeta: obj.(*cisv1.VirtualServerList).ListMeta}
	for _, item := range obj.(*cisv1.VirtualServerList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested virtualServers.
func (c *FakeVirtualServers) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(virtualserversResource, c.ns, opts))

}

// Create takes the representation of a virtualServer and creates it.  Returns the server's representation of the virtualServer, and an error, if there is any.
func (c *FakeVirtualServers) Create(ctx context.Context, virtualServer *cisv1.VirtualServer, opts v1.CreateOptions) (result *cisv1.VirtualServer, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(virtualserversResource, c.ns, virtualServer), &cisv1.VirtualServer{})

	if obj == nil {
		return nil, err
	}
	return obj.(*cisv1.VirtualServer), err
}

// Update takes the representation of a virtualServer and updates it. Returns the server's representation of the virtualServer, and an error, if there is any.
func (c *FakeVirtualServers) Update(ctx context.Context, virtualServer *cisv1.VirtualServer, opts v1.UpdateOptions) (result *cisv1.VirtualServer, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(virtualserversResource, c.ns, virtualServer), &cisv1.VirtualServer{})

	if obj == nil {
		return nil, err
	}
	return obj.(*cisv1.VirtualServer), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeVirtualServers) UpdateStatus(ctx context.Context, virtualServer *cisv1.VirtualServer, opts v1.UpdateOptions) (*cisv1.VirtualServer, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(virtualserversResource, "status", c.ns, virtualServer), &cisv1.VirtualServer{})

	if obj == nil {
		return nil, err
	}
	return obj.(*cisv1.VirtualServer), err
}

// Delete takes name of the virtualServer and deletes it. Returns an error if one occurs.
func (c *FakeVirtualServers) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(virtualserversResource, c.ns, name), &cisv1.VirtualServer{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeVirtualServers) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(virtualserversResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &cisv1.VirtualServerList{})
	return err
}

// Patch applies the patch and returns the patched virtualServer.
func (c *FakeVirtualServers) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *cisv1.VirtualServer, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(virtualserversResource, c.ns, name, pt, data, subresources...), &cisv1.VirtualServer{})

	if obj == nil {
		return nil, err
	}
	return obj.(*cisv1.VirtualServer), err
}
