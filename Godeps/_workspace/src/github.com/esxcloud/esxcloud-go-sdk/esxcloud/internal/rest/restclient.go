package rest

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"path/filepath"
)

type Request struct {
	Method      string
	URL         string
	ContentType string
	Body        io.Reader
}

const appJson string = "application/json"

func Get(client *http.Client, url string) (res *http.Response, err error) {
	req := Request{"GET", url, "", nil}
	res, err = Do(client, &req)
	return
}

func Post(client *http.Client, url string, body io.Reader) (res *http.Response, err error) {
	req := Request{"POST", url, appJson, body}
	res, err = Do(client, &req)
	return
}

func Delete(client *http.Client, url string) (res *http.Response, err error) {
	req := Request{"DELETE", url, "", nil}
	res, err = Do(client, &req)
	return
}

func Do(client *http.Client, req *Request) (res *http.Response, err error) {
	r, err := http.NewRequest(req.Method, req.URL, req.Body)
	if err != nil {
		return
	}
	if req.ContentType != "" {
		r.Header.Add("Content-Type", req.ContentType)
	}
	res, err = client.Do(r)
	return
}

func MultipartUploadFile(client *http.Client, url, filePath string, params map[string]string) (res *http.Response, err error) {
	file, err := os.Open(filePath)
	if err != nil {
		return
	}
	return MultipartUpload(client, url, file, filepath.Base(filePath), params)
}

func MultipartUpload(client *http.Client, url string, reader io.Reader, filename string, params map[string]string) (res *http.Response, err error) {
	buf := &bytes.Buffer{}
	writer := multipart.NewWriter(buf)
	if params != nil {
		for key, val := range params {
			writer.WriteField(key, val)
		}
	}
	header := textproto.MIMEHeader{}
	header.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s";`, "file", filename))
	header.Set("Content-Type", writer.FormDataContentType())
	part, err := writer.CreatePart(header)
	if err != nil {
		return
	}
	_, err = io.Copy(part, reader)
	if err != nil {
		return
	}
	err = writer.Close()
	if err != nil {
		return
	}
	res, err = Do(client, &Request{"POST", url, writer.FormDataContentType(), buf})
	return
}
