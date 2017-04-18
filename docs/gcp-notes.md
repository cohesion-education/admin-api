

Click on the "Generate new P12 key" to generate and download a new private key. Once you download the P12 file, use the following command to convert it into a PEM file.

    $ openssl pkcs12 -in key.p12 -passin pass:notasecret -out key.pem -nodes

## Resources
* https://github.com/GoogleCloudPlatform/google-cloud-go
* https://godoc.org/cloud.google.com/go/datastore
* https://cloud.google.com/storage/docs/access-control/signed-urls
* https://cloud.google.com/storage/docs/access-control/create-signed-urls-program
