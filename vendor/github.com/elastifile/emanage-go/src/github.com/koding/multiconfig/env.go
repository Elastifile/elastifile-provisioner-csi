package multiconfig

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/fatih/camelcase"
	"github.com/fatih/structs"

	"github.com/go-errors/errors"
)

// EnvironmentLoader satisifies the loader interface. It loads the
// configuration from the environment variables in the form of
// STRUCTNAME_FIELDNAME.
type EnvironmentLoader struct {
	// Prefix prepends given string to every environment variable
	// {STRUCTNAME}_FIELDNAME will be {PREFIX}_FIELDNAME
	Prefix string

	// CamelCase adds a seperator for field names in camelcase form. A
	// fieldname of "AccessKey" would generate a environment name of
	// "STRUCTNAME_ACCESSKEY". If CamelCase is enabled, the environment name
	// will be generated in the form of "STRUCTNAME_ACCESS_KEY"
	CamelCase bool
}

func (e *EnvironmentLoader) getPrefix(s *structs.Struct) string {
	if e.Prefix != "" {
		return e.Prefix
	}

	return s.Name()
}

// Load loads the source into the config defined by struct s
func (e *EnvironmentLoader) Load(s interface{}) error {
	strct := structs.New(s)

	prefix := e.getPrefix(strct)

	for _, field := range strct.Fields() {
		if err := e.processField(prefix, field); err != nil {
			return errors.Errorf("EnvironmentLoader: %v", err)
		}
	}

	return nil
}

// processField gets leading name for the env variable and combines the current
// field's name and generates environemnt variable names recursively
func (e *EnvironmentLoader) processField(prefix string, field *structs.Field) error {
	fieldName := e.generateFieldName(prefix, field)

	switch field.Kind() {
	case reflect.Struct:
		for _, f := range field.Fields() {
			if err := e.processField(fieldName, f); err != nil {
				return err
			}
		}
	default:
		v := os.Getenv(fieldName)
		_, ok := field.Value().(FieldSetter)

		if v == "" && !ok {
			return nil
		}
		if v == "" {
			v = fieldName
		}

		if err := FieldSet(field, v, nil); err != nil {
			return errors.Errorf("%v=%v: %v", fieldName, v, err)
		}
	}

	return nil
}

// PrintEnvs prints the generated environment variables to the std out.
func (e *EnvironmentLoader) PrintEnvs(s interface{}) {
	strct := structs.New(s)

	prefix := e.getPrefix(strct)

	for _, field := range strct.Fields() {
		e.printField(prefix, field)
	}
}

// printField prints the field of the config struct for the flag.Usage
func (e *EnvironmentLoader) printField(prefix string, field *structs.Field) {
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

type env struct {
	entries []string
}

// Return a []string with KEY=value pairs (e.g. for Docker)
func (e *EnvironmentLoader) GetEnvironment(s interface{}) []string {
	strct := structs.New(s)

	prefix := e.getPrefix(strct)
	env := &env{}

	for _, field := range strct.Fields() {
		e.appendFieldValue(env, prefix, field)
	}

	return env.entries
}

type EnvMarshaller interface {
	EnvMarshal(fieldName string) []string
}

func (e *EnvironmentLoader) appendFieldValue(env *env, prefix string, field *structs.Field) {
	fieldName := e.generateFieldName(prefix, field)

	if em, ok := field.Value().(EnvMarshaller); ok {
		for _, entry := range em.EnvMarshal(fieldName) {
			env.entries = append(env.entries, entry)
		}
		return
	}
	switch field.Kind() {
	case reflect.Struct:
		for _, f := range field.Fields() {
			e.appendFieldValue(env, fieldName, f)
		}
	default:
		if false {
			if field.IsZero() {
				return
			}
		}

		var stringValue string
		switch field.Kind() {
		case reflect.Slice:
			switch s := field.Value().(type) {
			case []string:
				stringValue = strings.Join(s, ",")
			case []int:
				var parts []string
				for _, v := range s {
					parts = append(parts, fmt.Sprintf("%v", v))
				}
				stringValue = strings.Join(parts, ",")
			}
		default:
			stringValue = fmt.Sprintf("%v", field.Value())
		}

		entry := fmt.Sprintf("%v=%v", fieldName, stringValue)
		env.entries = append(env.entries, entry)
	}
}

// generateFieldName generates the fiels name combined with the prefix and the
// struct's field name
func (e *EnvironmentLoader) generateFieldName(prefix string, field *structs.Field) string {
	fieldName := strings.ToUpper(field.Name())
	if e.CamelCase {
		fieldName = strings.ToUpper(strings.Join(camelcase.Split(field.Name()), "_"))
	}

	return strings.ToUpper(prefix) + "_" + fieldName
}
