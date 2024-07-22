package versions

import (
	"fmt"
	"sort"
)

// A node in an N-ary tree where children are sorted by the value
type V struct {
	Value uint16
	M     map[uint16]int
	Sub   []V
}

func (v *V) leaf() bool {
	return len(v.Sub) == 0
}

func (v *V) Insert(values []uint16) {
	if len(values) == 0 {
		return
	}
	value := values[0]

	i := sort.Search(len(v.Sub), func(i int) bool {
		return v.Sub[i].Value >= value
	})
	if i >= len(v.Sub) || v.Sub[i].Value != value {
		n := V{
			Value: value,
			M:     make(map[uint16]int),
		}
		// Insert at i
		v.Sub = append(v.Sub[:i], append(([]V{n}), v.Sub[i:]...)...)
		for j, sub := range v.Sub[i:] {
			v.M[sub.Value] = i + j
		}

	}

	v.Sub[i].Insert(values[1:])
}

func (v *V) Remove(values []uint16) {
	if len(values) == 0 {
		return
	}
	value := values[0]

	i, ex := v.M[value]
	if !ex {
		return
	}

	v.Sub[i].Remove(values[1:])

	if v.Sub[i].leaf() {
		v.Sub = append(v.Sub[:i], v.Sub[i+1:]...)
		delete(v.M, value)
		for j, sub := range v.Sub[i:] {
			v.M[sub.Value] = i + j
		}
	}
}

func (v *V) Export() (values []string) {
	values = []string{}
	for _, sub := range v.Sub {
		export := sub.Export()
		if len(export) == 0 {
			values = append(values, fmt.Sprintf("%d", sub.Value))
			continue
		}
		for _, subExp := range export {
			values = append(values, fmt.Sprintf("%d.%s", sub.Value, subExp))
		}
	}

	return
}

type Versions struct {
	Tree *V
}

func New() *Versions {
	return &Versions{
		Tree: &V{
			Value: 0, // not used (root)
			M:     make(map[uint16]int),
		},
	}
}

// Stores a version
func (sv *Versions) Add(version string) error {
	major, minor, patch, suffix, err := Parse(version)
	if err != nil {
		return err
	}
	if suffix != "" {
		// Not implemented
		return nil
	}

	sv.Tree.Insert([]uint16{major, minor, patch})

	return nil
}

// Removes a version
func (sv *Versions) Remove(version string) error {
	major, minor, patch, suffix, err := Parse(version)
	if err != nil {
		return err
	}
	if suffix != "" {
		// Not implemented
		return nil
	}

	sv.Tree.Remove([]uint16{major, minor, patch})

	return nil
}

// Export all versions. Only for testing
func (sv *Versions) Export() []string {
	return sv.Tree.Export()
}

func (sv *Versions) Empty() bool {
	return len(sv.Tree.Sub) == 0
}

// For the current versions it returns the latest version that would not trigger a major upgrade.
func (sv *Versions) LatestCompatible(currentVersion string) (latestVersion string, err error) {
	latestVersion = currentVersion

	major, minor, patch, _, err := Parse(currentVersion)
	if err != nil {
		return
	}

	var latestMinor uint16
	var latestPatch uint16
	if major == 0 {
		if len(sv.Tree.Sub) == 0 || sv.Tree.Sub[0].Value != 0 {
			return
		}
		minorIndex, ex := sv.Tree.Sub[0].M[minor]
		if !ex {
			return
		}
		minorEntry := sv.Tree.Sub[0].Sub[minorIndex]

		latestMinor = minor
		if len(minorEntry.Sub) == 0 {
			panic("no patch entries under a minor entry")
		}
		latestPatch = minorEntry.Sub[len(minorEntry.Sub)-1].Value
	} else {
		majorIndex, ex := sv.Tree.M[major]
		if !ex {
			return
		}
		majorEntry := sv.Tree.Sub[majorIndex]

		if len(majorEntry.Sub) == 0 {
			panic("no minor entries under a major entry")
		}
		minorEntry := majorEntry.Sub[len(majorEntry.Sub)-1]

		latestMinor = minorEntry.Value
		if len(minorEntry.Sub) == 0 {
			panic("no patch entries under a minor entry")
		}
		latestPatch = minorEntry.Sub[len(minorEntry.Sub)-1].Value
	}

	if latestMinor < minor || latestPatch < patch {
		return
	}

	latestVersion = fmt.Sprintf("%d.%d.%d", major, latestMinor, latestPatch)
	return
}
