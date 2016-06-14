package main

type Picture struct {
	Path  string `json:"path"`
	Width int
}

type PictureMap map[string]Picture
