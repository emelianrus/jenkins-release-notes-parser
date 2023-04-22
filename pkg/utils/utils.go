package utils

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetLevel(logrus.DebugLevel)
	// TODO: enable for windows only
	logrus.SetFormatter(&logrus.TextFormatter{ForceColors: true})
	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	logrus.SetOutput(os.Stdout)
}

// function to check if file exists
func IsFileExist(filePath string) bool {
	logrus.Debugf("Checking is file exist: %s", filePath)
	_, err := os.Stat(filePath)

	if os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

func IsUrl(line string) bool {
	if strings.HasPrefix(line, "http://") || strings.HasPrefix(line, "https://") {
		logrus.Debugf("Line: '%s' is URL", line)
		return true
	} else {
		logrus.Debugf("Line: '%s' is not URL", line)
		return false
	}
}

func DoRequestGet(url string) *http.Response {
	var client = &http.Client{Timeout: 100 * time.Second}
	// TODO: DRY
	logrus.Debugf("Downloading: %s \n", url)
	response, err := client.Get(url)
	if err != nil || response.StatusCode != 200 {
		fmt.Println("download error")
		logrus.Error(err)
		os.Exit(1)
	}
	return response
}

func GetNextVersion(versions []string, currentVersion string) (string, error) {
	sort.Slice(versions, func(i, j int) bool {
		return !IsNewerThan(versions[i], versions[j])
	})

	vlen := len(versions)
	for idx, value := range versions {
		if idx+2 > vlen {
			logrus.Warn("end of versions list reached")
			return "", errors.New("end of the slice reached, assume latest version")
		}
		if value == currentVersion {
			return versions[idx+1], nil
		}
	}
	// NOTE: should not reach here
	return "", errors.New("no version found, should not reach this return")
}