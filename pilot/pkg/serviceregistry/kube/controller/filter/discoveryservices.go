package filter

import (
	"sync"

	"istio.io/api/annotation"
	"istio.io/istio/pkg/config/visibility"
	corev1 "k8s.io/api/core/v1"
)

type DiscoveryServicesFilter interface {
	// return true if the input object resides in a namespace selected for discovery
	Filter(obj interface{}) bool
}

type discoveryServicesFilter struct {
	lock sync.RWMutex
}

func NewDiscoveryServicesFilter() DiscoveryServicesFilter {
	discoveryServicesFilter := &discoveryServicesFilter{}
	return discoveryServicesFilter
}

func (d *discoveryServicesFilter) Filter(obj interface{}) bool {
	d.lock.RLock()
	defer d.lock.RUnlock()

	switch svc := obj.(type) {
	case *corev1.Service:
		// permit if service is exposed
		return isExposed(svc)
	default:
		return true
	}
}

func isExposed(svc *corev1.Service) bool {
	exportTo, ok := svc.Annotations[annotation.NetworkingExportTo.Name]
	if !ok {
		return true
	}
	// For now, we only care about none visibility
	if visibility.Instance(exportTo) == visibility.None {
		return false
	}
	return true
}
