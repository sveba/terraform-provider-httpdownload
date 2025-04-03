package client

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type HttpDownloadClient struct {
	client http.Client
}

func NewHttpDownloadClient() *HttpDownloadClient {
	c := HttpDownloadClient{
		client: http.Client{Timeout: 10 * time.Second},
	}
	return &c
}

func (c *HttpDownloadClient) DownloadFileToDest(ctx context.Context, filename string, remote_url string) error {
	tflog.Info(ctx, "Downloading file", map[string]any{"url": remote_url, "filename": filename})

	// Download the file
	httpResp, err := http.Get(remote_url)
	if err != nil {
		tflog.Error(ctx, err.Error())
		return err
	}
	defer httpResp.Body.Close()

	// Check server response
	if httpResp.StatusCode != http.StatusOK {
		err = fmt.Errorf("file download error %s", httpResp.Status)
		tflog.Error(ctx, err.Error())
		return err
	}

	// Create the file
	file, err := os.Create(filename)
	if err != nil {
		tflog.Error(ctx, err.Error())
		return err
	}
	defer file.Close()

	// Write the body to file
	_, err = io.Copy(file, httpResp.Body)
	if err != nil {
		tflog.Error(ctx, err.Error())
		return err
	}

	return nil
}
