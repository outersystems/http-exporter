#!/bin/bash
set -eEuo pipefail

stacktrace () {
	echo "FAIL: unhandled error. stacktrace:"
	i=0
	while caller $i; do
		i=$((i+1))
	done
}

init() {
	trap "stacktrace" ERR
	registry="outersystems"
	context="$(cd "$(dirname "${BASH_SOURCE[0]}" )" && pwd)"
	dirname="$(basename ${context})"
	image="${registry:+${registry}/}${dirname}"
	bin="${dirname}"
}

clean() {
	#Cleaning go app
	rm -f ${context}/${bin}
}

compile() {
	echo "Compiling ${bin}..."
	docker run --rm -ti -w ${context} -v ${context}:${context} cell/cvim \
		go build -tags netgo
}

build() {
	local action=${1:-}
	local tag=${2:-latest}
	local status=0

	echo "Building ${image}:${tag} ..."
	if [ "${action}" == "release" ]; then
		docker tag ${image}:n ${image}:n-1 &>/dev/null || true
		docker rmi ${image}:latest ${image}:n &>/dev/null || true
	fi

	out=$(docker build -t ${image}:n ${context} 2>&1 || status=$?)
	if [ ${status} -ne 0 ]; then
		echo -e ${out} >&2
		exit $status
	fi

	if [ "${action}" == "release" ]; then
		docker tag ${image}:n ${image}:${tag}
	fi
	docker tag ${image}:n ${image}:latest
}

all() {
	local githash=
	if [ $(cd ${context} && git diff-index  HEAD -- | wc -l) -eq 0 ]; then
		githash="$(cd ${context} && git rev-parse --verify HEAD)"
	fi

	if [ -n "${githash}" ]; then
		local already=0
		docker inspect ${image}:${githash} &>/dev/null || already=$?
		if [ $already -eq 0 ]; then
			echo "Already built: ${image}:${githash}"
			exit 1
		fi

		echo "Release"
		compile
		build release $githash
	else
		echo "Snapshot"
		compile
		build snapshot
	fi

	clean
}

main() {

	if [ -n "${1:-}" -a "$(type -t ${1:-})" != "function" ]; then
		echo -e "Invalid action!\n" >&2
		exit 1
	fi

	init
	${@:-all}
}

main $@
