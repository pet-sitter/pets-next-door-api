package bucketinfra

import (
	"io"
	"path/filepath"

	"github.com/aws/aws-sdk-go/private/protocol/rest"
	"github.com/google/uuid"
	pnd "github.com/pet-sitter/pets-next-door-api/api"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/rs/zerolog/log"
)

type FileUploader interface {
	UploadFile(file io.ReadSeeker, fileName string) (url string, appError *pnd.AppError)
}

type S3Client struct {
	s3Client   *s3.S3
	bucketName string
}

func NewS3Client(keyID, key, endpoint, region, bucketName string) (*S3Client, error) {
	s3Config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(keyID, key, ""),
		Endpoint:         aws.String(endpoint),
		Region:           aws.String(region),
		S3ForcePathStyle: aws.Bool(true),
	}
	newSession, err := session.NewSession(s3Config)
	if err != nil {
		return nil, err
	}

	s3Client := s3.New(newSession)

	return &S3Client{
		s3Client:   s3Client,
		bucketName: bucketName,
	}, nil
}

func (c *S3Client) UploadFile(file io.ReadSeeker, fileName string) (string, *pnd.AppError) {
	randomFileName := generateRandomFileName(fileName)
	fullPath := "media/" + randomFileName

	if _, err := c.uploadToS3(file, fullPath); err != nil {
		return "", pnd.ErrUnknown(err)
	}

	req, _ := c.GetFileRequest(fullPath)
	rest.Build(req)
	if err := req.Send(); err != nil {
		return "", pnd.ErrUnknown(err)
	}

	return req.HTTPRequest.URL.String(), nil
}

func (c *S3Client) uploadToS3(file io.ReadSeeker, fullPath string) (*s3.PutObjectOutput, error) {
	result, err := c.s3Client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(c.bucketName),
		Key:    aws.String(fullPath),
		Body:   file,
	})
	if err != nil {
		log.Error().Err(err).Msg("error uploading file")
		return nil, err
	}
	return result, nil
}

func (c *S3Client) GetFileRequest(
	fileName string,
) (req *request.Request, output *s3.GetObjectOutput) {
	return c.s3Client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(c.bucketName),
		Key:    aws.String(fileName),
	})
}

func generateRandomFileName(originalFileName string) string {
	extension := filepath.Ext(originalFileName)
	return uuid.New().String() + extension
}
