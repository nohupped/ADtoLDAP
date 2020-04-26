# [Thank you for the HowTo reference.](https://technicalnotes.wordpress.com/2014/04/19/openldap-setup-with-memberof-overlay/)

The documentation in the above URL is included here for reference.

## openldap setup with memberof overlay

The post summarises steps executed to setup openldap with memberof overlay on Ubuntu 12.04.

### Background

Post-installation, this is how our cn=config looked-

```ldif
ubuntu@PS6226:~/openldap/memberof$ sudo ldapsearch -Q -LLL -Y EXTERNAL -H ldapi:/// -b cn=config dn
dn: cn=config
dn: cn=module{0},cn=config
dn: cn=schema,cn=config
dn: cn={0}core,cn=schema,cn=config
dn: cn={1}cosine,cn=schema,cn=config
dn: cn={2}nis,cn=schema,cn=config
dn: cn={3}inetorgperson,cn=schema,cn=config
dn: olcBackend={0}hdb,cn=config
dn: olcDatabase={-1}frontend,cn=config
dn: olcDatabase={0}config,cn=config
dn: olcDatabase={1}hdb,cn=config
```

### Steps for setting up memberof overlay

1. Load memberof module and configure memberof overlay

    ```bash
    ubuntu@PS6226:~/openldap/memberof$ cat memberof_load_configure.ldif
    dn: cn=module{1},cn=config
    cn: module{1}
    objectClass: olcModuleList
    olcModuleLoad: memberof
    olcModulePath: /usr/lib/ldap

    dn: olcOverlay={0}memberof,olcDatabase={1}hdb,cn=config
    objectClass: olcConfig
    objectClass: olcMemberOf
    objectClass: olcOverlayConfig
    objectClass: top
    olcOverlay: memberof
    olcMemberOfDangling: ignore
    olcMemberOfRefInt: TRUE
    olcMemberOfGroupOC: groupOfNames
    olcMemberOfMemberAD: member
    olcMemberOfMemberOfAD: memberOf
    psl@PS6226:~/openldap/memberof$

    ubuntu@PS6226:~/openldap/memberof$ sudo ldapadd -Q -Y EXTERNAL -H ldapi:/// -f memberof_load_configure.ldif
    adding new entry “cn=module{1},cn=config”

    adding new entry “olcOverlay={0}memberof,olcDatabase={1}hdb,cn=config”

    ```

2. Add referential integrity to the ldap config

    * Modify the cn=module entry to load refint

        ```bash
        ubuntu@PS6226:~/openldap/memberof$ cat 1refint.ldif
        dn: cn=module{1},cn=config
        add: olcmoduleload
        olcmoduleload: refint
        ubuntu@PS6226:~/openldap/memberof$

        ubuntu@PS6226:~/openldap/memberof$ sudo ldapmodify -Q -Y EXTERNAL -H ldapi:/// -f 1refint.ldif
        modifying entry “cn=module{1},cn=config”

        ```

    * Configure refint module

        ```bash
        ubuntu@PS6226:~/openldap/memberof$ cat 2refint.ldif
        dn: olcOverlay={1}refint,olcDatabase={1}hdb,cn=config
        objectClass: olcConfig
        objectClass: olcOverlayConfig
        objectClass: olcRefintConfig
        objectClass: top
        olcOverlay: {1}refint
        olcRefintAttribute: memberof member manager owner
        ubuntu@PS6226:~/openldap/memberof$

        ubuntu@PS6226:~/openldap/memberof$ sudo ldapadd -Q -Y EXTERNAL -H ldapi:/// -f 2refint.ldif
        adding new entry “olcOverlay={1}refint,olcDatabase={1}hdb,cn=config”

        ```

        The system is configured to use memberof attribute for groups!

3. Create groups and add members to the group

    This is how our domain setup looked

    ```bash
    ubuntu@PS6226:~/openldap/memberof$ ldapsearch -x -LLL -H ldap:/// -b dc=example,dc=com dn
    dn: dc=example,dc=com
    dn: cn=admin,dc=example,dc=com
    dn: ou=people,dc=example,dc=com
    dn: ou=groups,dc=example,dc=com
    dn: uid=john,ou=people,dc=example,dc=com
    dn: uid=mahesh,ou=people,dc=example,dc=com
    ```

    Add group

    ```bash
    ubuntu@PS6226:~/openldap/memberof$ cat addgroup-groupofnames.ldif
    dn: cn=peas_dev,ou=groups,dc=example,dc=com
    objectClass: groupofnames
    cn: peas_dev
    description: All users
    # add the group members all of which are
    # assumed to exist under people
    member: uid=john,ou=people,dc=example,dc=com
    member: uid=mahesh,ou=people,dc=example,dc=com

    ubuntu@PS6226:~/openldap/memberof$

    ubuntu@PS6226:~/openldap/memberof$ ldapadd -x -D cn=admin,dc=example,dc=com -W -f groupofnames.ldif
    Enter LDAP Password:
    adding new entry “cn=peas_dev,ou=groups,dc=example,dc=com”

    ```

    Check group membership

    ```bash
    ubuntu@PS6226:~/openldap/memberof$ ldapsearch -x -LLL -H ldap:/// -b uid=mahesh,ou=people,dc=example,dc=com dn
    dn: uid=mahesh,ou=people,dc=example,dc=com

    ubuntu@PS6226:~/openldap/memberof$ ldapsearch -x -LLL -H ldap:/// -b uid=mahesh,ou=people,dc=example,dc=com dn memberof
    dn: uid=mahesh,ou=people,dc=example,dc=com
    memberOf: cn=peas_dev,ou=groups,dc=example,dc=com

    ```
