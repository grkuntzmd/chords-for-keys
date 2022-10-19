package main

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type (
	note struct {
		name           string
		accidental     string
		includeInScale bool
	}

	scale struct {
		name      string
		intervals []int
	}

	noteIndex struct {
		note  note
		index int
	}

	state struct {
		key        noteIndex
		scale      scale
		scaleNotes []string
	}
)

const (
	sharp = "\u266f"
	flat  = "\u266d"
)

var (
	notes [][]note = [][]note{
		{{"B", sharp, false}, {"C", "", true}},
		{{"C", sharp, true}, {"D", flat, true}},
		{{"D", "", true}},
		{{"D", sharp, true}, {"E", flat, true}},
		{{"E", "", true}, {"F", flat, false}},
		{{"E", sharp, false}, {"F", "", true}},
		{{"F", sharp, true}, {"G", flat, true}},
		{{"G", "", true}},
		{{"G", sharp, true}, {"A", flat, true}},
		{{"A", "", true}},
		{{"A", sharp, true}, {"B", flat, true}},
		{{"B", "", true}, {"C", flat, false}},
	}

	scales = []scale{
		{"Major", []int{2, 2, 1, 2, 2, 2}},
		{"Natural Minor", []int{2, 1, 2, 2, 1, 2}},
	}
)

func init() {
}

func main() {
	var keyNames []string
	keyNotes := make(map[string]noteIndex)
	for i, g := range notes {
		for _, n := range g {
			if n.includeInScale {
				noteName := n.name + n.accidental
				keyNames = append(keyNames, noteName)
				keyNotes[noteName] = noteIndex{n, i}
			}
		}
	}

	var scaleNames []string
	scaleIntervals := make(map[string]scale)
	for _, s := range scales {
		scaleNames = append(scaleNames, s.name)
		scaleIntervals[s.name] = s
	}

	a := app.New()
	a.Settings().SetTheme(&myTheme{})

	state := state{
		key:   noteIndex{notes[0][0], 0},
		scale: scales[0],
	}
	state.scaleNotes = enumerateScale(state.key.note.name, state.key.note.accidental, state.key.index, state.scale.intervals)

	scaleLabel := widget.NewLabel(makeScaleLabel(state.scaleNotes))

	keySelector := widget.NewSelect(keyNames, func(s string) {
		state.key = keyNotes[s]
		state.scaleNotes = enumerateScale(state.key.note.name, state.key.note.accidental, state.key.index, state.scale.intervals)
		scaleLabel.SetText(makeScaleLabel(state.scaleNotes))
	})
	keySelector.SetSelectedIndex(0)

	scaleSelector := widget.NewSelect(scaleNames, func(s string) {
		state.scale = scaleIntervals[s]
		state.scaleNotes = enumerateScale(state.key.note.name, state.key.note.accidental, state.key.index, state.scale.intervals)
		scaleLabel.SetText(makeScaleLabel(state.scaleNotes))
	})
	scaleSelector.SetSelectedIndex(0)

	w := a.NewWindow("Chords for Keys")
	// w.Resize(fyne.NewSize(1000, 600))
	w.SetContent(
		container.NewVBox(
			container.NewHBox(
				widget.NewLabel("Key"),
				keySelector,
				widget.NewLabel("Scale"),
				scaleSelector,
				scaleLabel,
				layout.NewSpacer(),
			),
			widget.NewSeparator(),
			container.NewVBox(
				widget.NewCard("", "Triads", widget.NewLabel("")),
				widget.NewCard("", "Sevenths", widget.NewLabel("")),
			),
		),
	)
	w.ShowAndRun()
}

func makeScaleLabel(scaleNotes []string) string {
	return fmt.Sprintf("Scale Notes: %s", strings.Join(scaleNotes, " "))
}

func enumerateScale(note, accidental string, index int, intervals []int) []string {
	var scaleNotes []string
	scaleNotes = append(scaleNotes, note+accidental)

	lastNote := note
	ind := index
	notesLength := len(notes)
	for _, interval := range intervals {
		ind += interval

		length := len(notes[index])
		currentNote := notes[index][length-1]

		for _, n := range notes[ind%notesLength] {
			if n.name == lastNote {
				continue
			}
			currentNote = n
			break
		}
		lastNote = currentNote.name
		scaleNotes = append(scaleNotes, currentNote.name+currentNote.accidental)
	}

	return scaleNotes
}
