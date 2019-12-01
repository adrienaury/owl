# Owl

Owl is a platform agnostic set of tools to manage realms of units, users and groups.

## Concepts

Owl is very opiniated on how to manage user accounts, but it is also highly customizable.

Only 4 types of object are manipulated.

Each object has a unique identifier and a set of prefefined properties that can be multivalued. Additional properties can be configured.

### Realms

Property | Description
--       | --
ID       | Unique realm identifier
URL      | Location of the realm
Username | Used as login account to the realm backend

### Units

Property | Description
--       | --
ID       | Unique unit identifier

### Users

Property    | Description
--          | --
ID          | Unique user identifier
FirstNames  | First names [multivalued property]
LastNames   | Last names [multivalued property]
Emails      | E-mails owned by the user [multivalued property]

### Groups

Property | Description
--       | --
ID       | Unique group identifier
Name     | Name of the group
Members  | Ids of users in the group [multivalued property]

## Tools

### Owl CLI

Owl CLI respect the UNIX philosophy :

* Write programs that do one thing and do it well.
* Write programs to work together.
* Write programs to handle text streams, because that is a universal interface.

How ?

* manage realms, units, users and groups - that's it
* read inputs as json on the standard input by default
* write outputs as json on the standard output by default
* write logs on the standard error by default

#### Examples

```text
$ owl realm set dev ldap://dev.my-company.com/dc=example,dc=com cn=admin,dc=example,dc=com
Set realm 'dev' to 'ldap://dev.my-company.com/dc=example,dc=com'.

$ owl realm list -o table
ID    Username                    URL
dev   cn=admin,dc=example,dc=com  ldap://dev.my-company.com/dc=example,dc=com
prod  cn=admin,dc=example,dc=com  ldap://prod.my-company.com/dc=example,dc=com

$ owl realm login dev
Password :
Connected to realm 'dev' as user 'admin'.

$ owl unit list -o table
The realm contains no unit.

$ owl unit create <<< '{"ID": "my-unit", "Description": "Test unit"}'
Created unit 'my-unit' in realm 'dev'.

$ owl unit list -o table
ID       Description
my-unit  Test unit

$ owl unit use my-unit
Using unit 'my-unit' for next commands.

$ owl user create <<< '{"ID": "batman", "FirstNames": ["Bruce"], "LastNames": ["Wayne"]}'
Created user 'batman' in unit 'my-unit' of realm 'dev'.

$ owl user apply <<< '{"ID": "batman", "Emails": ["bruce.wayne@gotham.dc"]}'
Modified user 'batman' in unit 'my-unit' of realm 'dev'.

$ owl user apply joker firstname=Arthur lastname=Flake email=arthur.flake@gotham.dc
Created user 'joker' in unit 'my-unit' of realm 'dev'.

$ owl user append joker firstname="Jack"
Modifier user 'joker' in unit 'my-unit' of realm 'dev'.

$ owl group create bad-guys member=joker member=batman
Created group 'bad-guys' in unit 'my-unit' of realm 'dev'.

$ owl group member remove bad-guys batman
Modified group 'bad-guys' in unit 'my-unit' of realm 'dev'.

$ owl user list -o table
ID      First Names   Last Names  E-mails
batman  Bruce         Wayne       bruce.wayne@gotham.dc
joker   Arthur, Jack  Flake       arthur.flake@gotham.dc

$ owl user list -o json
{"Users": [{"ID": "batman", "FirstNames": ["Bruce"], "LastNames": ["Wayne"], "Emails": ["bruce.wayne@gotham.dc"]}, {"ID": "joker", "FirstNames": ["Arthur", "Jack"], "LastNames": ["Flake"], "Emails": ["arthur.flake@gotham.dc"]}]}

$ owl user list -o json | jq
{
    "Users": [
        {
            "ID": "batman",
            "FirstNames": [
                "Bruce"
            ],
            "LastNames": [
                "Wayne"
            ],
            "Emails": [
                "bruce.wayne@gotham.dc"
            ]
        },
        {
            "ID": "joker",
            "FirstNames": [
                "Arthur",
                "Jack"
            ],
            "LastNames": [
                "Flake"
            ],
            "Emails": [
                "arthur.flake@gotham.dc"
            ]
        }
    ]
}

$ owl user ls -o json | jq ".Users | [.[].ID]"
[
  "batman",
  "joker"
]

$ owl user list -o json | owl import --realm=prod --unit=organization
Imported 2 users in unit 'organization' of realm 'prod'.
```

#### Installation

Download the latest version for your OS from the [release page](https://github.com/adrienaury/owl/releases).

### Owl REST Server

TODO

### Owl Web GUI

TODO

## Contribute

Contributions to this project are very welcome.

If you want to contribute, please check CONTRIBUTING.md

## Links

* Issue Tracker: github.com/adrienaury/mailmock/issues
* Source Code: github.com/adrienaury/mailmock

## Support

If you are having issues, please let me know.
I'm Adrien and my mail is adrien.aury@gmail.com

## License

### Main license

The project is licensed under the [MIT License](https://opensource.org/licenses/MIT).
