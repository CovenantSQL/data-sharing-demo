package api

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"io"
	"net/http"
	"path/filepath"

	"github.com/Sirupsen/logrus"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/labstack/echo"
)

const (
	bucket = "ds-demo"
	region = "cn-northwest-1"
)

func Upload(c echo.Context) (err error) {
	token := c.Request().Header.Get("Authorization")
	if token == "" {
		logrus.Error("Upload: no token got")
		return c.JSON(http.StatusBadRequest, "no token got")
	} else {
		db := getDB(token)
		if db == nil {
			return c.JSON(http.StatusForbidden, "invalid token")
		}
	}

	// Source
	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	defer src.Close()

	buff := bytes.NewBuffer(make([]byte, 0, file.Size))
	_, err = io.Copy(buff, src)
	if err != nil {
		logrus.Errorf("read upload file failed: %v", err)
		return c.JSON(http.StatusBadRequest, err)
	}
	checkSum := md5.Sum(buff.Bytes())

	//select Region to use.
	conf := aws.Config{Region: aws.String(region)}
	sess, err := session.NewSession(&conf)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, err)
	}
	svc := s3manager.NewUploader(sess)

	logrus.Debug("Uploading file to S3")
	result, err := svc.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(file.Filename),
		Body:   buff,
	})
	if err != nil {
		logrus.Errorf("upload to s3 failed: %v", err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	logrus.Debugf("Successfully uploaded %s to %s\n", file.Filename, result.Location)

	return c.JSON(http.StatusOK, H{
		"attach_uri": filepath.Base(result.Location),
		"attach_sum": hex.EncodeToString(checkSum[:]),
	})
}

func Download(c echo.Context) (err error) {
	//token := c.Request().Header.Get("Authorization")
	//if token == "" {
	//	logrus.Error("Upload: no token got")
	//	return c.JSON(http.StatusBadRequest, "no token got")
	//} else {
	//	db := getDB(token)
	//	if db == nil {
	//		return c.JSON(http.StatusForbidden, "invalid token")
	//	}
	//}

	conf := aws.Config{Region: aws.String(region)}
	sess, err := session.NewSession(&conf)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, err)
	}
	downloader := s3manager.NewDownloader(sess)
	buff := &aws.WriteAtBuffer{}
	objName := c.Param("file")
	numBytes, err := downloader.Download(buff,
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(objName),
		})
	if err != nil {
		logrus.Errorf("Unable to download item %q, %v", objName, err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	c.Response().WriteHeader(http.StatusOK)
	numBytes, err = io.Copy(c.Response().Writer, bytes.NewReader(buff.Bytes()))

	logrus.Debugf("Downloaded %s from S3: %d", objName, numBytes)

	return
}
