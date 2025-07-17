package handlers

import (
	"imageboard/config"
	"strings"
)

func ExtractRatingsAndMap(queryParams []config.QueryParam) ([]config.Rating, map[string]bool) {
	var ratings []config.Rating
	ratingsMap := map[string]bool{}
	for _, param := range queryParams {
		if param.Key == "rating" {
			switch strings.ToLower(param.Value) {
			case "safe":
				ratings = append(ratings, config.RatingSafe)
				ratingsMap["Safe"] = true
			case "questionable":
				ratings = append(ratings, config.RatingQuestionable)
				ratingsMap["Questionable"] = true
			case "sensitive":
				ratings = append(ratings, config.RatingSensitive)
				ratingsMap["Sensitive"] = true
			case "explicit":
				ratings = append(ratings, config.RatingExplicit)
				ratingsMap["Explicit"] = true
			}
		}
	}
	if len(ratings) == 0 {
		ratings = []config.Rating{
			config.RatingSafe,
			config.RatingQuestionable,
			config.RatingSensitive,
		}
		ratingsMap["Safe"] = true
		ratingsMap["Questionable"] = true
		ratingsMap["Sensitive"] = true
	}
	return ratings, ratingsMap
}

func ExtractQueryTags(queryParams []config.QueryParam) (string, []string) {
	for _, param := range queryParams {
		if param.Key == "tags" {
			tags := strings.TrimSpace(param.Value)
			if tags == "" {
				return "", nil
			}
			tagList := strings.Split(tags, ",")
			for i := range tagList {
				tagList[i] = strings.TrimSpace(tagList[i])
			}
			return tags, tagList
		}
	}
	return "", nil
}
