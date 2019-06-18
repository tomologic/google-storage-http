# Serve Google Storage files via HTTP

Source: https://fale.io/blog/2018/04/12/an-http-server-to-serve-gcs-files/

## Configuration
Following environment variables can be set:

* PORT is a required port number to serve HTTP on.
* BUCKET is a required bucket name to serve files from.
* LOGGING=true enables log output for every request.

## Running in docker
If not overridden by environment variables, Dockerfile assumes LOGGING=true
and PORT=8080.

    docker build -t google-storage-http .
    docker run -e BUCKET=bucket_name \
        -e GOOGLE_APPLICATION_CREDENTIALS=/etc/gcloud/service_account_key.json \
        -v "$HOME/.config/gcloud/bucket-service-account.json:/etc/gcloud/service_account_key.json" \
        google-storage-http
