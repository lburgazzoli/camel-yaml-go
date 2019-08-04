package main

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

type Property struct {
	ID    string
	Value interface{}
}

type Definition struct {
	ID         string
	Value      string
	Parameters []Property
	Attributes []Property
	Outputs    []Definition
}

func (d *Definition) MarshalYAML() (interface{}, error) {
	definition, err := d.encode()
	if err != nil {
		return nil, err
	}

	answer := make(map[interface{}]interface{})
	answer[d.ID] = definition

	for _, p := range d.Attributes {
		answer[p.ID] = p.Value
	}

	return answer, nil
}

func (d *Definition) UnmarshalYAML(unmarshal func(interface{}) error) error {
	data := make(map[interface{}]interface{})
	if err := unmarshal(&data); err != nil {
		return err
	}

	d.Parameters = make([]Property, 0)
	d.Attributes = make([]Property, 0)

	for k, v := range data {
		switch t := v.(type) {
		case map[interface{}]interface{}:
			d.ID = k.(string)

			if err := d.decode(t); err != nil {
				return err
			}

			delete(data, k)
		default:
			d.Attributes = append(d.Attributes, Property{ID: k.(string), Value: v})
		}
	}

	return nil
}

func (d *Definition) encode() (map[interface{}]interface{}, error) {
	if d.Value != "" {
		return map[interface{}]interface{}{d.ID: d.Value}, nil
	}

	definition := make(map[interface{}]interface{})

	for _, p := range d.Parameters {
		definition[p.ID] = p.Value
	}

	if len(d.Outputs) > 0 {
		steps := make([]map[string]interface{}, 0, len(d.Outputs))

		for _, s := range d.Outputs {
			sd, err := s.encode()
			if err != nil {
				return nil, err
			}

			steps = append(steps, map[string]interface{}{s.ID: sd})
		}

		definition["steps"] = steps

	}

	return definition, nil
}

func (d *Definition) decode(data map[interface{}]interface{}) error {
	for k, v := range data {
		if k != "steps" {
			d.Parameters = append(d.Parameters, Property{ID: k.(string), Value: v})
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

func main() {
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
