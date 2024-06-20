package service

import (
	"errors"
	"log"

	"github.com/sibeur/gotaro/core/common"
	driver_lib "github.com/sibeur/gotaro/core/common/driver"
	"github.com/sibeur/gotaro/core/entity"
	"github.com/sibeur/gotaro/core/repository"
)

type MediaService struct {
	repo          *repository.Repository
	DriverManager *driver_lib.DriverManager
}

func NewMediaService(repo *repository.Repository, driverManager *driver_lib.DriverManager) *MediaService {
	return &MediaService{repo: repo, DriverManager: driverManager}
}

func (u *MediaService) FindAll() ([]*entity.Media, error) {
	result, err := u.repo.Media.FindAll()
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (u *MediaService) Upload(ruleSlug, fileName string, isCommit bool) (*entity.Media, error) {
	tempFilePath := common.TemporaryFolder + "/" + fileName

	rule, err := u.repo.Rule.FindBySlug(ruleSlug)
	if err != nil {
		log.Printf("Error finding rule: %v", err)
		return nil, err
	}

	if rule == nil {
		log.Printf("Error finding rule")
		return nil, errors.New(common.ErrRuleNotFoundMsg)
	}

	driver, err := u.repo.Driver.FindByID(rule.DriverID)
	if err != nil {
		log.Printf("Error finding driver: %v", err)
		return nil, err
	}

	if driver == nil {
		log.Printf("Error finding driver")
		return nil, errors.New(common.ErrDriverNotFoundMsg)
	}

	driverClient := u.DriverManager.GetDriver(driver.Slug)
	if driverClient == nil {
		log.Printf("Error finding driver client")
		return nil, errors.New(common.ErrDriverClientNotFoundMsg)
	}

	fileAliasName := common.GetFileNameUnique(fileName)

	fileMetaData, err := common.GetFileMetaData(tempFilePath)
	if err != nil {
		log.Printf("Error getting file meta data: %v", err)
		return nil, err
	}

	// validate file size
	fileSizeKB := fileMetaData.FileSize / 1024
	if fileSizeKB > rule.MaxSize {
		log.Printf("File size exceeded max size: %v", fileSizeKB)
		return nil, errors.New(common.ErrFileSizeExceededMsg)
	}

	// validate file mime
	if !common.IsMimeValid(rule.Mimes, fileMetaData.FileMime) {
		log.Printf("File mime invalid: %v", fileMetaData.FileMime)
		return nil, errors.New(common.ErrFileMimeInvalidMsg)
	}

	folder := driver.GetDefaultFolder()
	targetFilePath := folder + "/" + fileAliasName
	if folder == "/" {
		targetFilePath = fileAliasName
	}

	filePathFromDriver := driver.GetFilePathFromDriver(targetFilePath)

	mediaLink, err := driverClient.UploadFile(tempFilePath, targetFilePath)
	if err != nil {
		log.Printf("Error uploading file: %v", err)
		return nil, err
	}

	isPublic, _ := driverClient.IsStorageAssetPublic()

	media := entity.Media{
		RuleSlug:           ruleSlug,
		DriverSlug:         driver.Slug,
		FileOriginalName:   fileName,
		FileAliasName:      fileAliasName,
		FileExt:            fileMetaData.FileExt,
		FileMime:           fileMetaData.FileMime,
		FileSize:           fileMetaData.FileSize,
		FilePath:           mediaLink,
		FilePathFromDriver: filePathFromDriver,
		IsCommit:           isCommit,
		IsPublic:           isPublic,
	}
	err = u.repo.Media.Create(&media)

	if err != nil {
		log.Printf("Error creating media: %v", err)
		return nil, err
	}

	return &media, nil
}

func (u *MediaService) Delete(ruleSlug, fileAliasName string) error {
	return u.repo.Media.Delete(ruleSlug, fileAliasName)
}

func (u *MediaService) FindMedia(ruleSlug, fileAliasName string) (*entity.Media, error) {
	media, err := u.repo.Media.FindMedia(ruleSlug, fileAliasName)
	if err != nil {
		log.Printf("Error finding media: %v", err)
		return nil, err
	}

	if media == nil {
		return nil, nil
	}
	if !media.IsPublic {
		signedUrl, err := u.repo.Media.GetCachedSignedUrl(ruleSlug, fileAliasName)
		if err != nil {
			log.Printf("Error getting cached signed url: %v", err)
		}

		if signedUrl != "" {
			media.FilePath = signedUrl
		}

		if signedUrl == "" {
			driver := u.DriverManager.GetDriver(media.DriverSlug)
			if driver == nil {
				log.Printf("Error finding driver client")
				return nil, errors.New(common.ErrDriverNotFoundMsg)
			}
			newSignedUrl, err := driver.GetSignedUrl(media.FileAliasName)
			if err != nil {
				log.Printf("Error getting signed url: %v", err)
				return nil, err
			}
			media.FilePath = newSignedUrl
			go func(ruleSlug, fileAliasName, newSignedUrl string) {
				u.repo.Media.SetCachedSignedUrl(ruleSlug, fileAliasName, newSignedUrl)
				err = u.repo.Media.SetSignedUrl(media.RuleSlug, media.FileAliasName, newSignedUrl)
				if err != nil {
					log.Printf("Error setting signed url: %v", err)
				}
			}(ruleSlug, fileAliasName, newSignedUrl)
		}
	}
	return media, nil
}
