package confmanager

import (
	goflag "flag"
	"log"

	"github.com/caarlos0/env"
	flag "github.com/spf13/pflag"
)

func ParseCmdFlagsSecondary() {
	flag.CommandLine.AddGoFlagSet(goflag.CommandLine) // for future
	// Заново парсим аргументы командной строки, чтобы переопределить то, что было задано по дефолту и конфигом
	goflag.Parse()
}

func parseEnvVariables(cfg any) {
	err := env.Parse(cfg)
	if err != nil {
		log.Fatal(err)
	}
}
