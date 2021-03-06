// Copyright 2018 The Hugo Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package hugolib

import (
	"fmt"

	"github.com/gohugoio/hugo/common/text"

	"github.com/bep/mapstructure"
	"github.com/pkg/errors"
)

type refArgs struct {
	Path         string
	Lang         string
	OutputFormat string
}

func (p *Page) decodeRefArgs(args map[string]interface{}) (refArgs, *Site, error) {
	var ra refArgs
	err := mapstructure.WeakDecode(args, &ra)
	if err != nil {
		return ra, nil, nil
	}
	s := p.s

	if ra.Lang != "" && ra.Lang != p.Lang() {
		// Find correct site
		found := false
		for _, ss := range p.s.owner.Sites {
			if ss.Lang() == ra.Lang {
				found = true
				s = ss
			}
		}

		if !found {
			p.s.siteRefLinker.logNotFound(ra.Path, fmt.Sprintf("no site found with lang %q", ra.Lang), p, text.Position{})
			return ra, nil, nil
		}
	}

	return ra, s, nil
}

func (p *Page) Ref(argsm map[string]interface{}) (string, error) {
	return p.ref(argsm, p)
}

func (p *Page) ref(argsm map[string]interface{}, source interface{}) (string, error) {
	args, s, err := p.decodeRefArgs(argsm)
	if err != nil {
		return "", errors.Wrap(err, "invalid arguments to Ref")
	}

	if s == nil {
		return p.s.siteRefLinker.notFoundURL, nil
	}

	if args.Path == "" {
		return "", nil
	}

	return s.refLink(args.Path, source, false, args.OutputFormat)

}

func (p *Page) RelRef(argsm map[string]interface{}) (string, error) {
	return p.relRef(argsm, p)
}

func (p *Page) relRef(argsm map[string]interface{}, source interface{}) (string, error) {
	args, s, err := p.decodeRefArgs(argsm)
	if err != nil {
		return "", errors.Wrap(err, "invalid arguments to Ref")
	}

	if s == nil {
		return p.s.siteRefLinker.notFoundURL, nil
	}

	if args.Path == "" {
		return "", nil
	}

	return s.refLink(args.Path, source, true, args.OutputFormat)

}
