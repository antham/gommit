version:
	git stash -u
	sed -i "s/[[:digit:]]\+\.[[:digit:]]\+\.[[:digit:]]\+/$(v)/g" gommit/version.go
	git add -A
	git commit -m "feat(version) : "$(v)
	git tag v$(v) master

compile:
	git stash -u
	gox -output "build/{{.Dir}}_{{.OS}}_{{.Arch}}"

fmt:
	find ! -path "./vendor/*" -name "*.go" -exec gofmt -s -w {} \;

run-tests:
	./test.sh

run-quick-tests:
	go test -v $(shell glide nv)

test-all: run-tests

test-package:
	go test -race -cover -coverprofile=/tmp/gommit github.com/antham/gommit/$(pkg)
	go tool cover -html=/tmp/gommit
