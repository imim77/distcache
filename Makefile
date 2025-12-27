build:
	go build -o bin/distcache

run: build
	 ./bin/distcache
runfollower: build
	 ./bin/distcache --listenaddr :4000 --leaderaddr :3000
test:
	go test -v ./...