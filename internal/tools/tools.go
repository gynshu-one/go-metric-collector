package tools

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

func Contains(sl []string, s string) bool {
	for _, v := range sl {
		if v == s {
			return true
		}
	}
	return false
}
func Int64Ptr(i int64) *int64 {
	return &i
}
func Float64Ptr(f float64) *float64 {
	return &f
}

func GetProjectRoot() string {
	cwd, err := os.Getwd()
	if err != nil {
		return ""
	}

	for {
		if _, err = os.Stat(filepath.Join(cwd, "go.mod")); err == nil {
			return cwd
		}
		if cwd == "/" {
			break
		}
		cwd = filepath.Dir(cwd)
	}
	panic(fmt.Errorf("error getting project root"))
}
func ParseDuration(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		log.Fatalf("Error parsing duration %q: %s", s, err)
	}
	return d
}
