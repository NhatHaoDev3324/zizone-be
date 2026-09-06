package utils

import (
	"bytes"
	"context"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"
	"sync"

	"github.com/NhatHaoDev3324/zizone-be/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/deepteams/webp"
	"github.com/disintegration/imaging"
	"github.com/google/uuid"
)

var allowedImageExts = map[string]string{
	".jpg":  "image/jpeg",
	".jpeg": "image/jpeg",
	".png":  "image/png",
	".gif":  "image/gif",
	".webp": "image/webp",
}

var allowedVideoExts = map[string]string{
	".mp4":  "video/mp4",
	".mov":  "video/quicktime",
	".avi":  "video/x-msvideo",
	".mkv":  "video/x-matroska",
	".webm": "video/webm",
}

// webpQuality là chất lượng nén webp (0-100). 80 cân bằng tốt giữa dung lượng và chất lượng.
const webpQuality = 80

// maxImageDimension là kích thước tối đa (px) của cạnh dài nhất.
// Ảnh lớn hơn sẽ được thu nhỏ để giảm dung lượng tối đa khi upload/tải.
const maxImageDimension = 1920

func UploadR2Image(file *multipart.FileHeader, userID uuid.UUID) (string, error) {
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if _, ok := allowedImageExts[ext]; !ok {
		return "", fmt.Errorf("định dạng ảnh không được hỗ trợ: %s", ext)
	}

	// Chuyển đổi ảnh sang webp trước khi upload lên R2
	webpBytes, err := convertImageToWebp(file)
	if err != nil {
		return "", err
	}

	key := fmt.Sprintf("images/%s/%s.webp", userID, uuid.New().String())
	return uploadToR2(bytes.NewReader(webpBytes), key, "image/webp")
}

// convertImageToWebp đọc file ảnh, decode và encode lại thành định dạng webp.
func convertImageToWebp(file *multipart.FileHeader) ([]byte, error) {
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("không thể mở file: %w", err)
	}
	defer src.Close()

	img, _, err := image.Decode(src)
	if err != nil {
		return nil, fmt.Errorf("không thể đọc ảnh: %w", err)
	}

	// Thu nhỏ nếu ảnh quá lớn (giữ tỉ lệ, cạnh dài nhất = maxImageDimension).
	// Đây là bước tiết kiệm dung lượng hiệu quả nhất.
	if b := img.Bounds(); b.Dx() > maxImageDimension || b.Dy() > maxImageDimension {
		img = imaging.Fit(img, maxImageDimension, maxImageDimension, imaging.Lanczos)
	}

	var buf bytes.Buffer
	if err := webp.Encode(&buf, img, &webp.EncoderOptions{Quality: float32(webpQuality), Method: 4}); err != nil {
		return nil, fmt.Errorf("không thể chuyển đổi sang webp: %w", err)
	}

	return buf.Bytes(), nil
}

const maxConcurrentUploads = 20

func UploadMultipleR2Images(files []*multipart.FileHeader, userID uuid.UUID) ([]string, error) {
	type result struct {
		index int
		url   string
		err   error
	}

	numFiles := len(files)
	urls := make([]string, numFiles)
	results := make(chan result, numFiles)
	sem := make(chan struct{}, maxConcurrentUploads)
	var wg sync.WaitGroup

	for i, file := range files {
		wg.Add(1)
		go func(i int, file *multipart.FileHeader) {
			defer wg.Done()
			sem <- struct{}{}        // chiếm slot
			defer func() { <-sem }() // giải phóng slot

			url, err := UploadR2Image(file, userID)
			results <- result{index: i, url: url, err: err}
		}(i, file)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	var errs []error
	for res := range results {
		if res.err != nil {
			errs = append(errs, fmt.Errorf("file %d: %w", res.index, res.err))
		} else {
			urls[res.index] = res.url
		}
	}

	if len(errs) > 0 {
		return urls, fmt.Errorf("lỗi upload: %w", joinErrors(errs))
	}

	return urls, nil
}

func joinErrors(errs []error) error {
	if len(errs) == 0 {
		return nil
	}
	msgs := make([]string, len(errs))
	for i, e := range errs {
		msgs[i] = e.Error()
	}
	return fmt.Errorf("%s", strings.Join(msgs, "; "))
}

func UploadR2Video(file *multipart.FileHeader, userID uuid.UUID) (string, error) {
	ext := strings.ToLower(filepath.Ext(file.Filename))
	contentType, ok := allowedVideoExts[ext]
	if !ok {
		return "", fmt.Errorf("định dạng video không được hỗ trợ: %s", ext)
	}

	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("không thể mở file: %w", err)
	}
	defer src.Close()

	key := fmt.Sprintf("videos/%s/%s%s", userID, uuid.New().String(), ext)
	return uploadToR2(src, key, contentType)
}

func uploadToR2(body io.Reader, key, contentType string) (string, error) {
	bucket := config.GetR2BucketName()

	_, err := config.R2Client.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(key),
		Body:        body,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", fmt.Errorf("không thể upload file lên R2: %w", err)
	}

	publicURL := config.GetR2PublicURL()
	return fmt.Sprintf("%s/%s", publicURL, key), nil
}
