package minio

import (
	"context"
	"fmt"
	"mime/multipart"
	"strings"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioStorage struct {
	Client     *minio.Client
	BucketName string
	BaseURL    string
}

func NewMinioStorage(endpoint, accessKey, secretKey, bucketName, baseURL string, useSSL bool) (*MinioStorage, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}

	exists, err := client.BucketExists(context.Background(), bucketName)
	if err != nil {
		return nil, err
	}
	if !exists {
		if err := client.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{}); err != nil {
			return nil, err
		}
	}

	return &MinioStorage{
		Client:     client,
		BucketName: bucketName,
		BaseURL:    strings.TrimRight(baseURL, "/"),
	}, nil
}

func (m *MinioStorage) AddFiles(ctx context.Context, userID string, files []*multipart.FileHeader) ([]string, error) {
	urls := make([]string, 0, len(files))

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			return nil, err
		}
		defer file.Close()

		objectName := fmt.Sprintf("%s/%s-%s", userID, uuid.New().String(), fileHeader.Filename)

		_, err = m.Client.PutObject(ctx, m.BucketName, objectName, file, fileHeader.Size, minio.PutObjectOptions{
			ContentType: fileHeader.Header.Get("Content-Type"),
		})
		if err != nil {
			return nil, err
		}

		urls = append(urls, fmt.Sprintf("%s/%s/%s", m.BaseURL, m.BucketName, objectName))
	}

	return urls, nil
}

func (m *MinioStorage) DeleteUserFiles(ctx context.Context, userID string) error {
	objectsCh := make(chan minio.ObjectInfo)

	go func() {
		defer close(objectsCh)
		for object := range m.Client.ListObjects(ctx, m.BucketName, minio.ListObjectsOptions{
			Prefix:    userID + "/",
			Recursive: true,
		}) {
			if object.Err != nil {
				continue
			}
			objectsCh <- object
		}
	}()

	for object := range objectsCh {
		if err := m.Client.RemoveObject(ctx, m.BucketName, object.Key, minio.RemoveObjectOptions{}); err != nil {
			return err
		}
	}

	return nil
}

func (m *MinioStorage) ReplaceUserFiles(ctx context.Context, userID string, files []*multipart.FileHeader) ([]string, error) {
	if err := m.DeleteUserFiles(ctx, userID); err != nil {
		return nil, err
	}
	return m.AddFiles(ctx, userID, files)
}

func (m *MinioStorage) GetUserFiles(ctx context.Context, userID string) ([]string, error) {
	urls := []string{}

	for object := range m.Client.ListObjects(ctx, m.BucketName, minio.ListObjectsOptions{
		Prefix:    userID + "/",
		Recursive: true,
	}) {
		if object.Err != nil {
			return nil, object.Err
		}
		urls = append(urls, fmt.Sprintf("%s/%s/%s", m.BaseURL, m.BucketName, object.Key))
	}

	return urls, nil
}
