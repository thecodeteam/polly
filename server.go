// package main
//
// import (
// 	"bytes"
// 	"fmt"
// 	"io"
//
// 	"github.com/akutz/gofig"
// 	"github.com/emccode/libstorage"
// )
//
// var config gofig.Config
//
// const (
// 	/*
// 	   libStorageConfigBase is the base config for tests
// 	   01 - the host address to server and which the client uses
// 	   02 - the executors directory
// 	   03 - the client TLS section. use an empty string if TLS is disabled
// 	   04 - the server TLS section. use an empty string if TLS is disabled
// 	   05 - the services
// 	*/
// 	libStorageConfigBase = `
// libstorage:
//   host: %[1]s
//   profiles:
//     enabled: true
//     groups:
//     - local=127.0.0.1%[3]s
//   server:
//     endpoints:
//       localhost:
//         address: %[1]s%[4]s
//     services:%[5]s
// `
//
// 	libStorageConfigService = `
//       %[1]s:
//         libstorage:
//           driver: %[2]s
// `
// )
//
// func startServer() (io.Closer, <-chan error) {
// 	host := "tcp://127.0.0.1:7979"
// 	driverName := "mock"
// 	serviceName := "mock"
//
// 	config = gofig.New()
//
// 	services := &bytes.Buffer{}
// 	services.WriteString(
// 		fmt.Sprintf(libStorageConfigService, serviceName, driverName))
//
// 	configYaml := fmt.Sprintf(
// 		libStorageConfigBase,
// 		host, "/tmp/libstorage/executors",
// 		"", "", //clientTLS, serverTLS
// 		services.String())
//
// 	//fmt.Print(configYaml)
//
// 	configYamlBuf := []byte(configYaml)
// 	if err := config.ReadConfig(bytes.NewReader(configYamlBuf)); err != nil {
// 		panic(err)
// 	}
//
// 	return libstorage.Serve(config)
// }

package main

import (
	"fmt"
	"os"

	"github.com/akutz/gofig"

	"github.com/emccode/libstorage/cli/servers"
)

const (
	/*
	   libStorageConfigBase is the base config for tests

	   01 - the host address to server and which the client uses
	   02 - the executors directory
	   03 - the client TLS section. use an empty string if TLS is disabled
	   04 - the server TLS section. use an empty string if TLS is disabled
	   05 - the services
	*/
	libStorageConfigBase = `
libstorage:
  host: %[1]s
  profiles:
    enabled: true
    groups:
    - local=127.0.0.1%[3]s
  server:
    endpoints:
      localhost:
        address: %[1]s%[4]s
    services:%[5]s
`

	libStorageConfigService = `
      %[1]s:
        libstorage:
          driver: %[2]s
`
	libStorageConfigClientTLS = `
    tls:
      serverName: libstorage-server
      certFile: %s
      keyFile: %s
      trustedCertsFile: %s
`

	libStorageConfigServerTLS = `
        tls:
          serverName: libstorage-server
          certFile: %s
          keyFile: %s
          trustedCertsFile: %s
          clientCertRequired: true
`
)

var (
	tlsPath = fmt.Sprintf(
		"%s/src/github.com/emccode/libstorage/.tls", os.Getenv("GOPATH"))
	serverCrt    = fmt.Sprintf("%s/libstorage-server.crt", tlsPath)
	serverKey    = fmt.Sprintf("%s/libstorage-server.key", tlsPath)
	clientCrt    = fmt.Sprintf("%s/libstorage-client.crt", tlsPath)
	clientKey    = fmt.Sprintf("%s/libstorage-client.key", tlsPath)
	trustedCerts = fmt.Sprintf("%s/libstorage-ca.crt", tlsPath)
)

var config gofig.Config

func startServer() (gofig.Config, error) {
	cfg, _, errs := servers.Run("", false, "mock", "mock")
	go func() {
		err := <-errs
		panic(err)
	}()

	return cfg, nil
}
