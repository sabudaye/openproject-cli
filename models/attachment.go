package models

type Attachment struct {
	Id          uint64
	FileName    string
	FileSize    int64
	Status      string
	ContentType string
	DigestAlgo  string
	DigestHash  string
	CreatedAt   string
	// DownloadHref is the raw downloadLocation link from the API: a path
	// relative to the instance (local storage) or an absolute URL
	// (remote storage such as S3).
	DownloadHref string
}
