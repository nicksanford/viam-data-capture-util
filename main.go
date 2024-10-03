package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	v1 "go.viam.com/api/app/datasync/v1"
	"go.viam.com/rdk/data"
)

var recursive bool

func init() {
	flag.BoolVar(&recursive, "r", false, "decode all .capture files into .json files recursively from the provided directory")
}

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) != 1 {
		log.Fatalf("usage: %s [-r] (viam-data-capture-file.capture | directory_for_recursive_option)", os.Args[0])
	}

	if recursive {
		filepath.Walk(args[0], func(path string, info fs.FileInfo, err error) error {
			if err != nil || info.IsDir() || filepath.Ext(path) != ".capture" {
				return nil
			}

			f, err := os.Create(strings.ReplaceAll(path, ".capture", ".json"))
			if err != nil {
				panic(err.Error())
			}

			convert(path, f)
			f.Close()
			return nil
		})
	} else {
		convert(args[0], os.Stdout)
	}
}

func convert(captureFileName string, out io.Writer) {
	f, err := os.Open(captureFileName)
	if err != nil {
		log.Fatal(err.Error())
	}

	cf, err := data.ReadCaptureFile(f)
	if err != nil {
		f.Close()
		log.Fatal(err.Error())
	}
	defer cf.Close()

	sd, err := data.SensorDataFromCaptureFile(cf)
	if err != nil {
		log.Fatal(err.Error())
	}

	j, err := json.Marshal(&v1.DataCaptureUploadRequest{
		Metadata:       DataCaptureMetadataToUploadMetadata(cf.ReadMetadata(), captureFileName),
		SensorContents: sd,
	})

	if err != nil {
		log.Fatal(err.Error())
	}

	w := bufio.NewWriter(out)
	if _, err := w.WriteString(string(j)); err != nil {
		log.Fatal(err.Error())
	}
	w.Flush()
}

func DataCaptureMetadataToUploadMetadata(dcm *v1.DataCaptureMetadata, filename string) *v1.UploadMetadata {
	return &v1.UploadMetadata{
		ComponentName: dcm.ComponentName,
		ComponentType: dcm.ComponentType,
		// This is not correct for images (and possibly for other binary files)
		// TODO: Do what viam-server does in this case
		FileExtension:    dcm.FileExtension,
		FileName:         filename,
		MethodName:       dcm.MethodName,
		MethodParameters: dcm.MethodParameters,
		Tags:             dcm.Tags,
		Type:             dcm.Type,
	}
}
