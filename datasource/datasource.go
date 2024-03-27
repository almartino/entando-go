package datasource

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"strconv"
	"strings"
)

// Configuration is a normalized database configuration
// built by parsing the Entando injected datasource environment variables:
//
// https://developer.entando.com/next/docs/curate/bundle-details.html#microservices-environment-variables
type Configuration struct {
	Host     string
	Port     int
	Username string
	Password string
	// the database name
	Name string
}

type dbEnv struct {
	URL      string `envconfig:"DATASOURCE_URL" default:"jdbc:postgresql://localhost:5432/conferencems"`
	Username string `envconfig:"DATASOURCE_USERNAME" default:"conferencems"`
	Password string `envconfig:"DATASOURCE_PASSWORD" default:"root"`
}

// NewConfiguration parses the following env variables:
//   - SPRING_DATASOURCE_URL: the full jdbc url, e.g. jdbc:postgresql://db:5432/dbname
//   - SPRING_DATASOURCE_USERNAME
//   - SPRING_DATASOURCE_PASSWORD
//
// After parsing those variables, the datasource url is destructured in a Configuration alongside
// the username and password variables.
func NewConfiguration() Configuration {
	var c dbEnv
	envconfig.MustProcess("SPRING", &c)
	url := strings.SplitN(c.URL, "//", 2)
	if len(url) != 2 {
		panic("unexpected database url")
	}
	url = strings.SplitN(url[1], "/", 2)
	if len(url) != 2 {
		panic("unexpected db url len")
	}
	hostPort := strings.SplitN(url[0], ":", 2)
	if len(hostPort) != 2 {
		panic("unexpected host port len")
	}
	port, err := strconv.Atoi(hostPort[1])
	if err != nil {
		panic(fmt.Sprintf("could not parse database port: %v", err))
	}
	return Configuration{
		Host:     hostPort[0],
		Port:     port,
		Username: c.Username,
		Password: c.Password,
		Name:     url[1],
	}
}
