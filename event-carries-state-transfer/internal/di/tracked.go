package di

import "strings"

type tracked map[string]int

func (t tracked) add(info depInfo) tracked {
	newList := make(tracked, len(t))

	for k, v := range t {
		newList[k] = v
	}

	newList[info.key] = len(newList)

	return newList
}

func (s tracked) ordered() []string {
	keys := make([]string, len(s))
	for k, v := range s {
		keys[v] = k
	}

	return keys
}

func (s tracked) String() string {
	return strings.Join(s.ordered(), ",")
}
