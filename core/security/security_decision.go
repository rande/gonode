// Copyright © 2014-2018 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package security

type DecisionVoter interface {
	Support(o interface{}) bool
	Decide(t SecurityToken, attrs Attributes, o interface{}) bool
}

type AffirmativeDecision struct {
	Voters                     []Voter
	AllowIfAllAbstainDecisions bool
}

func (d *AffirmativeDecision) Support(o interface{}) bool {
	return true
}

func (d *AffirmativeDecision) Decide(t SecurityToken, attrs Attributes, o interface{}) bool {
	deny := 0

	for _, v := range d.Voters {

		if !v.Support(o) {
			continue
		}

		r, _ := v.Vote(t, o, attrs)

		switch r {
		case ACCESS_GRANTED:
			return true
		case ACCESS_DENIED:
			deny++
		}
	}

	if deny > 0 {
		return false
	}

	return d.AllowIfAllAbstainDecisions
}
