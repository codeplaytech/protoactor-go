root_dir:=$(CURDIR)
project_name:=$(shell basename $(CURDIR))
# ifeq ($(project_name),)
# $(error "missing project_name")
# endif 

cp:=consul
ttl:=10s
loops:=10000
clients:=10
interval:=0ms
env=prod


start:
	tmux new-session -d -s eg
	tmux setenv -t eg PROTO_ACTOR_ENV $(env)
	tmux split-window -t "eg:0"   -v
	tmux split-window -t "eg:0.0" -h -p 66
	tmux split-window -t "eg:0.1" -h -p 50
	tmux select-pane -t "eg:0.3"
	tmux send-keys -t "eg:0.0" "go run ./server/ --provider $(cp) --ttl $(ttl) --port 9991" Enter
	tmux send-keys -t "eg:0.1" "go run ./server/ --provider $(cp) --ttl $(ttl) --port 9992" Enter
	tmux send-keys -t "eg:0.2" "go run ./server/ --provider $(cp) --ttl $(ttl) --port 9993" Enter
	tmux send-keys -t "eg:0.3" "sleep 2 && go run ./client/ --provider $(cp) --clients $(clients) --loops $(loops) --interval $(interval) " Enter
	# tmux send-keys -t "eg:0.2" "go run member/main.go" Enter
	tmux attach -t eg
	tmux kill-session -t eg


stop:
	tmux kill-session -t eg


start-with-etcd:
	make start cp=etcd


debug:
	PROTO_ACTOR_ENV=dev make start cp=etcd


mock-clients-10w:
	make start cp=etcd clients=100000 loops=10 interval=100ms


mock-clients-20w:
	make start cp=etcd clients=200000 loops=10 interval=100ms


protoc:=docker run --rm -it \
	-v $(root_dir):/$(project_name) \
	-w /$(project_name)  \
	cupen/protoc:3.19.4-20220424


proto:
	$(protoc) -I=. --go_out=. --go_opt=paths=source_relative shared/*.proto
	$(protoc) -I=. --gograinv2_out=. shared/*.proto

