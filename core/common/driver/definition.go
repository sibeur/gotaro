package driver

type StorageDriverType uint32

const (
	// Google Cloud Storage Driver
	GCSDriverType StorageDriverType = 1
)

var AllowedDrivers []StorageDriverType = []StorageDriverType{
	GCSDriverType,
}

type UploadFileOpts struct {
	Mime string
}
