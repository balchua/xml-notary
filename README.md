# Xml Notary

Xml Digital Signature using goxmldsig

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


## Usage

Starting the server

``` shell
sign-xml-doc serve -c ./generated-certs/client.pem -k ./generated-certs/client-key.pem 

2023-01-22T09:05:33.976+0800    INFO    cmd/root.go:56  application started
2023-01-22T09:05:33.976+0800    INFO    cmd/serve.go:44 list all arguments      {"port": 5000, "certFile": "./generated-certs/client.pem", "keyFile": "./generated-certs/client-key.pem"}

 ┌───────────────────────────────────────────────────┐ 
 │                    Xml Notary                     │ 
 │                   Fiber v2.41.0                   │ 
 │               http://127.0.0.1:5000               │ 
 │       (bound on host 0.0.0.0 and port 5000)       │ 
 │                                                   │ 
 │ Handlers ............. 1  Processes ........... 1 │ 
 │ Prefork ....... Disabled  PID ............. 26435 │ 
 └───────────────────────────────────────────────────┘ 

```

### Endpoints

#### Sign an xml document

``` shell
curl --location --request POST 'http://localhost:5000/api/sign' \
--header 'Content-Type: application/xml' \
--data-raw '<test>
<hello/>
</test>'

```

Response

``` xml
<test>
<hello/>
<ds:Signature xmlns:ds="http://www.w3.org/2000/09/xmldsig#"><ds:SignedInfo><ds:CanonicalizationMethod Algorithm="http://www.w3.org/2006/12/xml-c14n11"/><ds:SignatureMethod Algorithm="http://www.w3.org/2001/04/xmldsig-more#rsa-sha256"/><ds:Reference URI=""><ds:Transforms><ds:Transform Algorithm="http://www.w3.org/2000/09/xmldsig#enveloped-signature"/><ds:Transform Algorithm="http://www.w3.org/2006/12/xml-c14n11"/></ds:Transforms><ds:DigestMethod Algorithm="http://www.w3.org/2001/04/xmlenc#sha256"/><ds:DigestValue>yCSuuZIFBfaUQKisRpp0UVq1020iROjq1J+z0YPr8kA=</ds:DigestValue></ds:Reference></ds:SignedInfo><ds:SignatureValue>LSmSIkesXixWQV1bUsW9Kh7tYgn221lhVjpWcAwItpkGPXVma2NQ+J4CIQ2gJXo141TSsNN658czSHQ5tmZHrN/9DjZ2vE90+GyR4TTjF+AjoeJvv8T4uGRYRcWJrBmuV1hglTR+m26dx+eOxwXZDpaltJSfRudbM2DTCD52kNRZ44gU7hiU1FC5eqt3Cb/rHm0Fk8xT6G7MpgCIeh2nIR96gJsaQCIsXgxnTsJdTTSYNxKPyCy+80qxBK1NH3n0nkErL7cvPF1QoY1E1KuOFGtamtdkaAnB1yVjyNidn7diycYymWuU4LmFMU7AmiGzFPY8Uc/+W6hwoUbVaTCpfQ==</ds:SignatureValue><ds:KeyInfo><ds:X509Data><ds:X509IssuerSerial><ds:X509IssuerName>CN=Punggol CA,OU=Main Office Root CA,O=Main Office,L=Singapore,C=SG</ds:X509IssuerName><ds:X509SerialNumber/></ds:X509IssuerSerial><ds:X509Certificate>MIID0DCCArigAwIBAgIUOEJIe8cP/hj0elNle9Cvg3fcyckwDQYJKoZIhvcNAQELBQAwajELMAkGA1UEBhMCU0cxEjAQBgNVBAcTCVNpbmdhcG9yZTEUMBIGA1UEChMLTWFpbiBPZmZpY2UxHDAaBgNVBAsTE01haW4gT2ZmaWNlIFJvb3QgQ0ExEzARBgNVBAMTClB1bmdnb2wgQ0EwHhcNMjMwMTIxMDQ1MTAwWhcNMjYwMTIwMjI1MTAwWjBjMQswCQYDVQQGEwJTRzEQMA4GA1UEBxMHUHVuZ2dvbDERMA8GA1UEChMIS29waXRpYW0xFjAUBgNVBAsTDVdlc3Rlcm4gc3RhbGwxFzAVBgNVBAMTDkNsaWVudCBDb25uZWN0MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAoLu5efRobGbr9AuIimfHdTCxLeQ7h0ulzJsvAuOyU0YfBZW5X0sEVFwRrBChbDFoDQp1/PArKiHmZddlRXhjW8Yx3aTQm6xRYDkHZjFMjhc1bVofQudZxwdcgle/RJOrJblm6gYCUVigaURAAW4EBpUMaHcLnRnfb1Kj0FkCOJ8Ge92NpZCRvSVvU2GwZW9SekOE5Ncs1aZ7WFjliyuUYKpRtRvyOnSD/QXqobf7xjKrk3GuC63hqVIleov7Q9Vn89ExFY9FZsJ6MM6QHBqo+nTC8RcW3OoOfGfLvReYK7rapCVF477ESYCUUnL2DyRYW2vaQcJUX7Y4OvHwgDMG6wIDAQABo3UwczAOBgNVHQ8BAf8EBAMCAqQwEwYDVR0lBAwwCgYIKwYBBQUHAwIwDAYDVR0TAQH/BAIwADAdBgNVHQ4EFgQU7TrcDmUU9UMEavPfoPMdamN8TEEwHwYDVR0jBBgwFoAUDcLavYdcRtJ9shg6LiohwKsebfYwDQYJKoZIhvcNAQELBQADggEBAIymbbwVsxUftUXO52J+DbdpkmAbr/Nh1QVboKhfe3mhley9W/TYTUmA7baoP4gTzpHnZh7bG34Ez+nQ8iW6jZdzL7+FNix49KxlpOFA9EzEmfYyN5mSctT2DgtuY9Ix1rO3exv/75eATUiX2bo0GfegRFZmTwoaug00GTFwK54QDKSNemSfO7yq1H7l/vyqilvEWJmxMusZ/4PriEIC4ud+K/ULGgS1eScHI08gTTpoRnbHh0bPzAuxjonuvKdjohIiLgGvolMPpzbwLFbJDjlhu9q7dacoid6tnHxGNq5fUnIiiubW4AjmBKdcoyXvkaMBE8mXqkEEKcUMK7CwR1c=</ds:X509Certificate></ds:X509Data></ds:KeyInfo></ds:Signature></test>
```