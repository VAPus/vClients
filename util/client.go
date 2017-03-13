package util

import (
	"path/filepath"
	"os"
	"io/ioutil"
	"bytes"
	"archive/zip"
	"strings"
	"crypto/md5"
	"fmt"
)

type Client struct {
	Name string
	Win winClient
	Linux linuxClient
}

type winClient struct {
	Enabled bool
	Size int
	CheckSum string
}

type linuxClient struct {
	Enabled bool
	Size int
	CheckSum string
}

// GetClientList retrieves a client list from the given path
func GetClientList(path string) ([]Client, error) {
	list := []Client{}

	// Read main directory
	f, err := ioutil.ReadDir(path)

	if err != nil {
		return nil, err
	}

	for _, file := range f {

		if !file.IsDir() {
			continue
		}

		client, err := getClient(
			filepath.Join(path, file.Name()),
		)

		if err != nil {
			return nil, err
		}

		list = append(list, client)
	}

	return list, nil
}

func getClient(path string) (Client, error) {
	client := Client{}
	client.Win = winClient{}
	client.Linux = linuxClient{}

	// Check windows files
	_, err := os.Stat(filepath.Join(path, "win"))

	if err == nil {
		client.Win.Enabled = true

		// Create file buffer
		buff, err := createClientZip(filepath.Join(path, "win"))

		if err != nil {
			return client, err
		}

		client.Win.Size = len(buff.Bytes())
		client.Win.CheckSum = fmt.Sprintf("%x", md5.Sum(buff.Bytes()))

		if err := ioutil.WriteFile(filepath.Join(path, "win.zip"), buff.Bytes(), os.ModePerm); err != nil {
			return client, err
		}
	}

	// Check linux versions
	_, err = os.Stat(filepath.Join(path, "linux"))

	if err == nil {
		client.Linux.Enabled = true

		// Create file buffer
		buff, err := createClientZip(filepath.Join(path, "linux"))

		if err != nil {
			return client, err
		}

		client.Linux.Size = len(buff.Bytes())
		client.Linux.CheckSum = fmt.Sprintf("%x", md5.Sum(buff.Bytes()))

		if err := ioutil.WriteFile(filepath.Join(path, "linux.zip"), buff.Bytes(), os.ModePerm); err != nil {
			return client, err
		}
	}

	return client, nil
}

func createClientZip(p string) (*bytes.Buffer, error) {
	// Create buffer
	buff := &bytes.Buffer{}

	// Create zip writer
	w := zip.NewWriter(buff)

	if err := filepath.Walk(p, func(path string, info os.FileInfo, e error) error {

		if info.IsDir() {
			return nil
		}

		f, err := w.Create(strings.Replace(path, p+"\\", "", 1))

		if err != nil {
			return err
		}

		// Read file
		content, err := ioutil.ReadFile(path)

		if err != nil {
			return err
		}

		f.Write(content)

		return nil
	}); err != nil {
		return nil, err
	}

	// Close zip writer
	w.Close()

	return buff, nil
}