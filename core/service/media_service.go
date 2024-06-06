package service

import (
	"errors"

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

func (u *MediaService) Upload(ruleSlug, fileName string, isCommit bool) (string, error) {
	tempFilePath := common.TemporaryFolder + "/" + fileName

	rule, err := u.repo.Rule.FindBySlug(ruleSlug)
	if err != nil {
		return "", err
	}

	if rule == nil {
		return "", errors.New(common.ErrRuleNotFoundMsg)
	}

	driver, err := u.repo.Driver.FindByID(rule.DriverID)
	if err != nil {
		return "", err
	}

	if driver == nil {
		return "", errors.New(common.ErrDriverNotFoundMsg)
	}

	driverClient := u.DriverManager.GetDriver(driver.Slug)
	if driverClient == nil {
		return "", errors.New(common.ErrDriverClientNotFoundMsg)
	}

	fileAliasName := common.GetFileNameUnique(fileName)

	fileMetaData, err := common.GetFileMetaData(tempFilePath)
	if err != nil {
		return "", err
	}

	// validate file size
	fileSizeKB := fileMetaData.FileSize / 1024
	if fileSizeKB > rule.MaxSize {
		return "", errors.New(common.ErrFileSizeExceededMsg)
	}

	// validate file mime
	if !common.IsMimeValid(rule.Mimes, fileMetaData.FileMime) {
		return "", errors.New(common.ErrFileMimeInvalidMsg)
	}

	folder := driver.GetDefaultFolder()
	targetFilePath := folder + "/" + fileAliasName
	if folder == "/" {
		targetFilePath = fileAliasName
	}

	filePathFromDriver := driver.GetFilePathFromDriver(targetFilePath)

	mediaLink, err := driverClient.UploadFile(tempFilePath, targetFilePath)
	if err != nil {
		return "", err
	}

	err = u.repo.Media.Create(&entity.Media{
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
	})

	if err != nil {
		return "", err
	}

	return mediaLink, nil
}

func (u *MediaService) Delete(ruleSlug, fileAliasName string) error {
	return u.repo.Media.Delete(ruleSlug, fileAliasName)
}

func (u *MediaService) FindMedia(ruleSlug, fileAliasName string) (*entity.Media, error) {
	return u.repo.Media.FindMedia(ruleSlug, fileAliasName)
}
