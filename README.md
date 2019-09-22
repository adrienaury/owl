# Owl

Owl is a command line interface to interact with a Ldap database.

## Example

```bash
owl connect :url
owl use :basedn
owl search :filter
owl add
owl delete
owl modify
owl rename
owl disconnect
```

## Architecture

## Algo

### Login

1. Load session
2. Complete user input
   1. Server URL
      1. Valeur passée en argument ? => oui: 2.2
      2. Valeur passée par flag ? => oui: 2.2
      3. Valeur présente dans session ? => oui: 2.2
      4. Terminal ? => oui: 2.2
      5. Erreur : server doit être spécifié
   2. Username
      1. Valeur passée par flag ? => oui: 2.3
      2. Valeur présente dans session ? => oui: 2.3
      3. Terminal ? => oui: 2.3
      4. Erreur : username doit être spécifié
   3. Password
      1. Valeur passée par flag ? => oui: 3
      2. Valeur présente dans coffre-fort ? => oui: 3
      3. Terminal ? => oui: 3
      4. Erreur : password doit être spécifié
3. Normalize inputs
   1. Server URL
4. Validate inputs
   1. Nombre d'arguments
   2. Serveur passé qu'une seule fois
5. Essayer de s'authentifier
   1. Invalid credentials et terminal ? => oui: 6
   2. Invalid credentials et pas terminal ? => oui: 7
   3. Autre erreur ? => oui: 7
   4. OK ? => oui: 8
6. Retry avec nouvelles valeurs
   1. Demander username et password => 5
7. Afficher message erreur
8. Mettre à jour
   1. Session
   2. Coffre-fort
