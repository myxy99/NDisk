package model

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type (
	File struct {
		gorm.Model
		Name         string
		Size         uint64
		Type         string
		Md5          string
		Sha1         string
		Sha256       string
		FileSystem   uint8
		FileRealPath string
		BlockSize    uint64
	}
)

func (f *File) TableName() string {
	return "db_common_file"
}

func (f *File) SetHash(ht, hc string) (err error) {
	switch strings.ToLower(ht) {
	case "md5":
		f.Md5 = hc
	case "sha1":
		f.Sha1 = hc
	case "sha256":
		f.Sha256 = hc
	default:
		err = fmt.Errorf("can not support hashType %s", ht)
	}
	return
}

func (f *File) Add() (err error) {
	err = MainDB().Create(f).Error
	return
}

func (f *File) SliceSize() uint64 {
	return (f.Size + f.BlockSize - 1) / f.BlockSize
}

func (f *File) FullPath(idx int) string {
	return fmt.Sprintf("%s/%d/%d", f.FileRealPath, f.ID, idx)
}
