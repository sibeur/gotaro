package driver

import (
	"encoding/base64"
	"encoding/json"

	"go.mongodb.org/mongo-driver/bson"
)

type GCSDriverConfig struct {
	ProjectID      string `json:"project_id" bson:"project_id"`
	BucketName     string `json:"bucket_name" bson:"bucket_name"`
	DefaultFolder  string `json:"default_folder" bson:"default_folder"`
	ServiceAccount string `json:"service_account" bson:"service_account"`
}

func NewGCSDriverConfig(projectID, bucketName, defaultFolder string, serviceAccount []byte) *GCSDriverConfig {
	encodedServiceAccount := base64.StdEncoding.EncodeToString(serviceAccount)
	return &GCSDriverConfig{
		ProjectID:      projectID,
		BucketName:     bucketName,
		DefaultFolder:  defaultFolder,
		ServiceAccount: encodedServiceAccount,
	}
}

func (conf *GCSDriverConfig) GetDecodedServiceAccount() ([]byte, error) {
	decodedBytes, err := base64.StdEncoding.DecodeString(conf.ServiceAccount)
	if err != nil {
		return nil, err
	}
	return decodedBytes, nil
}

func (conf *GCSDriverConfig) ToJSONBytes() ([]byte, error) {
	jsonBytes, err := json.Marshal(conf)
	if err != nil {
		return nil, err
	}
	return jsonBytes, nil
}

func (conf *GCSDriverConfig) ToBSON() (bson.M, error) {
	jsonBytes, err := conf.ToJSONBytes()
	if err != nil {
		return nil, err
	}
	var bsonM bson.M
	err = json.Unmarshal(jsonBytes, &bsonM)
	if err != nil {
		return nil, err
	}
	return bsonM, nil
}
