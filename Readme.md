simpler HTTP server

The server working with Post struct.

The server response on the simple request as POST, GET, PUT and DELETE.

If you make POST, PUT, PATCH or DELETE requests, changes will be automatically and safely saved to the memory Storage.
Your request body JSON should be object enclosed, just like the GET output. (for example {"Author": "Foo", "Message":"Foo"})
Id values are not mutable. Any id value in the body of your PUT or PATCH request will be ignored. Only a value set in a POST request will be respected, but only if not already taken.
A POST, PUT or PATCH request should include a Content-Type: application/json header to use the JSON in the request body. 