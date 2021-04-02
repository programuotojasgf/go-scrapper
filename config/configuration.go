package config

import (
	"fmt"
	"github.com/tkanos/gonfig"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

// Configuration struct for the scrapper
type Configuration struct {
	ConnectionString               string
	DatabaseName                   string
	ReviewsURLFirstPage            string
	WebsiteVisitorParallelismLimit int
}

var (
	// Config provides static access of configuration
	Config Configuration
)

func init() {
	err := gonfig.GetConf(getFileName(), &Config)
	if err != nil {
		fmt.Println(err)
		os.Exit(500)
	}
}

func getFileName() string {
	env := os.Getenv("ENV")
	if len(env) == 0 {
		env = "development"
	}
	filename := []string{"config.", env, ".json"}
	_, dirname, _, _ := runtime.Caller(0)
	filePath := path.Join(filepath.Dir(dirname), strings.Join(filename, ""))

	return filePath
}
