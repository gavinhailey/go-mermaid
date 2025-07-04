package flowchart

import (
	"reflect"
	"strings"
	"testing"

	"github.com/gavinhailey/go-mermaid/diagrams/utils/basediagram"
)

// MockIDGenerator is a simple ID generator for testing
type MockIDGenerator struct {
	currentID uint64
}

func (m *MockIDGenerator) NextID() uint64 {
	current := m.currentID
	m.currentID++
	return current
}

func TestNewFlowchart(t *testing.T) {
	tests := []struct {
		name string
		want *Flowchart
	}{
		{
			name: "Create new flowchart with default settings",
			want: &Flowchart{
				BaseDiagram: basediagram.NewBaseDiagram(NewFlowchartConfigurationProperties()),
				Direction:   FlowchartDirectionTopToBottom,
				CurveStyle:  CurveStyleNone,
				classes:     make([]*Class, 0),
				nodes:       make([]*Node, 0),
				subgraphs:   make([]*Subgraph, 0),
				links:       make([]*Link, 0),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewFlowchart()

			// Remove the comparison of idGenerator as it's an interface
			tt.want.idGenerator = got.idGenerator

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewFlowchart() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFlowchart_String(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(*Flowchart)
		contains []string
	}{
		{
			name: "Empty flowchart",
			setup: func(f *Flowchart) {
				f.SetDirection(FlowchartDirectionTopToBottom)
			},
			contains: []string{
				"flowchart TB",
			},
		},
		{
			name: "Flowchart with class",
			setup: func(f *Flowchart) {
				f.AddClass("myClass")
			},
			contains: []string{
				"flowchart TB",
			},
		},
		{
			name: "Flowchart with subgraph",
			setup: func(f *Flowchart) {
				f.AddSubgraph("My Subgraph")
			},
			contains: []string{
				"flowchart TB",
				"subgraph 0 [My Subgraph]",
			},
		},
		{
			name: "Flowchart with multiple elements",
			setup: func(f *Flowchart) {
				f.AddClass("myClass")
				node1 := f.AddNode("My Node 1")
				node2 := f.AddNode("My Node 2")
				f.AddSubgraph("My Subgraph")
				f.AddLink(node1, node2)
			},
			contains: []string{
				"flowchart TB",
				"0@{ shape: rect, label: \"My Node 1\"}",
				"1@{ shape: rect, label: \"My Node 2\"}",
				"subgraph 2 [My Subgraph]",
				"0 --> 1",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flowchart := NewFlowchart()
			tt.setup(flowchart)
			got := flowchart.String()
			for _, want := range tt.contains {
				if !strings.Contains(got, want) {
					t.Errorf("String() missing expected content %q in:\n%s", want, got)
				}
			}
		})
	}
}

func TestFlowchart_AddNode(t *testing.T) {
	tests := []struct {
		name      string
		flowchart *Flowchart
		text      string
		wantNode  *Node
	}{
		{
			name:      "Add simple node",
			flowchart: NewFlowchart(),
			text:      "Test Node",
			wantNode: &Node{
				ID:    "0",
				Text:  "Test Node",
				Shape: NodeShapeProcess,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.flowchart.AddNode(tt.text)

			if !reflect.DeepEqual(got, tt.wantNode) {
				t.Errorf("AddNode() = %v, want %v", got, tt.wantNode)
			}

			if len(tt.flowchart.nodes) != 1 || !reflect.DeepEqual(tt.flowchart.nodes[0], got) {
				t.Errorf("Node not added to flowchart correctly")
			}
		})
	}
}

func TestFlowchart_AddLink(t *testing.T) {
	tests := []struct {
		name      string
		flowchart *Flowchart
		setup     func(*Flowchart) (*Node, *Node)
		wantLink  *Link
	}{
		{
			name:      "Add simple link",
			flowchart: NewFlowchart(),
			setup: func(f *Flowchart) (*Node, *Node) {
				from := f.AddNode("Start")
				to := f.AddNode("End")
				return from, to
			},
			wantLink: &Link{
				Shape:  LinkShapeOpen,
				Head:   LinkArrowTypeArrow,
				Tail:   LinkArrowTypeNone,
				Length: 0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			from, to := tt.setup(tt.flowchart)
			got := tt.flowchart.AddLink(from, to)

			// Update expected link with actual nodes
			tt.wantLink.From = from
			tt.wantLink.To = to

			if !reflect.DeepEqual(got, tt.wantLink) {
				t.Errorf("AddLink() = %v, want %v", got, tt.wantLink)
			}

			if len(tt.flowchart.links) != 1 || !reflect.DeepEqual(tt.flowchart.links[0], got) {
				t.Errorf("Link not added to flowchart correctly")
			}
		})
	}
}

func TestFlowchart_AddSubgraph(t *testing.T) {
	flowchart := NewFlowchart()

	subgraph := flowchart.AddSubgraph("Test Subgraph")
	if subgraph == nil {
		t.Error("AddSubgraph() returned nil")
	}

	if len(flowchart.subgraphs) != 1 {
		t.Errorf("AddSubgraph() resulted in %d subgraphs, want 1", len(flowchart.subgraphs))
	}

	if subgraph.Title != "Test Subgraph" {
		t.Errorf("AddSubgraph() created subgraph with title %q, want %q", subgraph.Title, "Test Subgraph")
	}

	// Test ID generation
	subgraph2 := flowchart.AddSubgraph("Second Subgraph")
	if subgraph2.ID <= subgraph.ID {
		t.Errorf("Second subgraph ID %s should be greater than first subgraph ID %s", subgraph2.ID, subgraph.ID)
	}
}

func TestFlowchart_AddClass(t *testing.T) {
	flowchart := NewFlowchart()

	class := flowchart.AddClass("TestClass")
	if class == nil {
		t.Error("AddClass() returned nil")
	}

	if len(flowchart.classes) != 1 {
		t.Errorf("AddClass() resulted in %d classes, want 1", len(flowchart.classes))
	}

	if class.Name != "TestClass" {
		t.Errorf("AddClass() created class with name %q, want %q", class.Name, "TestClass")
	}

	// Test multiple classes
	class2 := flowchart.AddClass("SecondClass")
	if class2.Name != "SecondClass" {
		t.Errorf("AddClass() created second class with name %q, want %q", class2.Name, "SecondClass")
	}

	if len(flowchart.classes) != 2 {
		t.Errorf("After adding second class, got %d classes, want 2", len(flowchart.classes))
	}
}

func TestFlowchart_SetDirection(t *testing.T) {
	tests := []struct {
		name      string
		direction flowchartDirection
	}{
		{"Top to Bottom", FlowchartDirectionTopToBottom},
		{"Left to Right", FlowchartDirectionLeftRight},
		{"Right to Left", FlowchartDirectionRightLeft},
		{"Bottom Up", FlowchartDirectionBottomUp},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flowchart := NewFlowchart()
			result := flowchart.SetDirection(tt.direction)

			if flowchart.Direction != tt.direction {
				t.Errorf("SetDirection() = %v, want %v", flowchart.Direction, tt.direction)
			}

			if result != flowchart {
				t.Error("SetDirection() should return flowchart for chaining")
			}
		})
	}
}

func TestLink_SetText(t *testing.T) {
	link := NewLink(nil, nil)
	result := link.SetText("Test Text")

	if link.Text != "Test Text" {
		t.Errorf("SetText() = %v, want %v", link.Text, "Test Text")
	}

	if result != link {
		t.Error("SetText() should return link for chaining")
	}
}

func TestLink_SetShape(t *testing.T) {
	tests := []struct {
		name  string
		shape linkShape
	}{
		{"Open", LinkShapeOpen},
		{"Dotted", LinkShapeDotted},
		{"Thick", LinkShapeThick},
		{"Invisible", LinkShapeInvisible},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			link := NewLink(nil, nil)
			result := link.SetShape(tt.shape)

			if link.Shape != tt.shape {
				t.Errorf("SetShape() = %v, want %v", link.Shape, tt.shape)
			}

			if result != link {
				t.Error("SetShape() should return link for chaining")
			}
		})
	}
}

func TestLink_SetLength(t *testing.T) {
	link := NewLink(nil, nil)
	result := link.SetLength(5)

	if link.Length != 5 {
		t.Errorf("SetLength() = %v, want %v", link.Length, 5)
	}

	if result != link {
		t.Error("SetLength() should return link for chaining")
	}
}

func TestLink_SetHead(t *testing.T) {
	tests := []struct {
		name      string
		arrowType linkArrowType
	}{
		{"None", LinkArrowTypeNone},
		{"Arrow", LinkArrowTypeArrow},
		{"Left Arrow", LinkArrowTypeLeftArrow},
		{"Bullet", LinkArrowTypeBullet},
		{"Cross", LinkArrowTypeCross},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			link := NewLink(nil, nil)
			result := link.SetHead(tt.arrowType)

			if link.Head != tt.arrowType {
				t.Errorf("SetHead() = %v, want %v", link.Head, tt.arrowType)
			}

			if result != link {
				t.Error("SetHead() should return link for chaining")
			}
		})
	}
}

func TestLink_SetTail(t *testing.T) {
	tests := []struct {
		name      string
		arrowType linkArrowType
	}{
		{"None", LinkArrowTypeNone},
		{"Arrow", LinkArrowTypeArrow},
		{"Left Arrow", LinkArrowTypeLeftArrow},
		{"Bullet", LinkArrowTypeBullet},
		{"Cross", LinkArrowTypeCross},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			link := NewLink(nil, nil)
			result := link.SetTail(tt.arrowType)

			if link.Tail != tt.arrowType {
				t.Errorf("SetTail() = %v, want %v", link.Tail, tt.arrowType)
			}

			if result != link {
				t.Error("SetTail() should return link for chaining")
			}
		})
	}
}
