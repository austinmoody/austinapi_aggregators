package main

import (
	"fmt"
	"strings"
)

type TypeChoices struct {
	Options []string
	Value   string
}

func (s *TypeChoices) Set(value string) error {
	if !contains(s.Options, value) {
		return fmt.Errorf("Invalid value, must be one of: %s", strings.Join(s.Options, ", "))
	}
	s.Value = value
	return nil
}

func (s *TypeChoices) String() string {
	return s.Value
}

func contains(slice []string, value string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}
