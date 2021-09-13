package templates

const (
	DockerignoreTemplate = `# More info: https://docs.docker.com/engine/reference/builder/#dockerignore-file
# Ignore build and test binaries.
bin/
testbin/
`
	GitignoreTemplate = `
# Binaries for programs and plugins
*.exe
*.exe~
*.dll
*.so
*.dylib
bin
testbin/*
# Test binary, build with ` + "`go test -c`" + `
*.test
# Output of the go coverage tool, specifically when used with LiteIDE
*.out
# Kubernetes Generated files - skip generated files, except for vendored files
!vendor/**/zz_generated.*
# editor and IDE paraphernalia
.idea
.tmp/*
*.swp
*.swo
*~
`
)
