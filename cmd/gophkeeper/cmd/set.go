package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"

	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	filePath     string
	meta         string
	credName     string
	credPassword string
	credTitle    string
	cardNumber   string
	cardHolder   string
	cardCvv      int
)

func init() {
	setCmd.PersistentFlags().StringVarP(&meta, "meta", "m", "", "metadata")

	setFileCmd.Flags().StringVar(&filePath, "path", "", "path to file")
	setFileCmd.MarkFlagRequired("path")

	setCredCmd.Flags().StringVar(&credName, "name", "", "user name")
	setCredCmd.Flags().StringVar(&credPassword, "password", "", "user password")
	setCredCmd.Flags().StringVar(&credTitle, "title", "", "credentials record name")
	setCredCmd.MarkFlagRequired("name")
	setCredCmd.MarkFlagRequired("password")
	setCredCmd.MarkFlagRequired("title")

	setBankCmd.Flags().StringVar(&cardNumber, "number", "", "card number")
	setBankCmd.Flags().StringVar(&cardHolder, "holder", "", "card holder")
	setBankCmd.Flags().IntVar(&cardCvv, "cvv", 0, "card cvv")
	setBankCmd.MarkFlagRequired("number")
	setBankCmd.MarkFlagRequired("holder")
	setBankCmd.MarkFlagRequired("cvv")

	setCmd.AddCommand(setFileCmd)
	setCmd.AddCommand(setCredCmd)
	setCmd.AddCommand(setBankCmd)
	rootCmd.AddCommand(setCmd)
}

var setCmd = &cobra.Command{
	Use:   "set",
	Short: "set value to remote storage",
}

var setFileCmd = &cobra.Command{
	Use:   "file",
	Short: "store binary data",
	RunE: func(cmd *cobra.Command, args []string) error {
		pr, pw := io.Pipe()
		writer := multipart.NewWriter(pw)
		ct := writer.FormDataContentType()
		go func() {
			fileName := filepath.Base(filePath)
			file, err := os.OpenFile(filePath, os.O_RDONLY, os.ModePerm)
			if err != nil {
				pw.CloseWithError(err)
				return
			}
			defer file.Close()
			part, err := writer.CreateFormFile("file", fileName)
			if err != nil {
				pw.CloseWithError(err)
				return
			}
			_, err = io.Copy(part, file)
			if err != nil {
				pw.CloseWithError(err)
				return
			}
			pw.CloseWithError(writer.Close())
		}()

		req, err := http.NewRequest(http.MethodPost, upstreamURL+"/api/keeper/file", pr)
		if err != nil {
			return err
		}
		setAuthToken(req)
		req.Header.Set("Content-Type", ct)

		resp, err := httpClient.Do(req)
		if err != nil {
			return err
		}
		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("failed to upload file, http status code:'%d' and body:'%s'", resp.StatusCode, body)
		}
		return nil
	},
}

type credentialsRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	Title    string `json:"title"`
	Meta     string `json:"meta"`
}

var setCredCmd = &cobra.Command{
	Use:   "cred",
	Short: "store credentials",
	RunE: func(cmd *cobra.Command, args []string) error {
		body, err := json.Marshal(credentialsRequest{
			Name:     credName,
			Password: credPassword,
			Title:    credTitle,
			Meta:     meta,
		})
		if err != nil {
			return err
		}

		req, err := http.NewRequest(http.MethodPost, upstreamURL+"/api/keeper/credentials", bytes.NewBuffer(body))
		if err != nil {
			return err
		}
		req.Header.Set("Content", "application/json")
		setAuthToken(req)
		resp, err := httpClient.Do(req)
		if err != nil {
			return err
		}

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("failed to store credentials, http status code:'%d' and body:'%s'", resp.StatusCode, body)
		}
		return nil
	},
}

type setBankRequest struct {
	Title  string
	Meta   string
	Number string
	Holder string
	Cvv    int
}

var setBankCmd = &cobra.Command{
	Use:   "bank",
	Short: "set bank account",
	RunE: func(cmd *cobra.Command, args []string) error {
		body, err := json.Marshal(setBankRequest{
			Number: cardNumber,
			Holder: cardHolder,
			Cvv:    cardCvv,
			Title:  credTitle,
			Meta:   meta,
		})
		if err != nil {
			return err
		}

		req, err := http.NewRequest(http.MethodPost, upstreamURL+"/api/keeper/bank", bytes.NewBuffer(body))
		if err != nil {
			return err
		}
		req.Header.Set("Content", "application/json")
		setAuthToken(req)
		resp, err := httpClient.Do(req)
		if err != nil {
			return err
		}

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("failed to store credentials, http status code:'%d' and body:'%s'", resp.StatusCode, body)
		}
		return nil
	},
}
