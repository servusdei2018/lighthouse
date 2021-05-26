/*
Lighthouse
Copyright (C) 2021  Nathanael Bracy

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published
by the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package main

import (
	"log"
	"math/rand"
	"time"
)

const (
	PORT = ":1990"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	var lighthouse = NewMud()
	log.Fatal(lighthouse.ListenAndServe(PORT))
}
