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

package v1

import (
	"context"
	"time"

	v1 "github.com/F5Networks/k8s-bigip-ctlr/v2/config/apis/cis/v1"
	scheme "github.com/F5Networks/k8s-bigip-ctlr/v2/config/client/clientset/versioned/scheme"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// TLSProfilesGetter has a method to return a TLSProfileInterface.
// A group's client should implement this interface.
type TLSProfilesGetter interface {
	TLSProfiles(namespace string) TLSProfileInterface
}

// TLSProfileInterface has methods to work with TLSProfile resources.
type TLSProfileInterface interface {
	Create(ctx context.Context, tLSProfile *v1.TLSProfile, opts metav1.CreateOptions) (*v1.TLSProfile, error)
	Update(ctx context.Context, tLSProfile *v1.TLSProfile, opts metav1.UpdateOptions) (*v1.TLSProfile, error)
	Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1.TLSProfile, error)
	List(ctx context.Context, opts metav1.ListOptions) (*v1.TLSProfileList, error)
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.TLSProfile, err error)
	TLSProfileExpansion
}

// tLSProfiles implements TLSProfileInterface
type tLSProfiles struct {
	client rest.Interface
	ns     string
}

// newTLSProfiles returns a TLSProfiles
func newTLSProfiles(c *CisV1Client, namespace string) *tLSProfiles {
	return &tLSProfiles{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the tLSProfile, and returns the corresponding tLSProfile object, and an error if there is any.
func (c *tLSProfiles) Get(ctx context.Context, name string, options metav1.GetOptions) (result *v1.TLSProfile, err error) {
	result = &v1.TLSProfile{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("tlsprofiles").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of TLSProfiles that match those selectors.
func (c *tLSProfiles) List(ctx context.Context, opts metav1.ListOptions) (result *v1.TLSProfileList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1.TLSProfileList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("tlsprofiles").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested tLSProfiles.
func (c *tLSProfiles) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("tlsprofiles").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a tLSProfile and creates it.  Returns the server's representation of the tLSProfile, and an error, if there is any.
func (c *tLSProfiles) Create(ctx context.Context, tLSProfile *v1.TLSProfile, opts metav1.CreateOptions) (result *v1.TLSProfile, err error) {
	result = &v1.TLSProfile{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("tlsprofiles").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(tLSProfile).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a tLSProfile and updates it. Returns the server's representation of the tLSProfile, and an error, if there is any.
func (c *tLSProfiles) Update(ctx context.Context, tLSProfile *v1.TLSProfile, opts metav1.UpdateOptions) (result *v1.TLSProfile, err error) {
	result = &v1.TLSProfile{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("tlsprofiles").
		Name(tLSProfile.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(tLSProfile).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the tLSProfile and deletes it. Returns an error if one occurs.
func (c *tLSProfiles) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("tlsprofiles").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *tLSProfiles) DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("tlsprofiles").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched tLSProfile.
func (c *tLSProfiles) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.TLSProfile, err error) {
	result = &v1.TLSProfile{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("tlsprofiles").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
