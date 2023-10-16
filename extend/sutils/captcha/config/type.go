// Copyright 2009  The "config" Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package config

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// Substitutes values, calculated by callback, on matching regex
func (c *Config) computeVar(beforeValue *string, regx *regexp.Regexp, headsz, tailsz int, withVar func(*string) string) (*string, error) {
	var i int
	computedVal := beforeValue
	for i = 0; i < _DEPTH_VALUES; i++ { // keep a sane depth

		vr := regx.FindStringSubmatchIndex(*computedVal)
		if len(vr) == 0 {
			break
		}

		varname := (*computedVal)[vr[headsz]:vr[headsz+1]]
		varVal := withVar(&varname)
		if varVal == "" {
			return &varVal, errors.New(fmt.Sprintf("Option not found: %s", varname))
		}

		// substitute by new value and take off leading '%(' and trailing ')s'
		//  %(foo)s => headsz=2, tailsz=2
		//  ${foo}  => headsz=2, tailsz=1
		newVal := (*computedVal)[0:vr[headsz]-headsz] + varVal + (*computedVal)[vr[headsz+1]+tailsz:]
		computedVal = &newVal
	}

	if i == _DEPTH_VALUES {
		retVal := ""
		return &retVal,
			errors.New(
				fmt.Sprintf(
					"Possible cycle while unfolding variables: max depth of %d reached",
					strconv.Itoa(_DEPTH_VALUES),
				),
			)
	}

	return computedVal, nil
}

// Bool has the same behaviour as String but converts the response to bool.
// See "boolString" for string values converted to bool.
func (c *Config) Bool(section string, option string) (value bool, err error) {
	sv, err := c.String(section, option)
	if err != nil {
		return false, err
	}

	value, ok := boolString[strings.ToLower(sv)]
	if !ok {
		return false, errors.New("could not parse bool value: " + sv)
	}

	return value, nil
}

// BoolMuti has the same behaviour as String but converts the response to bool.
// See "boolString" for string values converted to bool.
func (c *Config) BoolMuti(section string, option string) (retValue []bool, err error) {
	retValue = []bool{}
	sv, err := c.StringMuti(section, option)
	if err != nil {
		return retValue, err
	}

	for i := 0; i < len(sv); i++ {
		value, ok := boolString[strings.ToLower(sv[i])]
		if !ok {
			return retValue, errors.New("could not parse bool value: " + sv[i])
		}

		retValue = append(retValue, value)
	}

	return retValue, nil
}

// Float has the same behaviour as String but converts the response to float.
func (c *Config) Float(section string, option string) (value float64, err error) {
	sv, err := c.String(section, option)
	if err == nil {
		value, err = strconv.ParseFloat(sv, 64)
	}

	return value, err
}

// FloatMuti has the same behaviour as String but converts the response to float.
func (c *Config) FloatMuti(section string, option string) (retValue []float64, err error) {
	retValue = []float64{}
	sv, err := c.StringMuti(section, option)
	for i := 0; i < len(sv); i++ {
		var value float64
		if err == nil {
			value, err = strconv.ParseFloat(sv[i], 64)
		} else {
			value = 0
		}
		retValue = append(retValue, value)
	}

	return retValue, err
}

// Int has the same behaviour as String but converts the response to int.
func (c *Config) Int(section string, option string) (value int, err error) {
	sv, err := c.String(section, option)
	if err == nil {
		value, err = strconv.Atoi(sv)
	}

	return value, err
}

// Int has the same behaviour as String but converts the response to int.
func (c *Config) IntMuti(section string, option string) (retValue []int, err error) {
	retValue = []int{}
	sv, err := c.StringMuti(section, option)
	for i := 0; i < len(sv); i++ {
		var value int
		if err == nil {
			value, err = strconv.Atoi(sv[i])
		} else {
			value = 0
		}
		retValue = append(retValue, value)
	}

	return retValue, err
}

// RawString gets the (raw) string value for the given option in the section.
// The raw string value is not subjected to unfolding, which was illustrated in
// the beginning of this documentation.
//
// It returns an error if either the section or the option do not exist.
func (c *Config) RawString(section string, option string) (value string, err error) {
	if _, ok := c.data[section]; ok {
		if tValue, ok := c.data[section][option]; ok {
			return tValue.v, nil
		}
	}
	return c.RawStringDefault(option)
}

// RawStringMuti gets the (raw) string value for the given option in the section.
// The raw string value is not subjected to unfolding, which was illustrated in
// the beginning of this documentation.
//
// It returns an error if either the section or the option do not exist.
func (c *Config) RawStringMuti(section string, option string) (value []string, err error) {
	if _, ok := c.data[section]; ok {
		if tValue, ok := c.data[section][option]; ok {
			return tValue.vMuti, nil
		}
	}
	return c.RawStringMutiDefault(option)
}

// RawStringDefault gets the (raw) string value for the given option from the
// DEFAULT section.
//
// It returns an error if the option does not exist in the DEFAULT section.
func (c *Config) RawStringDefault(option string) (value string, err error) {
	if tValue, ok := c.data[DEFAULT_SECTION][option]; ok {
		return tValue.v, nil
	}
	return "", OptionError(option)
}

// RawStringMutiDefault gets the (raw) string value for the given option from the
// DEFAULT section.
//
// It returns an error if the option does not exist in the DEFAULT section.
func (c *Config) RawStringMutiDefault(option string) (value []string, err error) {
	if tValue, ok := c.data[DEFAULT_SECTION][option]; ok {
		return tValue.vMuti, nil
	}
	return []string{}, OptionError(option)
}

// String gets the string value for the given option in the section.
// If the value needs to be unfolded (see e.g. %(host)s example in the beginning
// of this documentation), then String does this unfolding automatically, up to
// _DEPTH_VALUES number of iterations.
//
// It returns an error if either the section or the option do not exist, or the
// unfolding cycled.
func (c *Config) String(section string, option string) (value string, err error) {
	value, err = c.RawString(section, option)
	if err != nil {
		return "", err
	}

	// % variables
	computedVal, err := c.computeVar(&value, varRegExp, 2, 2, func(varName *string) string {
		lowerVar := *varName
		// search variable in default section as well as current section
		varVal, _ := c.data[DEFAULT_SECTION][lowerVar]
		if _, ok := c.data[section][lowerVar]; ok {
			varVal = c.data[section][lowerVar]
		}
		return varVal.v
	})
	value = *computedVal

	if err != nil {
		return value, err
	}

	// $ environment variables
	computedVal, err = c.computeVar(&value, envVarRegExp, 2, 1, func(varName *string) string {
		return os.Getenv(*varName)
	})
	value = *computedVal
	return value, err
}

// StringMuti gets the string value for the given option in the section.
// If the value needs to be unfolded (see e.g. %(host)s example in the beginning
// of this documentation), then String does this unfolding automatically, up to
// _DEPTH_VALUES number of iterations.
//
// It returns an error if either the section or the option do not exist, or the
// unfolding cycled.
func (c *Config) StringMuti(section string, option string) (retValue []string, err error) {
	retValue = []string{}
	values, err := c.RawStringMuti(section, option)
	if err != nil {
		return []string{}, err
	}

	for i := 0; i < len(values); i++ {

		value := values[i]
		// % variables
		computedVal, err := c.computeVar(&value, varRegExp, 2, 2, func(varName *string) string {
			lowerVar := *varName
			// search variable in default section as well as current section
			varVal, _ := c.data[DEFAULT_SECTION][lowerVar]
			if _, ok := c.data[section][lowerVar]; ok {
				varVal = c.data[section][lowerVar]
			}
			return varVal.v
		})
		value = *computedVal

		if err != nil {
			return retValue, err
		}

		// $ environment variables
		computedVal, err = c.computeVar(&value, envVarRegExp, 2, 1, func(varName *string) string {
			return os.Getenv(*varName)
		})
		value = *computedVal

		retValue = append(retValue, value)
	}
	return retValue, err
}
