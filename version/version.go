// Copyright 2017 Joan Llopis. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// version package holds the latest version information following the
// semver spec (http://semver.org)
package version

import (
	"fmt"

	"github.com/coreos/go-semver/semver"
)

var (
	Name       = "unknown"
	Version    = "0.0.1"
	APIVersion = "unknown"
)

func init() {
	v, err := semver.NewVersion(Version)
	if err == nil {
		APIVersion = fmt.Sprintf("%d.%d", v.Major, v.Minor)
	}
}
