package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

var dataID string

func init() {
	getCmd.PersistentFlags().StringVar(&dataID, "id", "", "data identificator")
	getCmd.MarkPersistentFlagRequired("id")
	getCmd.AddCommand(getCredCmd)
	getCmd.AddCommand(getFileCmd)
	getCmd.AddCommand(getBankCmd)
	rootCmd.AddCommand(getCmd)
}

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "get value from storage",
}

var getFileCmd = &cobra.Command{
	Use:   "file",
	Short: "get binary data",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := http.Client{}
		req, err := http.NewRequest(http.MethodGet, upstreamURL+"/api/keeper/file/"+dataID, nil)
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
			return fmt.Errorf("failed to download file with http code: %d", resp.StatusCode)
		}

		contentDisposition := resp.Header.Get("Content-Disposition")
		_, params, err := mime.ParseMediaType(contentDisposition)
		if err != nil {
			return err
		}
		filename := params["filename"]
		if len(filename) == 0 {
			return fmt.Errorf("failed to get filename")
		}

		out, err := os.Create(filename)
		if err != nil {
			return err
		}
		defer out.Close()

		_, err = io.Copy(out, resp.Body)
		if err != nil {
			return err
		}
		return nil
	},
}

type credentialsResponse struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	Title    string `json:"title"`
	Meta     string `json:"meta"`
}

var getCredCmd = &cobra.Command{
	Use:   "cred",
	Short: "get credentials",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := http.Client{}
		req, err := http.NewRequest(http.MethodGet, upstreamURL+"/api/keeper/credentials/"+dataID, nil)
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
			return fmt.Errorf("failed to get credentials with http code: %d", resp.StatusCode)
		}

		var bodyResp credentialsResponse
		decoder := json.NewDecoder(resp.Body)
		err = decoder.Decode(&bodyResp)
		if err != nil {
			return err
		}
		fmt.Printf("UserName: %s Password: %s\n", bodyResp.Name, bodyResp.Password)
		return nil
	},
}

type bankResponse struct {
	Title  string
	Meta   string
	Number string
	Holder string
	Cvv    int
}

var getBankCmd = &cobra.Command{
	Use:   "bank",
	Short: "get bank data",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := http.Client{}
		req, err := http.NewRequest(http.MethodGet, upstreamURL+"/api/keeper/bank/"+dataID, nil)
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
			return fmt.Errorf("failed to get credentials with http code: %d", resp.StatusCode)
		}

		var bodyResp bankResponse
		decoder := json.NewDecoder(resp.Body)
		err = decoder.Decode(&bodyResp)
		if err != nil {
			return err
		}
		fmt.Printf("Card number: %s card holder: %s cvv:%d\n", bodyResp.Number, bodyResp.Holder, bodyResp.Cvv)
		return nil
	},
}
