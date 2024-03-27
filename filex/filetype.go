package filex

const (
	Other    uint8  = 1
	Image    uint8  = 2
	Video    uint8  = 3
	Audio    uint8  = 4
	Model    uint8  = 5
	OtherStr string = "other"
	ImageStr string = "image"
	VideoStr string = "video"
	AudioStr string = "audio"
	ModelStr string = "model"
)

func fileTypeToUint8(fileType string) uint8 {
	switch fileType {
	case OtherStr:
		return Other
	case ImageStr:
		return Image
	case VideoStr:
		return Video
	case AudioStr:
		return Audio
	case ModelStr:
		return Model
	default:
		return Other
	}
}
