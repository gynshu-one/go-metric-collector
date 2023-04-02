package middlewares

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"fmt"
	"github.com/andybalholm/brotli"
	"github.com/gin-gonic/gin"
	"github.com/gynshu-one/go-metric-collector/internal/tools"
	"io"
	"net/http"
	"strings"
)

func MiscDecompress() gin.HandlerFunc {
	return func(c *gin.Context) {
		encoding := c.Request.Header.Get("Content-Encoding")
		encodings := strings.Split(encoding, ",")
		if len(encodings) > 0 {
			if tools.Contains(encodings, "gzip") {
				c.Request = deCompressGzip(c.Request)
			} else if tools.Contains(encodings, "deflate") {
				c.Request = deCompressDeflate(c.Request)
			} else if tools.Contains(encodings, "br") {
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
