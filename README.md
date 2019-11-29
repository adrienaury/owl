# Owl

Owl is a set of tools to manage realms of units, users and groups.

## Concepts

Each object has a unique identifier and a set of prefefined properties that can be multivalued.

### Realms

Property | Description
--       | --
id       | Unique realm identifier
url      | Location of the realm
username | Used as login account to the realm backend

### Units

Property | Description
--       | --
id       | Unique unit identifier

### Users

Property    | Description
--          | --
id          | Unique user identifier
first-name  | First names [multivalued property]
middle-name | Middle names [multivalued property]
last-name   | Last names [multivalued property]
email       | Email owned by the user [multivalued property]
group       | Ids of groups the user is member of [multivalued property]

### Groups

Property | Description
--       | --
id       | Unique group identifier
name     | Name of the group
member   | Ids of users in the group [multivalued property]

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
$ owl realm set dev ldap://dev.my-company.com/dc=example,dc=com
Set realm 'dev' to 'ldap://dev.my-company.com/dc=example,dc=com'.

$ owl realm list -o table
Identifier  Username  URL
dev         admin     ldap://dev.my-company.com/dc=example,dc=com
prod        admin     ldap://prod.my-company.com/dc=example,dc=com

$ owl realm login dev
Password :
Connected to realm 'dev' as user 'admin'.

$ owl unit list -o table
The realm contains no unit.

$ owl unit apply <<< '{"id": "my-unit"}'
Created unit 'my-unit' in realm 'dev'.

$ owl unit list
{"units": [{"id": "my-unit"}]}

$ owl unit use my-unit
Using unit 'my-unit' for next commands.

$ owl user create <<< '{"id": "batman", "first-name": ["Bruce"], "last-name": ["Wayne"]}'
Created user 'batman' in unit 'my-unit' of realm 'dev'.

$ owl user apply <<< '{"id": "batman", "email": "bruce.wayne@gotham.dc"}'
Modified user 'batman' in unit 'my-unit' of realm 'dev'.

$ owl user apply joker first-name=Arthur last-name=Flake email=arthur.flake@gotham.dc
Created user 'joker' in unit 'my-unit' of realm 'dev'.

$ owl user add first-name="Jack"
Modifier user 'joker' in unit 'my-unit' of realm 'dev'.

$ owl group create bad-guys member=joker member=batman
Created group 'bad-guys' in unit 'my-unit' of realm 'dev'.

$ owl group remove member=batman
Modified group 'bad-guys' in unit 'my-unit' of realm 'dev'.

$ owl user list -o table
Identifier  First-Name    Middle-Name  Last-Name  Email                   Group
batman      Bruce                      Wayne      bruce.wayne@gotham.dc
joker       Arthur, Jack               Flake      arthur.flake@gotham.dc  bad-guys

$ owl user list
{"users": [{"id": "batman", "first-name": ["Bruce"], "last-name": ["Wayne"], "email": ["bruce.wayne@gotham.dc"]}, {"id": "joker", "first-name": ["Arthur", "Jack"], "last-name": ["Flake"], "email": ["arthur.flake@gotham.dc"]}]}

$ owl user list | jq
{
    "users": [
        {
            "id": "batman",
            "first-name": [
                "Bruce"
            ],
            "last-name": [
                "Wayne"
            ],
            "email": [
                "bruce.wayne@gotham.dc"
            ]
        },
        {
            "id": "joker",
            "first-name": [
                "Arthur",
                "Jack"
            ],
            "last-name": [
                "Flake"
            ],
            "email": [
                "arthur.flake@gotham.dc"
            ]
        }
    ]
}

$ owl user list | owl import --realm=prod --unit=organization
Imported 2 users in unit 'organization' of realm 'prod'.
```

#### Installation

Download the latest version for your OS from the [release page](https://github.com/adrienaury/owl/releases).

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
