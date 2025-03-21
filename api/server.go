package api

import (
	"fmt"
	"net/http"
)

func RunServer() {

	if !ApplicationProperties.ServerProperties.UseHttps && ApplicationProperties.ServerProperties.RedirectHttps {
		panic("https is disabled, cannot create redirect")
	}

	if ApplicationProperties.ServerProperties.UseHttps {
		go runHttps()
	}

	runHttp()
}

func runHttps() {
	address := fmt.Sprintf("%s:%d", ApplicationProperties.ServerProperties.Address, ApplicationProperties.ServerProperties.HttpsPort)

	err := http.ListenAndServeTLS(
		address,
		ApplicationProperties.ServerProperties.TlsCertificatePath,
		ApplicationProperties.ServerProperties.TlsKeyPath,
		nil,
	)

	if err != nil {
		panic(fmt.Sprintf("cannot start https server on %s, cause: %s", address, err.Error()))
	}
}

type RedirectHandler struct {
}

func (RedirectHandler) ServeHTTP(writter http.ResponseWriter, request *http.Request) {

	redirectTo := fmt.Sprintf(
		"https://%s:%d/%s",
		request.URL.Hostname(),
		ApplicationProperties.ServerProperties.HttpsPort,
		request.URL.RequestURI(),
	)

	http.Redirect(writter, request, redirectTo, http.StatusMovedPermanently)
}

func runHttp() {
	address := fmt.Sprintf("%s:%d", ApplicationProperties.ServerProperties.Address, ApplicationProperties.ServerProperties.HttpPort)

	var handler http.Handler

	if ApplicationProperties.ServerProperties.RedirectHttps {
		handler = &RedirectHandler{}
	}

	err := http.ListenAndServe(address, handler)
	if err != nil {
		panic(fmt.Sprintf("cannot start http server on %s, cause: %s", address, err.Error()))
	}

}
