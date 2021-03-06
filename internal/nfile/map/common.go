package _map

import (
	"errors"
	"fmt"
	"strings"
)

const (
	DefaultPage     = 1
	DefaultPageSize = 10
)

var DefaultPageRequest = PageList{
	Page:     DefaultPage,
	PageSize: DefaultPageSize,
}

type (
	PageList struct {
		Page     int    `json:"page" form:"page" validate:"required,number" label:"页码"`
		PageSize int    `json:"page_size" form:"page_size" validate:"required,number" label:"页码大小"`
		Keyword  string `json:"keyword" form:"keyword"`
		IsDelete bool   `json:"is_delete" form:"is_delete"`
	}

	IdMap struct {
		Id uint `uri:"id" json:"id" validate:"required,number,min=1" label:"id"`
	}

	Header struct {
		FileId     uint   `header:"File-Id"`
		SliceIndex int    `header:"Slice-Index"`
		Size       uint   `header:"Size"`
		HashType   string `header:"Hash-Type"`
		HashCode   string `header:"Hash-Code"`
	}
)

func (h *Header) Validate() error {
	if hType := h.HashType; len(hType) > 0 {
		switch strings.ToUpper(hType) {
		case "md5":
		case "sha1":
		case "sha256":
		default:
			return fmt.Errorf("hash type %s is not support", hType)
		}
		if len(h.HashCode) == 0 {
			return errors.New("if hashType not null,hashCode is required")
		}
	}
	return nil
}

func NewHeader() *Header {
	return &Header{}
}
