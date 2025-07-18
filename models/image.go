package models

import (
	"fmt"
	"imageboard/config"
	"imageboard/utils/format"
	"imageboard/utils/math"
	"strings"

	"gorm.io/gorm"
)

type ImageSize struct {
	gorm.Model
	ImageID  uint                 `gorm:"not null;index" json:"-"`
	Image    Image                `gorm:"foreignKey:ImageID" json:"-"`
	SizeType config.ImageSizeType `gorm:"not null;size:50" json:"size_type"`
	Width    int                  `gorm:"not null" json:"width"`
	Height   int                  `gorm:"not null" json:"height"`
	FileSize int64                `gorm:"not null" json:"file_size"`
}

func (s *ImageSize) BeforeCreate(tx *gorm.DB) error {
	if s.Width <= 0 || s.Height <= 0 {
		return fmt.Errorf("image dimensions must be greater than zero")
	}

	if s.FileSize <= 0 {
		return fmt.Errorf("file size must be greater than zero")
	}

	return nil
}

func (s *ImageSize) GetURL() string {
	if config.S3.PublicURL == "" {
		return ""
	}
	return fmt.Sprintf("%s/%s/%s/%s/%s", config.S3.PublicURL, config.S3.BucketName, config.S3.FolderPath, s.SizeType, s.Image.FileName)
}

func (s *ImageSize) GetAspectRatio() float64 {
	if s.Height == 0 {
		return 0
	}
	return float64(s.Width) / float64(s.Height)
}

func (s *ImageSize) GetDimensions() string {
	return fmt.Sprintf("%dx%d", s.Width, s.Height)
}

func (s *ImageSize) GetFileSizeFormatted() string {
	return format.FileSize(s.FileSize)
}

type Image struct {
	gorm.Model
	FileName       string                  `gorm:"not null;size:255" json:"file_name"`
	ContentType    config.ImageContentType `gorm:"not null;size:100" json:"content_type"`
	MD5Hash        string                  `gorm:"not null;size:32" json:"md5_hash"`
	Title          string                  `gorm:"default:'';size:255" json:"title"`
	Description    string                  `gorm:"default:'';type:text" json:"description"`
	SourceURL      string                  `gorm:"default:'';size:500" json:"source_url"`
	Rating         config.Rating           `gorm:"not null;default:'safe';size:15" json:"rating"`
	IsApproved     bool                    `gorm:"not null;default:true" json:"is_approved"`
	IsDeleted      bool                    `gorm:"not null;default:false" json:"is_deleted"`
	ThreadLocked   bool                    `gorm:"not null;default:false" json:"thread_locked"`
	UploaderID     uint                    `gorm:"not null;index" json:"-"`
	Uploader       User                    `gorm:"foreignKey:UploaderID" json:"uploader"`
	ApproverID     *uint                   `gorm:"index" json:"-"`
	Approver       *User                   `gorm:"foreignKey:ApproverID" json:"approver,omitempty"`
	RelatedImages  []Image                 `gorm:"many2many:image_relationships;joinForeignKey:image_id;joinReferences:related_image_id" json:"related_images,omitempty"`
	ViewCount      int64                   `gorm:"not null;default:0" json:"view_count"`
	FavouriteCount int64                   `gorm:"not null;default:0" json:"favorite_count"`
	CommentCount   int64                   `gorm:"not null;default:0" json:"comment_count"`
	Sizes          []ImageSize             `gorm:"foreignKey:ImageID" json:"sizes,omitempty"`
	Tags           []Tag                   `gorm:"many2many:image_tags;joinForeignKey:image_id;joinReferences:tag_id" json:"tags,omitempty"`
	FavoritedBy    []User                  `gorm:"many2many:user_favorites" json:"favorited_by,omitempty"`
	Comments       []Comment               `gorm:"foreignKey:ImageID" json:"comments,omitempty"`
}

func (i *Image) BeforeCreate(tx *gorm.DB) error {
	i.FileName = strings.TrimSpace(i.FileName)
	i.Title = strings.TrimSpace(i.Title)
	i.Description = strings.TrimSpace(i.Description)

	if i.FileName == "" {
		return fmt.Errorf("file name cannot be empty")
	}

	if len(i.MD5Hash) != 32 {
		return fmt.Errorf("MD5 hash must be exactly 32 characters long")
	}

	return nil
}

func (i *Image) BeforeDelete(tx *gorm.DB) error {
	return tx.Exec(`UPDATE tags SET count = count - 1 WHERE id IN (
		SELECT tag_id FROM image_tags WHERE image_id = ?
	) AND count > 0`, i.ID).Error
}

func (i *Image) GetURL(sizeType config.ImageSizeType) string {
	for _, size := range i.Sizes {
		if size.SizeType == sizeType {
			return size.GetURL()
		}
	}

	return ""
}

func (i *Image) GetSize(sizeType config.ImageSizeType) *ImageSize {
	for _, size := range i.Sizes {
		if size.SizeType == sizeType {
			return &size
		}
	}
	return nil
}

func (i *Image) GetSizeByString(sizeType string) *ImageSize {
	for _, size := range i.Sizes {
		if string(size.SizeType) == sizeType {
			return &size
		}
	}
	return nil
}

func (i *Image) GetOriginalDimensions() string {
	if fullSize := i.GetSize(config.ImageSizeTypeOriginal); fullSize != nil {
		return fullSize.GetDimensions()
	}
	return "Unknown"
}

func (i *Image) GetOriginalSize() *ImageSize {
	return i.GetSize(config.ImageSizeTypeOriginal)
}

func (i *Image) GetSmallSize() *ImageSize {
	return i.GetSize(config.ImageSizeTypeSmall)
}

func (i *Image) GetMediumSize() *ImageSize {
	return i.GetSize(config.ImageSizeTypeMedium)
}

func (i *Image) GetLargeSize() *ImageSize {
	return i.GetSize(config.ImageSizeTypeLarge)
}

func (i *Image) GetThumbnailSize() *ImageSize {
	return i.GetSize(config.ImageSizeTypeThumbnail)
}

func (i *Image) GetAspectRatio() string {
	if fullSize := i.GetSize(config.ImageSizeTypeOriginal); fullSize != nil {
		if fullSize.Height == 0 {
			return "Unknown"
		}

		width := fullSize.Width
		height := fullSize.Height

		divisor := math.GCD(width, height)
		simplifiedWidth := width / divisor
		simplifiedHeight := height / divisor

		return fmt.Sprintf("%d:%d", simplifiedWidth, simplifiedHeight)
	}
	return "Unknown"
}

func (i *Image) AddSize(tx *gorm.DB, sizeType config.ImageSizeType, width, height int, fileSize int64) (*ImageSize, error) {
	if width <= 0 || height <= 0 {
		return nil, fmt.Errorf("image dimensions must be greater than zero")
	}

	if fileSize <= 0 {
		return nil, fmt.Errorf("file size must be greater than zero")
	}

	size := &ImageSize{
		ImageID:  i.ID,
		SizeType: sizeType,
		Width:    width,
		Height:   height,
		FileSize: fileSize,
	}

	if err := tx.Create(size).Error; err != nil {
		return nil, fmt.Errorf("failed to create image size: %v", err)
	}

	i.Sizes = append(i.Sizes, *size)
	return size, nil
}

func (i *Image) AddRelatedImage(tx *gorm.DB, relatedImage *Image) error {
	if relatedImage.ID == 0 {
		return fmt.Errorf("related image must be saved before adding relationship")
	}

	if relatedImage.IsDeleted {
		return fmt.Errorf("cannot add deleted image as related image")
	}

	if i.ID == relatedImage.ID {
		return fmt.Errorf("cannot relate an image to itself")
	}

	// If the relationship already exists, do nothing
	var count int64
	if err := tx.Table("image_relationships").Where("image_id = ? AND related_image_id = ?", i.ID, relatedImage.ID).Count(&count).Error; err != nil {
		return fmt.Errorf("failed to check existing relationship: %v", err)
	}

	if count > 0 {
		return nil
	}

	// Create bi-directional relationship
	if err := tx.Model(&i).Association("RelatedImages").Append(relatedImage); err != nil {
		return fmt.Errorf("failed to add related image: %v", err)
	}
	if err := tx.Model(&relatedImage).Association("RelatedImages").Append(i); err != nil {
		return fmt.Errorf("failed to add related image: %v", err)
	}

	return nil
}

func (i *Image) RemoveRelatedImage(tx *gorm.DB, relatedImage *Image) error {
	if relatedImage.ID == 0 {
		return fmt.Errorf("related image must be saved before removing relationship")
	}

	if i.ID == relatedImage.ID {
		return fmt.Errorf("cannot remove self from related images")
	}

	// Remove bi-directional relationship
	if err := tx.Model(&i).Association("RelatedImages").Delete(relatedImage); err != nil {
		return fmt.Errorf("failed to remove related image: %v", err)
	}
	if err := tx.Model(&relatedImage).Association("RelatedImages").Delete(i); err != nil {
		return fmt.Errorf("failed to remove related image: %v", err)
	}

	return nil
}

func (i *Image) GetRelatedImages(tx *gorm.DB) ([]Image, error) {
	var relatedImages []Image
	if err := tx.Model(&i).Association("RelatedImages").Find(&relatedImages); err != nil {
		return nil, fmt.Errorf("failed to get related images: %v", err)
	}
	return relatedImages, nil
}

func (i *Image) AddTag(tx *gorm.DB, tag *Tag) error {
	if i.IsDeleted || tag.IsDeleted {
		return fmt.Errorf("cannot add tag to deleted image or add deleted tag")
	}

	// Check if already associated
	var count int64
	if err := tx.Table("image_tags").Where("image_id = ? AND tag_id = ?", i.ID, tag.ID).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	// Add association
	if err := tx.Model(i).Association("Tags").Append(tag); err != nil {
		return err
	}

	// Update tag count
	return tx.Model(tag).UpdateColumn("count", gorm.Expr("count + ?", 1)).Error
}

func (i *Image) RemoveTag(tx *gorm.DB, tag *Tag) error {
	// Check if associated
	var count int64
	if err := tx.Table("image_tags").Where("image_id = ? AND tag_id = ?", i.ID, tag.ID).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		return nil // Not associated
	}

	// Remove association
	if err := tx.Model(i).Association("Tags").Delete(tag); err != nil {
		return err
	}

	// Update tag count
	return tx.Model(tag).UpdateColumn("count", gorm.Expr("GREATEST(count - ?, 0)", 1)).Error
}

func (i *Image) GetTags(tx *gorm.DB) ([]Tag, error) {
	var tags []Tag
	if err := tx.Model(i).Association("Tags").Find(&tags); err != nil {
		return nil, fmt.Errorf("failed to get image tags: %v", err)
	}
	return tags, nil
}

func (i *Image) DeleteImage(tx *gorm.DB) error {
	if i.IsDeleted {
		return fmt.Errorf("image is already deleted")
	}
	i.IsDeleted = true
	return tx.Save(i).Error
}
