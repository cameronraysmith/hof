package workflows

import (
	"github.com/hofstadter-io/ghacue"
	"github.com/hofstadter-io/hof/ci/gha/common"
)

ghacue.#Workflow & {
	name: "dagger"
	on:   _ | *["push"]
	env: {
		HOF_TELEMETRY_DISABLED: "1"
		HOF_FMT_VERSION:        "v0.6.9-rc.2"
	}
	jobs: {
		inception: {
			"runs-on": "ubuntu-latest"
			concurrency: {
				group:                "${{ github.workflow }}-inception-${{ github.ref_name }}"
				"cancel-in-progress": true
			}

			steps: [
				common.Steps.go.setup & {#ver: "1.23.x"},
				common.Steps.checkout,
				common.Steps.vars,
				common.Steps.go.deps,
				common.Steps.docker.compat,

				{
					name: "dagger-in-dagger"
					run:  "go run ./test/dagger/main/dagger-in-dagger.go"
				},
				{
					name: "dockerd-in-dagger"
					run:  "go run ./test/dagger/main/dockerd-in-dagger.go"
				},
			]
		}

		hof: {
			environment: "hof mod testing"
			"runs-on":   "ubuntu-latest"
			concurrency: {
				group:                "${{ github.workflow }}-hof-${{ github.ref_name }}"
				"cancel-in-progress": true
			}

			steps: [
				common.Steps.go.setup,
				common.Steps.dagger.cache,
				common.Steps.checkout,
				common.Steps.vars,
				common.Steps.go.deps,
				common.Steps.docker.compat,

				{
					name: "hof-in-dagger"
					run:  "go run ./test/dagger/main/hof.go"
					env: {
						GITHUB_TOKEN:       "${{secrets.GITHUB_TOKEN}}"
						GITLAB_TOKEN:       "${{secrets.GITLAB_TOKEN}}"
						BITBUCKET_USERNAME: "hofstadter"
						BITBUCKET_PASSWORD: "${{secrets.BITBUCKET_TOKEN}}"
					}
				},
			]
		}
	}
}
