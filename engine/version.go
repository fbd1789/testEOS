package engine

import (
	"strconv"
	"strings"
)

// CompareVersions compare deux versions EOS ("4.28.0F" vs "4.34.0F").
// Retourne -1 si v1 < v2, 0 si égal, 1 si v1 > v2.
func CompareVersions(v1, v2 string) int {
	// Enlève le "F" à la fin et sépare les numéros
	strip := func(v string) []int {
		v = strings.Split(v, "F")[0]
		parts := strings.Split(v, ".")
		var out []int
		for _, p := range parts {
			n, err := strconv.Atoi(p)
			if err != nil {
				out = append(out, 0)
			} else {
				out = append(out, n)
			}
		}
		return out
	}

	a := strip(v1)
	b := strip(v2)

	// Compare chaque segment
	for i := 0; i < len(a) || i < len(b); i++ {
		var x, y int
		if i < len(a) {
			x = a[i]
		}
		if i < len(b) {
			y = b[i]
		}
		if x < y {
			return -1
		}
		if x > y {
			return 1
		}
	}
	return 0
}
