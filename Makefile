redis-docker:
	docker run -d -p 6379:6379 --name redis redis

go-run-poller:
	go run poller/main.go

go-run-scheduler:
	go run scheduler/main.go