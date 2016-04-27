# ADtoLDAP
This program will Gather Results based from Active Directory, or even another openldap server based on attributes specified in /etc/ldapsync.ini, and sync it to the second ldap server. For Active directory to LDAP syncing, we need to make sure that the schema of the openldap server is prepared to accomodate the additional attibutes AD incorporates (an example would be the `memberOf:` attribute)

#### How to:
Enable `memberOf` attribute in ldap, to accomodate AD field, by using the 3 ldif files included in the repo.
(Thanks to https://technicalnotes.wordpress.com/2014/04/19/openldap-setup-with-memberof-overlay/)

```
ldapadd -Q -Y EXTERNAL -H ldapi:/// -f memberof_load_configure.ldif 
ldapmodify -Q -Y EXTERNAL -H ldapi:/// -f 1refint.ldif
ldapadd -Q -Y EXTERNAL -H ldapi:/// -f 2refint.ldif
```


Sample /etc/ldapsync.ini file for syncing from Active directory to an openldap server):
```
[ADServer]
ADHost = <AD Server IP>
ADPort = 389
#Page the result size to prevent possible OOM error and crash
ADPage = 500
#AD Connection Timeout in seconds (Defaults to 10)
ADConnTimeOut = 10
username = cn=ADUser,cn=Users,dc=example,dc=com
password = somepassword
basedn = ou=someOu,dc=example,dc=com
#Attributes required to be pulled
attr = comment, givenName, unixHomeDirectory, sn, loginShell, memberOf, dn, o, uid, objectclass, cn, displayName, cn, uidNumber, gidNumber, member
#ldap filter
filter = (cn=*)

[LDAPServer]
LDAPHost = <ldap server ip>
LDAPPort = 389
#Page LDAP result
LDAPPage = 500
#LDAP Connection Timeout in seconds (Defaults to 10)
LDAPConnTimeOut = 10
username = cn=someldapuser,dc=example,dc=com
password = someldappasswd
basedn = ou=someOu,dc=example,dc=com
attr = comment, givenName, homeDirectory, sn, loginShell, memberOf, dn, o, uid, objectclass, cn, displayName, cn, uidNumber, gidNumber, memberUid
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

`./go-sync --init`

The results can be verified, before the sync can be run continuously

`./go-sync --sync`

This is not daemonised from within (see http://stackoverflow.com/questions/10067295/how-to-start-a-go-program-as-a-daemon-in-ubuntu) and has to use an upstart or init script.(TODO)


##### Sample /etc/ldapsync.ini for syncing from one openldap server to another openldap server

```
[ADServer]
ADHost = <LDAP server1 IP>
ADPort = 389
#Page the result size to prevent possible OOM error and crash
ADPage = 500
#AD Connection Timeout in seconds (Defaults to 10)
ADConnTimeOut = 10
username = cn=someuser,dc=example,dc=com
password = somepassword1
basedn = ou=SomeOU,dc=example,dc=com
#Attributes required to be pulled
#attr = comment, givenName, unixHomeDirectory, sn, loginShell, memberOf, dn, o, uid, objectclass, cn, displayName, cn, uidNumber, gidNumber, member
attr = comment, givenName, homeDirectory, sn, loginShell, memberOf, dn, o, uid, objectclass, cn, displayName, cn, uidNumber, gidNumber, memberUid
#ldap filter
filter = (cn=*)

[LDAPServer]
LDAPHost = 127.0.0.1
LDAPPort = 389
#Page LDAP result
LDAPPage = 500
#LDAP Connection Timeout in seconds (Defaults to 10)
LDAPConnTimeOut = 10
username = cn=SomeUser1,dc=example,dc=com
password = somepassword2
basedn = ou=SomeOU,dc=example,dc=com
#attr = "distinguishedName, comment, givenName, primaryGroupID, unixHomeDirectory, sn, loginShell, memberOf, dn, o, uid, objectclass, cn, displayName, cn, uidNumber, gidNumber"
#attr = *
attr = comment, givenName, homeDirectory, sn, loginShell, memberOf, dn, o, uid, objectclass, cn, displayName, cn, uidNumber, gidNumber, memberUid
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

