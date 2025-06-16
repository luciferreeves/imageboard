package models

type UserLevel int

const (
	UserLevelMember UserLevel = iota
	UserLevelContributor
	UserLevelJanitor
	UserLevelModerator
	UserLevelAdmin
	UserLevelSuperAdmin
)

func (l UserLevel) String() string {
	switch l {
	case UserLevelMember:
		return "Member"
	case UserLevelContributor:
		return "Contributor"
	case UserLevelJanitor:
		return "Janitor"
	case UserLevelModerator:
		return "Moderator"
	case UserLevelAdmin:
		return "Admin"
	default:
		return "Unknown"
	}
}

func (l UserLevel) Color() string {
	switch l {
	case UserLevelMember:
		return "#8B9DC3" // Soft periwinkle blue
	case UserLevelContributor:
		return "#7FCDAE" // Mint green
	case UserLevelJanitor:
		return "#9BB5FF" // Light electric blue
	case UserLevelModerator:
		return "#FF9F9B" // Coral pink
	case UserLevelAdmin:
		return "#C39BD3" // Lavender purple
	case UserLevelSuperAdmin:
		return "#FFD93D" // Electric yellow
	default:
		return "#B0B0B0" // Neutral gray
	}
}

type Rating string

const (
	RatingSafe         Rating = "Safe"
	RatingQuestionable Rating = "Questionable"
	RatingSensitive    Rating = "Sensitive"
	RatingExplicit     Rating = "Explicit"
)

type ImageContentType string

const (
	ImageContentTypeJPEG    ImageContentType = "image/jpeg"
	ImageContentTypePNG     ImageContentType = "image/png"
	ImageContentTypeGIF     ImageContentType = "image/gif"
	ImageContentTypeWebP    ImageContentType = "image/webp"
	ImageContentTypeAVIF    ImageContentType = "image/avif"
	ImageContentTypeSVG     ImageContentType = "image/svg+xml"
	ImageContentTypeBMP     ImageContentType = "image/bmp"
	ImageContentTypeTIFF    ImageContentType = "image/tiff"
	ImageContentTypeICO     ImageContentType = "image/x-icon"
	ImageContentTypeHEIC    ImageContentType = "image/heic"
	ImageContentTypeHEIF    ImageContentType = "image/heif"
	ImageContentTypeUnknown ImageContentType = "application/octet-stream"
)

type ImageSizeType string

const (
	ImageSizeTypeIcon      ImageSizeType = "icon"
	ImageSizeTypeThumbnail ImageSizeType = "thumbnail"
	ImageSizeTypeSmall     ImageSizeType = "small"
	ImageSizeTypeMedium    ImageSizeType = "medium"
	ImageSizeTypeLarge     ImageSizeType = "large"
	ImageSizeTypeOriginal  ImageSizeType = "original"
)

type TagType string

const (
	TagTypeGeneral   TagType = "general"
	TagTypeArtist    TagType = "artist"
	TagTypeCopyright TagType = "copyright"
	TagTypeCharacter TagType = "character"
	TagTypeMeta      TagType = "meta"
)

func (t TagType) Color() string {
	switch t {
	case TagTypeGeneral:
		return "#4ECDC4" // Turquoise cyan
	case TagTypeArtist:
		return "#FF6B9D" // Hot pink
	case TagTypeCopyright:
		return "#A8E6CF" // Mint green
	case TagTypeCharacter:
		return "#FFB347" // Peach orange
	case TagTypeMeta:
		return "#DDA0DD" // Plum purple
	default:
		return "#E6E6FA" // Light lavender
	}
}

type EmailTokenType string

const (
	EmailTokenTypeVerification  EmailTokenType = "verification"
	EmailTokenTypePasswordReset EmailTokenType = "password_reset"
	EmailTokenTypeChangeEmail   EmailTokenType = "change_email"
)
