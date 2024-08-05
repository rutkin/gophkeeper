package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var userName string

func init() {
	registerCmd.Flags().StringVarP(&userName, "username", "u", "", "username required")
	registerCmd.MarkFlagRequired("username")
	rootCmd.AddCommand(registerCmd)
}

type registerRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "register new user",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(userName) == 0 {
			return fmt.Errorf("empty username is not allowed")
		}
		fmt.Println("Enter password:")
		password, err := term.ReadPassword(0)
		if err != nil {
			return err
		}

		req, err := json.Marshal(registerRequest{
			Name:     userName,
			Password: string(password),
		})
		if err != nil {
			return err
		}

		resp, err := http.Post(upstreamURL+"/api/register", "application/json", bytes.NewBuffer(req))
		if err != nil {
			return err
		}
		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("failed to register user, http status code:'%d' and body:'%s'", resp.StatusCode, body)
		}
		fmt.Println("User registered")
		return nil
	},
}
