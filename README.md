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

##### Manage realms

Create or modify realms with `owl realm set` command.

```shell
$ owl realm set dev ldap://dev.my-company.com/dc=example,dc=com cn=admin,dc=example,dc=com
Set realm 'dev' to 'ldap://dev.my-company.com/dc=example,dc=com'.
```

List created realms with `owl realm list` command.

```shell
$ owl realm list -o table
ID    Username                    URL
dev   cn=admin,dc=example,dc=com  ldap://dev.my-company.com/dc=example,dc=com
prod  cn=admin,dc=example,dc=com  ldap://prod.my-company.com/dc=example,dc=com
```

Login into a realm with `owl realm login` command.

```shell
$ owl realm login dev
Password :
Connected to realm 'dev' as user 'admin'.
```

##### Manage organizational units

Create a new unit with `owl unit create` command.

```shell
$ owl unit create <<< '{"ID": "my-unit", "Description": "Test unit"}'
Created unit 'my-unit' in realm 'dev'.
```

List existing units with `owl unit list` command.

```shell
$ owl unit list -o table
ID       Description
my-unit  Test unit
```

To create users and groups, you first need to select a unit with `owl unit use` command.

```shell
$ owl unit use my-unit
Using unit 'my-unit' for next commands.
```

##### Manage users

To create a user, use `owl user create` command.

```shell
$ owl user create <<< '{"ID": "batman", "FirstNames": ["Bruce"], "LastNames": ["Wayne"]}'
Created user 'batman' in unit 'my-unit' of realm 'dev'.
```

You can also create or replace an existing user with `owl user apply` command.

```shell
$ owl user apply <<< '{"ID": "batman", "Emails": ["bruce.wayne@gotham.dc"]}'
Modified user 'batman' in unit 'my-unit' of realm 'dev'.

$ owl user apply joker firstname=Arthur lastname=Flake email=arthur.flake@gotham.dc
Created user 'joker' in unit 'my-unit' of realm 'dev'.
```

To only add a single attribute, use `owl user append` command.

```shell
$ owl user append joker firstname="Jack"
Modifier user 'joker' in unit 'my-unit' of realm 'dev'.
```

List user with `owl user list` command.

```shell
$ owl user list -o table
ID      First Names   Last Names  E-mails
batman  Bruce         Wayne       bruce.wayne@gotham.dc
joker   Arthur, Jack  Flake       arthur.flake@gotham.dc
```

##### Manage groups

You guessed it, use `owl group create` command to create a group.

```shell
$ owl group create bad-guys member=joker member=batman
Created group 'bad-guys' in unit 'my-unit' of realm 'dev'.
```

Member list can be modified withe `owl group member` sub-commands.

```shell
$ owl group member remove bad-guys batman
Modified group 'bad-guys' in unit 'my-unit' of realm 'dev'.

$ owl group member add good-guys batman
Modified group 'good-guys' in unit 'my-unit' of realm 'dev'.
```

##### Advanded commands

All commands can output results in JSON or YAML format, thanks to the `--output` (short `-o`) flag.

```shell
$ owl user list -o json
{"Users": [{"ID": "batman", "FirstNames": ["Bruce"], "LastNames": ["Wayne"], "Emails": ["bruce.wayne@gotham.dc"]}, {"ID": "joker", "FirstNames": ["Arthur", "Jack"], "LastNames": ["Flake"], "Emails": ["arthur.flake@gotham.dc"]}]}
```

This universal interface enable the use of other programs, for example `jq`.

```shell
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
```

Owl also understand JSON if passed throught stdin.

```shell
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
