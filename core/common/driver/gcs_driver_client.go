package driver

import (
	"context"
	"io"
	"log"
	"os"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

type GCPDriverClientUseCase interface {
	// GetBucket returns the bucket object
	GetBucket() *storage.BucketHandle
	GetClient() *storage.Client
	GetDriverConfig() *GCSDriverConfig
	GetObjectNames() ([]string, error)
	UploadFile(filePath string, targetFilePath string) (string, error)
	Close()
}

type GCPDriverClient struct {
	driverConfig *GCSDriverConfig
	client       *storage.Client
	ctx          context.Context
}

func NewGCPDriverClient(driverConfig *GCSDriverConfig) (*GCPDriverClient, error) {
	// generate gcp client
	ctx := context.Background()
	jsonAuth, err := driverConfig.GetDecodedServiceAccount()
	if err != nil {
		return nil, err
	}
	authOpt := option.WithCredentialsJSON(jsonAuth)
	client, err := storage.NewClient(ctx, authOpt)
	if err != nil {
		return nil, err
	}
	return &GCPDriverClient{
		driverConfig: driverConfig,
		client:       client,
	}, nil
}

func (gcp *GCPDriverClient) GetBucket() *storage.BucketHandle {
	return gcp.client.Bucket(gcp.driverConfig.BucketName)
}

func (gcp *GCPDriverClient) GetClient() *storage.Client {
	return gcp.client
}

func (gcp *GCPDriverClient) GetDriverConfig() *GCSDriverConfig {
	return gcp.driverConfig
}

func (gcp *GCPDriverClient) GetObjectNames() ([]string, error) {
	objects := []string{}
	it := gcp.GetBucket().Objects(gcp.ctx, nil)
	for {
		obj, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		objects = append(objects, obj.Name)
	}
	return objects, nil
}

func (gcp *GCPDriverClient) UploadFile(filePath string, targetFilePath string) (string, error) {
	client := gcp.client

	bucketName := gcp.driverConfig.BucketName

	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	ctx := context.Background()
	obj := client.Bucket(bucketName).Object(targetFilePath)
	wc := obj.NewWriter(ctx)
	if _, err := io.Copy(wc, file); err != nil {
		return "", err
	}
	if err := wc.Close(); err != nil {
		return "", err
	}

	attrs, err := obj.Attrs(ctx)
	if err != nil {
		return "", err
	}

	return attrs.MediaLink, nil
}

func (gcp *GCPDriverClient) Close() {
	if err := gcp.client.Close(); err != nil {
		log.Printf("Error closing gcp driver client %v", err)
	}
}
