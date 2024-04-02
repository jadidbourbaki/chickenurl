# chickenurl

- Data is persistent with a configurable SLA of 1 second
- O(1) lookups for hash and url
- Production-level logging
- Stress tested with curl
- Integration tested with curl

To build: 

```
go build .
```

To deploy, just run the built binary:


Unit tests:

```
go test 
```

Integration tests:

```
test_curl_put.bash
```

Then use the output of the above script, let's say it's `ABCDEF`:

```
test_curl_get.bash ABCDEF
```

and 

```
test_curl_delete.bash ABCDEF
```
