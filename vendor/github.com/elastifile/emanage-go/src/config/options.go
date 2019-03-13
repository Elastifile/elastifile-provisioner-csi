package config

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/fatih/camelcase"
	"github.com/fatih/structs"
	"github.com/go-errors/errors"
	"github.com/koding/multiconfig"
)

// OptionLoader satisfies the multiconfig.Loader interface.
// It creates on the fly options based on the field names and parses them to load into the given pointer of struct s.
type OptionLoader struct {
	Options   []string
	CamelCase bool

	optionMap map[string]string
}

func (e *OptionLoader) getPrefix(s *structs.Struct) string {
	return ""
}

func (e *OptionLoader) parseOptions() error {
	logger.Debug("Parsing options")

	e.optionMap = map[string]string{}

	for _, opt := range e.Options {
		parts := strings.Split(opt, "=")
		if len(parts) != 2 {
			return errors.Errorf("Invalid option format '%v', expected 'key=value'", opt)
		}

		key, value := parts[0], parts[1]
		logger.Debug("Option", "key", key, "value", value)
		e.optionMap[key] = value
	}

	logger.Debug("optionMap", "optionMap", e.optionMap)

	return nil
}

// Load loads the source into the config defined by struct s
func (e *OptionLoader) Load(s interface{}) error {
	err := e.parseOptions()
	if err != nil {
		return err
	}

	strct := structs.New(s)
	prefix := e.getPrefix(strct)

	for _, field := range strct.Fields() {
		if err := e.processField(prefix, field); err != nil {
			return fmt.Errorf("OptionLoader: %v", err)
		}
	}

	if len(e.optionMap) > 0 {
		return errors.Errorf("Unrecognized options were specified: %v", e.optionMap)
	}

	return nil
}

func (e *OptionLoader) processField(prefix string, field *structs.Field) error {
	fieldName := e.generateFieldName(prefix, field)

	switch field.Kind() {
	case reflect.Struct:
		for _, f := range field.Fields() {
			if err := e.processField(fieldName, f); err != nil {
				return err
			}
		}
	default:
		// Yea, though I walk through the valley of the shadow of
		// death, I shall fear no evil, for thou art with me, thy rod
		// and thy staff, they comfort me.
		v, ok := e.optionMap[fieldName]
		_, isFieldSetter := field.Value().(multiconfig.FieldSetter)

		if !ok && !isFieldSetter {
			return nil
		}

		delete(e.optionMap, fieldName)

		if v == "" && !isFieldSetter {
			return nil
		}

		if v == "" {
			v = fieldName
		}

		defer func() {
			for k := range e.optionMap {
				if strings.HasPrefix(k, fieldName+".") {
					delete(e.optionMap, k)
				}
			}
		}()
		if err := multiconfig.FieldSet(field, v, e.optionMap); err != nil {
			return err
		}
	}

	return nil
}

func (e *OptionLoader) PrintOptions(s interface{}) {
	strct := structs.New(s)

	prefix := e.getPrefix(strct)

	for _, field := range strct.Fields() {
		e.printField(prefix, field)
	}
}

func (e *OptionLoader) printField(prefix string, field *structs.Field) {
	fieldName := e.generateFieldName(prefix, field)

	switch field.Kind() {
	case reflect.Struct:
		for _, f := range field.Fields() {
			e.printField(fieldName, f)
		}
	default:
		fmt.Println("  ", fieldName)
	}
}

// generateFieldName generates the field name combined with the prefix and the
// struct's field name
func (e *OptionLoader) generateFieldName(prefix string, field *structs.Field) string {
	fieldName := strings.ToLower(strings.Join(camelcase.Split(field.Name()), "_"))

	var parts []string
	if prefix != "" {
		parts = append(parts, strings.ToLower(prefix))
	}
	parts = append(parts, fieldName)
	return strings.Join(parts, ".")
}
