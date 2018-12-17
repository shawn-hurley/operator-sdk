// Copyright 2018 The Operator-SDK Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ansible

import (
	"path/filepath"

	"github.com/operator-framework/operator-sdk/pkg/scaffold/input"
)

const MoleculeDefaultInstallFile = "INSTALL.rst"

type MoleculeDefaultInstall struct {
	input.Input
}

// GetInput - gets the input
func (m *MoleculeDefaultInstall) GetInput() (input.Input, error) {
	if m.Path == "" {
		m.Path = filepath.Join(MoleculeDefaultDir, MoleculeDefaultInstallFile)
	}
	m.TemplateBody = moleculeDefaultInstallAnsibleTmpl

	return m.Input, nil
}

const moleculeDefaultInstallAnsibleTmpl = `*******
Docker driver installation guide
*******

Requirements
============

* General molecule dependencies (see https://molecule.readthedocs.io/en/latest/installation.html)
* Docker Engine
* docker-py
* docker

Install
=======

    $ sudo pip install docker-py
`
