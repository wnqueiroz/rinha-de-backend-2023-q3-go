#!/bin/sh 
# Exemplos de requests
# curl -v -XPOST -H "content-type: application/json" -d '{"apelido" : "xpto", "nome" : "xpto xpto", "nascimento" : "2000-01-01", "stack": null}' "http://localhost:9999/pessoas"
# curl -v -XGET "http://localhost:9999/pessoas/1"
# curl -v -XGET "http://localhost:9999/pessoas?t=xpto"
# curl -v "http://localhost:9999/contagem-pessoas"

docker compose rm -f
docker compose down -v --rmi local --remove-orphans
docker compose up -d --build

GATLING_BIN_DIR=$HOME/gatling/3.9.5/bin
WORKSPACE=$(pwd)/stress-test

sh $GATLING_BIN_DIR/gatling.sh -rm local -s RinhaBackendSimulation \
    -rd "DESCRICAO" \
    -rf $WORKSPACE/user-files/results \
    -sf $WORKSPACE/user-files/simulations \
    -rsf $WORKSPACE/user-files/resources \

sleep 5

curl -v "http://localhost:9999/contagem-pessoas"

go tool pprof -http=localhost: http://localhost:9999/debug/pprof/heap