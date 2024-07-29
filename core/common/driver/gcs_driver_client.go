package driver

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"cloud.google.com/go/storage"
	"github.com/sibeur/gotaro/core/common"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

type GCPDriverClientUseCase interface {
	// GetBucket returns the bucket object
	GetBucket() *storage.BucketHandle
	GetClient() *storage.Client
	GetDriverConfig() *GCSDriverConfig
	GetObjectNames() ([]string, error)
	UploadFile(filePath string, targetFilePath string, opts ...*UploadFileOpts) (string, error)
	GetSignedUrl(filePath string) (string, error)
	IsStorageAssetPublic() (bool, error)
	IsStorageBucketExist() (bool, error)
	IsHasStorageAdminPrivilage() (bool, error)
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

func (gcp *GCPDriverClient) UploadFile(filePath string, targetFilePath string, opts ...*UploadFileOpts) (string, error) {
	opt := &UploadFileOpts{}

	if len(opts) > 0 {
		opt = opts[0]
	}

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
	if opt.Mime != "" {
		wc.ContentType = opt.Mime
	}
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

	url := attrs.MediaLink

	isPublic, _ := gcp.IsStorageAssetPublic()
	if !isPublic {
		signedUrl, _ := gcp.GetSignedUrl(targetFilePath)
		url = signedUrl
	}

	return url, nil
}

func (gcp *GCPDriverClient) GetSignedUrl(filePath string) (string, error) {
	opts := &storage.SignedURLOptions{
		Scheme:  storage.SigningSchemeV4,
		Method:  "GET",
		Expires: time.Now().Add(common.DefaultSignedURLTTL),
	}
	signedUrl, err := gcp.client.Bucket(gcp.driverConfig.BucketName).SignedURL(filePath, opts)
	if err != nil {
		return "", err
	}
	return signedUrl, nil
}

func (gcp *GCPDriverClient) IsStorageAssetPublic() (bool, error) {
	ctx := context.Background()

	bucketName := gcp.driverConfig.BucketName
	bucket := gcp.client.Bucket(bucketName)

	// Check if the bucket is public
	attrs, err := bucket.Attrs(ctx)
	if err != nil {
		return false, err
	}
	isPublic := false
	for _, entity := range attrs.ACL {
		if entity.Entity == "allUsers" && entity.Role == storage.RoleReader {
			isPublic = true
			break
		}
	}
	return isPublic, nil
}

func (gcp *GCPDriverClient) IsStorageBucketExist() (bool, error) {

	ctx := context.Background()

	bucketName := gcp.driverConfig.BucketName
	_, err := gcp.client.Bucket(bucketName).Attrs(ctx)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (gcp *GCPDriverClient) IsHasStorageAdminPrivilage() (bool, error) {

	ctx := context.Background()

	bucketName := gcp.driverConfig.BucketName

	bucket := gcp.client.Bucket(bucketName)

	// Check if the bucket is public
	attrs, err := bucket.Attrs(ctx)
	if err != nil {
		if errors.Is(err, storage.ErrBucketNotExist) {
			fmt.Println("IsHasStorageAdminPrivilage error", err)
			return false, errors.New(common.ErrBucketNotExistMsg)
		}
		return false, err
	}
	isHasPrivilage := false
	for _, entity := range attrs.ACL {
		if entity.Entity == "allUsers" && entity.Role == storage.RoleOwner {
			isHasPrivilage = true
			break
		}
	}
	return isHasPrivilage, nil
}

func (gcp *GCPDriverClient) Close() {
	if err := gcp.client.Close(); err != nil {
		log.Printf("Error closing gcp driver client %v", err)
	}
}
