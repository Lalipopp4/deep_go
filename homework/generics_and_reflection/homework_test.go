package main

import (
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type Person struct {
	Name    string `properties:"name"`
	Address string `properties:"address,omitempty"`
	Age     int    `properties:"age"`
	Married bool   `properties:"married"`
}

func Serialize(person Person) string {
	var (
		tp    = reflect.TypeOf(person)
		value = reflect.ValueOf(person)
		res   = make([]string, 0, value.NumField())
	)

	for i := range value.NumField() {
		field := value.Field(i)
		propTag, ok := tp.Field(i).Tag.Lookup("properties")
		if !ok {
			continue
		}

		tags := strings.Split(propTag, ",")
		canOmit := len(tags) == 2 && tags[1] == "omitempty"
		switch field.Kind() {
		case reflect.String:
			if !canOmit || field.String() != "" {
				res = append(res, tags[0]+"="+field.String())
			}
		case reflect.Int:
			if !canOmit || field.Int() != 0 {
				res = append(res, tags[0]+"="+strconv.FormatInt(field.Int(), 10))
			}
		case reflect.Bool:
			if !canOmit || field.Bool() != false {
				res = append(res, tags[0]+"="+strconv.FormatBool(field.Bool()))
			}
		}

	}

	return strings.Join(res, "\n")
}

func TestSerialization(t *testing.T) {
	tests := map[string]struct {
		person Person
		result string
	}{
		"test case with empty fields": {
			result: "name=\nage=0\nmarried=false",
		},
		"test case with fields": {
			person: Person{
				Name:    "John Doe",
				Age:     30,
				Married: true,
			},
			result: "name=John Doe\nage=30\nmarried=true",
		},
		"test case with omitempty field": {
			person: Person{
				Name:    "John Doe",
				Age:     30,
				Married: true,
				Address: "Paris",
			},
			result: "name=John Doe\naddress=Paris\nage=30\nmarried=true",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result := Serialize(test.person)
			assert.Equal(t, test.result, result)
		})
	}
}
