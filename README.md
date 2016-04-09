# ADtoLDAP
Gather AD Results based on attributes and sync to LDAP


Sample ini file:
```
[ADServer]
ADHost = <AD Host IP>
ADPort = 389
#Page the result size to prevent possible OOM error and crash
ADPage = 500
#AD Connection Timeout in seconds (Defaults to 10)
ADConnTimeOut = 10
username = cn=linuxuser,cn=Users,dc=example,dc=com
password = somepasswd
basedn = ou=something,dc=example,dc=com
#Attributes required to be pulled
attr = comment, givenName, unixHomeDirectory, sn, loginShell, memberOf, dn, o, uid, objectclass, cn, displayName, cn, uidNumber, gidNumber
#ldap filter
filter = (cn=*)

[LDAPServer]
LDAPHost = <ldap server ip>
LDAPPort = 389
#Page LDAP result
LDAPPage = 500
#LDAP Connection Timeout in seconds (Defaults to 10)
LDAPConnTimeOut = 10
#username = cn=linuxauth,cn=Users,dc=internal,dc=media,dc=net
username = cn=someotheruser,dc=example,dc=com
password = someotherpasswd
basedn = ou=something,dc=example,dc=com
attr = distinguishedName, comment, givenName, primaryGroupID, unixHomeDirectory, sn, loginShell, memberOf, dn, o, uid, objectclass, cn, displayName, cn, uidNumber, gidNumber
filter = (cn=*)


[Replace]
userObjectClass = posixAccount,top,inetOrgPerson
groupObjectClass = top,posixGroup
[Map]
#format is to map ADAttribute = ldapattribute
unixHomeDirectory = homeDirectory
```
