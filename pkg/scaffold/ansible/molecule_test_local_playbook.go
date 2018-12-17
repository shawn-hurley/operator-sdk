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

	"github.com/operator-framework/operator-sdk/pkg/scaffold"
	"github.com/operator-framework/operator-sdk/pkg/scaffold/input"
)

const MoleculeTestLocalPlaybookFile = "playbook.yml"

type MoleculeTestLocalPlaybook struct {
	input.Input
	Resource scaffold.Resource
}

// GetInput - gets the input
func (m *MoleculeTestLocalPlaybook) GetInput() (input.Input, error) {
	if m.Path == "" {
		m.Path = filepath.Join(MoleculeTestLocalDir, MoleculeTestLocalPlaybookFile)
	}
	m.TemplateBody = moleculeTestLocalPlaybookAnsibleTmpl

	return m.Input, nil
}

const moleculeTestLocalPlaybookAnsibleTmpl = `---

- name: Build Operator in Kubernetes docker container
  hosts: k8s
  vars:
    image_name: {{.Resource.FullGroup}}/{{.ProjectName}}:testing
  tasks:
  # using command so we don't need to install any dependencies
  - name: Get existing image hash
    command: docker images -q {{"{{image_name}}"}}
    register: prev_hash
    changed_when: false

  - name: Build Operator Image
    command: docker build -f /build/build/Dockerfile -t {{"{{ image_name }}"}} /build
    register: build_cmd
    changed_when: not prev_hash.stdout or (prev_hash.stdout and prev_hash.stdout not in ''.join(build_cmd.stdout_lines[-2:]))

- name: Converge
  hosts: localhost
  connection: local
  vars:
    ansible_python_interpreter: '{{ "{{ ansible_playbook_python }}" }}'
    deploy_dir: "{{"{{ lookup('env', 'MOLECULE_PROJECT_DIRECTORY') }}/deploy"}}"
    image_name: {{.Resource.FullGroup}}/{{.ProjectName}}:testing
  tasks:

  - name: Create the Operator Deployment
    k8s:
      namespace: '{{ "{{ namespace }}" }}'
      definition: "{{"{{"}} lookup('template', '/'.join([deploy_dir, 'operator.yaml'])) {{"}}"}}"
    vars:
      pull_policy: Never
      REPLACE_IMAGE: '{{ "{{ image_name }}" }}'

  - name: Create the {{.Resource.FullGroup}}/{{.Resource.Version}}.{{.Resource.Kind}}
    k8s:
      namespace: '{{ "{{ namespace }}" }}'
      definition: "{{"{{"}} lookup('file', '/'.join([deploy_dir, 'crds/{{.Resource.Group}}_{{.Resource.Version}}_{{.Resource.LowerKind}}_cr.yaml'])) {{"}}"}}"

  - name: Get the newly created Custom Resource
    debug:
      msg: "{{"{{"}} lookup('k8s', group='{{.Resource.FullGroup}}', api_version='{{.Resource.Version}}', kind='{{.Resource.Kind}}', namespace=namespace, resource_name=cr.metadata.name) {{"}}"}}"
    vars:
      cr: "{{"{{"}} lookup('file', '/'.join([deploy_dir, 'crds/{{.Resource.Group}}_{{.Resource.Version}}_{{.Resource.LowerKind}}_cr.yaml'])) | from_yaml {{"}}"}}"

- import_playbook: '{{"{{ playbook_dir }}/../default/asserts.yml"}}'
`
