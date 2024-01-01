package transactions

type FilterItem struct {
	Folder  string
	Effects []string
}

type FilterItems []FilterItem
type Filter map[string]map[string]bool

func NewFilter(filters FilterItems) Filter {
	filterMap := make(map[string]map[string]bool)
	for _, filter := range filters {
		effectMap := make(map[string]bool)
		for _, effect := range filter.Effects {
			effectMap[effect] = false
		}
		filterMap[filter.Folder] = effectMap
	}
	return filterMap
}

func (filter Filter) IsSelected(selections ...string) bool {
	if len(filter) == 0 {
		return true
	}

	selectionLength := len(selections)
	if selectionLength == 0 {
		return true
	}

	effectMap, ok := filter[selections[0]]
	if !ok {
		return false
	}

	if selectionLength < 2 {
		return true
	}

	if len(effectMap) == 0 {
		return true
	}

	_, ok = effectMap[selections[1]]
	return ok
}
