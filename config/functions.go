package config

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"time"
)

func getEnv(key, defaultVal string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultVal
}

func getEnvBool(key string, defaultVal bool) bool {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.ParseBool(value); err == nil {
			return parsed
		}
	}
	return defaultVal
}

func getEnvDuration(key string, defaultVal time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if parsed, err := time.ParseDuration(value); err == nil {
			return parsed
		}
	}
	return defaultVal
}

func getEnvInt64(key string, defaultVal int64) int64 {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.ParseInt(value, 10, 64); err == nil {
			return parsed
		}
	}
	return defaultVal
}

func getEnvFloat64(key string, defaultVal float64) float64 {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.ParseFloat(value, 64); err == nil {
			return parsed
		}
	}
	return defaultVal
}

func setFieldFromEnv(field reflect.Value, envKey, defaultVal string) {
	switch field.Kind() {
	case reflect.String:
		field.SetString(getEnv(envKey, defaultVal))
	case reflect.Bool:
		defaultBool, _ := strconv.ParseBool(defaultVal)
		field.SetBool(getEnvBool(envKey, defaultBool))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		defaultInt, _ := strconv.ParseInt(defaultVal, 10, 64)
		field.SetInt(getEnvInt64(envKey, defaultInt))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		defaultUint, _ := strconv.ParseUint(defaultVal, 10, 64)
		setUintField(field, envKey, defaultUint)
	case reflect.Float32, reflect.Float64:
		defaultFloat, _ := strconv.ParseFloat(defaultVal, 64)
		field.SetFloat(getEnvFloat64(envKey, defaultFloat))
	default:
		setDurationField(field, envKey, defaultVal)
	}
}

func setUintField(field reflect.Value, envKey string, defaultVal uint64) {
	if value := os.Getenv(envKey); value != "" {
		if parsed, err := strconv.ParseUint(value, 10, 64); err == nil {
			field.SetUint(parsed)
			return
		}
	}
	field.SetUint(defaultVal)
}

func setDurationField(field reflect.Value, envKey, defaultVal string) {
	if field.Type() == reflect.TypeOf(time.Duration(0)) {
		defaultDuration, _ := time.ParseDuration(defaultVal)
		field.Set(reflect.ValueOf(getEnvDuration(envKey, defaultDuration)))
	}
}

func setFieldDefault(field reflect.Value, defaultVal string) {
	switch field.Kind() {
	case reflect.String:
		field.SetString(defaultVal)
	case reflect.Bool:
		if defaultBool, err := strconv.ParseBool(defaultVal); err == nil {
			field.SetBool(defaultBool)
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if defaultInt, err := strconv.ParseInt(defaultVal, 10, 64); err == nil {
			field.SetInt(defaultInt)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if defaultUint, err := strconv.ParseUint(defaultVal, 10, 64); err == nil {
			field.SetUint(defaultUint)
		}
	case reflect.Float32, reflect.Float64:
		if defaultFloat, err := strconv.ParseFloat(defaultVal, 64); err == nil {
			field.SetFloat(defaultFloat)
		}
	default:
		if field.Type() == reflect.TypeOf(time.Duration(0)) {
			if defaultDuration, err := time.ParseDuration(defaultVal); err == nil {
				field.Set(reflect.ValueOf(defaultDuration))
			}
		}
	}
}

func validateConfigInput(config any) (reflect.Value, reflect.Type, error) {
	v := reflect.ValueOf(config)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return reflect.Value{}, nil, fmt.Errorf("config must be a pointer to struct")
	}
	elem := v.Elem()
	return elem, elem.Type(), nil
}

func Parse(config any) error {
	elem, t, err := validateConfigInput(config)
	if err != nil {
		return err
	}

	for i := range elem.NumField() {
		field := elem.Field(i)
		fieldType := t.Field(i)

		if !field.CanSet() {
			continue
		}

		envKey := fieldType.Tag.Get("env")
		defaultVal := fieldType.Tag.Get("default")

		if envKey == "" {
			continue
		}

		setFieldFromEnv(field, envKey, defaultVal)
	}

	return nil
}

func Defaults[T any](config *T) *T {
	v := reflect.ValueOf(config)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return config
	}

	elem := v.Elem()
	t := elem.Type()
	newStruct := reflect.New(t)
	newElem := newStruct.Elem()

	for i := range elem.NumField() {
		field := newElem.Field(i)
		fieldType := t.Field(i)

		if !field.CanSet() {
			continue
		}

		defaultVal := fieldType.Tag.Get("default")
		if defaultVal == "" {
			continue
		}

		setFieldDefault(field, defaultVal)
	}

	return newStruct.Interface().(*T)
}
