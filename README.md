# FreePBX LDAP Directory
A simple LDAP server to serve a searchable address book of internal extensions from the FreePBX DB

## How it works
It starts the LDAP service on port 10389 and responds to directory search requests by translating them into a SQL query against the "asterisk.users" table in MySQL.

Since we aren't working with sensitive information or trying to implement authentication, but most phones require a bind request with a username & password before they'll search, it'll respond as success to any bind request without checking credentials.

This means the address list will always be up-to-date, as there is no import/export.

Two fields are returned for each result, "displayName" and "telephoneNumber".

MySQL to LDAP mapping is:
* "name" in MySQL maps to "displayName" in LDAP
* "extension" in MySQL maps to "telephoneNumber" in LDAP

## Build & Usage
To build, you will need the Go runtime and to build you just need to run:

```
go build
```

## Recommended Install Procedure
```
# mkdir -p /opt/freepbx-ldap
# cp <freepbx-ldap binary location> /opt/freepbx-ldap/freepbx-ldap
# chown -R asterisk:asterisk /opt/freepbx-ldap
# chmod +x /opt/freepbx-ldap/freepbx-ldap
# cp <systemd/freepbx-ldap.service location> /etc/systemd/system/freepbx-ldap.service
# systemctl daemon-reload
# systemctl enable freepbx-ldap
# systemctl start freepbx-ldap
```

## Phone Configuration
You'll need to configure your IP phones to look up against the LDAP server.

See examples below:

### Snom720 - snom720-main.htm
```xml
<?xml version="1.0" encoding="utf-8"?>
<settings>
        <phone-settings>
                *** Other Settings ***

                <ldap_server perm="">***server_ip***</ldap_server>
                <ldap_port perm="">10389</ldap_port>
                <ldap_base perm="">dc=asterisk</ldap_base>
                <ldap_username perm="">asterisk</ldap_username>
                <ldap_max_hits perm="">100</ldap_max_hits>
                <ldap_search_filter perm="">(&(telephoneNumber=*)(displayName=%))</ldap_search_filter>
                <ldap_number_filter perm="">(&(telephoneNumber=%)(displayName=*))</ldap_number_filter>
                <ldap_name_attributes perm="">displayName</ldap_name_attributes>
                <ldap_number_attributes perm="">telephoneNumber</ldap_number_attributes>
                <ldap_display_name perm="">%displayName</ldap_display_name>

                <gui_fkey1 perm="">keyevent F_DIRECTORY_SEARCH</gui_fkey1>
        </phone-settings>
</settings>
```

### Snom300 - snom300-main.htm
```xml
<?xml version="1.0" encoding="utf-8"?>
<settings>
        <phone-settings>
                *** Other Settings ***

                <ldap_server perm="">***server_ip***</ldap_server>
                <ldap_port perm="">10389</ldap_port>
                <ldap_base perm="">dc=asterisk</ldap_base>
                <ldap_username perm="">asterisk</ldap_username>
                <ldap_max_hits perm="">25</ldap_max_hits>
                <ldap_search_filter perm="">(&(telephoneNumber=*)(displayName=%))</ldap_search_filter>
                <ldap_number_filter perm="">(&(telephoneNumber=%)(displayName=*))</ldap_number_filter>
                <ldap_name_attributes perm="">displayName</ldap_name_attributes>
                <ldap_number_attributes perm="">telephoneNumber</ldap_number_attributes>
                <ldap_display_name perm="">%displayName</ldap_display_name>

                <idle_cancel_key_action perm="">keyevent F_DIRECTORY_SEARCH</idle_cancel_key_action>
        </phone-settings>
        <functionKeys e="2">
                <fkey idx="3" context="active" label="" perm="">keyevent F_DIRECTORY_SEARCH</fkey>
        </functionKeys>
</settings>
```

### Polycom SoundPoint IP - sip.cfg (must be firmware UC 4+)
```xml
<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<localcfg>
  *** Other Settings ***

  <dir>
      <dir.corp
          dir.corp.address="ldap://***server_ip***"
          dir.corp.port="10389"
          dir.corp.transport="TCP"
          dir.corp.baseDN="dc=asterisk"
          dir.corp.scope="sub"
          dir.corp.filterPrefix=""
          dir.corp.user="uid=asterisk,dc=asterisk"
          dir.corp.pageSize="32"
          dir.corp.password="supersecret"
          dir.corp.cacheSize="128"
          dir.corp.leg.pageSize="8"
          dir.corp.leg.cacheSize="32"
          dir.corp.autoQuerySubmitTimeout="1"
          dir.corp.viewPersistence="0"
          dir.corp.leg.viewPersistence="0"
          dir.corp.sortControl="0">
          <dir.corp.attribute
              dir.corp.attribute.1.name="displayName"
              dir.corp.attribute.1.label="Display Name"
              dir.corp.attribute.1.type="first_name"
              dir.corp.attribute.1.searchable="1"
              dir.corp.attribute.1.filter=""
              dir.corp.attribute.1.sticky="0"
              dir.corp.attribute.2.name="telephoneNumber"
              dir.corp.attribute.2.label="phone number"
              dir.corp.attribute.2.type="phone_number"
              dir.corp.attribute.2.filter=""
              dir.corp.attribute.2.sticky="0"
              dir.corp.attribute.2.searchable="1">
          </dir.corp.attribute>
          <dir.corp.backGroundSync
              dir.corp.backGroundSync.period="3600">
          </dir.corp.backGroundSync>
          <dir.corp.vlv
              dir.corp.vlv.allow="1"
              dir.corp.vlv.sortOrder="displayName telephoneNumber">
          </dir.corp.vlv>
      </dir.corp>
  </dir>

  <feature feature.corporateDirectory.enabled="1"/>
  <softkey softkey.feature.directories="1"/>
</localcfg>
