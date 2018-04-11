// Copyright © 2018 Microsoft Corporation and contributors
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/eventgrid/eventgrid"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var startConfig = viper.New()

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		http.ListenAndServe(fmt.Sprintf(":%d", viper.GetInt("port")), nil)
	},
}

func init() {
	rootCmd.AddCommand(startCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// startCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// startCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	startCmd.Flags().IntP("port", "p", 80, "The port that should be used to host the server emulating an EventGrid Topic.")
	startConfig.BindPFlag("port", startCmd.Flags().Lookup("port"))

	http.HandleFunc("/api/events", ProcessEventsHandler)
	http.HandleFunc("/subscribe", RegisterSubscriberHandler)
}

// ProcessEventsHandler reads an HTTP Request that informs an Event Grid topic of an Event.
// It then relays that message to all subscribers who have not filtered out messages
// of this type and subject.
func ProcessEventsHandler(resp http.ResponseWriter, req *http.Request) {
	const MaxPayloadSize = 1024 * 1024
	const MaxEventSize = 64 * 1024

	limitedBody := io.LimitReader(req.Body, MaxPayloadSize)

	var payload eventgrid.Event

	if contents, err := ioutil.ReadAll(limitedBody); err == nil {
		if err = json.Unmarshal(contents, &payload); err != nil {
			resp.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(resp)
			return
		}
	} else {
		resp.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(resp, "unable to parse request body:", err)
		return
	}
}

// RegisterSubscriberHandler mimics the ARM behavior of adding a subscriber to an Event Grid Topic.
//
// Note: It is important to understand that in a production app, this
func RegisterSubscriberHandler(resp http.ResponseWriter, req *http.Request) {

}
