package main

import (
	"encoding/json"
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
	parameters map[string]interface{}
	steps      []Definition
}

func (d *Definition) MarshalJSON() ([]byte, error) {
	return nil, nil
}

func (d *Definition) UnmarshalJSON(b []byte) error {
	return nil
}

func (d *Definition) MarshalYAML() (interface{}, error) {
	return nil, nil
}

func (d *Definition) UnmarshalYAML(unmarshal func(interface{}) error) error {
	d.parameters = make(map[string]interface{})
	d.steps = make([]Definition, 0)

	data := make(map[string]interface{})
	if err := unmarshal(&data); err != nil {
		return err
	}

	return d.decode(data)
}

func (d *Definition) decode(data map[string]interface{}) error {
	for k, v := range data {
		if k != "steps" {
			d.parameters[k] = v
		} else {
			switch t := v.(type) {
			case []interface{}:
				for _, s := range t {
					switch m := s.(type) {
					case map[string]interface{}:
						var step Definition
						if err := step.decode(m); err != nil {
							return err
						}

						d.steps = append(d.steps, step)
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
	raw map[string]Definition
}

func (r *Route) MarshalJSON() ([]byte, error) {
	return json.Marshal(&r.raw)
}

func (r *Route) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, &r.raw)
}

func (r *Route) MarshalYAML() (interface{}, error) {
	return yaml.Marshal(&r.raw)
}

func (r *Route) UnmarshalYAML(unmarshal func(interface{}) error) error {
	r.raw = make(map[string]Definition)
	return unmarshal(&r.raw)
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

	b, err := yaml.Marshal(&r)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", r)
	fmt.Printf("%s\n", string(b))
}
