package dtos

import (
	"github.com/sabudaye/openproject-cli/models"
)

type AttachmentDigestDto struct {
	Algorithm string `json:"algorithm,omitempty"`
	Hash      string `json:"hash,omitempty"`
}

type AttachmentLinksDto struct {
	Self             *LinkDto `json:"self,omitempty"`
	Container        *LinkDto `json:"container,omitempty"`
	Author           *LinkDto `json:"author,omitempty"`
	DownloadLocation *LinkDto `json:"downloadLocation,omitempty"`
}

type AttachmentDto struct {
	Id          int64                `json:"id,omitempty"`
	FileName    string               `json:"fileName,omitempty"`
	FileSize    int64                `json:"fileSize,omitempty"`
	Filesize    int64                `json:"filesize,omitempty"`
	Status      string               `json:"status,omitempty"`
	ContentType string               `json:"contentType,omitempty"`
	Digest      *AttachmentDigestDto `json:"digest,omitempty"`
	CreatedAt   string               `json:"createdAt,omitempty"`
	Links       *AttachmentLinksDto  `json:"_links,omitempty"`
}

type attachmentElements struct {
	Elements []*AttachmentDto `json:"elements"`
}

type AttachmentCollectionDto struct {
	Embedded attachmentElements `json:"_embedded"`
	Type     string             `json:"_type"`
	Total    int64              `json:"total"`
	Count    int64              `json:"count"`
}

/////////////// MODEL CONVERSION ///////////////

func (dto *AttachmentDto) Convert() *models.Attachment {
	fileSize := dto.FileSize
	if fileSize == 0 && dto.Filesize != 0 {
		// The APIv3 spec declares the attribute as "fileSize" but its own
		// examples (and some instances) serialize "filesize"; accept both.
		fileSize = dto.Filesize
	}

	attachment := &models.Attachment{
		Id:          uint64(dto.Id),
		FileName:    dto.FileName,
		FileSize:    fileSize,
		Status:      dto.Status,
		ContentType: dto.ContentType,
		CreatedAt:   dto.CreatedAt,
	}
	if dto.Digest != nil {
		attachment.DigestAlgo = dto.Digest.Algorithm
		attachment.DigestHash = dto.Digest.Hash
	}
	if dto.Links != nil && dto.Links.DownloadLocation != nil {
		attachment.DownloadHref = dto.Links.DownloadLocation.Href
	}

	return attachment
}

func (dto *AttachmentCollectionDto) Convert() []*models.Attachment {
	var attachments = make([]*models.Attachment, len(dto.Embedded.Elements))

	for idx, a := range dto.Embedded.Elements {
		attachments[idx] = a.Convert()
	}

	return attachments
}
