package model

import (
	"chaincode/pkg/utils"
	"errors"
)

func ValidateUser(user User) error {
	// User ID: 3-16 characters, only letters, numbers, and underscores
	// User Name: 3-16 characters, only letters, numbers, and underscores

	if !utils.ValidateLength(user.ID, 3, 16) {
		return errors.New("User ID must be between 3 and 16 characters")
	}
	if !utils.ValidateLength(user.Name, 3, 16) {
		return errors.New("User Name must be between 3 and 16 characters")
	}

	if !utils.ValidateName(user.ID) {
		return errors.New("User ID must contain only letters, numbers, and underscores")
	}
	if !utils.ValidateName(user.Name) {
		return errors.New("User Name must contain only letters, numbers, and underscores")
	}

	return nil
}

func ValidateFile(file File) error {
	// File Size:  0-2GB as Bytes
	// File Hash: SHA-256

	if !utils.ValidateRange64(file.Size, 0, int64(2)*1024*1024*1024) {
		return errors.New("File Size must be between 0 and 2GB")
	}
	if !utils.ValidateSHA256(file.Hash) {
		return errors.New("File Hash must be a SHA-256 hash")
	}

	return nil
}

func ValidateDatasetFile(datasetFile DatasetFile) error {
	// File Name: 1-64 characters [not contain invalid characters]
	// File Hash: SHA-256

	if !utils.ValidateLength(datasetFile.FileName, 1, 64) {
		return errors.New("File Name must be between 1 and 64 characters")
	}
	if !utils.ValidateFileName(datasetFile.FileName) {
		return errors.New("File Name must not contain invalid characters")
	}
	if !utils.ValidateSHA256(datasetFile.Hash) {
		return errors.New("File Hash must be a SHA-256 hash")
	}

	return nil
}

func ValidateVersion(version Version) error {
	// Files: list of DatasetFile
	// Rows: non-negative integer
	// Creation Time: ISO 8601
	// Change Log: 0-1024 characters

	for _, file := range version.Files {
		if err := ValidateDatasetFile(file); err != nil {
			return err
		}
	}
	if version.Rows < 0 {
		return errors.New("Rows must be a non-negative integer")
	}
	if !utils.ValidateTime(version.CreationTime) {
		return errors.New("Creation Time must be an ISO 8601 timestamp")
	}
	if !utils.ValidateLength(version.ChangeLog, 0, 1024) {
		return errors.New("Changelog must be between 0 and 1024 characters")
	}

	return nil
}

func ValidateDataset(dataset Dataset) error {
	// Owner ID: existing user [3-16 characters, only letters, numbers, and underscores]
	// Dataset Name: 3-64 characters [only letters, numbers, and underscores]

	if !utils.ValidateLength(dataset.Owner, 3, 16) {
		return errors.New("Owner ID must be between 3 and 16 characters")
	}
	if !utils.ValidateLength(dataset.Name, 3, 64) {
		return errors.New("Dataset Name must be between 3 and 64 characters")
	}
	if !utils.ValidateName(dataset.Owner) {
		return errors.New("Owner ID must contain only letters, numbers, and underscores")
	}
	if !utils.ValidateName(dataset.Name) {
		return errors.New("Dataset Name must contain only letters, numbers, and underscores")
	}

	for _, version := range dataset.Versions {
		if err := ValidateVersion(version); err != nil {
			return err
		}
	}

	return nil
}

func ValidateRecord(record Record) error {
	// Owner ID: existing user [3-16 characters, only letters, numbers, and underscores]
	// Dataset Name: existing dataset [3-64 characters, only letters, numbers, and underscores]
	// User ID: existing user [3-16 characters, only letters, numbers, and underscores]
	// Files: list of DatasetFile
	// Time: ISO 8601

	if !utils.ValidateLength(record.DatasetOwner, 3, 16) {
		return errors.New("Dataset Owner must be between 3 and 16 characters")
	}
	if !utils.ValidateLength(record.DatasetName, 3, 64) {
		return errors.New("Dataset Name must be between 3 and 64 characters")
	}
	if !utils.ValidateLength(record.User, 3, 16) {
		return errors.New("User ID must be between 3 and 16 characters")
	}
	if !utils.ValidateName(record.DatasetOwner) {
		return errors.New("Dataset Owner must contain only letters, numbers, and underscores")
	}
	if !utils.ValidateName(record.DatasetName) {
		return errors.New("Dataset Name must contain only letters, numbers, and underscores")
	}
	if !utils.ValidateName(record.User) {
		return errors.New("User ID must contain only letters, numbers, and underscores")
	}
	if !utils.ValidateTime(record.Time) {
		return errors.New("Time must be an ISO 8601 timestamp")
	}
	for _, file := range record.Files {
		if err := ValidateDatasetFile(file); err != nil {
			return err
		}
	}

	return nil
}
