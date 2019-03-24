package s3

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/yansal/gallery/storage"
)

type Storage struct {
	s3     *s3.S3
	bucket string
}

func New(bucket string) (*Storage, error) {
	s, err := session.NewSession()
	if err != nil {
		return nil, err
	}
	return &Storage{
		bucket: bucket,
		s3:     s3.New(s),
	}, nil
}

func (s *Storage) List(prefix string) (res storage.ListResult, err error) {
	in := &s3.ListObjectsInput{
		Bucket:    aws.String(s.bucket),
		Delimiter: aws.String("/"),
		Prefix:    aws.String(prefix),
	}
	out, err := s.s3.ListObjects(in)
	if err != nil {
		return res, err
	}

	for i := range out.CommonPrefixes {
		res.Results = append(res.Results, *out.CommonPrefixes[i].Prefix)
	}
	for i := range out.Contents {
		res.Results = append(res.Results, *out.Contents[i].Key)
	}
	return res, nil
}
