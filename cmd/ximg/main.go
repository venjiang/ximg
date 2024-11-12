package main

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"

	flag "github.com/spf13/pflag"
	"github.com/venjiang/ximg"
)

var (
	output string
)

func main() {
	flag.Parse()
	// check arguments
	if len(os.Args) < 2 {
		fmt.Println("Usage: ximg <file>")
		os.Exit(1)
	}
	file := os.Args[1]
	data, err := getFileContent(file)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// check image type
	var buf bytes.Buffer
	mime := http.DetectContentType(data)
	switch mime {
	case "image/jpeg":
		buf.WriteString("data:image/jpeg;base64,")
	case "image/png":
		buf.WriteString("data:image/png;base64,")
	default:
		slog.Error("unsupported mime type", "mime", mime)
		os.Exit(1)
	}
	// base64 encode
	buf.Write(ximg.Base64Encode(data))
	// output
	if output == "" {
		fmt.Println(buf.String())
		return
	}
	slog.Info("output to file", "file", output)
	if err := os.WriteFile(output, buf.Bytes(), 0644); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func getFileContent(file string) (data []byte, err error) {
	if isRemoteFile(file) {
		slog.Info("download image...", "file", file)
		data, err = downloadImage(file)
		if err != nil {
			slog.Error("download image failed: %s", "error", err)
			return nil, err
		}
		slog.Info("downloaded completed", "len", len(data))
	} else {
		data, err = os.ReadFile(file)
		if err != nil {
			slog.Error("read file failed", "error", err)
			return nil, err
		}
	}
	return
}

func isRemoteFile(file string) bool {
	return strings.HasPrefix(file, "http://") || strings.HasPrefix(file, "https://")
}

func downloadImage(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func init() {
	flag.StringVarP(&output, "output", "o", "", "output file")
}
