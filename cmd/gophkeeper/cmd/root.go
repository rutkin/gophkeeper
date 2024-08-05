package cmd

import (
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/theherk/viper"
)

var upstreamURL = "http://127.0.0.1:8080"

var rootCmd = &cobra.Command{
	Use:   "gophkeeper",
	Short: "gophkeeper application for store secrets",
	Long:  "gophkeeper is a command line client for store secrets",
}

func makeRequest(url string) error {
	client := http.Client{}
	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return err
	}
	setAuthToken(req)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send request %s with http status: %d", url, resp.StatusCode)
	}
	return nil
}

func setAuthToken(req *http.Request) {
	token := viper.GetString("token")
	req.Header.Set("authorization", fmt.Sprintf("bearer %s", token))
}

func Execute() error {
	return rootCmd.Execute()
}
