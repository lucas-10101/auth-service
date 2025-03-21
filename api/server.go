package api

import (
	"fmt"
	"net/http"
)

func Run() {
	http.ListenAndServeTLS(
		fmt.Sprintf("%s:%d", ApplicationProperties.ServerProperties.Address, ApplicationProperties.ServerProperties.Port),
		ApplicationProperties.ServerProperties.TlsCertificatePath,
		ApplicationProperties.ServerProperties.TlsKeyPath,
		nil,
	)
}
