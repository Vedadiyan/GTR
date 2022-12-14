# GTR 

GTR (being short for Go Tiny Router) is a minimalistic router
that was initially developed to identify RESTful call signatures
for the purpose of caching them (based on parameters).

[Usage Examples]
GTR parses URLs based on Express style templates, for example:

	`http://www.abcdefg.com/api/v1/users/:username/details`

The provided URL can be successfully matched against the following
URLs:

	`http://www.abcdefg.com/api/v1/users/ken/details`
	`http://www.abcdefg.com/api/v1/users/dennis/details`

Alongside with route parameter, GTR also supports query string
specification, in such a way that if a query parameter is specified
a match will be successfull only and only if the target URL also
specifies that query parameter and that it is of the same value as
originally specified. For example:

	`http://www.abcdefg.com/api/v1/users/:username/details?type=cached`

The provided URL can be successfully matched against the following
URLs:

	`http://www.abcdefg.com/api/v1/users/ken/details?type=cached&format=JSON`
	`http://www.abcdefg.com/api/v1/users/dennis/details?type=cached`

However, the following URLs will NOT be successfullt matched:

	`http://www.abcdefg.com/api/v1/users/ken/details?format=JSON`
	`http://www.abcdefg.com/api/v1/users/dennis/details?`

This behavior has been designed intentional to serve the original purpose
of the library.
