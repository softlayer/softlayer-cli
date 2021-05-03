package version

import "strconv"

var major, minor, build string

// plugin version
const PLUGIN_VERSION = "0.0.1"

// plugin major version
var PLUGIN_MAJOR_VERSION = convert(major)

// plugin minor version
var PLUGIN_MINOR_VERSION = convert(minor)

// plugin build version
var PLUGIN_BUILD_VERSION = convert(build)

func convert(numstring string) int {
	if i, err := strconv.Atoi(numstring); err == nil {
		return i
	} else {
		return 0
	}
}
