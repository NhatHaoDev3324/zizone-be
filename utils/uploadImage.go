package utils

import (
	"context"
	"io"
	"mime/multipart"
	"sync"

	"github.com/NhatHaoDev3324/zizone-be/config"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

func UploadImage(file *multipart.FileHeader) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	return UploadImageFromReader(src)
}

func UploadImageFromReader(reader io.Reader) (string, error) {
	resp, err := config.Cloud.Upload.Upload(context.Background(), reader, uploader.UploadParams{
		Folder: "zizone",
		Format: "png",
	})
	if err != nil {
		return "", err
	}

	return resp.SecureURL, nil
}

func UploadMultipleImages(files []*multipart.FileHeader) ([]string, error) {
	if len(files) == 0 {
		return []string{}, nil
	}

	var wg sync.WaitGroup
	urls := make([]string, len(files))
	errs := make([]error, len(files))

	semaphore := make(chan struct{}, 5)

	for i, file := range files {
		wg.Add(1)
		go func(i int, f *multipart.FileHeader) {
			defer wg.Done()

			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			url, err := UploadImage(f)
			if err != nil {
				errs[i] = err
				return
			}
			urls[i] = url
		}(i, file)
	}

	wg.Wait()

	for _, err := range errs {
		if err != nil {
			return nil, err
		}
	}

	return urls, nil
}
