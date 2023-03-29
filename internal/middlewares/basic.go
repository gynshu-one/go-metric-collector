package middlewares

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"fmt"
	"github.com/andybalholm/brotli"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

func MiscDecompress() gin.HandlerFunc {
	return func(c *gin.Context) {
		encoding := c.Request.Header.Get("Content-Encoding")
		if encoding != "" {
			switch encoding {
			//case "gzip":
			//	c.Request = deCompressGzip(c.Request)
			case "deflate":
				c.Request = deCompressDeflate(c.Request)
			case "br":
				c.Request = deCompressBr(c.Request)
			}
		}
		c.Next()
	}
}

func deCompressDeflate(r *http.Request) *http.Request {
	// Replace body bytes with compressed bytes
	bt := new(bytes.Buffer)
	_, err := bt.ReadFrom(r.Body)
	if err != nil {
		return nil
	}
	dc, err := decompressDeflate(bt.Bytes())
	if err != nil {
		return nil
	}
	r.Body = io.NopCloser(bytes.NewReader(dc))
	return r
}
func deCompressGzip(r *http.Request) *http.Request {
	// Replace body bytes with compressed bytes
	bt := new(bytes.Buffer)
	_, err := bt.ReadFrom(r.Body)
	if err != nil {
		return nil
	}
	dc, err := decompressGzip(bt.Bytes())
	if err != nil {
		return nil
	}
	r.Body = io.NopCloser(bytes.NewReader(dc))
	return r
}

func deCompressBr(r *http.Request) *http.Request {
	// Replace body bytes with compressed bytes
	bt := new(bytes.Buffer)
	_, err := bt.ReadFrom(r.Body)
	if err != nil {
		return nil
	}
	dc, err := decompressBr(bt.Bytes())
	if err != nil {
		return nil
	}
	r.Body = io.NopCloser(bytes.NewReader(dc))
	return r
}
func decompressBr(data []byte) ([]byte, error) {
	r := brotli.NewReader(bytes.NewReader(data))
	var b bytes.Buffer
	_, err := b.ReadFrom(r)
	if err != nil {
		return nil, fmt.Errorf("failed decompress data: %v", err)
	}

	return b.Bytes(), nil
}
func decompressDeflate(data []byte) ([]byte, error) {
	r := flate.NewReader(bytes.NewReader(data))
	defer r.Close()

	var b bytes.Buffer
	_, err := b.ReadFrom(r)
	if err != nil {
		return nil, fmt.Errorf("failed decompress data: %v", err)
	}

	return b.Bytes(), nil
}

func decompressGzip(data []byte) ([]byte, error) {
	r, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed init decompress reader: %v", err)
	}
	defer r.Close()

	var b bytes.Buffer
	_, err = b.ReadFrom(r)
	if err != nil {
		return nil, fmt.Errorf("failed decompress data: %v", err)
	}

	return b.Bytes(), nil
}
func compressDeflate(data []byte) ([]byte, error) {
	var b bytes.Buffer
	w, err := flate.NewWriter(&b, flate.BestCompression)
	if err != nil {
		return nil, fmt.Errorf("failed init compress writer: %v", err)
	}
	_, err = w.Write(data)
	if err != nil {
		return nil, fmt.Errorf("failed write data to compress temporary buffer: %v", err)
	}
	err = w.Close()
	if err != nil {
		return nil, fmt.Errorf("failed compress data: %v", err)
	}
	return b.Bytes(), nil
}
