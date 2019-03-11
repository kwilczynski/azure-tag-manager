# Tag manager 

Currently the software works for Azure only. 

## Prerequisites

For Azure you need to create service principal.


## Download

You can download a binary for your architecture from https://bitbucket.org/nordcloud/tagmanager/downloads/ 

### Azure
Create service principal file:

```bash
az ad sp create-for-rbac --sdk-auth > my.auth
```

and export path to the authorizer:

```bash
export AZURE_AUTH_LOCATION=my.auth
```

## How it works ?

Tag rewriter accepts the payload where rules are defined. Each rule has a list of conditions and a lits of actions. If all conditions evaluate to true for a resource, all actions are executed. 

```json
{
  "dryrun": true,
  "rules":  [
    {
        "name": "Tag me stuff", 
        "conditions": [
            {"type": "tagEqual", "tag": "darek", "value" : "example"},
            {"type": "tagExists", "tag": "darek7"},
            {"type": "tagNotExists", "tag": "env"}
        ], 
        "actions": [
            {"type": "addTag", "tag": "mucha", "value": "zoo" },
            {"type": "addTag", "tag": "mucha3", "value": "zoo" }
        ]
      }
    ]
}
```

As of today, the tag rewriter accepts four kinds of conditions (all are case senstive):

* `noTags` - checks if there are no tags set 
* `tagEqual` - checks if a `tag` has a `value` set 
* `tagNotEqual` - checks if a `tag` has a value set different than `value` 
* `tagExists` - checks if a tag with key `tag` exists
* `tagNotExists` - same as above but negative
* `regionEqual` - checks if resource is in key `region` (aka location in azure)
* `regionNotEqual` - same as above but negative
* `rgEqual` - match resource group in a key `resourceGroup`
* `rgNotEqual` - match not resource group 


## Running 

`./tagmanager mapping.json` 

### Example run

```

➜  pantageusz git:(master) ✗ ./tagmanager mapping.json
INFO[0002] 👍  Conditions are true for (darek33) with ID = /subscriptions/6690b014-bdbd-4496-98ee-f2f255699f70/resourceGroups/darek/providers/Microsoft.Storage/storageAccounts/darek33 
INFO[0002]      🔥  DryRun Firing action addTag on resource /subscriptions/6690b014-bdbd-4496-98ee-f2f255699f70/resourceGroups/darek/providers/Microsoft.Storage/storageAccounts/darek33 
INFO[0002]      🔥  DryRun Firing action addTag on resource /subscriptions/6690b014-bdbd-4496-98ee-f2f255699f70/resourceGroups/darek/providers/Microsoft.Storage/storageAccounts/darek33 

```

## Changelog

0.2

* support for named rules (you must use `name` to add a name for a rule)
* support for a new condition check `rgEqual` and `rgNotEqual` to match resource groups   the syntax is ` { "type": "rgEqual", "resourceGroup": "myRg" }`
* support for checking for no tags `noTags` 
* a less verbose debug level


## Todo 

* Azure ARM policy setting 
* Better CLI with flags
* AWS support for EC2

## Licence 

Dariusz Dwornikowski, 

Nordcloud PAT team (Hold my 🍺, not your my 🐴🐴)

