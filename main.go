package main

import (
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

const (
	sharp = "\u266f"
	flat  = "\u266d"
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

	chord struct {
		name     string
		position string
		notes    []string
	}

	model struct {
		key        note
		keyIndex   int
		scale      scale
		scaleNotes []string

		keySelector   *widget.Select
		scaleSelector *widget.Select
		scaleLabel    *widget.Label

		triadGrid   *fyne.Container
		seventhGrid *fyne.Container
	}
)

var (
	allNotes [][]note = [][]note{
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

	allScales = []scale{
		{"Major", []int{2, 2, 1, 2, 2, 2}},
		{"Minor", []int{2, 1, 2, 2, 1, 2}},
	}

	keyNames       []string
	keyNotes       = make(map[string]noteIndex)
	scaleNames     []string
	scaleIntervals = make(map[string]scale)
)

func init() {
	for i, g := range allNotes {
		for _, n := range g {
			if n.includeInScale {
				noteName := n.name + n.accidental
				keyNames = append(keyNames, noteName)
				keyNotes[noteName] = noteIndex{n, i}
			}
		}
	}

	for _, s := range allScales {
		scaleNames = append(scaleNames, s.name)
		scaleIntervals[s.name] = s
	}
}

func main() {
	a := app.New()
	a.Settings().SetTheme(&myTheme{})

	key := allNotes[0][0]
	keyIndex := 0
	scale := allScales[0]
	m := model{
		key:        key,
		scale:      scale,
		scaleNotes: enumerateScale(key.name, key.accidental, keyIndex, scale.intervals),
	}

	ui := buildUI(&m)

	w := a.NewWindow("Chords for Keys")
	w.SetContent(ui)
	w.ShowAndRun()
}

func enumerateScale(note, accidental string, index int, intervals []int) []string {
	var scaleNotes []string
	scaleNotes = append(scaleNotes, note+accidental)

	lastNote := note
	ind := index
	notesLength := len(allNotes)
	for _, interval := range intervals {
		ind += interval

		length := len(allNotes[index])
		currentNote := allNotes[index][length-1]

		for _, n := range allNotes[ind%notesLength] {
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

func buildUI(m *model) *fyne.Container {
	m.scaleLabel = widget.NewLabel(strings.Join(m.scaleNotes, " "))

	m.triadGrid = container.NewGridWithColumns(7)
	m.seventhGrid = container.NewGridWithColumns(7)

	m.keySelector = widget.NewSelect(keyNames, func(s string) {
		k := keyNotes[s]
		m.key = k.note
		m.keyIndex = k.index
		refreshUI(m)
	})
	m.keySelector.SetSelectedIndex(0)

	m.scaleSelector = widget.NewSelect(scaleNames, func(s string) {
		m.scale = scaleIntervals[s]
		refreshUI(m)
	})
	m.scaleSelector.SetSelectedIndex(0)

	fillChordGrid(buildTriads(m), m.triadGrid)
	fillChordGrid(buildSevenths(m), m.seventhGrid)

	return container.NewPadded(
		container.NewVBox(
			container.NewHBox(
				widget.NewLabel("Key"),
				m.keySelector,
				layout.NewSpacer(),
				widget.NewLabel("Scale"),
				m.scaleSelector,
				layout.NewSpacer(),
				widget.NewLabel("Scale Notes"),
				m.scaleLabel,
			),
			widget.NewSeparator(),
			container.NewVBox(
				widget.NewCard("", "Triads", m.triadGrid),
				widget.NewCard("", "Sevenths", m.seventhGrid),
			),
		),
	)
}

func buildChords(m *model, pattern []int, suffixes []string, positionNames []string) []chord {
	patLen := len(pattern)
	var chords []chord
	for i, n := range m.scaleNotes {
		c := chord{
			position: positionNames[i],
			name:     n + suffixes[i],
			notes:    make([]string, 0, patLen),
		}
		for _, p := range pattern {
			c.notes = append(c.notes, m.scaleNotes[(i+p)%len(m.scaleNotes)])
		}
		chords = append(chords, c)
	}

	return chords
}

func buildTriads(m *model) []chord {
	pattern := []int{0, 2, 4}
	var suffixes []string
	switch m.scale.name {
	case "Major":
		suffixes = []string{"", "m", "m", "", "", "m", "dim"}
	case "Minor":
		suffixes = []string{"m", "dim", "", "m", "m", "", ""}
	}
	positionNames := []string{"I", "II", "III", "IV", "V", "VI", "VII"}

	return buildChords(m, pattern, suffixes, positionNames)
}

func buildSevenths(m *model) []chord {
	pattern := []int{0, 2, 4, 6}
	var suffixes []string
	switch m.scale.name {
	case "Major":
		suffixes = []string{"M7", "m7", "m7", "M7", "7", "m7", "m7\u266d5"}
	case "Minor":
		suffixes = []string{"m7", "m7\u266d5", "M7", "m7", "m7", "M7", "7"}
	}
	positionNames := []string{"I", "II", "III", "IV", "V", "VI", "VII"}

	return buildChords(m, pattern, suffixes, positionNames)
}

func refreshUI(m *model) {
	m.scaleNotes = enumerateScale(m.key.name, m.key.accidental, m.keyIndex, m.scale.intervals)
	m.scaleLabel.SetText(strings.Join(m.scaleNotes, " "))

	fillChordGrid(buildTriads(m), m.triadGrid)
	fillChordGrid(buildSevenths(m), m.seventhGrid)
}

func fillChordGrid(chords []chord, grid *fyne.Container) {
	grid.RemoveAll()
	for _, c := range chords {
		card := widget.NewCard(c.name, c.position, widget.NewLabel(strings.Join(c.notes, " ")))
		grid.Add(card)
	}
}
