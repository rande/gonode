// Copyright © 2014-2016 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package security

func EnsureRoles(roles []string, rls ...string) []string {

	var found bool
	for _, rl := range rls {
		found = false

		for _, role := range roles {
			if role == rl {
				found = true

				break
			}
		}

		if found == false {
			roles = append(roles, rl)
		}
	}

	return roles
}
