//+build cover_main

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

package tests

/*
Go toolchain lacks the ability to instrument arbitrary executables with
coverage counters.

This file wraps the maddy executable into a minimal layer of "test" logic to
make 'go test' work for it and produce the coverage report.

Use ./build_cover.sh to compile it into ./maddy.cover.

References:
https://stackoverflow.com/questions/43381335/how-to-capture-code-coverage-from-a-go-binary
https://blog.cloudflare.com/go-coverage-with-external-tests/
https://github.com/albertito/chasquid/blob/master/coverage_test.go
*/

import (
	"os"
	"testing"

	"github.com/jmcarbo/maddy"
)

func TestMain(m *testing.M) {
	// -test.* flags are registered somewhere in init() in "testing" (?)
	// so calling flag.Parse() in maddy.Run() catches them up.

	// maddy.Run changes the working directory, we need to change it back so
	// -test.coverprofile writes out profile in the right location.
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	code := maddy.Run()

	if err := os.Chdir(wd); err != nil {
		panic(err)
	}

	// Silence output produced by "testing" runtime.
	_, w, err := os.Pipe()
	if err == nil {
		os.Stderr = w
		os.Stdout = w
	}

	// Even though we do not have any tests to run, we need to call out into
	// "testing" to make it process flags and produce the coverage report.
	m.Run()

	// TestMain doc says we have to exit with a sensible status code on our
	// own.
	os.Exit(code)
}
