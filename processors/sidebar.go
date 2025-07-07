package processors

import (
	"fmt"
	"imageboard/database"
	"imageboard/models"

	"github.com/gofiber/fiber/v2"
)

type SiteStats struct {
	Posts    string
	Tags     string
	Today    string
	Storage  string
	Comments string
}

func SidebarContextProcessor(ctx *fiber.Ctx) error {
	popularTags, popularTagsErr := database.GetPopularTags(15)
	if popularTagsErr != nil || len(popularTags) == 0 {
		mockTags := []models.Tag{
			{Name: "anime", Type: models.TagTypeGeneral, Count: 1523},
			{Name: "manga", Type: models.TagTypeGeneral, Count: 892},
			{Name: "kawaii", Type: models.TagTypeGeneral, Count: 756},
			{Name: "retro", Type: models.TagTypeMeta, Count: 634},
			{Name: "y2k", Type: models.TagTypeMeta, Count: 511},
			{Name: "aesthetic", Type: models.TagTypeGeneral, Count: 445},
			{Name: "sakura", Type: models.TagTypeArtist, Count: 389},
			{Name: "studio_ghibli", Type: models.TagTypeCopyright, Count: 312},
			{Name: "totoro", Type: models.TagTypeCharacter, Count: 298},
			{Name: "sailor_moon", Type: models.TagTypeCharacter, Count: 267},
			{Name: "pokemon", Type: models.TagTypeCopyright, Count: 234},
			{Name: "pixiv", Type: models.TagTypeMeta, Count: 198},
			{Name: "digital_art", Type: models.TagTypeMeta, Count: 176},
			{Name: "watercolor", Type: models.TagTypeGeneral, Count: 145},
			{Name: "minimalist", Type: models.TagTypeGeneral, Count: 123},
		}
		ctx.Locals("PopularTags", mockTags)
	} else {
		ctx.Locals("PopularTags", popularTags)
	}

	recentTags, recentTagsErr := database.GetRecentTags(10)
	if recentTagsErr != nil || len(recentTags) == 0 {
		mockRecentTags := []models.Tag{
			{Name: "cyberpunk", Type: models.TagTypeGeneral, Count: 23},
			{Name: "vaporwave", Type: models.TagTypeMeta, Count: 45},
			{Name: "synthwave", Type: models.TagTypeGeneral, Count: 12},
			{Name: "retrocomputing", Type: models.TagTypeMeta, Count: 8},
			{Name: "neon", Type: models.TagTypeGeneral, Count: 67},
			{Name: "glitch", Type: models.TagTypeMeta, Count: 34},
			{Name: "pixel_art", Type: models.TagTypeGeneral, Count: 89},
			{Name: "lo_fi", Type: models.TagTypeGeneral, Count: 56},
		}
		ctx.Locals("RecentTags", mockRecentTags)
	} else {
		ctx.Locals("RecentTags", recentTags)
	}

	postsCount, postsErr := database.GetTotalPostsCount()
	tagsCount, tagsCountErr := database.GetTotalTagsCount()
	commentsCount, commentsErr := database.GetTotalCommentsCount()
	todayCount, todayErr := database.GetTodayPostsCount()
	storageSize, storageErr := database.GetTotalStorageSize()

	var stats SiteStats

	if postsErr == nil {
		stats.Posts = fmt.Sprintf("%d", postsCount)
	} else {
		stats.Posts = "0"
	}
	if tagsCountErr == nil {
		stats.Tags = fmt.Sprintf("%d", tagsCount)
	} else {
		stats.Tags = "0"
	}
	if commentsErr == nil {
		stats.Comments = fmt.Sprintf("%d", commentsCount)
	} else {
		stats.Comments = "0"
	}
	if todayErr == nil {
		stats.Today = fmt.Sprintf("%d new", todayCount)
	} else {
		stats.Today = "0 new"
	}
	if storageErr == nil {
		stats.Storage = storageSize
	} else {
		stats.Storage = "0 B"
	}

	ctx.Locals("SiteStats", stats)

	return ctx.Next()
}
