version:
	git stash -u
	sed -i "s/[[:digit:]]\+\.[[:digit:]]\+\.[[:digit:]]\+/$(v)/g" gommit/version.go
	git add -A
	git commit -m "feat(version) : "$(v)
	git tag v$(v) master

fmt:
	find ! -path "./vendor/*" -name "*.go" -exec go fmt {} \;

gometalinter:
	gometalinter -D gotype --vendor --deadline=240s --dupl-threshold=200 -e '_string' -j 5 ./...

run-tests:
	./test.sh

test-all: gometalinter run-tests

test-package:
	go test -race -cover -coverprofile=/tmp/gommit github.com/antham/gommit/$(pkg)
	go tool cover -html=/tmp/gommit
