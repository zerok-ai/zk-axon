# axon

This is a golang project to query the collected scenarios, traces and spans. This service will run on client cluster. It uses postgres to collect the data.

All the authentications are done in zk-Auth service and the query are directed to Axon only after the authentication is successful.