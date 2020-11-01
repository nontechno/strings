// Copyright 2020 The NonTechno Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package strings

import (
	"fmt"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

type Resolver func(justKey, fullMatch string, logger *log.Entry) (string, error)

// this func expands constructs like ${foo} into specified replacement.
// Q: why is it better than os.Expand?
// A:	- you can customize what prefix (e.g. "${") /postfix (e.g. "}") are
//		- no "unintended" behavior, like expanding "$foo" in addition to ${foo}
//		- unresolved entries can be left unchanged
//		- "foo" can include previously disallowed chanacters (e.g. ".")
//
// warning: this is a non-recursive expander - it will not resolve expanded values
func Expand(source string, prefix, postfix string, resolver Resolver, logger *log.Entry) (string, error) {
	already := ""
	for {
		start := strings.Index(source, prefix)
		if start < 0 {
			return already + source, nil
		}

		keyStart := start + len(prefix)
		keyLen := strings.Index(source[keyStart:], postfix)

		if keyLen < 0 {
			// there is no postfix found - return as it is
			logger.Warningf("found `prefix` but not `postfix` - possible but unlikely schenario...")
			return already + source, nil
		}

		end := keyStart + keyLen + len(postfix) - 1
		key := source[keyStart : keyStart+keyLen]

		val, err := resolver(key, prefix+key+postfix, logger)
		if err != nil {
			logger.WithError(err).Errorf("the provided `resolver` failed to find a match for the key (%s)", key)
			return "", err
		}

		already += source[:start] + val
		source = source[end+1:]
	}
}

// resolve using env vars, if not found - leave the original intact
func EnvironmentResolverIntact(justKey, fullMatch string, logger *log.Entry) (string, error) {
	if value, found := os.LookupEnv(justKey); found {
		return value, nil
	}
	logger.Warningf("failed to find environment variable (%s)", justKey)
	return fullMatch, nil
}

// resolve using env vars, if not found - substitute with empty string
func EnvironmentResolver(justKey, fullMatch string, logger *log.Entry) (string, error) {
	if value, found := os.LookupEnv(justKey); found {
		return value, nil
	}
	logger.Warningf("failed to find environment variable (%s)", justKey)
	return "", nil
}

// resolve using env vars, if not found - return error
func EnvironmentResolverFail(justKey, fullMatch string, logger *log.Entry) (string, error) {
	if value, found := os.LookupEnv(justKey); found {
		return value, nil
	}
	err := fmt.Errorf("failed to find env.var (%s)", justKey)
	logger.WithError(err).Errorf("env var not found")
	return "", err
}
