#!/bin/bash

#based on https://github.com/clakech/couchbase-elastic

wait_for_start() {
    "$@"
    while [ $? -ne 0 ]
    do
        echo 'waiting for couchbase to start'
        sleep 1
        "$@"
    done
}

echo "launch couchbase"
/entrypoint.sh couchbase-server &

# wait for couchbase to be up - `couchbase-cli server-info` is broken / n/a
wait_for_start wget -q localhost:8091

echo $CLUSTER_TYP + " setup"

if [ $CLUSTER_TYPE == "MASTER" ]
  then
    # init the cluster - get_all_the_things
    couchbase-cli cluster-init -c 127.0.0.1:8091 -u Administrator -p password \
    --cluster-username=Administrator --cluster-password=password \
    --cluster-ramsize=600 --cluster-index-ramsize=512 \
    --index-storage-setting=default --services=data,index,query,fts

    # create bucket
    couchbase-cli bucket-create -c 127.0.0.1:8091 -u Administrator -p password \
    --bucket=company --bucket-type=couchbase --bucket-ramsize=600 --enable-flush=1 --wait
elif [ $CLUSTER_TYPE == "REPLICA" ]; then
    # add server to cluster
    couchbase-cli server-add -c couchbase_master:8091 -u Administrator -p password \
    --server-add="$(hostname -I | xargs):8091" --server-add-username=Administrator \
    --server-add-password=password --services=data,index,query,fts

    couchbase-cli rebalance -c couchbase_master:8091 -u Administrator -p password
fi


wait