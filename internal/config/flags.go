package config

import (
	"flag"
)

type Flags struct {
	FlagRunAddr           string
	FlagRunShorternerAddr string
	FlagLog               bool
}

func ParseFlags(args []string) Flags {
	fs := flag.NewFlagSet("shortner", flag.ContinueOnError)

	flags := Flags{
		FlagRunAddr:           "",
		FlagRunShorternerAddr: "",
		FlagLog:               false,
	}

	fs.StringVar(&flags.FlagRunAddr, "a", "localhost:8080", "Адрес запуска http-сервера")

	fs.StringVar(&flags.FlagRunShorternerAddr, "b", "http://localhost:8080", "Базовый адрес результирующего сокращённого URL ")

	fs.BoolVar(&flags.FlagLog, "l", false, "Вывод логов в консоль")

	fs.Parse(args)

	return flags
}
