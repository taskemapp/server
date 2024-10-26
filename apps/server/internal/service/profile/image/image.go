package image

type Processing interface {
	// ConvertToWebp from any valid img format
	ConvertToWebp(b []byte) error
}
