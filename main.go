package main

import (
	"os"

	"github.com/Jeiwan/opscript/cmd"
	"github.com/sirupsen/logrus"
)

func main() {
	if os.Getenv("DEBUG") != "" {
		logrus.SetLevel(logrus.DebugLevel)
	}

	if err := cmd.New().Execute(); err != nil {
		logrus.Fatalln(err)
	}
}
