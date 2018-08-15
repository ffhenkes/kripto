package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/NeowayLabs/logger"
	"github.com/benthor/gocli"
	"github.com/ffhenkes/kripto/auth"
	"github.com/ffhenkes/kripto/model"
)

var logK = logger.Namespace("kripto.cli")

func main() {

	passphrase := os.Getenv("PHRASE")

	if "" == passphrase {
		logK.Fatal("Missing passphrase! Export <PHRASE> before continue!")
	}

	cli := gocli.MkCLI("Welcome to Kripto CLI! Type help for valid commands.")

	err := cli.AddOption("help", "prints this help message", cli.Help)
	if err != nil {
		logK.Fatal("Critical failure!")
	}

	err = cli.AddOption("exit", "exits the CLI", cli.Exit)
	if err != nil {
		logK.Fatal("Critical failure!")
	}

	err = cli.AddOption("quit", "", cli.Exit)
	if err != nil {
		logK.Fatal("Critical failure!")
	}

	err = cli.AddOption("add", "Creates a valid user for Kripto! Type: username@password", func(args []string) string {
		res := ""

		size := len(args)

		if size == 0 {
			res = "Missing value! <username@password>"
			return res
		}

		if size > 1 {
			res = "Too many arguments! Use: username@password"
			return res
		}

		input := strings.Split(args[0], "@")

		if len(input) < 2 {
			res = "Bad input format! Use: username@password"
			return res
		}

		password := normalizePassword(input)

		if "" == password {
			res = "Password must not be empty!"
			return res
		}

		c := model.Credentials{
			Username: input[0],
			Password: password,
		}

		login := auth.NewLogin(&c)
		ok := login.AddCredentials(passphrase)
		if ok != nil {
			res = "Error adding new credentials!!"
			return res
		}

		res = "\"" + c.Username + "@***********\""
		return fmt.Sprintf("User added successfully %s", res)
	})
	if err != nil {
		logK.Fatal("Critical failure!")
	}

	cli.DefaultOption(func(args []string) string {
		return fmt.Sprintf("%s: command not found, type 'help' for help", args[0])
	})

	cli.Loop("<kripto>::@ ")

	fmt.Println("Good bye! Thank you for using Kripto!")

}

func normalizePassword(a []string) string {
	normal := ""
	for k, v := range a {

		if k == 1 {
			normal = v
		}

		if k > 1 {
			normal += fmt.Sprintf("@%s", v)
		}

	}
	return normal
}
