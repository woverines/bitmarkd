// Copyright (c) 2014-2018 Bitmark Inc.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"
	"path"

	"github.com/urfave/cli"

	"github.com/bitmark-inc/bitmarkd/command/bitmark-cli/configuration"
	"github.com/bitmark-inc/bitmarkd/command/bitmark-cli/encrypt"
)

func runSetup(c *cli.Context) error {

	m := c.App.Metadata["config"].(*metadata)

	name, err := checkName(c.GlobalString("identity"))
	if nil != err {
		return err
	}

	testnet := c.Bool("testnet")
	livenet := c.Bool("livenet")

	if testnet == livenet {
		return fmt.Errorf("only one of testnet/livenet can be set")
	}

	connect, err := checkConnect(c.String("connect"))
	if nil != err {
		return err
	}

	description, err := checkDescription(c.String("description"))
	if nil != err {
		return err
	}

	// optional existing hex key value
	privateKey, err := checkOptionalKey(c.String("privateKey"))
	if nil != err {
		return err
	}

	if m.verbose {
		fmt.Fprintf(m.e, "config: %s\n", m.file)
		fmt.Fprintf(m.e, "testnet: %t\n", testnet)
		fmt.Fprintf(m.e, "livenet: %t\n", livenet)
		fmt.Fprintf(m.e, "connect: %s\n", connect)
		fmt.Fprintf(m.e, "identity: %s\n", name)
		fmt.Fprintf(m.e, "description: %s\n", description)
	}

	// Create the folder hierarchy for configuration if not existing
	configDir := path.Dir(m.file)
	d, err := checkFileExists(configDir)
	if nil != err {
		if err := os.MkdirAll(configDir, 0750); nil != err {
			return err
		}
	} else if !d {
		return fmt.Errorf("path: %q is not a directory", configDir)
	}

	configData := &configuration.Configuration{
		DefaultIdentity: name,
		TestNet:         testnet,
		Connect:         connect,
		Identities:      make([]encrypt.IdentityType, 0),
	}

	err = addIdentity(configData, name, description, privateKey, c.GlobalString("password"), testnet)
	if nil != err {
		return err
	}

	m.config = configData
	m.save = true

	return nil
}
