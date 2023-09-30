package utils

// module compares package versions in different formats

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

// can cause bug with non supported versions like 1.1.1-beta-1
// equal false
// older false
// newer true
func IsNewerThan(newVersion string, oldVersion string) bool {

	isNewStandartNewVersion, _ := regexp.MatchString(".*\\.v.*", newVersion)
	isNewStandartOldVersion, _ := regexp.MatchString(".*\\.v.*", oldVersion)

	// equal
	if isNewStandartNewVersion && isNewStandartOldVersion && newVersion == oldVersion {
		return false
	}

	// newer version stadart grater than any old version
	if isNewStandartNewVersion && !isNewStandartOldVersion {
		return true
	}
	if isNewStandartOldVersion && !isNewStandartNewVersion {
		return false
	}

	// compare new style
	if isNewStandartNewVersion && isNewStandartOldVersion {
		if len(strings.Split(newVersion, ".")) > 2 && len(strings.Split(oldVersion, ".")) > 2 {
			return compareByIter(strings.Split(newVersion, ".v")[0], strings.Split(oldVersion, ".v")[0])
		}

		// might raise error?
		newVersionSplit, _ := strconv.Atoi(strings.Split(newVersion, ".")[0])
		oldVersionSplit, _ := strconv.Atoi(strings.Split(oldVersion, ".")[0])

		if newVersionSplit > oldVersionSplit {
			return true
		} else {
			return false
		}

	}
	// compare old style
	if !isNewStandartNewVersion && !isNewStandartOldVersion {
		return compareByIter(newVersion, oldVersion)
	}
	logrus.Panicf("reached unexpected case new: %s old:%s\n", newVersion, oldVersion)
	return false
}

func compareByIter(newVersion string, oldVersion string) bool {
	isNewSemVer, _ := regexp.MatchString(".*-.*", newVersion)
	isOldSemVer, _ := regexp.MatchString(".*-.*", oldVersion)

	newVersionSplit := strings.Split(newVersion, ".")
	oldVersionSplit := strings.Split(oldVersion, ".")

	for i, new := range newVersionSplit {

		old := oldVersionSplit[i]

		newInt, _ := strconv.Atoi(new)
		oldInt, _ := strconv.Atoi(old)
		if newInt > oldInt {
			return true
		}
		// 4.5.5-3.0
		if newInt < oldInt {
			return false
		}
		// new 1.1.1.1
		// old 1.1.1
		if len(newVersionSplit) > len(oldVersionSplit) {
			if len(oldVersionSplit)-1 == i {
				return true
			}
		}

	}

	if isNewSemVer && isOldSemVer {
		firstPartNewVersion := strings.Split(newVersion, "-")[0]
		firstPartOldVersion := strings.Split(oldVersion, "-")[0]

		if len(firstPartNewVersion) > len(firstPartOldVersion) {
			return true

		} else if firstPartNewVersion == firstPartOldVersion {
			if compareByIter(firstPartNewVersion, firstPartOldVersion) {
				return true
			} else {
				return compareByIter(strings.Split(newVersion, "-")[1], strings.Split(oldVersion, "-")[1])
			}
		} else {
			return compareByIter(firstPartNewVersion, firstPartOldVersion)
		}
	}

	// new 1.1.1
	// old 1.1.1.1
	// if len new > len old and all number the same rely on lenght

	if len(newVersionSplit) > len(oldVersionSplit) {
		return true
	} else {
		return false
	}
}
