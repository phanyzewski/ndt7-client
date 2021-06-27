package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	ndt7 "github.com/m-lab/ndt7-client-go"
	"github.com/m-lab/ndt7-client-go/spec"
	"github.com/olekukonko/tablewriter"
)

// EmitterOutput is a human readable emitter. It emits the events generated
// by running a ndt7 test as pleasant stdout messages.
type EmitterOutput struct {
	out io.Writer
}

// NewEmitterOutput returns a new human readable emitter using the
// specified writer.
func NewEmitterOutput(w io.Writer) *EmitterOutput {
	return &EmitterOutput{w}
}

// Started handles the start event
func (e EmitterOutput) Started(test spec.TestKind) error {
	_, err := fmt.Fprintf(e.out, "\nstarted %s\n", test)
	return err
}

// Failed handles an error event
func (e EmitterOutput) Failed(test spec.TestKind, err error) error {
	_, failure := fmt.Fprintf(e.out, "\n%s failed: %s\n", test, err.Error())
	return failure
}

// Connected handles the connected event
func (e EmitterOutput) Connected(test spec.TestKind, fqdn string) error {
	_, err := fmt.Fprintf(e.out, "\n%s in progress with %s\n", test, fqdn)
	return err
}

// SpeedEvent handles the emitter output
func (e EmitterOutput) SpeedEvent(m *spec.Measurement) error {
	// The specification recommends that we show application level
	// measurements. Let's just do that in interactive mode. To this
	// end, we ignore any measurement coming from the server.
	if m.Origin != spec.OriginClient {
		return nil
	}
	if m.AppInfo == nil || m.AppInfo.ElapsedTime <= 0 {
		return errors.New("Missing m.AppInfo or invalid m.AppInfo.ElapsedTime")
	}
	elapsed := float64(m.AppInfo.ElapsedTime) / 1e06
	v := (8.0 * float64(m.AppInfo.NumBytes)) / elapsed / (1000.0 * 1000.0)
	_, err := fmt.Fprintf(e.out, "\r%7.1f Mbit/s", v)
	return err
}

// Completed posts a complete event notification
func (e EmitterOutput) Completed(test spec.TestKind) error {
	_, err := fmt.Fprintf(e.out, "\n%s: completed\n", test)

	return err
}

// Summary is a tabledized summary of the test activity
func (e EmitterOutput) Summary(client *ndt7.Client) {
	data := [][]string{}
	if dl, ok := client.Results()[spec.TestDownload]; ok {
		if dl.Client.AppInfo != nil && dl.Client.AppInfo.ElapsedTime > 0 {
			elapsed := float64(dl.Client.AppInfo.ElapsedTime) / 1e06
			downloaded := (8.0 * float64(dl.Client.AppInfo.NumBytes)) / elapsed / (1000.0 * 1000.0)

			data = append(data, []string{"Average Download Speed", fmt.Sprintf("%.2f Mbit/s", downloaded)})
		}

		if dl.Server.TCPInfo != nil {
			if dl.Server.TCPInfo.BytesSent > 0 {
				retrans := float64(dl.Server.TCPInfo.BytesRetrans) / float64(dl.Server.TCPInfo.BytesSent) * 100

				data = append(data, []string{"Retrans Percent", fmt.Sprintf("%.4f %%", retrans)})
			}

			minRTT := float64(dl.Server.TCPInfo.MinRTT) / 1000
			data = append(data, []string{"MinRTT", fmt.Sprintf("%.4f ms", minRTT)})
		}
	}

	if ul, ok := client.Results()[spec.TestUpload]; ok {
		if ul.Client.AppInfo != nil && ul.Client.AppInfo.ElapsedTime > 0 {
			elapsed := float64(ul.Client.AppInfo.ElapsedTime) / 1e06
			uploaded := (8.0 * float64(ul.Client.AppInfo.NumBytes)) / elapsed / (1000.0 * 1000.0)

			data = append(data, []string{"Average Upload Speed", fmt.Sprintf("%.2f Mbit/s", uploaded)})
		}
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Measurement", "Value"})

	for _, v := range data {
		table.Append(v)
	}

	fmt.Println()
	table.Render()
}
