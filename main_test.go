package main

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNoteIndexes(t *testing.T) {
	tests := []struct {
		note  string
		index int
	}{
		{"B♯", 0}, {"C", 0},
		{"C♯", 1}, {"D♭", 1},
		{"D", 2},
		{"D♯", 3}, {"E♭", 3},
		{"E", 4}, {"F♭", 4},
		{"E♯", 5}, {"F", 5},
		{"F♯", 6}, {"G♭", 6},
		{"G", 7},
		{"G♯", 8}, {"A♭", 8},
		{"A", 9},
		{"A♯", 10}, {"B♭", 10},
		{"B", 11}, {"C♭", 11},
	}

	for _, e := range tests {
		i, ok := noteIndexes[e.note]
		if !ok {
			t.Errorf("Missing note '%s' in noteIndexes", e.note)
		}
		if i != e.index {
			t.Errorf("Incorrect note index %d for note '%s' in noteIndexes, expected %d", i, e.note, e.index)
		}
	}
}

func TestEnumerateScaleMajor(t *testing.T) {
	tests := []struct {
		note  string
		scale []string
	}{
		{"C", []string{"C", "D", "E", "F", "G", "A", "B"}},
		{"C♯", []string{"C♯", "D♯", "E♯", "F♯", "G♯", "A♯", "B♯"}},
		{"D♭", []string{"D♭", "E♭", "F", "G♭", "A♭", "B♭", "C"}},
		{"D", []string{"D", "E", "F♯", "G", "A", "B", "C♯"}},
		{"D♯", []string{"D♯", "E♯", "G", "A♭", "B♭", "C", "D"}},
		{"E♭", []string{"E♭", "F", "G", "A♭", "B♭", "C", "D"}},
		{"E", []string{"E", "F♯", "G♯", "A", "B", "C♯", "D♯"}},
		{"F", []string{"F", "G", "A", "B♭", "C", "D", "E"}},
		{"F♯", []string{"F♯", "G♯", "A♯", "B", "C♯", "D♯", "E♯"}},
		{"G♭", []string{"G♭", "A♭", "B♭", "C♭", "D♭", "E♭", "F"}},
		{"G", []string{"G", "A", "B", "C", "D", "E", "F♯"}},
		{"G♯", []string{"G♯", "A♯", "B♯", "C♯", "D♯", "E♯", "G"}},
		{"A♭", []string{"A♭", "B♭", "C", "D♭", "E♭", "F", "G"}},
		{"A", []string{"A", "B", "C♯", "D", "E", "F♯", "G♯"}},
		{"A♯", []string{"A♯", "B♯", "D", "E♭", "F", "G", "A"}},
		{"B♭", []string{"B♭", "C", "D", "E♭", "F", "G", "A"}},
		{"B", []string{"B", "C♯", "D♯", "E", "F♯", "G♯", "A♯"}},
	}
	assert.Equal(t, len(chromatic), len(tests))

	intervals := scaleIntervals["Major"]
	for _, e := range tests {
		s := enumerateScale(e.note, intervals)
		if !reflect.DeepEqual(e.scale, s) {
			t.Errorf("Scales not equal: for scale %s expected %v, got %v", e.note, e.scale, s)
		}
	}
}

func TestEnumerateScaleMinor(t *testing.T) {
	tests := []struct {
		note  string
		scale []string
	}{
		{"C", []string{"C", "D", "E♭", "F", "G", "A♭", "B♭"}},
		{"C♯", []string{"C♯", "D♯", "E", "F♯", "G♯", "A", "B"}},
		{"D♭", []string{"D♭", "E♭", "F♭", "G♭", "A♭", "A", "B"}},
		{"D", []string{"D", "E", "F", "G", "A", "B♭", "C"}},
		{"D♯", []string{"D♯", "E♯", "F♯", "G♯", "A♯", "B", "C♯"}},
		{"E♭", []string{"E♭", "F", "G♭", "A♭", "B♭", "C♭", "D♭"}},
		{"E", []string{"E", "F♯", "G", "A", "B", "C", "D"}},
		{"F", []string{"F", "G", "A♭", "B♭", "C", "D♭", "E♭"}},
		{"F♯", []string{"F♯", "G♯", "A", "B", "C♯", "D", "E"}},
		{"G♭", []string{"G♭", "A♭", "A", "B", "C♯", "D", "E"}},
		{"G", []string{"G", "A", "B♭", "C", "D", "E♭", "F"}},
		{"G♯", []string{"G♯", "A♯", "B", "C♯", "D♯", "E", "F♯"}},
		{"A♭", []string{"A♭", "B♭", "C♭", "D♭", "E♭", "F♭", "G♭"}},
		{"A", []string{"A", "B", "C", "D", "E", "F", "G"}},
		{"A♯", []string{"A♯", "B♯", "C♯", "D♯", "E♯", "F♯", "G♯"}},
		{"B♭", []string{"B♭", "C", "D♭", "E♭", "F", "G♭", "A♭"}},
		{"B", []string{"B", "C♯", "D", "E", "F♯", "G", "A"}},
	}
	assert.Equal(t, len(chromatic), len(tests))

	intervals := scaleIntervals["Minor"]
	for _, e := range tests {
		s := enumerateScale(e.note, intervals)
		if !reflect.DeepEqual(e.scale, s) {
			t.Errorf("Scales not equal: for scale %s expected %v, got %v", e.note, e.scale, s)
		}
	}
}
