package ca

func SelfSignCaConf() string {
	return `
HOME			= .
openssl_conf = openssl_init
config_diagnostics = 1
oid_section = new_oids

[ new_oids ]
# Policies used by the TSA examples.
tsa_policy1 = 1.2.3.4.1
tsa_policy2 = 1.2.3.4.5.6
tsa_policy3 = 1.2.3.4.5.7

[openssl_init]
providers = provider_sect

[provider_sect]
default = default_sect

[default_sect]
# activate = 1
####################################################################
[ ca ]
default_ca	= CA_default		# The default ca section
####################################################################
[ CA_default ]
dir             = %s		# Where everything is kept
certs           = $dir/certs		# Where the issued certs are kept
crl_dir         = $dir/crl		# Where the issued crl are kept
database        = $dir/index.txt	# database index file.
#unique_subject = no			# Set to 'no' to allow creation of several certs with same subject.
new_certs_dir   = $dir/newcerts		# default place for new certs.
certificate     = $dir/cacert.pem 	# The CA certificate
serial          = $dir/serial 		# The current serial number
crlnumber       = $dir/crlnumber	# the current crl number must be commented out to leave a V1 CRL
crl             = $dir/crl.pem 		# The current CRL
private_key     = $dir/private/cakey.pem # The private key
x509_extensions = usr_cert		# The extensions to add to the cert
# Comment out the following two lines for the "traditional" (and highly broken) format.
name_opt        = ca_default		# Subject Name options
cert_opt        = ca_default		# Certificate field options
default_days    = 365			# how long to certify for
default_crl_days= 30			# how long before next CRL
default_md      = default		# use public key default MD
preserve        = no			# keep passed DN ordering
policy          = policy_match

[ policy_match ]
countryName         = match
stateOrProvinceName = match
organizationName    = match
organizationalUnitName = optional
commonName          = supplied
emailAddress        = optional


[ policy_anything ]
countryName         = optional
stateOrProvinceName = optional
localityName        = optional
organizationName    = optional
organizationalUnitName = optional
commonName          = supplied
emailAddress        = optional

####################################################################
[ req ]
prompt              = no
default_bits		= 2048
#default_keyfile 	= privkey.pem
distinguished_name	= req_distinguished_name
#attributes          = req_attributes
x509_extensions     = v3_ca  # The extensions to add to the self signed cert

# Passwords for private keys if not present they will be prompted for
# input_password = secret
# output_password = secret

string_mask         = utf8only

# req_extensions = v3_req # The extensions to add to a certificate request

[ req_distinguished_name ]
countryName             = %s
#countryName_default     = CN
#countryName_min			= 2
#countryName_max			= 2
stateOrProvinceName     = %s
#stateOrProvinceName_default	= Some-State
localityName            = %s
0.organizationName      = %s
#0.organizationName_default = Internet Widgits Pty Ltd
# we can do this but it is not needed normally :-)
#1.organizationName     = Second Organization Name (eg, company)
#1.organizationName_default = World Wide Web Pty Ltd
organizationalUnitName  = %s
#organizationalUnitName_default	=
commonName              = %s
#commonName_max          = 64
#emailAddress            = 
#emailAddress_max		= 64

# SET-ex3			= SET extension number 3

[ req_attributes ]
challengePassword           = A challenge password
challengePassword_min       = 4
challengePassword_max       = 20

unstructuredName            = An optional company name

[ usr_cert ]
# These extensions are added when 'ca' signs a request.
# This goes against PKIX guidelines but some CAs do it and some software
# requires this to avoid interpreting an end user certificate as a CA.

basicConstraints=CA:FALSE

# This is typical in keyUsage for a client certificate.
# keyUsage = nonRepudiation, digitalSignature, keyEncipherment

# PKIX recommendations harmless if included in all certificates.
subjectKeyIdentifier=hash
authorityKeyIdentifier=keyid,issuer

# This stuff is for subjectAltName and issuerAltname.
# Import the email address.
# subjectAltName=email:copy
# An alternative to produce certificates that aren't
# deprecated according to PKIX.
# subjectAltName=email:move

# Copy subject details
# issuerAltName=issuer:copy

# This is required for TSA certificates.
# extendedKeyUsage = critical,timeStamping

[ v3_req ]
# Extensions to add to a certificate request
basicConstraints = CA:FALSE
keyUsage = nonRepudiation, digitalSignature, keyEncipherment

[ v3_ca ]
# Extensions for a typical CA
# PKIX recommendation.
subjectKeyIdentifier=hash
authorityKeyIdentifier=keyid:always,issuer
basicConstraints = critical,CA:true

# Key usage: this is typical for a CA certificate. However since it will
# prevent it being used as an test self-signed certificate it is best
# left out by default.
# keyUsage = cRLSign, keyCertSign

# Include email address in subject alt name: another PKIX recommendation
# subjectAltName=email:copy
# Copy issuer details
# issuerAltName=issuer:copy

# DER hex encoding of an extension: beware experts only!
# obj=DER:02:03
# Where 'obj' is a standard or added object
# You can even override a supported extension:
# basicConstraints= critical, DER:30:03:01:01:FF

[ crl_ext ]
# CRL extensions.
# Only issuerAltName and authorityKeyIdentifier make any sense in a CRL.
# issuerAltName=issuer:copy
authorityKeyIdentifier=keyid:always

[ proxy_cert_ext ]
# These extensions should be added when creating a proxy certificate

# This goes against PKIX guidelines but some CAs do it and some software
# requires this to avoid interpreting an end user certificate as a CA.

basicConstraints=CA:FALSE

# This is typical in keyUsage for a client certificate.
# keyUsage = nonRepudiation, digitalSignature, keyEncipherment

# PKIX recommendations harmless if included in all certificates.
subjectKeyIdentifier=hash
authorityKeyIdentifier=keyid,issuer

# This stuff is for subjectAltName and issuerAltname.
# Import the email address.
# subjectAltName=email:copy
# An alternative to produce certificates that aren't
# deprecated according to PKIX.
# subjectAltName=email:move

# Copy subject details
# issuerAltName=issuer:copy

# This really needs to be in place for it to be a proxy certificate.
proxyCertInfo=critical,language:id-ppl-anyLanguage,pathlen:3,policy:foo

####################################################################
[ tsa ]

default_tsa = tsa_config1	# the default TSA section

[ tsa_config1 ]

# These are used by the TSA reply generation only.
dir		= ./demoCA		# TSA root directory
serial		= $dir/tsaserial	# The current serial number (mandatory)
crypto_device	= builtin		# OpenSSL engine to use for signing
signer_cert	= $dir/tsacert.pem 	# The TSA signing certificate
					# (optional)
certs		= $dir/cacert.pem	# Certificate chain to include in reply
					# (optional)
signer_key	= $dir/private/tsakey.pem # The TSA private key (optional)
signer_digest  = sha256			# Signing digest to use. (Optional)
default_policy	= tsa_policy1		# Policy if request did not specify it
					# (optional)
other_policies	= tsa_policy2, tsa_policy3	# acceptable policies (optional)
digests     = sha1, sha256, sha384, sha512  # Acceptable message digests (mandatory)
accuracy	= secs:1, millisecs:500, microsecs:100	# (optional)
clock_precision_digits  = 0	# number of digits after dot. (optional)
ordering		= yes	# Is ordering defined for timestamps?
				# (optional, default: no)
tsa_name		= yes	# Must the TSA name be included in the reply?
				# (optional, default: no)
ess_cert_id_chain	= no	# Must the ESS cert id chain be included?
				# (optional, default: no)
ess_cert_id_alg		= sha1	# algorithm to compute certificate
				# identifier (optional, default: sha1)

[insta] # CMP using Insta Demo CA
# Message transfer
server = pki.certificate.fi:8700
# proxy = # set this as far as needed, e.g., http://192.168.1.1:8080
# tls_use = 0
path = pkix/

# Server authentication
recipient = "/C=FI/O=Insta Demo/CN=Insta Demo CA" # or set srvcert or issuer
ignore_keyusage = 1 # potentially needed quirk
unprotected_errors = 1 # potentially needed quirk
extracertsout = insta.extracerts.pem

# Client authentication
ref = 3078 # user identification
secret = pass:insta # can be used for both client and server side

# Generic message options
cmd = ir # default operation, can be overridden on cmd line with, e.g., kur

# Certificate enrollment
subject = "/CN=openssl-cmp-test"
newkey = insta.priv.pem
out_trusted = apps/insta.ca.crt # does not include keyUsage digitalSignature
certout = insta.cert.pem

[pbm] # Password-based protection for Insta CA
# Server and client authentication
ref = $insta::ref # 3078
secret = $insta::secret # pass:insta

[signature] # Signature-based protection for Insta CA
# Server authentication
trusted = $insta::out_trusted # apps/insta.ca.crt

# Client authentication
secret = # disable PBM
key = $insta::newkey # insta.priv.pem
cert = $insta::certout # insta.cert.pem

[ir]
cmd = ir

[cr]
cmd = cr

[kur]
# Certificate update
cmd = kur
oldcert = $insta::certout # insta.cert.pem

[rr]
# Certificate revocation
cmd = rr
oldcert = $insta::certout # insta.cert.pem
`
}

var confReq = `[ req ]
prompt             = no
default_bits       = 2048
distinguished_name = req_distinguished_name
`
var confReqExt = `
req_extensions     = req_ext
`
var confReqDistName = `[ req_distinguished_name ]
countryName                    = %s
stateOrProvinceName            = %s
localityName                   = %s
organizationName               = %s
organizationalUnitName         = %s
commonName                     = %s
`

var confReqExtDetail = `[ req_ext ]
subjectAltName = @alt_names

[alt_names]
DNS.1 = %s
`

func SignServerCertConf() string {
	return confReq + confReqExt + confReqDistName + confReqExtDetail
}
func SignClientCertConf() string {
	return confReq + confReqDistName
}
