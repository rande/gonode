// Copyright © 2014-2016 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package security

type AuthorizationChecker interface {
	IsGranted(attrs Attributes, o interface{}) (bool, error)
}

type DefaultAuthorizationChecker struct {
	DecisionManager DecisionVoter
}

func (c *DefaultAuthorizationChecker) IsGranted(t SecurityToken, attrs Attributes, o interface{}) (bool, error) {

	if c.DecisionManager == nil {
		return false, nil
	}

	return c.DecisionManager.Decide(t, attrs, o), nil
}
