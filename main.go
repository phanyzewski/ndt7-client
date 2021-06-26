package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	ndt7 "github.com/m-lab/ndt7-client-go"
	"github.com/m-lab/ndt7-client-go/spec"
	"github.com/olekukonko/tablewriter"
)

const (
	clientName     = "ndt7-client-go-cmd"
	clientVersion  = "0.5.0"
	defaultTimeout = 55 * time.Second
)

func main() {
	client := ndt7.NewClient(clientName, clientVersion)
	parentCtx := context.Background()

	ctx, cancelFunc := context.WithTimeout(parentCtx, defaultTimeout)
	defer cancelFunc()

	dChan, err := client.StartDownload(ctx)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	h := NewHumanReadableWithWriter(os.Stdout)

	err = h.OnStarting(spec.TestDownload)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for event := range dChan {
		func(m *spec.Measurement) {
			h.OnDownloadEvent(&event)
		}(&event)
	}

	err = h.OnComplete(spec.TestDownload)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	uChan, err := client.StartUpload(ctx)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = h.OnStarting(spec.TestUpload)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for event := range uChan {
		func(m *spec.Measurement) {
			h.OnUploadEvent(&event)
		}(&event)
	}

	err = h.OnComplete(spec.TestUpload)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	data := [][]string{}
	if dl, ok := client.Results()[spec.TestDownload]; ok {
		if dl.Client.AppInfo != nil && dl.Client.AppInfo.ElapsedTime > 0 {
			elapsed := float64(dl.Client.AppInfo.ElapsedTime) / 1e06
			downloaded := (8.0 * float64(dl.Client.AppInfo.NumBytes)) / elapsed / (1000.0 * 1000.0)

			data = append(data, []string{"Average Download Speed", fmt.Sprintf("%v Mbit/s", downloaded)})
			// fmt.Printf("downloaded: %v Mbit/s\n", downloaded)
		}

		if dl.Server.TCPInfo != nil {
			if dl.Server.TCPInfo.BytesSent > 0 {
				retrans := float64(dl.Server.TCPInfo.BytesRetrans) / float64(dl.Server.TCPInfo.BytesSent) * 100

				data = append(data, []string{"Retrans Percent", fmt.Sprintf("%v", retrans)})
				// fmt.Printf("retrans: %v%% \n", retrans)
			}

			minRTT := float64(dl.Server.TCPInfo.MinRTT) / 1000
			data = append(data, []string{"MinRTT", fmt.Sprintf("%v", minRTT)})
			// fmt.Printf("minRTT: %vms \n", minRTT)
		}
	}

	if ul, ok := client.Results()[spec.TestUpload]; ok {
		if ul.Client.AppInfo != nil && ul.Client.AppInfo.ElapsedTime > 0 {
			elapsed := float64(ul.Client.AppInfo.ElapsedTime) / 1e06
			uploaded := (8.0 * float64(ul.Client.AppInfo.NumBytes)) / elapsed / (1000.0 * 1000.0)

			data = append(data, []string{"Average Upload Speed", fmt.Sprintf("%v Mbit/s", uploaded)})
			// fmt.Printf("uploaded: %v Mbit/s\n", uploaded)
		}
	}

	// data := [][]string{
	// 	[]string{"A", "The Good", "500"},
	// 	[]string{"B", "The Very very Bad Man", "288"},
	// 	[]string{"C", "The Ugly", "120"},
	// 	[]string{"D", "The Gopher", "800"},
	// }

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Measurement", "Value"})

	for _, v := range data {
		table.Append(v)
	}
	table.Render() // Send output
}

// HumanReadable is a human readable emitter. It emits the events generated
// by running a ndt7 test as pleasant stdout messages.
type HumanReadable struct {
	out io.Writer
}

// NewHumanReadableWithWriter returns a new human readable emitter using the
// specified writer.
func NewHumanReadableWithWriter(w io.Writer) HumanReadable {
	return HumanReadable{w}
}

// OnStarting handles the start event
func (h HumanReadable) OnStarting(test spec.TestKind) error {
	_, err := fmt.Fprintf(h.out, "\nstarting %s", test)
	return err
}

// OnError handles the error event
func (h HumanReadable) OnError(test spec.TestKind, err error) error {
	_, failure := fmt.Fprintf(h.out, "\n%s failed: %s\n", test, err.Error())
	return failure
}

// OnConnected handles the connected event
func (h HumanReadable) OnConnected(test spec.TestKind, fqdn string) error {
	_, err := fmt.Fprintf(h.out, "\n%s in progress with %s\n", test, fqdn)
	return err
}

// OnDownloadEvent handles an event emitted by the download test
func (h HumanReadable) OnDownloadEvent(m *spec.Measurement) error {
	return h.onSpeedEvent(m)
}

// OnUploadEvent handles an event emitted during the upload test
func (h HumanReadable) OnUploadEvent(m *spec.Measurement) error {
	return h.onSpeedEvent(m)
}

func (h HumanReadable) onSpeedEvent(m *spec.Measurement) error {
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
	_, err := fmt.Fprintf(h.out, "\nAvg. speed  : %7.1f Mbit/s", v)
	return err
}

// OnComplete handles the complete event
func (h HumanReadable) OnComplete(test spec.TestKind) error {
	_, err := fmt.Fprintf(h.out, "\n%s: complete\n", test)

	return err
}
