// Copyright Contributors to the Open Cluster Management project
package apply

import (
	"text/template"

	"github.com/openshift/library-go/pkg/operator/resource/resourceapply"
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/apimachinery/pkg/runtime"

	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
)

type Applier struct {
	kubeClient          kubernetes.Interface
	apiExtensionsClient apiextensionsclient.Interface
	dynamicClient       dynamic.Interface
	templateFuncMap     template.FuncMap
	scheme              *runtime.Scheme
	owner               runtime.Object
	cache               resourceapply.ResourceCache
	controller          *bool
	blockOwnerDeletion  *bool
}

//ApplierBuilder a builder to build the applier
type ApplierBuilder struct {
	applier Applier
}

type iApplierBuilder interface {
	//Build returns the builded applier
	Build() Applier
	//WithClient adds the several clients to the applier
	WithClient(
		kubeClient kubernetes.Interface,
		apiExtensionsClient apiextensionsclient.Interface,
		dynamicClient dynamic.Interface) *ApplierBuilder
	//WithTemplateFuncMap add template.FuncMap to the applier.
	WithTemplateFuncMap(fm template.FuncMap) *ApplierBuilder
	//WithOwner add an ownerref to the object
	WithOwner(owner runtime.Object, blockOwnerDeletion, controller bool, scheme *runtime.Scheme) *ApplierBuilder
	//WithCache add cache
	WithCache(cache resourceapply.ResourceCache) *ApplierBuilder
}

var _ iApplierBuilder = NewApplierBuilder()

//New ApplyBuilder
func NewApplierBuilder() *ApplierBuilder {
	return &ApplierBuilder{}
}

//Build returns the builded applier
func (a *ApplierBuilder) Build() Applier {
	if a.applier.cache == nil {
		a.applier.cache = NewResourceCache()
	}
	return a.applier
}

//WithClient adds the several clients to the applier
func (a *ApplierBuilder) WithClient(
	kubeClient kubernetes.Interface,
	apiExtensionsClient apiextensionsclient.Interface,
	dynamicClient dynamic.Interface) *ApplierBuilder {
	a.applier.kubeClient = kubeClient
	a.applier.apiExtensionsClient = apiExtensionsClient
	a.applier.dynamicClient = dynamicClient
	return a
}

//WithTemplateFuncMap add template.FuncMap to the applier.
func (a *ApplierBuilder) WithTemplateFuncMap(fm template.FuncMap) *ApplierBuilder {
	a.applier.templateFuncMap = fm
	return a
}

//WithOwner add an ownerref to the object
func (a *ApplierBuilder) WithOwner(owner runtime.Object, blockOwnerDeletion, controller bool, scheme *runtime.Scheme) *ApplierBuilder {
	a.applier.owner = owner
	a.applier.blockOwnerDeletion = &blockOwnerDeletion
	a.applier.controller = &controller
	a.applier.scheme = scheme
	return a
}

//WithCache set a the cache instead of using the default cache created on the Build()
func (a *ApplierBuilder) WithCache(cache resourceapply.ResourceCache) *ApplierBuilder {
	a.applier.cache = cache
	return a
}

func NewResourceCache() resourceapply.ResourceCache {
	return resourceapply.NewResourceCache()
}
