package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path"

	"github.com/rutkin/gophkeeper/cmd/gophkeeper/cmd"
	"github.com/theherk/viper"
)

func getUserHomeDir() string {
	usr, err := user.Current()
	if err != nil {
		fmt.Println("Can't get your home directory.")
		os.Exit(1)
	}

	return usr.HomeDir
}

func getConfigDir() string {
	return path.Join(getUserHomeDir(), ".config")
}

func getConfigPath() string {
	return path.Join(getConfigDir(), "pusher.json")
}

func main() {
	if _, err := os.Stat(getConfigDir()); os.IsNotExist(err) {
		err = os.Mkdir(getConfigDir(), os.ModeDir|0755)
		if err != nil {
			panic(err)
		}
	}

	if _, err := os.Stat(getConfigPath()); os.IsNotExist(err) {
		err = ioutil.WriteFile(getConfigPath(), []byte("{}"), 0600)
		if err != nil {
			panic(err)
		}
	}

	viper.SetConfigFile(getConfigPath())
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	cmd.Execute()
}
