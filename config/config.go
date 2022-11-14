package config

import (
	"fmt"
	"net/url"
	"os"
	"reflect"
	"strconv"
)

type Config struct {
	S3Endpoint        string `kind:"url" example:"http://localhost:8084" env:"SS3B_S3_ENDPOINT"` //has trailing slash
	FileSize          int    `kind:"int" example:"90000000" env:"SS3B_FILE_SIZE"`                //in bytes
	DatasetSize       int    `kind:"int" example:"50000000000" env:"SS3B_DATASET_SIZE"`          //in bytes
	MovingPopularSize int    `kind:"int" example:"9000000000" env:"SS3B_MOVING_POPULAR_SIZE"`    //in bytes
	TestDuration      int    `kind:"int" example:"180" env:"SS3B_TEST_DURATION"`                 //in seconds
	RequestPerSecond  int    `kind:"int" example:"50" env:"SS3B_REQUEST_PER_SECOND"`             //in seconds
}

func GetValues() Config {
	c := Config{}
	fields := reflect.VisibleFields(reflect.TypeOf(struct{ Config }{}))

	for _, field := range fields {
		switch field.Tag.Get("kind") {
		case "url":
			urlCleaned := checkAndCleanURL(getEnvValueString(field.Tag.Get("env")))

			if urlCleaned == "" {
				urlCleaned = field.Tag.Get("example")
			}

			reflect.ValueOf(&c).Elem().FieldByName(field.Name).SetString(urlCleaned)
		case "string":
			value := getEnvValueString(field.Tag.Get("env"))
			if value == "" {
				value = field.Tag.Get("example")
			}

			reflect.ValueOf(&c).Elem().FieldByName(field.Name).SetString(value)
		case "int":
			intValue := getEnvValueInt(field.Tag.Get("env"))
			if intValue == 0 {
				var err error
				intValue, err = strconv.Atoi(field.Tag.Get("example"))
				if err != nil {
					panic(err)
				}
			}
			reflect.ValueOf(&c).Elem().FieldByName(field.Name).SetInt(int64(intValue))
		}
	}

	return c
}

func checkAndCleanURL(urlDirty string) string {
	urlCleaned := urlDirty

	if urlCleaned == "" {
		return urlCleaned
	}

	// check if url is valid
	_, err := url.ParseRequestURI(urlCleaned)
	if err != nil {
		panic(err)
	}

	// remove trailing slash if present
	if urlCleaned[len(urlCleaned)-1:] == "/" {
		urlCleaned = urlCleaned[:len(urlCleaned)-1]
	}

	// check if url has protocol
	if urlCleaned[0:4] != "http" {
		urlCleaned = "https://" + urlCleaned
	}

	return urlCleaned
}

func getEnvValueString(env string) string {
	fmt.Println("---!>", os.Getenv(env), env)
	return os.Getenv(env)
}
func getEnvValueInt(env string) int {
	foo := os.Getenv(env)
	fmt.Println(foo, "--", env)
	if foo == "" {
		return 0
	}

	//  string to int
	bar, err := strconv.Atoi(foo)
	if err != nil {
		panic(err)
	}

	return bar
}

func GetCompleteURL(c Config, path string) string {
	return c.S3Endpoint + path
}
