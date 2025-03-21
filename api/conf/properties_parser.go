package conf

func GetPropertiesNameList() map[string]interface{} {
	return map[string]interface{}{

		// App Properties
		"application.name": "app-name",

		// Server Properties
		"application.server.address":        "127.0.0.1",
		"application.server.http-port":      80,
		"application.server.https-port":     443,
		"application.server.use-https":      false,
		"application.server.redirect-https": false,
		"application.server.tls-key-path":   "/path/to/key",
		"application.server.tls-cert-path":  "/path/to/cert",
	}
}

func LoadProperties() {

}

func WriteTemplate() {

}
