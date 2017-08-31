package old

type Picture struct {
	Path string `json:"path"`
	Type string `json:"type"`
}

type PictureMap map[string]Picture
