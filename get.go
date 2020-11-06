package main

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/pkg/errors"
)

// parseURI returns bucket and key from URIs like:
// - s3://bucket-name/dir
// - s3://bucket-name/dir/file.ext
func parseURI(uri string) (bucket, key string, err error) {
	if !strings.HasPrefix(uri, "s3://") {
		return "", "", fmt.Errorf("uri %s protocol is not s3", uri)
	}

	u, err := url.Parse(uri)
	if err != nil {
		return "", "", errors.Wrapf(err, "parse uri %s", uri)
	}

	bucket, key = u.Host, strings.TrimPrefix(u.Path, "/")
	return bucket, key, nil
}

// FetchRaw downloads the object from URI and returns it in the form of byte slice.
// uri must be in the form of s3 protocol: s3://bucket-name/key[...].
func FetchRaw(sess *session.Session, ctx context.Context, uri string) ([]byte, error) {
	bucket, key, err := parseURI(uri)
	if err != nil {
		return nil, err
	}

	buf := &aws.WriteAtBuffer{}
	_, err = s3manager.NewDownloader(sess).DownloadWithContext(
		ctx,
		buf,
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		})
	if err != nil {
		return nil, errors.Wrap(err, "fetch object from s3")
	}

	return buf.Bytes(), nil
}

func main() {
	//sess := session.New()
	so := session.Options{
		Config: aws.Config{
			S3ForcePathStyle: aws.Bool(true),
			Endpoint:         aws.String(os.Getenv("AWS_ENDPOINT")),
		},
		SharedConfigState: session.SharedConfigEnable,
	}
	sess, err := session.NewSessionWithOptions(so)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	/*
	   svc := s3.New(sess)

	   input := &s3.GetObjectInput{
	           Bucket: aws.String("bin-us-east-1.dev-cloudinfra.intelerad.com"),
	           Key:    aws.String("charts/ci/index.html"),
	   }

	   result, err := svc.GetObject(input)
	   if err != nil {
	           if aerr, ok := err.(awserr.Error); ok {
	                   switch aerr.Code() {
	                   case s3.ErrCodeNoSuchKey:
	                           fmt.Println(s3.ErrCodeNoSuchKey, aerr.Error())
	                   default:
	                           fmt.Println(aerr.Error())
	                   }
	           } else {
	                   // Print the error, cast err to awserr.Error to get the Code and
	                   // Message from an error.
	                   fmt.Println(err.Error())
	           }
	   }
	   fmt.Println(result)
	*/

	b, err := FetchRaw(sess, context.TODO(), "s3://bin-us-east-1.dev-cloudinfra.intelerad.com/charts/ci/index.yaml")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	fmt.Printf("%s err=%v\n", b, err)
}
