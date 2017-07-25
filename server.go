package main

import (
	"log"

	"flag"
	"fmt"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/spf13/viper"
	"github.com/team-tissue/labelizer-app/api/auth"
	"os"
	"path"
	"path/filepath"
)

const (
	CONFIG_PORT = "port"
)

func init() {
	initConfig()
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	cfgFilePath := *flag.String("config", ".labelizer.yml", "config file path")
	fmt.Println(cfgFilePath)
	flag.Parse()

	dir, _ := os.Getwd()

	if !filepath.IsAbs(cfgFilePath) {
		cfgFilePath = path.Join(dir, cfgFilePath)
	}

	if cfgFilePath == "" {
		// Use config file from the flag.
		cfgFilePath = path.Join(dir, ".labelizer.yml")
	}

	viper.SetConfigFile(cfgFilePath)
	viper.SetConfigType("yml")
	viper.AutomaticEnv() // read in environment variables that match
	viper.SetDefault(CONFIG_PORT, "8000")

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		fmt.Println(err)
	}
}

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Static("/assets/js", "app/dist")

	r := e.Group("/api")
	r.Use(middleware.JWT([]byte(auth.CONFIG_JWT_SECRET)))

	// Serve incoming routes from spa
	for _, route := range []string{"setup", "/"} {
		e.File(route, "public/index.html")
	}

	e.GET("/auth/github/signin", auth.Challenge)
	e.GET("/auth/github/code", auth.FetchCode)

	log.Fatal(e.Start(":" + viper.Get(CONFIG_PORT).(string)))
}
