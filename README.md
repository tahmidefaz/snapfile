# Snapfile
Snapfile is a REST file upload/download service that allows users to download files for a specific amount of time. The file is deleted once
the predermined download quota is met. Written in Go. The frontend is under development.

### Setup Instructions
1. Though it is optional, you should install docker and docker-compose. It helps with obtaining a PostgreSQL instance.
2. Run `make setup` to install go dependencies.
3. If you have docker/docker-compose installed, run `make up_db` to spin up the Postgres instance.
4. Go to the `Makefile` and verify your DB credentials under `run`.
5. Finally, run `make run` to start the service. The server listens and serves on port `3000` by default.

### API Interface
The API provides the following end points:
Endpoint|Method(s)|Instructions
--------|---------|------------
/v1/upload|`POST`|required `form-data`: `file` and `max_downloads`. optional: `preferred_url`. See below.
/v1/download/:url|`GET`|`url` will be provided to you when you upload a file.
/v1/ping|`GET`|simple api health check

#### form-data for the upload endpoint
In order to upload a file to the API, you need to provide the following `form-data` in the body of your `POST` request.
form-data key|Reuired/Optional|Explanation
-------------|----------------|-----------
`file`|Required|The file you want to upload.
`max_downloads`|Required|The number of times you want users to download the file before it gets deleted.
`preferred_url`|Optional|Your url preference. If a preferred url is not provided or is not available, an UUID is generated.
