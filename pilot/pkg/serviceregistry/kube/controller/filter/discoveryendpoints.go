package filter

import (
	"sync"

	corev1 "k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	listerv1 "k8s.io/client-go/listers/core/v1"
)

type DiscoveryEndpointsFilter interface {
	// return true if the input object resides in a namespace selected for discovery
	Filter(obj interface{}) bool
}

type discoveryEndpointsFilter struct {
	lock          sync.RWMutex
	serviceLister listerv1.ServiceLister
}

func NewDiscoveryEndpointsFilter(
	serviceLister listerv1.ServiceLister,
) DiscoveryEndpointsFilter {
	discoveryEndpointsFilter := &discoveryEndpointsFilter{
		serviceLister: serviceLister,
	}

	return discoveryEndpointsFilter
}

func (d *discoveryEndpointsFilter) Filter(obj interface{}) bool {
	d.lock.RLock()
	defer d.lock.RUnlock()

	switch ep := obj.(type) {
	case *corev1.Endpoints:
		// TODO(jwpark): How to handle endpoint slices?
		svc, err := d.serviceLister.Services(ep.GetNamespace()).Get(ep.Name)
		if err != nil {
			// If service is not found from lister, we should filter it.
			if kerrors.IsNotFound(err) {
				return false
			}
			// Handler error
			return true
		}
		return isExposed(svc)
	default:
		return true
	}
}
