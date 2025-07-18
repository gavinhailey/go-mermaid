package block

import (
	"reflect"
	"strings"
	"testing"

	"github.com/gavinhailey/go-mermaid/diagrams/utils/basediagram"
)

func TestNewDiagram(t *testing.T) {
	got := NewDiagram()
	want := &Diagram{
		BaseDiagram: basediagram.NewBaseDiagram(NewBlockConfigurationProperties()),
		Blocks:      make([]*Block, 0),
		Links:       make([]*Link, 0),
		Columns:     0,
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("NewDiagram() = %v, want %v", got, want)
	}
}

func TestDiagram_SetColumns(t *testing.T) {
	diagram := NewDiagram()
	result := diagram.SetColumns(3)

	if diagram.Columns != 3 {
		t.Errorf("SetColumns() = %v, want %v", diagram.Columns, 3)
	}

	if result != diagram {
		t.Error("SetColumns() should return diagram for chaining")
	}
}

func TestDiagram_AddBlock(t *testing.T) {
	diagram := NewDiagram()
	block := diagram.AddBlock("Test Block")

	if block == nil {
		t.Error("AddBlock() returned nil")
	}

	if block.Text != "Test Block" {
		t.Errorf("AddBlock() text = %v, want %v", block.Text, "Test Block")
	}

	if len(diagram.Blocks) != 1 {
		t.Errorf("AddBlock() resulted in %d blocks, want 1", len(diagram.Blocks))
	}
}

func TestDiagram_AddLink(t *testing.T) {
	diagram := NewDiagram()
	from := diagram.AddBlock("From")
	to := diagram.AddBlock("To")
	link := diagram.AddLink(from, to)

	if link == nil {
		t.Error("AddLink() returned nil")
	}

	if !reflect.DeepEqual(link.From, from) || !reflect.DeepEqual(link.To, to) {
		t.Error("AddLink() nodes do not match expected values")
	}

	if len(diagram.Links) != 1 {
		t.Errorf("AddLink() resulted in %d links, want 1", len(diagram.Links))
	}
}

func TestDiagram_String(t *testing.T) {
	tests := []struct {
		name     string
		diagram  *Diagram
		setup    func(*Diagram)
		contains []string
	}{
		{
			name:    "Empty diagram",
			diagram: NewDiagram(),
			contains: []string{
				"block-beta\n",
			},
		},
		{
			name:    "Diagram with columns",
			diagram: NewDiagram(),
			setup: func(d *Diagram) {
				d.SetColumns(3)
			},
			contains: []string{
				"block-beta",
				"columns 3",
			},
		},
		{
			name:    "Diagram with blocks and links",
			diagram: NewDiagram(),
			setup: func(d *Diagram) {
				b1 := d.AddBlock("Start")
				b2 := d.AddBlock("End")
				d.AddLink(b1, b2)
			},
			contains: []string{
				"block-beta\n",
				"[\"Start\"]",
				"[\"End\"]",
				"-->",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(tt.diagram)
			}

			got := tt.diagram.String()
			for _, want := range tt.contains {
				if !strings.Contains(got, want) {
					t.Errorf("Diagram.String() missing expected content: %v\nFull output:\n%v", want, got)
				}
			}
		})
	}
}

func TestDiagram_AddColumn(t *testing.T) {
	diagram := NewDiagram()
	initial := diagram.Columns

	result := diagram.AddColumn()

	if diagram.Columns != initial+1 {
		t.Errorf("AddColumn() = %v, want %v", diagram.Columns, initial+1)
	}

	if result != diagram {
		t.Error("AddColumn() should return diagram for chaining")
	}
}

func TestDiagram_RemoveColumn(t *testing.T) {
	diagram := NewDiagram()
	diagram.SetColumns(3)
	initial := diagram.Columns

	result := diagram.RemoveColumn()

	if diagram.Columns != initial-1 {
		t.Errorf("RemoveColumn() = %v, want %v", diagram.Columns, initial-1)
	}

	if result != diagram {
		t.Error("RemoveColumn() should return diagram for chaining")
	}
}

func TestDiagram_AddSpace(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(*Diagram)
		contains []string
	}{
		{
			name: "Add single space",
			setup: func(d *Diagram) {
				d.AddSpace()
			},
			contains: []string{
				"space",
			},
		},
		{
			name: "Add multiple spaces",
			setup: func(d *Diagram) {
				d.AddSpace()
				d.AddSpace()
			},
			contains: []string{
				"space",
				"space",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diagram := NewDiagram()
			tt.setup(diagram)
			got := diagram.String()
			for _, want := range tt.contains {
				if !strings.Contains(got, want) {
					t.Errorf("String() missing expected content %q in:\n%s", want, got)
				}
			}
		})
	}
}

func TestDiagram_AddSpaceWithWidth(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(*Diagram)
		contains []string
	}{
		{
			name: "Add space with width 2",
			setup: func(d *Diagram) {
				d.AddSpaceWithWidth(2)
			},
			contains: []string{
				"space:2",
			},
		},
		{
			name: "Add multiple spaces with different widths",
			setup: func(d *Diagram) {
				d.AddSpaceWithWidth(2)
				d.AddSpaceWithWidth(3)
			},
			contains: []string{
				"space:2",
				"space:3",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diagram := NewDiagram()
			tt.setup(diagram)
			got := diagram.String()
			for _, want := range tt.contains {
				if !strings.Contains(got, want) {
					t.Errorf("String() missing expected content %q in:\n%s", want, got)
				}
			}
		})
	}
}
