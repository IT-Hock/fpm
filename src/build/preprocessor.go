package build

import (
	"strconv"
	"strings"
)

var (
	Commit        = "unknown"
	BuildDate     = "unknown"
	GitBranch     = "unknown"
	GitDirty      = "unknown"
	GitCommitDate = "unknown"
	BuildVersion  = "unknown"

	GithubClientId     = ""
	GithubClientSecret = ""

	_Debug   = "false"
	_Version = "0.0.0"
)

var VersionMajor, _ = strconv.Atoi(strings.Split(_Version, ".")[0])
var VersionMinor, _ = strconv.Atoi(strings.Split(_Version, ".")[1])
var VersionPatch, _ = strconv.Atoi(strings.Split(_Version, ".")[2])

var VersionString = _Version

//goland:noinspection GoBoolExpressions
var Debug = _Debug == "1"
