package format

import (
	"bytes"
	"fmt"
	"image"
	_ "image/gif"
	"image/jpeg"
	"image/png"
	"strings"

	"golang.org/x/image/webp"
)

func init() {
	image.RegisterFormat("webp", "RIFF????WEBP", webp.Decode, webp.DecodeConfig)
}

func FileSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	val := float64(size) / float64(div)
	unitStr := "KMGTPE"[exp : exp+1]
	if val == float64(int64(val)) {
		return fmt.Sprintf("%d %sB", int64(val), unitStr)
	}
	return fmt.Sprintf("%.2f %sB", val, unitStr)
}

func RemoveExtension(fileName string) string {
	if fileName == "" {
		return fileName
	}
	parts := strings.Split(fileName, ".")
	if len(parts) <= 1 {
		return fileName
	}
	return strings.Join(parts[:len(parts)-1], ".")
}

func DecodeImage(imgData []byte) (image.Image, string, error) {
	img, formatName, err := image.Decode(bytes.NewReader(imgData))
	return img, formatName, err
}

func GetImageSizeAndData(img image.Image) (int64, []byte, error) {
	var buf bytes.Buffer
	switch img.(type) {
	case *image.NRGBA, *image.RGBA, *image.YCbCr:
		err := png.Encode(&buf, img)
		if err != nil {
			return 0, nil, err
		}
	default:
		err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: 85})
		if err != nil {
			return 0, nil, err
		}
	}
	return int64(buf.Len()), buf.Bytes(), nil
}
