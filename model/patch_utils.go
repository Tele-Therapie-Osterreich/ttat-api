package model

import "github.com/pkg/errors"

func stringUpdate(updates map[string]interface{}, k string, dst *string) error {
	v, ok := updates[k]
	if !ok {
		return nil
	}
	s, ok := v.(string)
	if !ok {
		return errors.New("non-string value for '" + k + "'")
	}
	*dst = s
	return nil
}

func optStringUpdate(updates map[string]interface{}, k string, dst **string) error {
	v, ok := updates[k]
	if !ok {
		return nil
	}
	s, ok := v.(string)
	if !ok {
		return errors.New("non-string value for '" + k + "'")
	}
	*dst = &s
	return nil
}

func boolUpdate(updates map[string]interface{}, k string, dst *bool) error {
	v, ok := updates[k]
	if !ok {
		return nil
	}
	b, ok := v.(bool)
	if !ok {
		return errors.New("non-boolean value for '" + k + "'")
	}
	*dst = b
	return nil
}
