package api

import (
	"fmt"
	"net/http"

	"github.com/lucas-10101/auth-service/api/conf"
	"github.com/lucas-10101/auth-service/api/utils"
)

func RunServer() {

	if !conf.ApplicationProperties.ServerProperties.UseHttps && conf.ApplicationProperties.ServerProperties.RedirectHttps {
		panic(utils.HTTPS_REDIRECT_CONFIGURATION_MISMATCH)
	}

	if conf.ApplicationProperties.ServerProperties.UseHttps {
		go runHttps()
	}

	runHttp()
}

func runHttps() {
	address := fmt.Sprintf("%s:%d", conf.ApplicationProperties.ServerProperties.Address, conf.ApplicationProperties.ServerProperties.HttpsPort)

	err := http.ListenAndServeTLS(
		address,
		conf.ApplicationProperties.ServerProperties.TlsCertificatePath,
		conf.ApplicationProperties.ServerProperties.TlsKeyPath,
		nil,
	)

	if err != nil {
		panic(utils.HTTPS_SERVER_START_FAILURE)
	}
}

type RedirectHandler struct {
}

func (RedirectHandler) ServeHTTP(writter http.ResponseWriter, request *http.Request) {

	redirectTo := fmt.Sprintf(
		"https://%s:%d/%s",
		request.URL.Hostname(),
		conf.ApplicationProperties.ServerProperties.HttpsPort,
		request.URL.RequestURI(),
	)

	http.Redirect(writter, request, redirectTo, http.StatusMovedPermanently)
}

func runHttp() {
	address := fmt.Sprintf("%s:%d", conf.ApplicationProperties.ServerProperties.Address, conf.ApplicationProperties.ServerProperties.HttpPort)

	var handler http.Handler

	if conf.ApplicationProperties.ServerProperties.RedirectHttps {
		handler = &RedirectHandler{}
	}

	err := http.ListenAndServe(address, handler)
	if err != nil {
		panic(utils.HTTP_SERVER_START_FAILURE)
	}

}
