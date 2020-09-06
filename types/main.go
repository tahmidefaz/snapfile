package types

import "gorm.io/gorm"

type DbModal struct {
	gorm.Model
	FileName		string	`json:"file_name"`
	MaxDownloads 	int		`json:"max_downloads"`
}

type FileInfo struct {
	FileName  		string	`json:"file_name" binding:"required"`
	MaxDownloads 	int		`json:"max_downloads" binding:"required"`
}