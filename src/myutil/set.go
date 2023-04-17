package myutil

type Set map[interface{}]bool

func NewSet(items ...interface{}) *Set {
	set := make(Set)
	for _, item := range items {
		set.Add(item)
	}
	return &set
}

// Add adds an element to the set
func (s Set) Add(item interface{}) {
	s[item] = true
}

// Remove removes an element from the set
func (s Set) Remove(item interface{}) {
	delete(s, item)
}

// Contains checks if an element is present in the set
func (s Set) Contains(item interface{}) bool {
	_, exists := s[item]
	return exists
}

// Size returns the number of elements in the set
func (s Set) Size() int {
	return len(s)
}

// Clear removes all elements from the set
func (s Set) Clear() {
	for k := range s {
		delete(s, k)
	}
}

// Elements returns a slice of all elements in the set
func (s Set) Elements() []interface{} {
	elements := make([]interface{}, 0, len(s))
	for k := range s {
		elements = append(elements, k)
	}
	return elements
}
