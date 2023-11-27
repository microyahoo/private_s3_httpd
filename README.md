private_s3_httpd
---------------

Private HTTP Server for S3 content.

Amazon S3 provides a public HTTP interface for accessing content, but what if you don't want publicly accessible files?

`private_s3_httpd` exposes a private HTTP/HTTPS endpoint for an S3 bucket so you can controll access to the data. This is ideal for accessing S3 via HTTP/HTTPS api as a backend data service, or for local http browsing of a private s3 bucket, or for use behind another authentication service (like [oauth2_proxy](https://github.com/bitly/oauth2_proxy)) to secure access.


```
Usage of bin/private_s3_httpd:
  -bucket string
    	S3 bucket name
  -listen string
    	address:port to listen on. (default ":8444")
  -log-requests
    	log HTTP requests (default true)
  -region string
    	AWS S3 Region (default "us-east-1")
  -s3-endpoint string
    	AWS S3 endpoint
  -ak string
    	AWS S3 access key
  -sk string
    	AWS S3 secret key
```
