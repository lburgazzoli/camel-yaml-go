package v3

import (
	"fmt"

	yamlv3 "gopkg.in/yaml.v3"
)

// ******************************************
//
//
//
// ******************************************

type Property struct {
	Key   string
	Value interface{}
}

type Definition struct {
	ID         string
	Value      string
	Parameters []Property
	Outputs    []Definition
}

func (d *Definition) MarshalYAML() (interface{}, error) {
	return nil, nil
}

func (d *Definition) UnmarshalYAML(node *yamlv3.Node) error {
	switch node.Kind {
	case yamlv3.MappingNode:
		for i := 0; i < len(node.Content); i += 2 {
			switch {
			case node.Content[i+1].Kind == yamlv3.MappingNode:
				d.ID = node.Content[i].Value
				d.decode(node.Content[i+1])
			case node.Content[i+1].Kind == yamlv3.ScalarNode:
				d.ID = node.Content[i].Value
				d.Value = node.Content[i+1].Value
			default:
				return fmt.Errorf("unknown node kind %v", node.Kind)
			}
		}
	default:
		return fmt.Errorf("unknown node kind %v", node.Kind)
	}

	return nil
}

func (d *Definition) decode(node *yamlv3.Node) error {
	for i := 0; i < len(node.Content); i += 2 {
		switch {
		case node.Content[i+1].Kind == yamlv3.SequenceNode && node.Content[i].Value == "steps":
			d.Outputs = make([]Definition, 0, len(node.Content[i+1].Content))
			err := node.Content[i+1].Decode(&d.Outputs)
			if err != nil {
				return err
			}
		case node.Content[i+1].Kind == yamlv3.ScalarNode:
			d.Parameters = append(d.Parameters, Property{Key: node.Content[i].Value, Value: node.Content[i+1].Value})
		default:
			return fmt.Errorf("unknown node kind %v", node.Kind)
		}
	}

	return nil
}

type Route struct {
	Definition
	Attributes []Property
}

func (r *Route) UnmarshalYAML(node *yamlv3.Node) error {
	switch node.Kind {
	case yamlv3.MappingNode:
		for i := 0; i < len(node.Content); i += 2 {
			switch {
			case node.Content[i+1].Kind == yamlv3.MappingNode:
				r.Definition.ID = node.Content[i].Value
				r.Definition.decode(node.Content[i+1])
			case node.Content[i+1].Kind == yamlv3.ScalarNode:
				r.Attributes = append(r.Attributes, Property{Key: node.Content[i].Value, Value: node.Content[i+1].Value})
			default:
				return fmt.Errorf("unknown node kind %v", node.Kind)
			}
		}
	default:
		return fmt.Errorf("unknown node kind %v", node.Kind)
	}

	return nil
}

const data string = `
id: test
from:
  uri: timer:tick?period=3s
  steps:
    - set-body:
        constant: Hello world!
    - to: "stream:out"
`

func Run() {
	var r Route

	err := yamlv3.Unmarshal([]byte(data), &r)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", r)
}
