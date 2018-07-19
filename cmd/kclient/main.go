package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/benthor/gocli"
	"github.com/ffhenkes/kripto/algo"
	"github.com/ffhenkes/kripto/fs"
)

const (
	path = "/data/authdb"
)

func main() {

	passphrase := os.Getenv("PHRASE")

	if "" == passphrase {
		fmt.Sprintf("Missing passphrase! Export <PHRASE> before continue!")
		os.Exit(0)
	}

	cli := gocli.MkCLI("Welcome Kripto CLI! Type help for valid commands.")

	// register help Option with cli.Help as callback
	cli.AddOption("help", "prints this help message", cli.Help)

	// register exit Option with cli.Exit as callback
	cli.AddOption("exit", "exits the CLI", cli.Exit)

	// register hidden quit Option with cli.Exit as callback. Should not appear in "help" list
	cli.AddOption("quit", "", cli.Exit)

	// demonstrate argument passing
	cli.AddOption("add", "Creates a valid user for Kripto! Type: username@password", func(args []string) string {
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

		user := input[0]
		passwd := string(algo.MakeSimpleHash(input[1]))
		user_string := fmt.Sprintf("%s@%s", user, passwd)

		symmetrical := algo.NewSymmetrical()
		data, err := symmetrical.Encrypt([]byte(user_string), passphrase)
		if err != nil {
			res = "Encryption error! Can`t continue!"
			return res
		}

		sys := fs.NewFileSystem(path)
		err = sys.MakeAuth(fmt.Sprintf(".%s", user), data)
		if err != nil {
			res = fmt.Sprintf("Error creating authentication for user: %s", user)
			return res
		}

		res = "\"" + user + "@***********\""
		return fmt.Sprintf("User added successfully %s", res)
	})

	cli.DefaultOption(func(args []string) string {
		return fmt.Sprintf("%s: command not found, type 'help' for help", args[0])
	})

	// run the main loop
	cli.Loop("<kripto>::@ ")

	fmt.Println("Good bye! Thank you for using Kripto!")

}
