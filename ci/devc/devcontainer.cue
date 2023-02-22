name: "hof cli"

image: "mcr.microsoft.com/devcontainers/universal:2"

postCreateCommand: """
make hof && hof mod tidy
echo "hallo! you can now work on hof code, just type 'make hof' to rebuild"
"""

customizations: {
	vscode: extensions: [
		"asdine.cue",
		"jallen7usa.vscode-cue-fmt",
	]
}
