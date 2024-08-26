package model

import (
	"errors"
	"regexp"
)

func ValidateUser(user User) error {
	// User ID: 3-16 characters, only letters, numbers, and underscores
	// User Name: 3-16 characters, only letters, numbers, and underscores

	if len(user.ID) < 3 || len(user.ID) > 16 {
		return errors.New("User ID must be between 3 and 16 characters")
	}
	if len(user.Name) < 3 || len(user.Name) > 16 {
		return errors.New("User Name must be between 3 and 16 characters")
	}
	if !regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString(user.ID) {
		return errors.New("User ID must contain only letters, numbers, and underscores")
	}
	if !regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString(user.Name) {
		return errors.New("User Name must contain only letters, numbers, and underscores")
	}
	return nil
}
func ValidateFile(file DatasetFile) error {
	// File Name: 3-64 characters
	// File Size: positive integer
	// File Hash: SHA-256

	if len(file.FileName) < 3 || len(file.FileName) > 64 {
		return errors.New("File Name must be between 3 and 64 characters")
	}
	if file.Size <= 0 {
		return errors.New("File Size must be a positive integer")
	}
	if !regexp.MustCompile(`^[a-f0-9]{64}$`).MatchString(file.Hash) {
		return errors.New("File Hash must be a SHA-256 hash")
	}
	return nil
}
func ValidateVersion(version DatasetVersion) error {
	// Creation Time: ISO 8601
	// Change Log: 0-1024 characters
	// Size: positive integer
	// Rows: non-negative integer
	// Files: list of SHA-256 hashes

	if !regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`).MatchString(version.CreationTime) {
		return errors.New("Creation Time must be an ISO 8601 timestamp")
	}
	if len(version.ChangeLog) < 0 || len(version.ChangeLog) > 1024 {
		return errors.New("Change Log must be between 0 and 1024 characters")
	}
	if version.Size <= 0 {
		return errors.New("Size must be a positive integer")
	}
	if version.Rows < 0 {
		return errors.New("Rows must be a non-negative integer")
	}
	for _, file := range version.Files {
		if !regexp.MustCompile(`^[a-f0-9]{64}$`).MatchString(file) {
			return errors.New("Files must be a list of SHA-256 hashes")
		}
	}
	return nil
}
func ValidateDataset(dataset Dataset) error {
	// Dataset Name: 1-64 characters
	// Owner ID: existing user [3-16 characters, only letters, numbers, and underscores]

	if len(dataset.Name) < 1 || len(dataset.Name) > 64 {
		return errors.New("Dataset Name must be between 1 and 64 characters")
	}
	if len(dataset.Owner) < 3 || len(dataset.Owner) > 16 {
		return errors.New("Owner ID must be between 3 and 16 characters")
	}
	if !regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString(dataset.Owner) {
		return errors.New("Owner ID must contain only letters, numbers, and underscores")
	}

	// Validate Version
	for _, version := range dataset.Versions {
		if err := ValidateVersion(version); err != nil {
			return err
		}
	}
	return nil
}
