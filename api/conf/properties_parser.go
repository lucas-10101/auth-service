package conf

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/lucas-10101/auth-service/api/utils"
)

func getPropertiesList() map[string]interface{} {

	return map[string]interface{}{

		// App Properties
		"ApplicationProperties.name": "app-name",

		// Server Properties
		"ApplicationProperties.ServerProperties.Address":            "127.0.0.1",
		"ApplicationProperties.ServerProperties.UseHttps":           false,
		"ApplicationProperties.ServerProperties.RedirectHttps":      false,
		"ApplicationProperties.ServerProperties.HttpPort":           80,
		"ApplicationProperties.ServerProperties.HttpsPort":          443,
		"ApplicationProperties.ServerProperties.TlsKeyPath":         "/path/to/key",
		"ApplicationProperties.ServerProperties.TlsCertificatePath": "/path/to/cert",

		// Logger Properties
		"ApplicationProperties.LoggerProperties.LogLevel":                int(slog.LevelInfo),
		"ApplicationProperties.LoggerProperties.LogDriver":               "file",
		"ApplicationProperties.LoggerProperties.AllowFallback":           true,
		"ApplicationProperties.LoggerProperties.FileName":                "application.log",
		"ApplicationProperties.LoggerProperties.MongoDbDatabaseName":     nil,
		"ApplicationProperties.LoggerProperties.MongoDbCollectionName":   nil,
		"ApplicationProperties.LoggerProperties.MongoDbConnectionString": nil,
		"ApplicationProperties.LoggerProperties.MongoDbMaxPoolSize":      nil,
		"ApplicationProperties.LoggerProperties.MongoDbMinPoolSize":      nil,
	}
}

func LoadProperties() {
	loadDinamicFrom(readPropertiesFile(), ApplicationProperties, "ApplicationProperties")
}

// cannot handle primitive types as pointers
//
// cannot handle sub structures as pointers
//
// make to deal with properties struct
//
// BFG-9000 property loader
func loadDinamicFrom(properties map[string]string, into any, withBasePath string) {
	if reflect.TypeOf(into).Kind() != reflect.Pointer && reflect.TypeOf(into).Elem().Kind() != reflect.Struct {
		panic("struct pointer required")
	}

	for _, visibleField := range reflect.VisibleFields(reflect.TypeOf(into).Elem()) {

		if visibleField.Type.Kind() == reflect.Struct {
			field := reflect.ValueOf(into).Elem().FieldByName(visibleField.Name)

			if field.CanSet() {
				loadDinamicFrom(properties, field.Addr().Interface(), fmt.Sprintf("%s.%s", withBasePath, visibleField.Name))
				continue
			}
		}

		if visibleField.Type.Kind() == reflect.Struct || visibleField.Type.Kind() == reflect.Pointer {
			continue
		}

		field := reflect.ValueOf(into).Elem().FieldByName(visibleField.Name)
		propertyValue, exists := properties[fmt.Sprintf("%s.%s", withBasePath, visibleField.Name)]
		if field.CanAddr() && exists {
			switch field.Kind() {
			case reflect.String:
				field.SetString(propertyValue)
			case reflect.Int:
				if intValue, parseErr := strconv.ParseInt(propertyValue, 10, 0); parseErr == nil {
					field.SetInt(intValue)
				}
			case reflect.Bool:
				if boolValue, parseErr := strconv.ParseBool(propertyValue); parseErr == nil {
					field.SetBool(boolValue)
				}
			default:
				continue
			}
		}
	}
}

func readPropertiesFile() map[string]string {

	file, err := os.OpenFile("application.properties", (os.O_RDONLY), 0644)

	if err != nil {
		panic(utils.PROPERTIES_FILE_READ_ERROR.WithErrorMessage(err))
	}

	defer file.Close()

	properties := map[string]string{}
	for lnNumber, scanner := 1, bufio.NewScanner(file); scanner.Scan(); lnNumber++ {
		line := scanner.Text()
		if pos := strings.Index(line, "#"); pos != -1 {
			line = line[0:pos]
		}

		if len(strings.TrimSpace(line)) == 0 {
			continue
		}

		parts := strings.Split(line, "=")

		if len(parts) != 2 {
			panic(fmt.Sprintf("%s, at line: %d", utils.PROPERTIES_ENTRY_BAD_FORMAT.ToString(), lnNumber))
		}

		properties[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
	}

	return properties
}

func WriteTemplate() {

	file, err := os.OpenFile("application.properties", (os.O_WRONLY | os.O_CREATE), 0644)

	if err != nil {
		panic(utils.PROPERTIES_FILE_WRITE_ERROR.WithErrorMessage(err))
	}

	writter := bufio.NewWriter(file)

	defer file.Close()
	defer writter.Flush()

	for propertyName, defaultValue := range getPropertiesList() {
		writter.WriteString(fmt.Sprintf("%s=%v\n", propertyName, defaultValue))
	}

}
