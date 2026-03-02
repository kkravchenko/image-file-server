package transport

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"
)

type Client struct {
	client *http.Client
}

func NewClient(timeout time.Duration) (*Client, error) {
	if timeout == 0 {
		return nil, errors.New("timeout can be zero")
	}

	return &Client{
		client: &http.Client{
			Timeout: timeout,
			Transport: &loggingRoundTripper{
				next: http.DefaultTransport,
			},
		},
	}, nil
}

func (c Client) DownloadFile(filepath string, url string) error {

	out, err := os.Create(filepath)
	if err != nil {
		slog.Error("Error create file "+filepath+": ", err.Error())
		return err
	}
	defer func() {
		err := out.Close()
		if err != nil {
			slog.Error("Error close file "+filepath+": ", err.Error())
		}
	}()

	resp, err := http.Get(url)
	if err != nil {
		slog.Error("Error download file "+filepath+": ", err.Error())
		return err
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			slog.Error("Error close file "+filepath+": ", err.Error())
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		slog.Error("Error copy to file: ", err.Error())
		return err
	}

	return nil
}

type loggingRoundTripper struct {
	next http.RoundTripper
}

func (l loggingRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	slog.Info(fmt.Sprintf("[%s] %s %s]n", time.Now().Format(time.ANSIC), r.Method, r.URL))
	return l.next.RoundTrip(r)
}
