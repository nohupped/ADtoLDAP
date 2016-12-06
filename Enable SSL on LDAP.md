# Enable SSL on an openldap server. 

(Referenced from https://www.server-world.info/en/note?os=Debian_8&p=openldap&f=4)
```
### Generate Certificate:
root@dlp:~# cd /etc/ssl/private 
root@dlp:/etc/ssl/private# openssl genrsa -aes128 -out server.key 2048 
Generating RSA private key, 2048 bit long modulus
...................+++
.....+++
e is 65537 (0x10001)
Enter pass phrase for server.key:# set passphrase
Verifying - Enter pass phrase for server.key:# confirm
##### remove passphrase from private key

root@dlp:/etc/ssl/private# openssl rsa -in server.key -out server.key 
Enter pass phrase for server.key:# input passphrase
writing RSA key

root@www:/etc/ssl/private# openssl req -new -days 3650 -key server.key -out server.csr 
You are about to be asked to enter information that will be incorporated
into your certificate request.
What you are about to enter is what is called a Distinguished Name or a DN.
There are quite a few fields but you can leave some blank
For some fields there will be a default value,
If you enter '.', the field will be left blank.

Country Name (2 letter code) [AU]:IN# country
State or Province Name (full name) [Some-State]:SomeState   # state
Locality Name (eg, city) []:SomeState# city
Organization Name (eg, company) [Internet Widgits Pty Ltd]:SomeCompany   # company
Organizational Unit Name (eg, section) []:SomeOU   # department
Common Name (e.g. server FQDN or YOUR name) []:mycompany.example.com   # server's FQDN
Email Address []:nohupped@gmail.com# email address
Please enter the following 'extra' attributes
to be sent with your certificate request
A challenge password []:
An optional company name []:
root@www:/etc/ssl/private# openssl x509 -in server.csr -out server.crt -req -signkey server.key -days 3650 
Signature ok
subject=/C=IN/ST=SomeState/L=SomeState/O=SomeCompany/OU=SomeOU/CN=mycompany.example.com/emailAddress=nohupped@gmail.com
Getting Private key
root@dlp:/etc/ssl/private# chmod 400 server.*
```
### Configure LDAP over TLS to make connection be secure

The mod_ssl.ldif file mentioned below is attached inside the `ldifs` directory.
```
root@dlp:~# cp /etc/ssl/private/server.key \
/etc/ssl/private/server.crt \
/etc/ssl/certs/ca-certificates.crt \
/etc/ldap/sasl2/ 
root@dlp:~# chown openldap. /etc/ldap/sasl2/server.key \
/etc/ldap/sasl2/server.crt \
/etc/ldap/sasl2/ca-certificates.crt
root@dlp:~# vi mod_ssl.ldif
# create new
 dn: cn=config
changetype: modify
add: olcTLSCACertificateFile
olcTLSCACertificateFile: /etc/ldap/sasl2/ca-certificates.crt
-
replace: olcTLSCertificateFile
olcTLSCertificateFile: /etc/ldap/sasl2/server.crt
-
replace: olcTLSCertificateKeyFile
olcTLSCertificateKeyFile: /etc/ldap/sasl2/server.key

root@dlp:~# ldapmodify -Y EXTERNAL -H ldapi:/// -f mod_ssl.ldif 
SASL/EXTERNAL authentication started
SASL username: gidNumber=0+uidNumber=0,cn=peercred,cn=external,cn=auth
SASL SSF: 0
modifying entry "cn=config"
```
Make sure you see `modifying entry "cn=config"` after the ldapmodify. Else the changes wouldn't be reflected.
```
root@dlp:~# vi /etc/default/slapd
# line 24: add
SLAPD_SERVICES="ldap:/// ldapi:/// ldaps:///"
root@dlp:~# systemctl restart slapd 
```