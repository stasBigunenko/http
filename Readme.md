simple HTTP server

The server automatically will connect with http connection on localhost (on 127.0.0.1:8080). If you want to change any of these values you can use enviroment variables "PORT" and/or "HOST".

The server working with Post struct in json format.
Post {"Id": "uuid.UUID", "Author": "string", "Message": "string"}

The server response on the simple request as POST, GET, PUT and DELETE.

If you make POST, PUT or DELETE requests, changes will be automatically and safely saved to the memory Storage.
Your request body JSON should be object enclosed, just like the GET output. (for example {"Author": "Foo", "Message":"Foo"})
Id values are not mutable. Any id value in the body of your PUT request will be ignored. 
A POST or PUT request should include a Content-Type: application/json header to use the JSON in the request body. 
