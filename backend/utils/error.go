package utils

import "github.com/mymi14s/goconda/models"

func LogError(error map[string]any) {

	errLog := &models.ErrorLog{}
	_, err := models.BaseModel{}.Create(errLog, map[string]any{
		"title":   error["title"],
		"error":   error["error"],
		"context": error["context"],
	})
	if err != nil {
		// handle error
	}
}
