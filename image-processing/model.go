package image_processing

type ImageResizerInput struct {
	OriginalFile     []byte
	OriginalFileName string
	Height           int
	Width            int
	Quality          int
}

type ImageResizerOutput struct {
	ImageContent []byte
	FileName     string
	Height       int
	Width        int
}

type processFileOutput struct {
	OutputFile []byte
	Height     uint
	Width      uint
}
