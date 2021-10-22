Simple CRUD API with additional functions: save data in and download data from *csv file.

HTTP Server working through gRPC Server which can be connected to differenet storages through enviroment variables.
You can connect to the following types:
- inmemory storage
- Redis storage
- PostgreSQL storage
- MongoDB storage
- ElastickSearch storage

HTTP server supports GraphQL service.

HTTP server has JWT auth. You can receive token by login to the server where you should indicate your name. You will receive the token in header's response which you should specify in all your requests.

Also HTTP server connected to monitoring system Prometheus and data visualization platform Grafana.

All services working through Docker with docker-compose file, or you can work through Kubernetes as well.

If the enviroment variables will be empty of any wanted storage, the connection will be using default credentials.

You can find lists of all enviroment variables in config files which are located in cmd dir to each server.
