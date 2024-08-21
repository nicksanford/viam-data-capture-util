package main

import (
	"bufio"
	"encoding/json"
	"log"
	"os"

	v1 "go.viam.com/api/app/datasync/v1"
	"go.viam.com/rdk/data"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("usage %s viam-data-capture-file.capture", os.Args[0])
	}

	f, err := os.Open(os.Args[1])
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
		Metadata:       DataCaptureMetadataToUploadMetadata(cf.ReadMetadata()),
		SensorContents: sd,
	})

	if err != nil {
		log.Fatal(err.Error())
	}

	w := bufio.NewWriter(os.Stdout)
	if _, err := w.WriteString(string(j)); err != nil {
		log.Fatal(err.Error())
	}

}

func DataCaptureMetadataToUploadMetadata(dcm *v1.DataCaptureMetadata) *v1.UploadMetadata {
	return &v1.UploadMetadata{
		ComponentName:    dcm.ComponentName,
		ComponentType:    dcm.ComponentType,
		FileExtension:    dcm.FileExtension,
		FileName:         os.Args[1],
		MethodName:       dcm.MethodName,
		MethodParameters: dcm.MethodParameters,
		PartId:           "unknown",
		Tags:             dcm.Tags,
		Type:             dcm.Type,
	}
}
