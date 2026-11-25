//go:build !darwin || ios

package main

import "errors"

func openBlackHole() (virtualOutput, error) { return nil, errors.New("BlackHole requires macOS") }
