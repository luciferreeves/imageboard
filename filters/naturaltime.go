package filters

import (
	"fmt"
	"time"

	"github.com/flosch/pongo2/v6"
)

func naturaltimeFilter(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	timeValue, ok := in.Interface().(time.Time)
	if !ok {
		return nil, &pongo2.Error{
			Sender:    "filter:naturaltime",
			OrigError: fmt.Errorf("input must be a time.Time, got %T", in.Interface()),
		}
	}

	now := time.Now()
	diff := now.Sub(timeValue)

	if diff < 0 {
		diff = -diff
		if diff < time.Minute {
			return pongo2.AsValue("in a few seconds"), nil
		} else if diff < time.Hour {
			minutes := int(diff.Minutes())
			if minutes == 1 {
				return pongo2.AsValue("in 1 minute"), nil
			}
			return pongo2.AsValue(fmt.Sprintf("in %d minutes", minutes)), nil
		} else if diff < 24*time.Hour {
			hours := int(diff.Hours())
			if hours == 1 {
				return pongo2.AsValue("in 1 hour"), nil
			}
			return pongo2.AsValue(fmt.Sprintf("in %d hours", hours)), nil
		} else {
			days := int(diff.Hours() / 24)
			if days == 1 {
				return pongo2.AsValue("in 1 day"), nil
			}
			return pongo2.AsValue(fmt.Sprintf("in %d days", days)), nil
		}
	}

	if diff < time.Minute {
		return pongo2.AsValue("just now"), nil
	} else if diff < time.Hour {
		minutes := int(diff.Minutes())
		if minutes == 1 {
			return pongo2.AsValue("1 minute ago"), nil
		}
		return pongo2.AsValue(fmt.Sprintf("%d minutes ago", minutes)), nil
	} else if diff < 24*time.Hour {
		hours := int(diff.Hours())
		if hours == 1 {
			return pongo2.AsValue("1 hour ago"), nil
		}
		return pongo2.AsValue(fmt.Sprintf("%d hours ago", hours)), nil
	} else if diff < 7*24*time.Hour {
		days := int(diff.Hours() / 24)
		if days == 1 {
			return pongo2.AsValue("1 day ago"), nil
		}
		return pongo2.AsValue(fmt.Sprintf("%d days ago", days)), nil
	} else if diff < 30*24*time.Hour {
		weeks := int(diff.Hours() / (24 * 7))
		if weeks == 1 {
			return pongo2.AsValue("1 week ago"), nil
		}
		return pongo2.AsValue(fmt.Sprintf("%d weeks ago", weeks)), nil
	} else if diff < 365*24*time.Hour {
		months := int(diff.Hours() / (24 * 30))
		if months == 1 {
			return pongo2.AsValue("1 month ago"), nil
		}
		return pongo2.AsValue(fmt.Sprintf("%d months ago", months)), nil
	} else {
		years := int(diff.Hours() / (24 * 365))
		if years == 1 {
			return pongo2.AsValue("1 year ago"), nil
		}
		return pongo2.AsValue(fmt.Sprintf("%d years ago", years)), nil
	}
}
