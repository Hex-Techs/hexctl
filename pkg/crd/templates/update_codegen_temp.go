package templates

// output-base 默认需要 go.mod repo 是三级目录

const UpdateCodeGenTemp = `#!/usr/bin/env bash

# set -o errexit
set -o nounset
set -o pipefail

MATCH=$(grep "^import" hack/tools.go)

if [[ ${MATCH} == "" ]];then
    sed -i".out" -e "s#// import _ \"k8s\.io/code-generator\"#import _ \"k8s\.io/code-generator\"#" hack/tools.go
fi

go mod vendor

SCRIPT_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
CODEGEN_PKG=${CODEGEN_PKG:-$(cd "${SCRIPT_ROOT}"; ls -d -1 ./vendor/k8s.io/code-generator 2>/dev/null || echo ../code-generator)}

API_BASE="${SCRIPT_ROOT}"/api/{{.Group}}/v1alpha1

mkdir -p ${API_BASE}
cp ${SCRIPT_ROOT}/api/v1alpha1/* ${API_BASE}

# generate the code with:
# --output-base    because this script should also be able to run inside the vendor dir of
#                  k8s.io/kubernetes. The output-base is needed for the generators to output into the vendor dir
#                  instead of the $GOPATH directly. For normal projects this can be dropped.
bash "${CODEGEN_PKG}"/generate-groups.sh "deepcopy,client,lister,informer" \
    {{.Repo}}/client {{.Repo}}/api \
    "{{.Group}}:{{.Version}}" \
    --go-header-file "${SCRIPT_ROOT}"/hack/boilerplate.go.txt \
    --output-base "$(dirname "${BASH_SOURCE[0]}")/../../../.."

if [ -d "${API_BASE}" ]; then
    cp ${SCRIPT_ROOT}/api/{{.Group}}/{{.Version}}/* ${SCRIPT_ROOT}/api/{{.Version}}/
    rm -rf ${SCRIPT_ROOT}/api/{{.Group}}
fi

find client -type f -name "*.go" | xargs sed -i".out" -e "s#{{.Repo}}/api/{{.Group}}/{{.Version}}#{{.Repo}}/api/{{.Version}}#g"
find client -type f -name "*go.out" | xargs rm -rf

sed -i".out" -e "s#import _ \"k8s\.io/code-generator\"#// import _ \"k8s\.io/code-generator\"#" hack/tools.go
find . -name tools.go.out | xargs rm -rf
go mod tidy

rm -rf vendor
`
