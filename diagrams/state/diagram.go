// Package state provides functionality for creating Mermaid state diagrams
package state

import (
	"strings"

	"github.com/gavinhailey/go-mermaid/diagrams/utils"
	"github.com/gavinhailey/go-mermaid/diagrams/utils/basediagram"
)

// Base string formats for state diagrams
const (
	baseDiagramType string = "stateDiagram-v2\n"
)

// Diagram represents a state diagram with states, transitions, and rendering options.
type Diagram struct {
	basediagram.BaseDiagram[StateConfigurationProperties]
	States      []*State
	Transitions []*Transition
}

// NewDiagram creates a new state diagram with default settings.
func NewDiagram() *Diagram {
	return &Diagram{
		BaseDiagram: basediagram.NewBaseDiagram(NewStateConfigurationProperties()),
		States:      make([]*State, 0),
		Transitions: make([]*Transition, 0),
	}
}

// AddState creates and adds a new state to the diagram.
func (d *Diagram) AddState(id, description string, stateType StateType) *State {
	state := NewState(id, description, stateType)
	d.States = append(d.States, state)
	return state
}

// AddTransition creates and adds a new transition between states.
func (d *Diagram) AddTransition(from, to *State, description string) *Transition {
	transition := NewTransition(from, to, description)
	d.Transitions = append(d.Transitions, transition)
	return transition
}

// String generates a Mermaid-formatted string representation of the state diagram.
func (d *Diagram) String() string {
	var sb strings.Builder

	sb.WriteString(baseDiagramType)

	for _, state := range d.States {
		sb.WriteString(state.String(""))
	}

	for _, transition := range d.Transitions {
		sb.WriteString(transition.String(""))
	}

	return d.BaseDiagram.String(sb.String())
}

// RenderToFile saves the diagram to a file at the specified path.
func (d *Diagram) RenderToFile(path string) error {
	return utils.RenderToFile(path, d.String())
}
