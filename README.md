# Xml Digital Signature using goxmldsig

This is a sample code to sign an verify xml digital signature

## create certificate and key

We will use `cfssl` to generate self signed certificates to sign the xml document.

The steps below will guide you through the process.

1. First we need to generate the Certificate Authority.  Prepare the certificate configuration.

``` json
{
  "CN": "Punggol CA",
  "key": {
    "algo": "rsa",
    "size": 2048
  },
  "names": [
  {
    "C": "SG",
    "L": "Singapore",
    "O": "Main Office",
    "OU": "Main Office Root CA"
  }
 ]
}
```

See [`ca.json`](templates/ca.json)

``` shell
cfssl gencert -initca templates/ca.json | cfssljson -bare generated-certs/ca
```

2. Create the `profile` to determine how it signs certificates.

``` json
{
  "signing": {
    "default": {
      "expiry": "26298h"
    },
    "profiles": {
      "client": {
        "usages": [
          "signing",
          "digital signature",
          "key encipherment", 
          "client auth"
        ],
        "expiry": "26298h"
      }
    }
  }
}
```

See [`cfssl.json`](templates/cfssl.json)

3. Generate the client certificates and certificate signing requests.

``` json
{
    "CN": "Client Connect",
    "key": {
      "algo": "rsa",
      "size": 2048
    },
    "names": [
      {
        "C":  "SG",
        "L":  "Punggol",
        "O":  "Kopitiam",
        "OU": "Western stall"
      }
    ],
    "ca": {
      "expiry": "42720h"
    }
}
```
See [`client.json`](templates/client.json)

``` shell
cfssl gencert -ca=generated-certs/ca.pem -ca-key=generated-certs/ca-key.pem -config=templates/cfssl.json -profile=client templates/client.json | cfssljson -bare generated-certs/client
```

4. Sign the CSR using the certificate authority certificates.

``` shell
cfssl sign -ca generated-certs/ca.pem -ca-key generated-certs/ca-key.pem -config templates/cfssl.json -profile client generated-certs/client.csr | cfssljson -bare generated-certs/client
```