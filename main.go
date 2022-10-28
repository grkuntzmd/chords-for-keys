package main

import (
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

const chromaticScaleLen = 12

type (
	scale struct {
		name      string
		intervals []int
	}

	model struct {
		key            string
		scale          string
		scaleNotes     []string
		scaleIntervals []int

		scaleLabel    *widget.Label
		keySelector   *widget.Select
		scaleSelector *widget.Select

		triadGrid         *fyne.Container
		seventhGrid       *fyne.Container
		secondaryDomGrid  *fyne.Container
		secondaryLeadGrid *fyne.Container
		tritoneSubGrid    *fyne.Container
	}

	chord struct {
		name     string
		position string
		notes    []string
	}
)

var (
	keyNames       = []string{"C", "C♯", "D♭", "D", "D♯", "E♭", "E", "F", "F♯", "G♭", "G", "G#", "A♭", "A", "A♯", "B♭", "B"}
	chromaticScale = []string{"C", "D♭", "D", "E♭", "E", "F", "G♭", "G", "A♭", "A", "B♭", "B"}

	noteEquivalents = [][]string{
		{"B♯", "C"},
		{"C♯", "D♭"},
		{"D"},
		{"D♯", "E♭"},
		{"E", "F♭"},
		{"E♯", "F"},
		{"F♯", "G♭"},
		{"G"},
		{"G♯", "A♭"},
		{"A"},
		{"A♯", "B♭"},
		{"B", "C♭"},
	}

	noteIndexes map[string]int

	majorIntervals = []int{2, 2, 1, 2, 2, 2}
	scaleNames     []string
	scaleIntervals map[string][]int
)

func init() {
	// Make the noteIndexes map.
	noteIndexes = make(map[string]int)
	for i, r := range noteEquivalents {
		for _, n := range r {
			noteIndexes[n] = i
		}
	}

	// Make scale names array and intervals map.
	scales := []scale{
		{"Major", majorIntervals},
		{"Minor", []int{2, 1, 2, 2, 1, 2}},
	}
	scaleIntervals = make(map[string][]int)
	for _, s := range scales {
		scaleNames = append(scaleNames, s.name)
		scaleIntervals[s.name] = s.intervals
	}
}

func main() {
	a := app.New()
	a.Settings().SetTheme(&myTheme{})

	key := keyNames[0]
	scale := scaleNames[0]
	m := model{
		key:            key,
		scale:          scale,
		scaleNotes:     enumerateScale(key, scaleIntervals[scale]),
		scaleIntervals: scaleIntervals[scale],
	}

	w := a.NewWindow("Chords for Keys")
	w.SetContent(m.buildUI())
	w.ShowAndRun()
}

func enumerateScale(note string, intervals []int) []string {
	var notes []string
	notes = append(notes, note)

	lastNatural := note[0:1]
	index := noteIndexes[note]
	numNotes := len(noteEquivalents)

	for _, interval := range intervals {
		index += interval
		currentNote := noteEquivalents[index%numNotes][0]

		for _, n := range noteEquivalents[index%numNotes] {
			if n[0:1] == lastNatural {
				continue
			}
			currentNote = n
			break
		}
		lastNatural = currentNote[0:1]
		notes = append(notes, currentNote)
	}

	return notes
}

func fillChordGrid(chords []chord, grid *fyne.Container) {
	grid.RemoveAll()
	for _, c := range chords {
		card := widget.NewCard(c.name, c.position, widget.NewLabel(strings.Join(c.notes, " ")))
		grid.Add(card)
	}
}

func (m *model) buildUI() *fyne.Container {
	m.scaleLabel = widget.NewLabel(strings.Join(m.scaleNotes, " "))

	m.triadGrid = container.NewGridWithColumns(7)
	m.seventhGrid = container.NewGridWithColumns(7)
	m.secondaryDomGrid = container.NewGridWithColumns(6)
	m.secondaryLeadGrid = container.NewGridWithColumns(6)
	m.tritoneSubGrid = container.NewGridWithColumns(1)

	m.keySelector = widget.NewSelect(keyNames, func(s string) {
		m.key = s
		m.refreshUI()
	})
	m.keySelector.SetSelectedIndex(0)

	m.scaleSelector = widget.NewSelect(scaleNames, func(s string) {
		m.scale = s
		m.scaleIntervals = scaleIntervals[s]
		m.refreshUI()
	})
	m.scaleSelector.SetSelectedIndex(0)

	m.refreshUI()

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
				widget.NewCard("", "Secondary Dominants", m.secondaryDomGrid),
				widget.NewCard("", "Secondary Lead Tones", m.secondaryLeadGrid),
				widget.NewCard("", "Tritone Substitution", m.tritoneSubGrid),
			),
		),
	)
}

func (m *model) refreshUI() {
	m.scaleNotes = enumerateScale(m.key, m.scaleIntervals)
	m.scaleLabel.SetText(strings.Join(m.scaleNotes, " "))

	fillChordGrid(m.buildTriads(), m.triadGrid)
	fillChordGrid(m.buildSevenths(), m.seventhGrid)
	fillChordGrid(m.buildSecondaryDoms(), m.secondaryDomGrid)
	fillChordGrid(m.buildSecondaryLeads(), m.secondaryLeadGrid)
	fillChordGrid(m.buildTritoneSubstition(), m.tritoneSubGrid)
}

func (m *model) buildChords(pattern []int, suffixes []string, positionNames []string) []chord {
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

func (m *model) buildTriads() []chord {
	pattern := []int{0, 2, 4}
	var suffixes []string
	switch m.scale {
	case "Major":
		suffixes = []string{"", "m", "m", "", "", "m", "°"}
	case "Minor":
		suffixes = []string{"m", "°", "", "m", "m", "", ""}
	}
	positionNames := []string{"I", "II", "III", "IV", "V", "VI", "VII"}

	return m.buildChords(pattern, suffixes, positionNames)
}

func (m *model) buildSevenths() []chord {
	pattern := []int{0, 2, 4, 6}
	var suffixes []string
	switch m.scale {
	case "Major":
		suffixes = []string{"M7", "m7", "m7", "M7", "7", "m7", "m7♭5"}
	case "Minor":
		suffixes = []string{"m7", "m7♭5", "M7", "m7", "m7", "M7", "7"}
	}
	positionNames := []string{"I⁷", "II⁷", "III⁷", "IV⁷", "V⁷", "VI⁷", "VII⁷"}

	return m.buildChords(pattern, suffixes, positionNames)
}

func (m *model) buildSecondaryDoms() []chord {
	pattern := []int{0, 2, 4, 6}
	positionNames := []string{"", "V⁷ / II", "V⁷ / III", "V⁷ / IV", "V⁷ / V", "V⁷ / VI", "V⁷ / VII"}
	var chords []chord
	for i, s := range m.scaleNotes {
		if i == 0 {
			continue
		}
		sec := enumerateScale(s, majorIntervals)
		fifth := sec[4]
		secLen := len(sec)
		c := chord{
			name:     fifth + "7",
			position: positionNames[i],
			notes:    []string{},
		}
		for _, p := range pattern {
			c.notes = append(c.notes, sec[(p+4)%secLen])
		}
		chords = append(chords, c)
	}

	return chords
}

func (m *model) buildSecondaryLeads() []chord {
	pattern := []int{0, 2, 4}
	positionNames := []string{"", "VII° / II", "VII° / III", "VII° / IV", "VII° / V", "VII° / VI", "VII° / VII"}
	var chords []chord
	for i, s := range m.scaleNotes {
		if i == 0 {
			continue
		}
		sec := enumerateScale(s, majorIntervals)
		seventh := sec[6]
		secLen := len(sec)
		c := chord{
			name:     seventh + "°",
			position: positionNames[i],
			notes:    []string{},
		}
		// fmt.Printf("note: %s, scale: %v, chord: %v\n", s, sec, c)
		for _, p := range pattern {
			c.notes = append(c.notes, sec[(p+6)%secLen])
		}
		chords = append(chords, c)
	}

	return chords
}

func (m *model) buildTritoneSubstition() []chord {
	pattern := []int{0, 4, 7, 10}
	index := noteIndexes[m.key] + 1
	c := chord{
		name:     chromaticScale[index%chromaticScaleLen] + "7",
		position: "sub VII⁷ / V⁷",
		notes:    []string{},
	}

	for _, p := range pattern {
		c.notes = append(c.notes, chromaticScale[(index+p)%chromaticScaleLen])
	}

	return []chord{c}
}
