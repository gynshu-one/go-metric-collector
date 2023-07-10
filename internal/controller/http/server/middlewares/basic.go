// Package middlewares contains the basic compression and decompression middlewares
package middlewares

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"github.com/andybalholm/brotli"
	"github.com/gin-gonic/gin"
	"github.com/gynshu-one/go-metric-collector/internal/tools"
	"github.com/rs/zerolog/log"
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
				log.Debug().Msg("Decompressing gzip")
				c.Request = deCompressGzip(c.Request)
			} else if tools.Contains(encodings, "deflate") {
				log.Debug().Msg("Decompressing deflate")
				c.Request = deCompressDeflate(c.Request)
			} else if tools.Contains(encodings, "br") {
				log.Debug().Msg("Decompressing br")
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
		log.Trace().Msgf("Failed to read body Deflate: %v", err)
		return nil
	}
	dc, err := decompressDeflate(bt.Bytes())
	if err != nil {
		log.Trace().Msgf("Failed to decompress Deflate: %v", err)
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
		log.Trace().Msgf("Failed to read body Gzip: %v", err)
		return nil
	}
	dc, err := decompressGzip(bt.Bytes())
	if err != nil {
		log.Trace().Msgf("Failed to decompress Gzip: %v", err)
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
		log.Trace().Msgf("Failed to read body Br: %v", err)
		return nil
	}
	dc, err := decompressBr(bt.Bytes())
	if err != nil {
		log.Trace().Msgf("Failed to decompress Br: %v", err)
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
		return nil, err
	}

	return b.Bytes(), nil
}
func decompressDeflate(data []byte) ([]byte, error) {
	r := flate.NewReader(bytes.NewReader(data))
	defer func() {
		err := r.Close()
		if err != nil {
			log.Trace().Msgf("Failed to close deflate reader: %v", err)
		}
	}()

	var b bytes.Buffer
	_, err := b.ReadFrom(r)
	if err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

func decompressGzip(data []byte) ([]byte, error) {
	r, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer func() {
		err = r.Close()
		if err != nil {
			log.Trace().Msgf("Failed to close gzip reader: %v", err)
		}
	}()

	var b bytes.Buffer
	_, err = b.ReadFrom(r)
	if err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}
