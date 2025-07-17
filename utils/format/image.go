package format

import (
	"imageboard/config"
	"strings"
)

func GetCDNURL() string {
	cdnURL := strings.TrimRight(config.S3.PublicURL, "/") + "/" + config.S3.BucketName
	if config.S3.FolderPath != "" {
		cdnURL += "/" + config.S3.FolderPath
	}
	return cdnURL
}
