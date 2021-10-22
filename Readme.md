Simple REST API for CRUD operations with additional functions: save data in and download data from *.csv file.

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

API Documentation.

To test the REST API , use any rest client like postman etc or cURL.

to login and receive a token:
curl --location --request GET 'http://localhost:8085/posts/login' \
--header 'Content-Type: application/json' \
--data-raw '{
    "Name": "stas"
}'

to Create post:
curl --location --request POST 'http://localhost:8085/posts/create' \
--header 'Authorization: Bearer <your_token>' \
--header 'Content-Type: application/json' \
--data-raw '{
    "Author": "Bob",
    "Message": "The second"
}'

to Get post:
curl --location --request GET 'http://localhost:8085/posts/<id_of_the_post>' \
--header 'Authorization: Bearer <your_token>' \
--header 'Content-Type: text/plain' \
--data-binary '@'

to Get all posts:
curl --location --request GET 'http://localhost:8085/posts/' \
--header 'Authorization: Bearer <your_token>' \
--data-raw ''

to Update post:
curl --location --request PUT 'http://localhost:8085/posts/<id_of_the_post>' \
--header 'Authorization: Bearer <your_token>' \
--header 'Content-Type: application/json' \
--data-raw '{
    "Author": "new_author",
    "Message": "new_message"
}'

to Delete post:
curl --location --request DELETE 'http://localhost:8085/posts/<id_of_the_post>' \
--header 'Authorization: Bearer <your_token>'

to Save data to the file:
curl --location --request GET 'http://localhost:8085/posts/download' \
--header 'Authorization: Bearer <your_token>'
--to 'file=@"path_name_of_the_file.csv"'

to Download data from file:
curl --location --request POST 'http://localhost:8085/posts/upload' \
--header 'Authorization: Bearer <your_token>' \
--form 'file=@"path_name_of_the_file.csv"'

to use GraphQL:

to Create post:
curl --location --request POST 'http://localhost:8085/posts/graphql/' \
--header 'Authorization: Bearer <your_token>' \
--header 'Content-Type: application/json' \
--data-raw '{"query":"mutation{\r\n  createPost(input:{author: \"andrew\", message:\"the first3\"}){\r\n    id,\r\n    author,\r\n    message\r\n  }\r\n}","variables":{}}'

to Get post:
curl --location --request POST 'http://localhost:8085/posts/graphql/' \
--header 'Authorization: Bearer <your_token>' \
--header 'Content-Type: application/json' \
--data-raw '{"query":"query{ getPost(id: \"<id_of_the_post>\"){\r\n    id,\r\n    author,\r\n    message\r\n  }\r\n}","variables":{}}'

to Get all posts:
curl --location --request POST 'http://localhost:8085/posts/graphql/' \
--header 'Authorization: Bearer <your_token>' \
--header 'Content-Type: application/json' \
--data-raw '{"query":"query {\r\n  getPosts{\r\n    id,\r\n    author,\r\n    message\r\n  }\r\n}","variables":{}}'

to Update post:
curl --location --request POST 'http://localhost:8085/posts/graphql/' \
--header 'Authorization: Bearer <your_token>' \
--header 'Content-Type: application/json' \
--data-raw '{"query":"mutation{\r\n  updatePost (input:{id: \"id_of_the_post>\", author: \"<new_name>\", message:\"<new_message>\"}){\r\n    id,\r\n    author,\r\n    message\r\n  }\r\n}","variables":{}}'

to Delete post:
curl --location --request POST 'http://localhost:8085/posts/graphql/' \
--header 'Authorization: Bearer <your_token>' \
--header 'Content-Type: application/json' \
--data-raw '{"query":"mutation{\r\n  deletePost (id: \"id_of_the_post>\")\r\n}","variables":{}}'
