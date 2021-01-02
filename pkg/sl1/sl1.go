/*
Copyright © 2020 Marco Ostaska

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/

// Package sl1 Package account implements utility routines for manipulating /api
package sl1

import (
	"fmt"
	"path"
	"path/filepath"

	"github.com/marco-ostaska/sl1cmd/pkg/sl1/httpcalls"
)

// These are predefined shortcuts for sl1 APIs sections.
const (
	AccountAPI = "/api/account"
)

// BasicInfo is an abstraction for basic api results
type BasicInfo []struct {
	URI         string `json:"URI"`
	Description string `json:"description"`
	ID          string
}

// Load loads BasicInfo from /api/xyz
func (bInfo *BasicInfo) Load(api string) error {
	var a httpcalls.APIData
	a.API = api
	err := a.NewRequest(&bInfo)
	if err != nil {
		return err
	}

	for i := 0; i < len(*bInfo); i++ {
		id, err := bInfo.Sl1ID((*bInfo)[i].Description)
		if err != nil {
			return err
		}
		(*bInfo)[i].ID = id
	}

	return nil
}

// SearchByURI search description by provided URI
func (bInfo *BasicInfo) SearchByURI(uri string) (string, error) {

	for _, u := range *bInfo {
		if uri == u.URI {
			return u.Description, nil
		}
	}
	return uri, fmt.Errorf("%s: no such uri", uri)

}

// IndexPosition searchs index position by provided BasicInfo.Description
func (bInfo *BasicInfo) IndexPosition(s string) (int, error) {

	for i, u := range *bInfo {
		if s == u.Description {
			return i, nil
		}
	}
	return 0, fmt.Errorf("%s: no such %v", s, path.Base(path.Dir((*bInfo)[0].URI)))

}

// Sl1ID returns sl1 id from BasicInfo.Description
func (bInfo *BasicInfo) Sl1ID(s string) (string, error) {
	i, err := bInfo.IndexPosition(s)
	if err != nil {
		return "", err
	}

	return path.Base((*bInfo)[i].URI), nil
}

// Println formats the output in a user friendly style
func (bInfo *BasicInfo) Println(args []string) {

	if len(args) == 0 {
		bInfo.printRange()
	}

	if len(args) > 0 {
		bInfo.printByDesc(args)
	}
}

// printRange formats the output for every BasicInfo range
func (bInfo *BasicInfo) printRange() {
	for _, i := range *bInfo {
		fmt.Printf("sl1id=%s(%s)\n", i.ID, i.Description)
	}
}

// printByDesc formats the output searching by Description
func (bInfo *BasicInfo) printByDesc(args []string) {
	for _, a := range args {
		i, err := bInfo.IndexPosition(a)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Printf("sl1id=%s(%s)\n", filepath.Base(((*bInfo)[i].URI)), ((*bInfo)[i].Description))
	}

}
