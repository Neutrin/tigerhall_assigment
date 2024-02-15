package inits

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/url"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

var (
	storageClient *storage.Client
	bucket        = "tigers_image_sightings"
)

// TODO make it something like to read from json object and not path
func InitStorageClient() {
	var err error
	ctx := context.Background()
	storageClient, err = storage.NewClient(ctx, option.WithCredentialsFile("/Users/nitinthakur-mbp/personal/tigerhall_assigment/core/inits/keys.json"))
	if err != nil {
		panic(err)
	}
	fmt.Println(" ********* init storage class intialiased succesfully *************")

}

func UploadFile(f multipart.File, uploadedFile *multipart.FileHeader) (*url.URL, error) {
	var (
		err error
		//url *url.URL
	)
	ctx := context.Background()
	sw := storageClient.Bucket(bucket).Object(uploadedFile.Filename).NewWriter(ctx)
	if _, err = io.Copy(sw, f); err != nil {
		return nil, err
	}
	if err = sw.Close(); err != nil {
		return nil, err
	}
	fmt.Println(" here ", sw.Attrs().Name)
	return url.Parse("/" + bucket + "/" + sw.Attrs().Name)

}
