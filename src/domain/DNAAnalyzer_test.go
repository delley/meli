package domain

import (
	"testing"
)

func TestDNAAnalyzer(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want bool
	}{
		{"Horizontal Analysis", []string{"CTGAGA", "CTATGC", "TATTGT", "AGAGGG", "CCCCTA", "TCACTG"}, true},
		{"Vertical Analysis", []string{"CTGAGA", "CTATGC", "TATTGT", "AGAGGG", "CCCATA", "TCACTG"}, true},
		{"Left diagonal analysis,", []string{"CTGACA", "CTATGC", "TATTGT", "AGAGGG", "CCCATA", "TCACTG"}, true},
		{"Right diagonal analysis", []string{"CTGATA", "CTATGC", "CATTGT", "CCAGGG", "CCCATA", "TCACTG"}, true},
		{"It's not a simian", []string{"ATGCGA", "CAGTGC", "TTATTT", "AGACGG", "GCGTCA", "TCACTG"}, false},
		{"Invalid structure", []string{"ATGCGA", "CAGTGC", "TTATTT", "AGACGG", "GCGTCA", "TCACTK"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := IsSimian(tt.args); got != tt.want {
				t.Errorf("isSimian() = %v, want %v", got, tt.want)
			}
		})
	}

}

func TestTransposer(t *testing.T) {
	var DNA = []string{"CTGAGA", "CTATGC", "TATTGT", "AGAGGG", "CCCCTA", "TCACTG"}
	var want = []string{"CCTACT", "TTAGCC", "GATACA", "ATTGCC", "GGGGTT", "ACTGAG"}
	got := transpose(DNA)
	if notEqual(got, want) {
		t.Errorf("transpose() = %v, want %v", got, want)
	}
}

func TestRotateCCW(t *testing.T) {
	var DNA = []string{"CTGAGA", "CTATGC", "TATTGT", "AGAGGG", "CCCCTA", "TCACTG"}
	var want = []string{"ACTGAG", "GGGGTT", "ATTGCC", "GATACA", "TTAGCC", "CCTACT"}
	got := rotateCCW(DNA)
	if notEqual(got, want) {
		t.Errorf("rotateCCW() = %v, want %v", got, want)
	}
}

func TestDiagonals(t *testing.T) {
	var DNA = []string{"CTGAGA", "CTATGC", "TATTGT", "AGAGGG", "CCCCTA", "TCACTG"}
	var want = []string{"CTTGTG", "CAACT", "TATGA", "TGCC", "GTGG", "ACA", "AGT", "CC", "GC", "T", "A"}
	got := diagonals(DNA)
	if notEqual(got, want) {
		t.Errorf("diagonals() = %v, want %v", got, want)
	}
}

func TestValidateStructure(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want bool
	}{
		{"A valid DNA", []string{"CTGAGA", "CTATGC", "TATTGT", "AGAGGG", "CCCCTA", "TCACTG"}, true},
		{"Not a valid DNA - invalid size", []string{"CTGAGA", "CTATGC", "TATTGT", "AGAGGG", "CCCATA", "TCACTG", "TCACTG"}, false},
		{"Not a valid DNA - invalid char", []string{"CTGACA", "CTATGC", "TATTGT", "AGAGGG", "CCCATA", "TCACTK"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := validateStructure(tt.args); got != tt.want {
				t.Errorf("isValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func notEqual(a, b []string) bool {
	return !equal(a, b)
}

func equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
