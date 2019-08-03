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

type Definition struct {
	id         string
	parameters map[string]interface{}
	steps      []Definition
}

func (d *Definition) encode() (map[interface{}]interface{}, error) {
	data := make(map[interface{}]interface{})

	for k, v := range d.parameters {
		data[k] = v
	}

	if len(d.steps) > 0 {
		steps := make([]map[string]interface{}, 0, len(d.steps))

		for _, s := range d.steps {
			sd, err := s.encode()
			if err != nil {
				return nil, err
			}

			steps = append(steps, map[string]interface{}{s.id: sd})
		}

		data["steps"] = steps

	}
	return data, nil
}

func (d *Definition) decode(data map[interface{}]interface{}) error {
	d.parameters = make(map[string]interface{})

	for k, v := range data {
		if k != "steps" {
			d.parameters[k.(string)] = v
		} else {
			switch t := v.(type) {
			case []interface{}:
				d.steps = make([]Definition, 0, len(t))

				for _, s := range t {
					switch m := s.(type) {
					case map[interface{}]interface{}:
						for k, v := range m {
							step := Definition{id: k.(string)}
							if err := step.decode(v.(map[interface{}]interface{})); err != nil {
								return err
							}

							d.steps = append(d.steps, step)
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
	parameters map[string]interface{}
	definition Definition
}

func (r *Route) MarshalYAML() (interface{}, error) {
	data := make(map[string]interface{})

	for k, v := range r.parameters {
		data[k] = v
	}

	if r.definition.id != "" {
		d, err := r.definition.encode()
		if err != nil {
			return nil, err
		}

		data[r.definition.id] = d
	}

	return yaml.Marshal(data)
}

func (r *Route) UnmarshalYAML(unmarshal func(interface{}) error) error {
	data := make(map[interface{}]interface{})
	if err := unmarshal(&data); err != nil {
		return err
	}

	r.parameters = make(map[string]interface{})

	for k, v := range data {
		switch t := v.(type) {
		case map[interface{}]interface{}:
			r.definition.id = k.(string)

			if err := r.definition.decode(t); err != nil {
				return err
			}

			delete(data, k)
		default:
			r.parameters[k.(string)] = v
		}
	}

	return nil
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
from:
  uri: timer:tick?period=3s
  steps:
    - set-body:
        constant: Hello world!
`

func main() {
	var r Route

	if err := yaml.Unmarshal([]byte(data), &r); err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", r)
}
