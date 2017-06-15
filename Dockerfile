FROM nordstrom/baseimage-alpine:3.2

MAINTAINER Marius Grigoriu

ADD build/kanban-stats kanban-stats

ENTRYPOINT /kanban-stats
