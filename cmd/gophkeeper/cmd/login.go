package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/theherk/viper"
	"golang.org/x/term"
)

var loginUserName string

func init() {
	loginCmd.Flags().StringVarP(&loginUserName, "username", "u", "", "username required")
	loginCmd.MarkFlagRequired("username")
	rootCmd.AddCommand(loginCmd)
}

type loginRequest struct {
	Name     string `json:"email"`
	Password string `json:"password"`
}

type loginResponse struct {
	AccessToken string `json:"token"`
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "login user",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(loginUserName) == 0 {
			return fmt.Errorf("empty username is not allowed")
		}
		fmt.Println("Enter password:")
		password, err := term.ReadPassword(0)
		if err != nil {
			return err
		}

		body, err := json.Marshal(loginRequest{
			Name:     loginUserName,
			Password: string(password),
		})
		if err != nil {
			return err
		}

		req, err := http.NewRequest(http.MethodPost, upstreamURL+"/api/login", bytes.NewBuffer(body))
		if err != nil {
			return err
		}

		resp, err := httpClient.Do(req)
		if err != nil {
			return err
		}
		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("failed to login user, http status code:'%d' and body:'%s'", resp.StatusCode, body)
		}
		var loginResp loginResponse
		decoder := json.NewDecoder(resp.Body)
		err = decoder.Decode(&loginResp)
		if err != nil {
			return err
		}
		viper.Set("token", loginResp.AccessToken)
		err = viper.WriteConfig()
		if err != nil {
			return err
		}
		fmt.Println("User authenticated")
		return nil
	},
}
