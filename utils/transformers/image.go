package transformers

import "imageboard/config"

func ConvertStringRatingToType(rating string) (config.Rating, error) {
	switch rating {
	case "safe":
		return config.RatingSafe, nil
	case "questionable":
		return config.RatingQuestionable, nil
	case "sensitive":
		return config.RatingSensitive, nil
	case "explicit":
		return config.RatingExplicit, nil
	default:
		return config.RatingSafe, nil
	}
}
