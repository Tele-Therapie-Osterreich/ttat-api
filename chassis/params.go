package chassis

import (
	"net/url"
	"strconv"
)

// IntParam extracts an integer URL query parameter.
func IntParam(qs url.Values, k string, dst *uint) error {
	s := qs.Get(k)
	if s != "" {
		i, err := strconv.Atoi(s)
		if err != nil {
			return err
		}
		*dst = uint(i)
	}
	return nil
}

// FloatParam extracts a floating point URL query parameter.
func FloatParam(qs url.Values, k string, dst **float64) error {
	s := qs.Get(k)
	if s != "" {
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return err
		}
		if f >= 0.0 {
			*dst = &f
		}
	}
	return nil
}

// StringParam extracts a string URL query parameter.
func StringParam(qs url.Values, k string, dst **string) {
	s := qs.Get(k)
	if s != "" {
		val := s
		*dst = &val
	}
}
