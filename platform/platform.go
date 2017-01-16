/*
Copyright (C) 2017 Andreas T Jonsson

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package platform

import (
	"log"
	"math"
	"path"
	"sync/atomic"
)

var (
	ConfigPath string
	idCounter  uint64
)

func CfgRootJoin(p ...string) string {
	return path.Clean(path.Join(ConfigPath, path.Join(p...)))
}

func NewId64() uint64 {
	if idCounter == math.MaxUint64 {
		log.Panicln("id space exhausted")
	}
	return atomic.AddUint64(&idCounter, 1) - 1
}
