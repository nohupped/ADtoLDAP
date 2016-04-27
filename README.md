# ADtoLDAP
Gather AD Results based on attributes and sync to LDAP

Enable `memberOf` attribute in ldap, to accomodate AD field, by using the 3 ldif files included in the repo.
(Thanks to https://technicalnotes.wordpress.com/2014/04/19/openldap-setup-with-memberof-overlay/)

```
ldapadd -Q -Y EXTERNAL -H ldapi:/// -f memberof_load_configure.ldif 
ldapmodify -Q -Y EXTERNAL -H ldapi:/// -f 1refint.ldif
ldapadd -Q -Y EXTERNAL -H ldapi:/// -f 2refint.ldif
```


Sample /etc/ldapsync.ini file:
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
