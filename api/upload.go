package api

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/labstack/echo"
)

func Upload(c echo.Context) error {
	//-----------
	// Read file
	//-----------

	// Source
	file, err := c.FormFile("file")
	if err != nil {
		return err
	}
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Destination
	dst, err := os.Create(file.Filename)
	if err != nil {
		return err
	}
	defer dst.Close()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, fmt.Sprintf("file %s uploaded successfully", file.Filename))
}
