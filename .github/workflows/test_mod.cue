package workflows

import "github.com/hofstadter-io/hof/.github/workflows/common"

common.#Workflow & {
	name: "test_mod"
	on: pull_request: { paths: ["lib/mod/**"] }
	jobs: test: {
		environment: "hof mod testing"
		steps: [ for step in common.#BuildSteps {step} ] + [{
			name: "Run mod tests"
			run: """
			hof test test.cue -s lib -t test -t mod
			"""
			env: {
				HOFMOD_SSHKEY: "${{secrets.HOFMOD_SSHKEY}}"
				GITHUB_TOKEN: "${{secrets.HOFMOD_TOKEN}}"
				GITLAB_TOKEN: "${{secrets.GITLAB_TOKEN}}"
				BITBUCKET_TOKEN: "${{secrets.BITBUCKET_TOKEN}}"
			}
		}]
	}
}

