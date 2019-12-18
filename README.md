# Owl

[![Go Report Card](https://goreportcard.com/badge/github.com/adrienaury/owl)](https://goreportcard.com/report/github.com/adrienaury/owl)
[![Github Release Card](https://img.shields.io/github/release/adrienaury/owl)](https://github.com/adrienaury/owl/releases)
[![codecov](https://codecov.io/gh/adrienaury/owl/branch/develop/graph/badge.svg)](https://codecov.io/gh/adrienaury/owl)
[![Build Status](https://travis-ci.org/adrienaury/owl.svg?branch=develop)](https://travis-ci.org/adrienaury/owl)
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fadrienaury%2Fowl.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2Fadrienaury%2Fowl?ref=badge_shield)

Owl is a platform agnostic set of tools to manage realms of units, users and groups. Thanks to the modular conception, any backend can theorically be used to store and access data (LDAP, MySQL, MongoDB, ...). For now, only LDAP is supported, please open an issue if another one is needed.

The project is composed of 3 tools :

* Owl CLI : manage your realms with a powerfull devops CLI
* Owl REST API Server : equivalent to the CLI but with exposed REST Endpoints
* Owl Web Administration GUI : graphical user interface in front of the REST API Server

## Concepts

Owl is opiniated on how to manage user accounts, but it is also highly customizable.

There is only 4 types of object manipulated.

Each object has a unique identifier and a set of prefefined properties that can be multivalued. Additional properties can be configured.

### Realms

Realms are associated with servers, instance, etc... where the data is persisted. Each realm is in isolation from other realms.

Property | Description
--       | --
ID       | Unique realm identifier
URL      | Location of the realm
Username | Used as login account to the realm backend

### Units

Units are logical grouping of users and groups, used to mimic real-world organization (like OU in LDAP).

Property    | Description
--          | --
ID          | Unique unit identifier
Description | Description of the unit

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

First principle : Owl CLI respect the UNIX philosophy.

> Write programs that do one thing and do it well.\
> Write programs to work together.\
> Write programs to handle text streams, because that is a universal interface.
>
> -- *Douglas McIlroy, inventor of Unix pipelines*

How ?

* it's simple, no complicated fanciness (like in LDAP for example)
* read inputs as json on the standard input by default
* write outputs as json on the standard output
* write logs on the standard error by default

Second principle : Owl is a devops tool.

Why ?

* it can be used by a human operator or automated by scripting
* local configuration is stored "as code" in local directory
* every object can be serialized as JSON or YAML, if you need to store them in a code repository (like Git)

#### Examples

##### Manage realms

Create or modify realms with `owl realm` command.

```console
$ owl realm dev ldap://dev.my-company.com/dc=example,dc=com cn=admin,dc=example,dc=com
Set realm 'dev' to 'ldap://dev.my-company.com/dc=example,dc=com'.
```

Login into a realm with `owl login` command. It is also possible to use the `--realm` flag on a specific command.

```console
$ owl login dev
Password :
Connected to realm 'dev' as user 'admin'.

$ owl login dev
Connected to realm 'dev' as user 'admin'.
```

List realms with `owl realms` command, current realm is highlighted with an asterisk.

```console
$ owl realms
ID     Username                    URL
dev    cn=admin,dc=example,dc=com  ldap://dev.my-company.com/dc=example,dc=com
prod*  cn=admin,dc=example,dc=com  ldap://prod.my-company.com/dc=example,dc=com
```

`TODO` Realm creation and login can be done in a single operation :

```console
$ owl login ldap://dev.my-company.com/dc=example,dc=com
Username : cn=admin,dc=example,dc=com
Password :
Name this realm : dev
Set realm 'dev' to 'ldap://dev.my-company.com/dc=example,dc=com'.
Connected to realm 'dev' as user 'admin'.

$ owl login ldap://dev.my-company.com/dc=example,dc=com
Username : cn=admin,dc=example,dc=com
Connected to realm 'dev' as user 'admin'.

$ owl login ldap://dev.my-company.com/dc=example,dc=com cn=admin,dc=example,dc=com
Connected to realm 'dev' as user 'admin'.
```

##### Manage organizational units

Create a new unit with `owl create unit` command.

```console
$ owl create unit my-unit "Test unit"
Created unit 'my-unit' in realm 'dev'.
```

The create command also read JSON on stdin, so these are other ways of doing :

```console
$ owl create unit <<< '{"ID": "my-unit", "Description": "Test unit"}'
Created unit 'my-unit' in realm 'dev'.

$ echo '{"ID": "my-unit", "Description": "Test unit"}' | owl create unit
Created unit 'my-unit' in realm 'dev'.
```

List existing units with `owl list unit` command.

```console
$ owl list unit
ID       Description
my-unit  Test unit
```

To create users and groups, you first need to select a unit with `owl unit` command. It is also possible to use the `--unit` flag on a specific command.

```console
$ owl unit my-unit
Using unit 'my-unit' for next commands.
```

Know which unit you're currently on with `owl unit` command.

```console
$ owl unit
Using unit 'my-unit'.
```

The special `default` unit is selected if `owl unit` is never used before. You can re-select the default unit at any time.

```console
$ owl unit default
Using default unit for next commands.

$ owl unit -
Using default unit for next commands.
```

##### Manage users

To create a user, use `owl create user` command.

```console
$ owl create user batman firstname=Bruce lastname=Wayne
Created user 'batman' in unit 'my-unit' of realm 'dev'.

$ owl create user <<< '{"ID": "batman", "FirstNames": ["Bruce"], "LastNames": ["Wayne"]}'
Created user 'batman' in unit 'my-unit' of realm 'dev'.
```

You can also create or replace an existing user with `owl apply user` command.

```console
$ owl apply user batman firstname=Bruce lastname=Wayne email=bruce.wayne@gotham.dc
Replaced user 'batman' in unit 'my-unit' of realm 'dev'.

$ owl apply user joker firstname=Arthur lastname=Flake email=arthur.flake@gotham.dc
Created user 'joker' in unit 'my-unit' of realm 'dev'.
```

To only add a single attribute, use `owl append user` command.

```console
$ owl append user joker firstname="Jack"
Modifier user 'joker' in unit 'my-unit' of realm 'dev'.
```

List user with `owl list user` command.

```console
$ owl list user
ID      First Names   Last Names  E-mails
batman  Bruce         Wayne       bruce.wayne@gotham.dc
joker   Arthur, Jack  Flake       arthur.flake@gotham.dc
```

Give user a random password with `owl password assign` command.

```console
$ owl password assign joker
Assigned new random password to user 'joker' in unit 'my-unit' of realm 'dev'.
```

##### Manage groups

You guessed it, use `owl create group` command to create a group.

```console
$ owl create group bad-guys member=joker member=batman
Created group 'bad-guys' in unit 'my-unit' of realm 'dev'.
```

Member list can be modified with `owl append` and `owl remove` commands.

```console
$ owl remove group bad-guys member=batman
Removed from group 'bad-guys' in unit 'my-unit' of realm 'dev'.

$ owl append group good-guys member=batman
Appended to group 'good-guys' in unit 'my-unit' of realm 'dev'.
```

##### Verbs

Here is a list of verbs available to manage objects, with `owl <verb> <object>` command structure.

Write verbs :

Verb   | Aliases        | If object already exists     | If object doesn't exist
--     | --             | --                           | --
create | insert, import | error                        | create object
apply  | replace, ap    | replace object               | create object
update | set            | replace specified attributes | error
upsert | -              | replace specified attributes | create object
append | add            | add attributes               | error
remove | rm             | remove attributes            | error
delete | del            | delete object                | nothing

Read verbs :

Verb | Aliases    | Description
--   | --         | --
list | ls, export | list all objects
get  | read       | read object with given ID

##### Export and import

List all objects with `owl list` without parameter. All write verbs can import a list of objects. Use them without parameter to mix different types.

```console
$ owl list user -o json | owl apply --realm=prod --unit organization
Replaced user 'batman' in unit 'organization' of realm 'prod'.
Replaced user 'joker' in unit 'organization' of realm 'prod'.
Created user 'robin' in unit 'organization' of realm 'prod'.
```

```console
$ owl export -o json | owl apply --realm=prod
Replaced unit 'my-unit' in realm 'prod'.
Replaced user 'batman' in unit 'my-unit' of realm 'prod'.
Replaced user 'joker' in unit 'my-unit' of realm 'prod'.
Created user 'robin' in unit 'my-unit' of realm 'prod'.
Replaced group 'good-guys' in unit 'my-unit' of realm 'prod'.
Replaced group 'bad-guys' in unit 'my-unit' of realm 'prod'.
```

##### Advanded commands

All commands can output results in JSON or YAML format, thanks to the `--output` (short `-o`) flag.

```console
$ owl list user -o json
{"Users": [{"ID": "batman", "FirstNames": ["Bruce"], "LastNames": ["Wayne"], "Emails": ["bruce.wayne@gotham.dc"]}, {"ID": "joker", "FirstNames": ["Arthur", "Jack"], "LastNames": ["Flake"], "Emails": ["arthur.flake@gotham.dc"]}]}
```

This universal interface enable the use of other programs, for example `jq`.

```console
$ owl list user -o json | jq
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

$ owl ls user -o json | jq ".Users | [.[].ID]"
[
  "batman",
  "joker"
]
```

Owl also understand JSON if passed throught stdin, this enables chaining of owl commands.

```console
$ owl list user -o json | owl apply user --realm=prod --unit=organization
Replaced user 'batman' in unit 'organization' of realm 'prod'.
Replaced user 'joker' in unit 'organization' of realm 'prod'.
```

#### Installation

Download the latest version for your OS from the [release page](https://github.com/adrienaury/owl/releases).

### Owl REST Server

`TODO`

### Owl Web GUI

`TODO`

## Contribute

Contributions to this project are very welcome.

If you want to contribute, please check CONTRIBUTING.md

## Links

* Issue Tracker: github.com/adrienaury/owl/issues
* Source Code: github.com/adrienaury/owl

## Support

If you are having issues, please let me know.
I'm Adrien and my mail is adrien.aury@gmail.com

## License

### Main license

The project is licensed under the [MIT License](https://opensource.org/licenses/MIT).

### Use of 3rd party librairies

Library                                     | Version | Licenses                        |
--------------------------------------------|---------|---------------------------------|
golang.org/pkg                              | v1.13   | [BSD-3-Clause](NOTICE.md#go)    |
github.com/spf13/cobra                      | v0.0.5  | [Apache-2.0](NOTICE.md#cobra)   |
github.com/mitchellh/go-homedir             | v1.1.0  | [MIT](NOTICE.md#go-homedir)     |
github.com/docker/docker-credential-helpers | v0.6.3  | [MIT](NOTICE.md#dch)            |
gopkg.in/yaml.v3                            | v3.0.0  | [Apache-2.0](NOTICE.md#yaml)    |
gopkg.in/ldap.v3                            | v3.1.0  | [MIT](NOTICE.md#ldap)           |
github.com/spf13/pflag                      | v1.0.3  | [BSD-3-Clause](NOTICE.md#pflag) |
gopkg.in/asn1-ber.v1                        | v1.0.0  | [MIT](NOTICE.md#asn1)           |

Check NOTICE.md for copyright notices.

[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fadrienaury%2Fowl.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Fadrienaury%2Fowl?ref=badge_large)
