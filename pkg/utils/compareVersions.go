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

	logrus.Debugf("check version compare new: %s old: %s\n", newVersion, oldVersion)
	logrus.Debugf("isNewStandartNewVersion %t,  isNewStandartOldVersion: %t\n", isNewStandartNewVersion, isNewStandartOldVersion)
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

func compareByIter(new string, old string) bool {
	isNewSemVer, _ := regexp.MatchString(".*-.*", new)
	isOldSemVer, _ := regexp.MatchString(".*-.*", old)

	if isNewSemVer && isOldSemVer {
		logrus.Debugln("checking semver")

		if len(strings.Split(new, "-")[0]) > len(strings.Split(old, "-")[0]) {
			return true

		} else if strings.Split(new, "-")[0] == strings.Split(old, "-")[0] {
			if compareByIter(strings.Split(new, "-")[0], strings.Split(old, "-")[0]) {
				return true
			} else {
				return compareByIter(strings.Split(new, "-")[1], strings.Split(old, "-")[1])
			}
		} else {
			return compareByIter(strings.Split(new, "-")[0], strings.Split(old, "-")[0])
		}
	}
	newVersionSplit := strings.Split(new, ".")
	oldVersionSplit := strings.Split(old, ".")

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
	// new 1.1.1
	// old 1.1.1.1
	// if len new > len old and all number the same rely on lenght

	if len(newVersionSplit) > len(oldVersionSplit) {
		return true
	} else {
		return false
	}
}
