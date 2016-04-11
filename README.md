# ADtoLDAP
Gather AD Results based on attributes and sync to LDAP

Enable `memberOf` attribute in ldap, to accomodate AD field, by using the 3 ldif files included in the repo.
```
ldapadd -Q -Y EXTERNAL -H ldapi:/// -f memberof_load_configure.ldif 
ldapmodify -Q -Y EXTERNAL -H ldapi:/// -f 1refint.ldif
ldapadd -Q -Y EXTERNAL -H ldapi:/// -f 2refint.ldif
```


Sample ini file:
```
[ADServer]
ADHost = <AD IP>
ADPort = 389
#Page the result size to prevent possible OOM error and crash
ADPage = 500
#AD Connection Timeout in seconds (Defaults to 10)
ADConnTimeOut = 10
username = cn=SomeUser,dc=example,dc=com
password = somepasswd
basedn = ou=someou,dc=example,dc=com
#Attributes required to be pulled
attr = comment, givenName, unixHomeDirectory, sn, loginShell, memberOf, dn, o, uid, objectclass, cn, displayName, cn, uidNumber, gidNumber, member
#ldap filter
filter = (cn=*)

[LDAPServer]
LDAPHost = <OpenLdapServerIP>
LDAPPort = 389
#Page LDAP result
LDAPPage = 500
#LDAP Connection Timeout in seconds (Defaults to 10)
LDAPConnTimeOut = 10
username = cn=somelinuxuser,dc=example,dc=com
password = someldappasswd
basedn = ou=someotherbasedn,dc=example,dc=com
attr = "distinguishedName, comment, givenName, primaryGroupID, unixHomeDirectory, sn, loginShell, memberOf, dn, o, uid, objectclass, cn, displayName, cn, uidNumber, gidNumber"
filter = (cn=*)


[Replace]
userObjectClass = posixAccount,top,inetOrgPerson
groupObjectClass = top,posixGroup
[Map]
#ADAttribute = ldapattribute
unixHomeDirectory = homeDirectory
#specify member mapping if you are selecting member attribute from *attr above
member = memberUid
```
