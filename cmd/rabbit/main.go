package main

import (
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/jmshal/rabbit/rabbit"
	cli "gopkg.in/urfave/cli.v1"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	app := cli.NewApp()
	app.Name = "rabbit"
	app.Usage = ""
	app.Version = rabbit.Version
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "config, c",
			Usage:  "Load configuration from `FILE`",
			Value:  "./config.json",
			EnvVar: "RABBIT_CONFIG_FILE",
		},
	}
	app.Commands = []cli.Command{
		cli.Command{
			Name:        "version",
			Usage:       "",
			Description: "Prints the version",
			Action: func(c *cli.Context) error {
				fmt.Println(rabbit.Version)
				return nil
			},
		},
	}
	app.Action = func(c *cli.Context) error {
		configPath := c.String("config")
		config, err := rabbit.LoadConfigFile(configPath)
		if err != nil {
			log.Fatalln(err)
		}
		app := rabbit.NewRabbit(config)
		log.Fatalln(app.Listen())
		return nil
	}
	app.Run(os.Args)
}
