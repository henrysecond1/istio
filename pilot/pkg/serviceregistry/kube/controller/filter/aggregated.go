package filter

type AggregatedFilter interface {
	Filter(obj interface{}) bool
}

type aggregatedFilter struct {
	filters []func(obj interface{}) bool
}

func NewAggregatedFilter(
	filterFuncs ...func(obj interface{}) bool,
) AggregatedFilter {
	filters := []func(obj interface{}) bool{}
	filters = append(filters, filterFuncs...)
	return &aggregatedFilter{filters: filters}
}

// All filters are ANDed
func (a *aggregatedFilter) Filter(obj interface{}) bool {
	for _, filter := range a.filters {
		if !filter(obj) {
			return false
		}
	}
	return true
}
