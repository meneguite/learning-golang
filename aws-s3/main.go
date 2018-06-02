package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}

func getNewSession() (*session.Session, error) {
	return session.NewSession(&aws.Config{
		Region: aws.String("sa-east-1"),
	})
}

func getListForBucket(ss *session.Session, bucket string) (*s3.ListObjectsOutput, error) {
	svc := s3.New(ss)
	return svc.ListObjects(&s3.ListObjectsInput{
		Bucket: aws.String(bucket),
	})
}

func downloadAllFiles(ss *session.Session, bucket string, list *s3.ListObjectsOutput) {
	var wg sync.WaitGroup
	wg.Add(len(list.Contents))

	downloader := getDownloader(ss)
	for _, item := range list.Contents {
		fmt.Println("Downloading file: ", *item.Key)
		go donwloadFile(&wg, downloader, bucket, *item.Key, "downloads")
	}
	wg.Wait()
}

func getDownloader(ss *session.Session) *s3manager.Downloader {
	return s3manager.NewDownloader(ss)
}

func donwloadFile(wg *sync.WaitGroup, d *s3manager.Downloader, bucket string, fileName string, path string) {
	defer wg.Done()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, 0666)
	}

	//time.Sleep(2 * time.Second)

	filePath := filepath.Join(path, fileName)
	file, err := os.Create(filePath)
	if err != nil {
		exitErrorf("Error on download files %q, %v", bucket, err)
	}

	defer file.Close()

	_, err = d.Download(file, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fileName),
	})
	if err != nil {
		exitErrorf("Error on download files %q, %v", bucket, err)
	}
}

func uploadAllFiles(ss *session.Session, bucket string) {
	files, err := ioutil.ReadDir("uploads")
	if err != nil {
		exitErrorf("Error on uploading files %q, %v", bucket, err)
	}

	var wg sync.WaitGroup
	wg.Add(len(files))

	uploader := getUploader(ss)
	for _, file := range files {
		fmt.Println("Uploading file: ", file.Name())
		go uploadFile(&wg, uploader, bucket, file.Name(), "uploads")
	}
	wg.Wait()
}

func getUploader(ss *session.Session) *s3manager.Uploader {
	return s3manager.NewUploader(ss)
}

func uploadFile(wg *sync.WaitGroup, u *s3manager.Uploader, bucket string, fileName string, path string) {
	defer wg.Done()
	filePath := filepath.Join(path, fileName)
	file, err := os.Open(filePath)
	if err != nil {
		exitErrorf("Error on uploading files %q, %v", bucket, err)
	}
	defer file.Close()

	//time.Sleep(2 * time.Second)

	_, err = u.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fileName),
		Body:   file,
	})
	if err != nil {
		exitErrorf("Error on uploading files %q, %v", bucket, err)
	}
}

func main() {
	if len(os.Args) != 2 {
		exitErrorf("Bucket name required")
	}
	bucket := os.Args[1]

	ss, err := getNewSession()
	if err != nil {
		exitErrorf("Erro ao abir a sessão", err)
	}

	list, err := getListForBucket(ss, bucket)
	if err != nil {
		exitErrorf("Unable to list items in bucket %q, %v", bucket, err)
	}

	downloadAllFiles(ss, bucket, list)
	uploadAllFiles(ss, bucket)
}
