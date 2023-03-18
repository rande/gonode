// Copyright © 2014-2023 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package security

type AuthorizationChecker interface {
	IsGranted(t SecurityToken, attrs Attributes, o interface{}) (bool, error)
}

type DefaultAuthorizationChecker struct {
	DecisionVoter DecisionVoter
}

func (c *DefaultAuthorizationChecker) IsGranted(t SecurityToken, attrs Attributes, o interface{}) (bool, error) {
	if c.DecisionVoter == nil {
		return false, nil
	}

	return c.DecisionVoter.Decide(t, attrs, o), nil
}
