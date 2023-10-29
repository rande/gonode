// Copyright © 2014-2023 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
)

type file struct {
	remoteURL string
	localFile string
}

var files = []file{
	{
		remoteURL: "https://unpkg.com/bootstrap/dist/css/bootstrap.min.css",
		localFile: "/../modules/bootstrap/static/bootstrap.min.css",
	},
	{
		remoteURL: "https://unpkg.com/bootstrap/dist/js/bootstrap.esm.js",
		localFile: "/../modules/bootstrap/static/bootstrap.esm.js",
	},
	{
		remoteURL: "https://unpkg.com/bootstrap/dist/js/bootstrap.min.js",
		localFile: "/../modules/bootstrap/static/bootstrap.min.js",
	},
	{
		remoteURL: "https://unpkg.com/bootstrap/dist/js/bootstrap.bundle.min.js",
		localFile: "/../modules/bootstrap/static/bootstrap.bundle.min.js",
	},
}

func main() {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		fmt.Println("Unable to get the current filename")
	}
	dirname := filepath.Dir(filename)
	fmt.Println(dirname)

	// Download each file from the files struct and save it to the local filesystem
	for _, f := range files {
		fmt.Println("Downloading", f.remoteURL, "to", dirname+f.localFile)
		if err := downloadFile(f.remoteURL, dirname+f.localFile); err != nil {
			fmt.Println(err)
			continue
		}
	}

	fmt.Println("Bootstrap files downloaded successfully!")
}

// Helper function to download a file from a URL and save it to a local file
func downloadFile(url string, filepath string) error {
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Download the file
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Write the file contents to disk
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
