package v2

import (
	"errors"
	"fmt"

	"gopkg.in/yaml.v2"
)

// ******************************************
//
//
//
// ******************************************

type Definition struct {
	ID         string        `yaml:"-"`
	Value      string        `yaml:"-"`
	Parameters yaml.MapSlice `yaml:"-"`
	Attributes yaml.MapSlice `yaml:"-"`
	Outputs    []Definition  `yaml:"steps"`
}

func (d *Definition) MarshalYAML() (interface{}, error) {
	definition, err := d.encode()
	if err != nil {
		return nil, err
	}

	var answer yaml.MapSlice

	if len(d.Attributes) > 0 {
		answer = append(answer, d.Attributes...)
	}
	if d.ID != "" {
		answer = append(answer, yaml.MapItem{
			Key:   d.ID,
			Value: definition,
		})
	}

	return answer, nil
}

func (d *Definition) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var data yaml.MapSlice
	if err := unmarshal(&data); err != nil {
		return err
	}

	d.Parameters = make([]yaml.MapItem, 0)
	d.Attributes = make([]yaml.MapItem, 0)

	for _, v := range data {
		switch t := v.Value.(type) {
		case map[interface{}]interface{}:
			d.ID = v.Key.(string)

			if err := d.decode(t); err != nil {
				return err
			}
		default:
			d.Attributes = append(d.Attributes, v)
		}
	}

	return nil
}

func (d *Definition) encode() (yaml.MapSlice, error) {
	if d.Value != "" {
		return []yaml.MapItem{
			{
				Key:   d.ID,
				Value: d.Value,
			},
		}, nil
	}

	var definition yaml.MapSlice

	if len(d.Parameters) > 0 {
		definition = append(definition, d.Parameters...)
	}
	if len(d.Outputs) > 0 {
		steps := make([]yaml.MapItem, 0, len(d.Outputs))

		for _, s := range d.Outputs {
			sd, err := s.encode()
			if err != nil {
				return nil, err
			}

			steps = append(steps, yaml.MapItem{
				Key:   s.ID,
				Value: sd,
			})
		}

		definition = append(definition, yaml.MapItem{
			Key:   "steps",
			Value: steps,
		})
	}

	return definition, nil
}

func (d *Definition) decode(data map[interface{}]interface{}) error {
	for k, v := range data {
		if k != "steps" {
			d.Parameters = append(d.Parameters, yaml.MapItem{Key: k.(string), Value: v})
		} else {
			switch t := v.(type) {
			case []interface{}:
				d.Outputs = make([]Definition, 0, len(t))

				for _, step := range t {
					switch definition := step.(type) {
					case map[interface{}]interface{}:
						for k, v := range definition {
							step := Definition{ID: k.(string)}
							switch content := v.(type) {
							case map[interface{}]interface{}:
								if len(content) != 1 {
									return errors.New("")
								}
								if err := step.decode(content); err != nil {
									return err
								}
							case string:
								step.Value = content
							default:
								return errors.New("")
							}

							d.Outputs = append(d.Outputs, step)
						}
					default:
						return errors.New("")
					}
				}
			default:
				return errors.New("")
			}
		}
	}

	return nil
}

// ******************************************
//
//
//
// ******************************************

// Route --
type Route struct {
	Definition
}

// ******************************************
//
//
//
// ******************************************

func (r *Route) DeepCopy() *Route {
	if r == nil {
		return nil
	}

	b, err := yaml.Marshal(r)
	if err != nil {
		panic(err)
	}

	var out Route

	err = yaml.Unmarshal(b, &out)
	if err != nil {
		panic(err)
	}

	return &out
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

	err := yaml.Unmarshal([]byte(data), &r)
	if err != nil {
		panic(err)
	}

	b, err := yaml.Marshal(&r)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", r)
	fmt.Printf("%s\n", string(b[:]))
}
