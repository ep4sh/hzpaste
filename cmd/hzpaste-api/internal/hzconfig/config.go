package hzconfig

import (
	"log"
	"os"
	"path"
	"strings"

	"github.com/tkanos/gonfig"
)

type Configuration struct {
	Host string `env:"HZPASTE_HOST"`
	Port string `env:"HZPASTE_PORT"`
}

func GetConfigFileName() string {
	env := os.Getenv("ENV")
	if len(env) == 0 {
		env = "development"
	}
	filename := []string{"config.", env, ".json"}
	filePath := path.Join(strings.Join(filename, ""))

	return filePath
}

func GetPortEnv() (string, string) {
	envVarName := "HZPASTE_PORT"
	// Specific environment variable override for heroku deployment only
	_, herokuMode := os.LookupEnv("HEROKU_DEPLOY")
	if herokuMode {
		envVarName = "PORT"
	}
	//
	return os.Getenv(envVarName), envVarName
}

func GetHostEnv() (string, string) {
	envVarName := "HZPASTE_HOST"
	return os.Getenv(envVarName), envVarName
}

func InitConfig() Configuration {
	PortEnvValue, PortEnvName := GetPortEnv()
	HostEnvValue, HostEnvName := GetHostEnv()
	configuration := Configuration{
		Port: PortEnvValue,
		Host: HostEnvValue,
	}
	if len(configuration.Port) == 0 || len(configuration.Host) == 0 {
		err := gonfig.GetConf(GetConfigFileName(), &configuration)
		if err != nil {
			log.Println("Cannot initialize configuration!")
			log.Println("Either provide configuration file", GetConfigFileName())
			log.Println("or set", HostEnvName, "and", PortEnvName, "environment variables")
			log.Fatal(err)
		}
	}

	return configuration
}
