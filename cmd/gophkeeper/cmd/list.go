package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/theherk/viper"
)

func init() {
	rootCmd.AddCommand(listCmd)
}

type itemResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type listItemsResponse struct {
	Items []itemResponse `json:"items"`
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list user items",
	RunE: func(cmd *cobra.Command, args []string) error {
		token := viper.GetString("token")
		client := http.Client{}
		req, err := http.NewRequest(http.MethodGet, upstreamURL+"/api/keeper", nil)
		if err != nil {
			return err
		}
		req.Header.Set("authorization", fmt.Sprintf("bearer %s", token))

		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("failed to list, http status code:'%d' and body:'%s'", resp.StatusCode, body)
		}

		var listResp listItemsResponse
		decoder := json.NewDecoder(resp.Body)
		err = decoder.Decode(&listResp)
		if err != nil {
			return err
		}

		fmt.Println("ID Name Type")

		for _, resp := range listResp.Items {
			fmt.Printf("%s %s %s\n", resp.ID, resp.Name, resp.Type)
		}
		return nil
	},
}
