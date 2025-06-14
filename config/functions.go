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

func Parse(config interface{}) error {
	v := reflect.ValueOf(config)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("config must be a pointer to struct")
	}

	v = v.Elem()
	t := v.Type()

	for i := range v.NumField() {
		field := v.Field(i)
		fieldType := t.Field(i)

		if !field.CanSet() {
			continue
		}

		envKey := fieldType.Tag.Get("env")
		defaultVal := fieldType.Tag.Get("default")

		if envKey == "" {
			continue
		}

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
			if value := os.Getenv(envKey); value != "" {
				if parsed, err := strconv.ParseUint(value, 10, 64); err == nil {
					field.SetUint(parsed)
					continue
				}
			}
			field.SetUint(defaultUint)

		case reflect.Float32, reflect.Float64:
			defaultFloat, _ := strconv.ParseFloat(defaultVal, 64)
			field.SetFloat(getEnvFloat64(envKey, defaultFloat))

		default:
			if field.Type() == reflect.TypeOf(time.Duration(0)) {
				defaultDuration, _ := time.ParseDuration(defaultVal)
				field.Set(reflect.ValueOf(getEnvDuration(envKey, defaultDuration)))
			} else {
				return fmt.Errorf("unsupported field type: %s", field.Kind())
			}
		}
	}

	return nil
}
