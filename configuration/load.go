package configuration

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/jinzhu/configor"
	"gopkg.in/yaml.v2"
)

func getPrefixForStruct(prefixes []string, fieldStruct *reflect.StructField) []string {
	if fieldStruct.Anonymous && fieldStruct.Tag.Get("anonymous") == "true" {
		return prefixes
	}
	return append(prefixes, fieldStruct.Name)
}

func addFlagByField(field reflect.Value, flagName string, defaultValue string, usageValue string) {
	switch field.Kind() {
	case reflect.Bool:
		value, _ := strconv.ParseBool(defaultValue)
		flag.Bool(flagName, value, usageValue)
	case reflect.String:
		flag.String(flagName, defaultValue, usageValue)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		value, _ := strconv.ParseUint(defaultValue, 10, int(field.Type().Size()))
		flag.Uint(flagName, uint(value), usageValue)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		value, _ := strconv.ParseInt(defaultValue, 10, int(field.Type().Size()))
		flag.Int(flagName, int(value), usageValue)
	case reflect.Float32, reflect.Float64:
		value, _ := strconv.ParseFloat(defaultValue, int(field.Type().Size()))
		flag.Float64(flagName, float64(value), usageValue)
	}
}

func isValidField(field reflect.Value) bool {
	return field.CanAddr() && field.CanInterface()
}

func isSimpleType(field reflect.Value) bool {
	return (field.Kind() != reflect.Slice) && (field.Kind() != reflect.Struct)
}

func getPathInStruct(fieldStruct reflect.StructField, prefixes ...string) string {
	return strings.Join(append(prefixes, fieldStruct.Name), ".")
}

func getFlagName(fieldStruct reflect.StructField, prefixes ...string) string {
	cliName := fieldStruct.Tag.Get("cli")
	if cliName == "" {
		structName := AddDelimiter(fieldStruct.Name)
		cliName = strings.ToLower(strings.Join(append(prefixes, structName), "-")) // DB_NAME
	}
	return cliName
}

func getFlagUsage(fieldStruct reflect.StructField) string {
	usageValue := fieldStruct.Tag.Get("usage")
	if required := fieldStruct.Tag.Get("required"); required != "" {
		return usageValue + " (required)"
	}
	return usageValue
}

func updateFieldValue(field reflect.Value, value string) error {
	return yaml.Unmarshal([]byte(value), field.Addr().Interface())
}

func setDefaultValue(field reflect.Value, fieldStruct reflect.StructField, prefixes ...string) error {
	if value := fieldStruct.Tag.Get("default"); value != "" {
		if err := updateFieldValue(field, value); err != nil {
			return err
		}
	} else if fieldStruct.Tag.Get("required") == "true" {
		// return error if it is required but blank
		return fmt.Errorf("%v (flag %v) is required, but blank", getPathInStruct(fieldStruct, prefixes...), getFlagName(fieldStruct, prefixes...))
	}
	return nil
}

func isBlankField(field reflect.Value) bool {
	return reflect.DeepEqual(field.Interface(), reflect.Zero(field.Type()).Interface())
}

func exportConfigToFlags(config interface{}, prefixes ...string) error {
	configValue := reflect.Indirect(reflect.ValueOf(config))
	if configValue.Kind() != reflect.Struct {
		return errors.New("invalid config, should be struct")
	}

	configType := configValue.Type()
	for i := 0; i < configType.NumField(); i++ {
		var (
			fieldStruct  = configType.Field(i)
			field        = configValue.Field(i)
			flagName     = getFlagName(fieldStruct, prefixes...) // read configuration from shell cli
			defaultValue = fieldStruct.Tag.Get("default")
			usageValue   = getFlagUsage(fieldStruct)
		)

		if !isValidField(field) {
			continue
		}

		for field.Kind() == reflect.Ptr {
			field = field.Elem()
		}

		addFlagByField(field, flagName, defaultValue, usageValue)

		if field.Kind() == reflect.Struct {
			if err := exportConfigToFlags(field.Addr().Interface(), getPrefixForStruct(prefixes, &fieldStruct)...); err != nil {
				return err
			}
		}
	}
	return nil
}

func setDefaultValuesToStruct(config interface{}, prefixes ...string) error {
	configValue := reflect.Indirect(reflect.ValueOf(config))
	if configValue.Kind() != reflect.Struct {
		return errors.New("invalid config, should be struct")
	}
	configType := configValue.Type()

	for i := 0; i < configType.NumField(); i++ {
		var (
			fieldStruct = configType.Field(i)
			field       = configValue.Field(i)
		)

		if !isValidField(field) {
			continue
		}

		for field.Kind() == reflect.Ptr {
			field = field.Elem()
		}
		setDefaultValue(field, fieldStruct, prefixes...)

		if field.Kind() == reflect.Struct {
			if err := setDefaultValuesToStruct(field.Addr().Interface(), getPrefixForStruct(prefixes, &fieldStruct)...); err != nil {
				return err
			}
		}
	}
	return nil
}

func isChangedFlag(flag *flag.Flag) bool {
	return flag.Value.String() != flag.DefValue
}

func exportChangesFlagsToConfig(config interface{}, prefixes ...string) error {
	configValue := reflect.Indirect(reflect.ValueOf(config))
	if configValue.Kind() != reflect.Struct {
		return errors.New("invalid config, should be struct")
	}
	configType := configValue.Type()

	for i := 0; i < configType.NumField(); i++ {
		var (
			fieldStruct = configType.Field(i)
			field       = configValue.Field(i)
			flagName    = getFlagName(fieldStruct, prefixes...)
		)

		if !isValidField(field) {
			continue
		}

		for field.Kind() == reflect.Ptr {
			field = field.Elem()
		}

		if isSimpleType(field) {
			cliFlag := flag.Lookup(flagName)
			if isChangedFlag(cliFlag) {
				if err := updateFieldValue(field, cliFlag.Value.String()); err != nil {
					return err
				}
			}
		}

		if isBlankField(field) {
			if err := setDefaultValue(field, fieldStruct, prefixes...); err != nil {
				return err
			}
		}

		if field.Kind() == reflect.Struct {
			if err := exportChangesFlagsToConfig(field.Addr().Interface(), getPrefixForStruct(prefixes, &fieldStruct)...); err != nil {
				return err
			}
		}
	}
	return nil
}

func setCommonFlags(defaultConfigFileName string) error {
	flag.Bool("export-config", false, "create default config file")
	flag.String("config-file", defaultConfigFileName, "set configuration file")
	return nil
}

func getConfigFileName(defaultFileName string) string {
	configFileName := defaultFileName
	cliFlag := flag.Lookup("config-file")
	if isChangedFlag(cliFlag) {
		configFileName = cliFlag.Value.String()
	}
	return configFileName
}

func exportDefaultConfigFile(config interface{}, configFileName string) error {
	setDefaultValuesToStruct(config)
	var buffer bytes.Buffer
	if err := toml.NewEncoder(&buffer).Encode(config); err != nil {
		return err
	}
	if err := ioutil.WriteFile(configFileName, buffer.Bytes(), 0644); err != nil {
		return err
	}
	fmt.Printf("default configuration saved to %s", configFileName)
	return nil
}

func checkCommonFlags(config interface{}, configFileName string) error {
	if cliFlag := flag.Lookup("export-config"); cliFlag.Value.String() == "true" {
		if err := exportDefaultConfigFile(config, configFileName); err != nil {
			fmt.Println(err)
		}
		os.Exit(0)
	}
	return nil
}

func normaliseFilePath(file string) string {
	binaryHomepath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return file
	}
	return filepath.Join(binaryHomepath, file)
}

func Load(config interface{}, file string) error {
	if err := exportConfigToFlags(config); err != nil {
		return err
	}

	configFilePath := normaliseFilePath(file)

	setCommonFlags(configFilePath)

	flag.Parse()
	if err := checkCommonFlags(config, getConfigFileName(configFilePath)); err != nil {
		return err
	}

	configor.Load(config, getConfigFileName(configFilePath))

	if err := exportChangesFlagsToConfig(config); err != nil {
		return err
	}

	return nil
}
