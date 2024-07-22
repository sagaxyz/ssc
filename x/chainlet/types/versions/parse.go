package versions

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

const MaxUint16 = ^uint16(0)

var semverRegexp = regexp.MustCompile(`^(?P<major>0|[1-9]\d*)\.(?P<minor>0|[1-9]\d*)\.(?P<patch>0|[1-9]\d*)(?:-(?P<prerelease>(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+(?P<buildmetadata>[0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$`)

func Check(ver string) bool {
	return semverRegexp.MatchString(ver)
}

func CheckUpgrade(old, new string) (major bool, err error) {
	oldMajor, oldMinor, _, _, err := Parse(old)
	if err != nil {
		err = fmt.Errorf("invalid version string '%s': %w", old, err)
		return
	}
	newMajor, newMinor, _, _, err := Parse(new)
	if err != nil {
		err = fmt.Errorf("invalid version string '%s': %w", new, err)
		return
	}

	// minor part acts as the major part when major is 0
	if oldMajor == 0 && newMajor == 0 {
		newMajor = newMinor
		oldMajor = oldMinor
	}
	if newMajor != oldMajor {
		if newMajor != oldMajor+1 {
			err = errors.New("major upgrades have to be increments of 1")
			return
		}
		major = true
	}

	return
}

func convertUint16(str string) (num uint16, err error) {
	i, err := strconv.Atoi(str)
	if err != nil {
		return
	}
	if i > int(MaxUint16) {
		err = errors.New("uint16 overflow")
		return
	}

	num = uint16(i)
	return
}

func Parse(ver string) (major, minor, patch uint16, suffix string, err error) {
	match := semverRegexp.FindStringSubmatch(ver)
	if match == nil || len(match) < 4 {
		err = errors.New("no regexp match")
		return
	}

	major, err = convertUint16(match[1])
	if err != nil {
		return
	}
	minor, err = convertUint16(match[2])
	if err != nil {
		return
	}
	patch, err = convertUint16(match[3])
	if err != nil {
		return
	}

	suffix = match[4]
	return
}
