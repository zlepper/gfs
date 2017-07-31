package internal

import (
	"flag"
	"github.com/zlepper/gfs"
	"log"
	"os"
)

func main() {
	configPath := flag.String("config", gfs.DefaultConfigPath, "The path to the config file.")
	persist := flag.Bool("persist", false, "Overwrite config file with options given as arguments.")
	username := flag.String("username", "", "The username of the user. Overrules whatever is in the config file.")
	password := flag.String("password", "", "The password of the user. Overrules whatever is in the config file.")
	port := flag.String("port", "", "The port to serve on. Overrules whatever is in the config file.")
	loginRequiredForRead := flag.Bool("loginRequiredForRead", false, "Enable to require login for being able to get directory listings, and downloading files.")
	serve := flag.String("serve", gfs.DefaultServePath, "The path that should be served by gfs.")

	flag.Parse()

	configs, err := gfs.GetConfigs(*configPath)
	if err != nil {
		log.Fatalln(err)
	}

	if *username != "" {
		configs.Username = *username
	}

	if *password != "" {
		passwordHash, err := gfs.CreatePassword(*password)
		if err != nil {
			log.Fatalln(err)
		}

		configs.Password = passwordHash
	}

	if *loginRequiredForRead {
		configs.LoginRequiredForRead = *loginRequiredForRead
	}

	if *port != "" {
		configs.Port = *port
	}

	if *serve != gfs.DefaultServePath {
		configs.Serve = *serve
	}

	if *persist {
		err := gfs.SaveConfigs(*configPath, configs)
		if err != nil {
			log.Panicln(err)
		}
		log.Println("Updated configs. Start gfs without the `persist` flag to start the program.")
		return
	}

	log.Println(*configs)

	os.MkdirAll(configs.Serve, os.ModePerm)

	gfs.RunServer(configs)
}
