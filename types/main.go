package types

import "gorm.io/gorm"

type DbModal struct {
	gorm.Model
	FileName		string	`json:"file_name" binding:"required"`
	Url				string	`json:"url" binding:"required"`
	MaxDownloads 	int		`json:"max_downloads" binding:"required"`
}

type FileInfo struct {
	FileName  		string	`json:"file_name" binding:"required"`
	MaxDownloads 	int		`json:"max_downloads" binding:"required"`
}
