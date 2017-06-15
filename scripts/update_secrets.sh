#!/bin/sh

kubectl delete --ignore-not-found=true secret kanban-stats
kubectl create secret generic kanban-stats --from-file=./secrets/TRELLO_KEY --from-file=./secrets/TRELLO_TOKEN --from-file=./secrets/INFLUX_PASS
