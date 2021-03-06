/*
Maddy Mail Server - Composable all-in-one email server.
Copyright © 2019-2020 Max Mazurov <fox.cpp@disroot.org>, Maddy Mail Server contributors

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package main

import (
	"fmt"

	"github.com/jmcarbo/maddy/internal/modify/dkim"
	"github.com/urfave/cli"
)

func dkimCommand(ctx *cli.Context) error {
	newKeyAlgo := ctx.String("algorithm")
	keyPath := ctx.String("keypath")
  fmt.Printf("Generating dkim keypair in %s with algorithm %s\n", keyPath, newKeyAlgo)
  _, err := dkim.GenerateAndWrite(keyPath, newKeyAlgo)
  if err != nil {
    fmt.Printf("Error [%s] generating dkim keypair in %s with algorithm %s\n", err.Error(), keyPath, newKeyAlgo)
  }
	return err
}
