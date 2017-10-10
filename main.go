package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"strings"

	"golang.org/x/crypto/ssh/terminal"
)

/* TODO:
 * readme.md
 * fancy colorized output based on how the balances compare to previous day/week/interval
 * use flag package to parse arguments
 * list all accounts
 * list details of one specific account
 * tests
 * gnuplot graphs?
 * use log package correctly
 */

func main() {
	bmc := new(BettermentAPIClient)
	config := GetConfig()
	bmc.Email = config.Email
	bmc.Password = config.Password
	bmc.Summary()
}

// Config ..
type Config struct {
	Email    string
	Password string
}

// GetConfig from ~/.bmc/bmc.json, or create it if it doesn't exist.
func GetConfig() Config {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}

	// Make directory, even if it already exists
	path := usr.HomeDir + "/.bmc"
	err = os.MkdirAll(path, 0755)
	if err != nil {
		panic(err)
	}

	// If file doesn't exist, prompt and create it
	fullpath := path + "/bmc.json"
	if _, err := os.Stat(fullpath); os.IsNotExist(err) {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Email address: ")
		email, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		// There is a \n by default
		email = strings.Trim(email, " \r\n")

		fmt.Print("Password: ")
		password, err := terminal.ReadPassword(0)
		if err != nil {
			panic(err)
		}

		fmt.Printf("\nWriting new config file to '%s'.\n", fullpath)
		conf := Config{email, string(password)}
		buf, err := json.Marshal(conf)
		if err != nil {
			panic(err)
		}
		err = ioutil.WriteFile(fullpath, buf, 0644)
		if err != nil {
			panic(err)
		}
		return conf
	}

	// Read existing config file
	file, _ := os.Open(fullpath)
	decoder := json.NewDecoder(file)
	config := Config{}
	err = decoder.Decode(&config)
	if err != nil {
		panic(err)
	}

	return config
}
