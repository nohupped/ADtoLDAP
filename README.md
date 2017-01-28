# ADtoLDAP
This program will gather results from Active Directory, or another openldap server based on the attributes specified in /etc/ldapsync.ini, and sync it to the second ldap server. The `basedn` which must be synced from, for eg:`basedn = ou=someOu,dc=example,dc=com` in the sample configuration below must be created on the destination server to acommodate the sync. For Active directory to LDAP syncing, we need to make sure that the schema of the openldap server is prepared to accomodate the additional attibutes AD incorporates, if we are syncing them. (an example would be the `memberOf:` attribute) Better - omit those unless required.
This can run over an encrypted connection as well. 

To use `TLS`, make sure to add the domain name for which the AD certificate is generated. If that fails, the program panics throwing 

```
panic: LDAP Result Code 200 "": x509: certificate is valid for example1.domain.com, example2.domain.com, EXAMPLE, not examples.domains.com
```

Using `TLS` will make it hard for decrypting the data transferred over wire. Without using TLS, the data can be viewed with a packet capturing program like tcpdump like
```
tcpdump -v -XX
```
####Requirements to set up TLS connection:

##### Get the pem file from the AD server:
From the windows server cmd, do 
```
certutil  -ca.cert ca_name.cer > ca.crt
```
This will generate the pem file, and will be saved in the working directory by the name ca.crt.
This pem file must be copied over from the master/AD server to the slave/openldap server, and the path to this file must be mentioned in the ldapsync.ini file to create a custom cert pool and use it as the Root CAs, so the DialTLS wouldn't panic with a `certificate signed by unknown authority` error.


#### How to:
##### install:
```
go get github.com/nohupped/ADtoLDAP
```
A custom Daemonizer is provided in the Daemonizer directory, that daemonize itself, forks again and runs the program, and capture any errors or panics that the program throws to the `STDOUT/STDERR`, and logs it to the syslog. Compile it as 
```
gcc  -W -Wall ./main.c ./src/ForkSelf.c -o daemonizer
```
The program can be daemonized as 
```
<path/to>/daemonizer <path/to>/ADtoLDAP
```

Enable `memberOf` attribute in ldap (required only if we are syncing it), to accomodate the equivalent AD field, by using the 3 ldif files included in the repo.
(Thanks to https://technicalnotes.wordpress.com/2014/04/19/openldap-setup-with-memberof-overlay/)

```
ldapadd -Q -Y EXTERNAL -H ldapi:/// -f memberof_load_configure.ldif 
ldapmodify -Q -Y EXTERNAL -H ldapi:/// -f 1refint.ldif
ldapadd -Q -Y EXTERNAL -H ldapi:/// -f 2refint.ldif
```

##### The permission of /etc/ldapsync.ini

The program checks if the file permissions for /etc/ldapsync.ini are too broad. If it is not 600, the program will report that, and will not start. This is checked to prevent a non-privileged read of the username and password used to bind to both the servers, which are stored in this configuration file. This can be over-ridden by running the program with the flag `--safe=false`.

##### 
/etc/ldapsync.ini file for syncing from Active directory to an openldap server:
```
[ADServer]
ADHost = <AD Server IP>
#ADPort = 389 for non ssl and 636 for ssl
ADPort = 636
UseTLS = true
# set InsecureSkipVerify to true for testing, to accept the certificate without any verification.
InsecureSkipVerify = false
#CRTValidFor will not be honored if InsecureSkipVerify is set to true.
CRTValidFor = example1.domain.com
#Path to the pem file, which is used to create the custom CA pool. Will not be honored if InsecureSkipVerify is set to true.
CRTPath = /etc/ldap.crt
#Page the result size to prevent possible OOM error and crash
ADPage = 500
#AD Connection Timeout in seconds (Defaults to 10)
ADConnTimeOut = 10
username = cn=ADUser,cn=Users,dc=example,dc=com
password = somepassword
basedn = ou=someOu,dc=example,dc=com
#Attributes required to be pulled
attr = givenName, unixHomeDirectory, sn, loginShell, memberOf, dn, o, uid, objectclass, cn, displayName, cn, uidNumber, gidNumber, member
#ldap filter
filter = (cn=*)

[LDAPServer]
LDAPHost = <ldap server ip>
LDAPPort = 636
UseTLS = true
InsecureSkipVerify = false
CRTValidFor = ldapserver.example.com
CRTPath = /etc/ldap/sasl2/server.crt
#Page LDAP result
LDAPPage = 500
#LDAP Connection Timeout in seconds (Defaults to 10)
LDAPConnTimeOut = 10
username = cn=someldapuser,dc=example,dc=com
password = someldappasswd
basedn = ou=someOu,dc=example,dc=com
attr = givenName, homeDirectory, sn, loginShell, memberOf, dn, o, uid, objectclass, cn, displayName, cn, uidNumber, gidNumber, memberUid
filter = (cn=*)


[Replace]
userObjectClass = posixAccount,top,inetOrgPerson
groupObjectClass = top,posixGroup
[Map]
#ADAttribute = ldapattribute, that is mapping AD attribute to the relevant ldap attribute
unixHomeDirectory = homeDirectory
#specify member mapping if you are selecting member attribute from *attr above
member = memberUid

[Sync]
#Add sleep time after each successful sync, in seconds.
sleepTime = 5

```
Do the initial run which will do the initial population 

`./ADtoLDAP --safe=false --sync=once --configfile=/etc/ldapsync.ini`

--safe=false will omit the config file permission checking.
--sync=once will do the initial run once and will exit.

The results can be verified, before the sync can be run continuously

`./ADtoLDAP`

The default options are:

--safe=true is the default behaviour, so the program will verify the config file's permission to make sure it is non-readable by groups and others.
--sync=daemon will continuously run the program in foreground only (see http://stackoverflow.com/questions/10067295/how-to-start-a-go-program-as-a-daemon-in-ubuntu). Use the daemonizer to daemonize it.
--config-file, if not specified, takes the default path to look for, as `/etc/ldapsync.ini`. This can be over-ridden with this argument.

##### Sample /etc/ldapsync.ini for syncing from one openldap server to another openldap server

```
[ADServer]
ADHost = <LDAP server1 IP>
ADPort = 636
UseTLS = true
# set InsecureSkipVerify to true for testing, to accept the certificate without any verification.
InsecureSkipVerify = false
#CRTValidFor will not be honored if InsecureSkipVerify is set to true.
CRTValidFor = example1.domain.com
#Path to the pem file, which is used to create the custom CA pool. Will not be honored if InsecureSkipVerify is set to true.
CRTPath = /etc/ldap.crt
#Page the result size to prevent possible OOM error and crash
ADPage = 500
#AD Connection Timeout in seconds (Defaults to 10)
ADConnTimeOut = 10
username = cn=someuser,dc=example,dc=com
password = somepassword1
basedn = ou=SomeOU,dc=example,dc=com
#Attributes required to be pulled
#attr = comment, givenName, unixHomeDirectory, sn, loginShell, memberOf, dn, o, uid, objectclass, cn, displayName, cn, uidNumber, gidNumber, member
attr = givenName, homeDirectory, sn, loginShell, memberOf, dn, o, uid, objectclass, cn, displayName, cn, uidNumber, gidNumber, memberUid
#ldap filter
filter = (cn=*)

[LDAPServer]
LDAPHost = 127.0.0.1
LDAPPort = 636
UseTLS = true
InsecureSkipVerify = false
CRTValidFor = ldapserver.example.com
CRTPath = /etc/ldap/sasl2/server.crt
#Page LDAP result
LDAPPage = 500
#LDAP Connection Timeout in seconds (Defaults to 10)
LDAPConnTimeOut = 10
username = cn=SomeUser1,dc=example,dc=com
password = somepassword2
basedn = ou=SomeOU,dc=example,dc=com
#attr = "distinguishedName, comment, givenName, primaryGroupID, unixHomeDirectory, sn, loginShell, memberOf, dn, o, uid, objectclass, cn, displayName, cn, uidNumber, gidNumber"
#attr = *
attr = givenName, homeDirectory, sn, loginShell, memberOf, dn, o, uid, objectclass, cn, displayName, cn, uidNumber, gidNumber, memberUid
filter = (cn=*)


[Replace]
userObjectClass = posixAccount,top,inetOrgPerson
groupObjectClass = top,posixGroup
[Map]
#ADAttribute = ldapattribute, that is mapping AD attribute to the relevant ldap attribute
#unixHomeDirectory = homeDirectory
#specify member mapping if you are selecting member attribute from *attr above
#member = memberUid

[Sync]
#Add sleep time after each successful sync, in seconds.
sleepTime = 60

```

Now we need to create index for the frequently accessed attributes in ldap. A sample ldif file with a few of the attributes are added into the ldif directory. Run the query

```ldapmodify -Q -Y EXTERNAL -H ldapi:/// -f addindex.ldif```

##### Monitoring
A monitoring module written in python 3 is provided in the `LdapSyncMonitor/SyncMonitor` directory, that seeks to the end of the log file, reads backwards until it finds another newline character(doing this because of the size of logs it can generate), and from the captured line, takes the timestamp, and do the math with the warning and critical thresholds that the class `Monitor` accepts, and exits with relevent exit code suitable for nagios. A sample usage can be found in `LdapSyncMonitor/monitorDaemon.py` script. 
