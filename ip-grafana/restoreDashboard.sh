curl -u admin:admin -H 'Content-Type: application/json' -X POST -d @trello.json http://monitoring.platform.prod.aws.cloud.nordstrom.net/api/datasources
curl -u admin:admin -H 'Content-Type: application/json' -X POST -d @teamStats.json http://monitoring.platform.prod.aws.cloud.nordstrom.net/api/dashboards/db
