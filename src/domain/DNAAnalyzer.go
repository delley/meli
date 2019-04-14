package domain

import (
	"fmt"
	"regexp"
	"strings"
)

const words = "ATCG"
const rg = "%c{4}"

// IsSimian determines whether a DNA is from a simian.
func IsSimian(dna []string) (bool, error) {

	if v := validateStructure(dna); !v {
		return false, fmt.Errorf("Invalid structure: %s", dna)
	}

	if m := match(dna); m {
		return true, nil
	}

	rdna := rotateCCW(dna)
	if m := match(rdna); m {
		return true, nil
	}

	dd := diagonals(dna)
	if m := match(dd); m {
		return true, nil
	}

	dr := diagonals(rdna)
	if m := match(dr); m {
		return true, nil
	}

	return false, nil
}

func match(dna []string) bool {
	for _, chain := range dna {
		for _, w := range words {
			matched, _ := regexp.MatchString(fmt.Sprintf(rg, w), chain)
			if matched {
				return true
			}
		}
	}
	return false
}

func transpose(s []string) []string {
	var t = []string{"", "", "", "", "", ""}
	for _, w := range s {
		for i, v := range w {
			t[i] = fmt.Sprintf("%s%c", t[i], v)
		}
	}

	return t
}

func rotateCCW(s []string) []string {
	var t = []string{"", "", "", "", "", ""}
	for _, w := range s {
		j := len(w) - 1
		for i, v := range w {
			t[j-i] = fmt.Sprintf("%s%c", t[j-i], v)
		}
	}

	return t
}

func diagonals(s []string) []string {
	var t = []string{"", "", "", "", "", "", "", "", "", "", ""}
	k := 0
	j := -2
	z := 1
	for _, w := range s {

		for i := 0; i < len(w); i++ {
			t[k] = fmt.Sprintf("%s%c", t[k], w[i])

			if k == 0 {
				j = j * -1
			}

			k = k + j
			if k < 0 {
				k = 0
			}
		}
		k = z
		z = z + j
		j = -2
	}

	return t
}

func validateStructure(dna []string) bool {
	if len(dna) != 6 {
		return false
	}

	c := strings.Join(dna, "")
	m, _ := regexp.MatchString("^[ATCG]*$", c)
	if !m {
		return false
	}

	return true
}
