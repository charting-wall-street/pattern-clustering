package pattern

import (
	"bytes"
	"strconv"
)

type EventSequence struct {
	Empty  bool
	Events []float32
}

type Pattern struct {
	Sequences map[string]EventSequence
}

func (p *Pattern) Length() int {
	for _, sequence := range p.Sequences {
		return len(sequence.Events)
	}
	return -1
}

func (p *Pattern) String() string {
	buffer := bytes.Buffer{}
	for key, sequence := range p.Sequences {
		buffer.WriteString(key)
		buffer.WriteString("   ")
		for _, event := range sequence.Events {
			if event != 0 {
				s := strconv.FormatFloat(float64(event), 'f', 2, 32)
				for len(s) < 7 {
					s = s + " "
					if len(s) < 7 {
						s = " " + s
					}
				}
				buffer.WriteString(s)
			} else {
				buffer.WriteString("   _   ")
			}
			buffer.WriteRune(' ')
		}
		buffer.WriteRune('\n')
	}
	return buffer.String()
}

func NewTracker(p *Pattern) *Tracker {
	return &Tracker{
		Pattern:   p,
		Frequency: 1,
	}
}

type Tracker struct {
	Pattern   *Pattern
	Frequency int
}

func NewCluster(ts []*Tracker) *Cluster {
	freq := 0
	for _, t := range ts {
		freq += t.Frequency
	}
	return &Cluster{
		Patterns:  ts,
		Frequency: freq,
	}
}

type Cluster struct {
	Patterns  []*Tracker
	Frequency int
}
