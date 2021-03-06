// Copyright 2016 Palantir Technologies, Inc.
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

package integration_test

import (
	"fmt"
	"testing"

	"github.com/nmiyake/pkg/gofiles"
	"github.com/palantir/godel/framework/pluginapitester"
	"github.com/palantir/godel/pkg/osarch"
	"github.com/palantir/godel/pkg/products/v2/products"
	"github.com/stretchr/testify/require"

	"github.com/palantir/distgo/publisher/publishertester"
)

func TestArtifactoryPublish(t *testing.T) {
	const godelYML = `exclude:
  names:
    - "\\..+"
    - "vendor"
  paths:
    - "godel"
`

	pluginPath, err := products.Bin("dist-plugin")
	require.NoError(t, err)

	publishertester.RunAssetPublishTest(t,
		pluginapitester.NewPluginProvider(pluginPath),
		nil,
		"artifactory",
		[]publishertester.TestCase{
			{
				Name: "publishes artifact and POM to Artifactory",
				Specs: []gofiles.GoFileSpec{
					{
						RelPath: "foo/foo.go",
						Src:     `package main; func main() {}`,
					},
				},
				ConfigFiles: map[string]string{
					"godel/config/godel.yml": godelYML,
					"godel/config/dist-plugin.yml": `
products:
  foo:
    build:
      main-pkg: ./foo
    dist:
      disters:
        type: os-arch-bin
    publish:
      group-id: com.test.group
      info:
        artifactory:
          config:
            url: http://artifactory.domain.com
            username: testUsername
            password: testPassword
            repository: testRepo
`,
				},
				Args: []string{
					"--dry-run",
				},
				WantOutput: func(projectDir string) string {
					return fmt.Sprintf(`[DRY RUN] Uploading out/dist/foo/1.0.0/os-arch-bin/foo-1.0.0-%s.tgz to http://artifactory.domain.com/artifactory/testRepo/com/test/group/foo/1.0.0/foo-1.0.0-%s.tgz
[DRY RUN] Uploading to http://artifactory.domain.com/artifactory/testRepo/com/test/group/foo/1.0.0/foo-1.0.0.pom
`, osarch.Current().String(), osarch.Current().String())
				},
			},
			{
				Name: "skips POM publish based on configuration",
				Specs: []gofiles.GoFileSpec{
					{
						RelPath: "foo/foo.go",
						Src:     `package main; func main() {}`,
					},
				},
				ConfigFiles: map[string]string{
					"godel/config/godel.yml": godelYML,
					"godel/config/dist-plugin.yml": `
products:
  foo:
    build:
      main-pkg: ./foo
    dist:
      disters:
        type: os-arch-bin
    publish:
      group-id: com.test.group
      info:
        artifactory:
          config:
            url: http://artifactory.domain.com
            username: testUsername
            password: testPassword
            repository: testRepo
            no-pom: true
`,
				},
				Args: []string{
					"--dry-run",
				},
				WantOutput: func(projectDir string) string {
					return fmt.Sprintf(`[DRY RUN] Uploading out/dist/foo/1.0.0/os-arch-bin/foo-1.0.0-%s.tgz to http://artifactory.domain.com/artifactory/testRepo/com/test/group/foo/1.0.0/foo-1.0.0-%s.tgz
`, osarch.Current().String(), osarch.Current().String())
				},
			},
			{
				Name: "can use flags to specify values",
				Specs: []gofiles.GoFileSpec{
					{
						RelPath: "foo/foo.go",
						Src:     `package main; func main() {}`,
					},
				},
				ConfigFiles: map[string]string{
					"godel/config/godel.yml": godelYML,
					"godel/config/dist-plugin.yml": `
products:
  foo:
    build:
      main-pkg: ./foo
    dist:
      disters:
        type: os-arch-bin
    publish:
      info:
        artifactory:
`,
				},
				Args: []string{
					"--dry-run",
					"--group-id", "com.test.group",
					"--url", "http://artifactory.domain.com",
					"--username", "testUsername",
					"--password", "testPassword",
					"--repository", "testRepo",
				},
				WantOutput: func(projectDir string) string {
					return fmt.Sprintf(`[DRY RUN] Uploading out/dist/foo/1.0.0/os-arch-bin/foo-1.0.0-%s.tgz to http://artifactory.domain.com/artifactory/testRepo/com/test/group/foo/1.0.0/foo-1.0.0-%s.tgz
[DRY RUN] Uploading to http://artifactory.domain.com/artifactory/testRepo/com/test/group/foo/1.0.0/foo-1.0.0.pom
`, osarch.Current().String(), osarch.Current().String())
				},
			},
			{
				Name: "can use flags to specify values including no-pom",
				Specs: []gofiles.GoFileSpec{
					{
						RelPath: "foo/foo.go",
						Src:     `package main; func main() {}`,
					},
				},
				ConfigFiles: map[string]string{
					"godel/config/godel.yml": godelYML,
					"godel/config/dist-plugin.yml": `
products:
  foo:
    build:
      main-pkg: ./foo
    dist:
      disters:
        type: os-arch-bin
    publish:
      info:
        artifactory:
`,
				},
				Args: []string{
					"--dry-run",
					"--group-id", "com.test.group",
					"--url", "http://artifactory.domain.com",
					"--username", "testUsername",
					"--password", "testPassword",
					"--repository", "testRepo",
					"--no-pom",
				},
				WantOutput: func(projectDir string) string {
					return fmt.Sprintf(`[DRY RUN] Uploading out/dist/foo/1.0.0/os-arch-bin/foo-1.0.0-%s.tgz to http://artifactory.domain.com/artifactory/testRepo/com/test/group/foo/1.0.0/foo-1.0.0-%s.tgz
`, osarch.Current().String(), osarch.Current().String())
				},
			},
			{
				Name: "publishes multiple artifacts and POM to Artifactory",
				Specs: []gofiles.GoFileSpec{
					{
						RelPath: "foo/foo.go",
						Src:     `package main; func main() {}`,
					},
				},
				ConfigFiles: map[string]string{
					"godel/config/godel.yml": godelYML,
					"godel/config/dist-plugin.yml": `
products:
  foo:
    build:
      main-pkg: ./foo
      os-archs:
        - os: darwin
          arch: amd64
        - os: linux
          arch: amd64
    dist:
      disters:
        type: os-arch-bin
        config:
          os-archs:
            - os: darwin
              arch: amd64
            - os: linux
              arch: amd64
    publish:
      group-id: com.test.group
      info:
        artifactory:
          config:
            url: http://artifactory.domain.com
            username: testUsername
            password: testPassword
            repository: testRepo
`,
				},
				Args: []string{
					"--dry-run",
				},
				WantOutput: func(projectDir string) string {
					return `[DRY RUN] Uploading out/dist/foo/1.0.0/os-arch-bin/foo-1.0.0-darwin-amd64.tgz to http://artifactory.domain.com/artifactory/testRepo/com/test/group/foo/1.0.0/foo-1.0.0-darwin-amd64.tgz
[DRY RUN] Uploading out/dist/foo/1.0.0/os-arch-bin/foo-1.0.0-linux-amd64.tgz to http://artifactory.domain.com/artifactory/testRepo/com/test/group/foo/1.0.0/foo-1.0.0-linux-amd64.tgz
[DRY RUN] Uploading to http://artifactory.domain.com/artifactory/testRepo/com/test/group/foo/1.0.0/foo-1.0.0.pom
`
				},
			},
		},
	)
}

func TestArtifactoryUpgradeConfig(t *testing.T) {
	pluginPath, err := products.Bin("dist-plugin")
	require.NoError(t, err)

	pluginapitester.RunUpgradeConfigTest(t,
		pluginapitester.NewPluginProvider(pluginPath),
		nil,
		[]pluginapitester.UpgradeConfigTestCase{
			{
				Name: `valid v0 config works`,
				ConfigFiles: map[string]string{
					"godel/config/dist-plugin.yml": `
products:
  foo:
    build:
      main-pkg: ./foo
      os-archs:
        - os: darwin
          arch: amd64
        - os: linux
          arch: amd64
    dist:
      disters:
        type: os-arch-bin
        config:
          os-archs:
            - os: darwin
              arch: amd64
            - os: linux
              arch: amd64
    publish:
      group-id: com.test.group
      info:
        artifactory:
          config:
            # comment
            url: http://artifactory.domain.com
            username: testUsername
            password: testPassword
            repository: testRepo
`,
				},
				WantOutput: ``,
				WantFiles: map[string]string{
					"godel/config/dist-plugin.yml": `
products:
  foo:
    build:
      main-pkg: ./foo
      os-archs:
        - os: darwin
          arch: amd64
        - os: linux
          arch: amd64
    dist:
      disters:
        type: os-arch-bin
        config:
          os-archs:
            - os: darwin
              arch: amd64
            - os: linux
              arch: amd64
    publish:
      group-id: com.test.group
      info:
        artifactory:
          config:
            # comment
            url: http://artifactory.domain.com
            username: testUsername
            password: testPassword
            repository: testRepo
`,
				},
			},
		},
	)
}
