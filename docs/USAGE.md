# How to use GuildOps on Discord

## Table of Contents

* [Introduction](#introduction)
* [Player actions](#player-actions)
    + [Link a player to a discord user](#link-a-player-to-a-discord-user)
    + [Create an absence](#create-an-absence)
    + [Delete an absence](#delete-an-absence)
    + [Get info about myself](#get-info-about-myself)
* [Guild Officer actions](#guild-officer-actions)
    + [Create a raid <a name="introduction"></a>](#create-a-raid--a-name--introduction----a-)
    + [Create a player](#create-a-player)
    + [Get info about a player](#get-info-about-a-player)
    + [Get info about a raid](#get-info-about-a-raid)
    + [List raids](#list-raids)
    + [Create a strike](#create-a-strike)
    + [List strikes on a player](#list-strikes-on-a-player)
    + [Delete a strike](#delete-a-strike)
    + [Create a fail](#create-a-fail)
    + [List fails on a player or on a raid](#list-fails-on-a-player-or-on-a-raid)
    + [Delete a fail](#delete-a-fail)
    + [Attribute a loot](#attribute-a-loot)
    + [Select a player to attribute a loot](#select-a-player-to-attribute-a-loot)
    + [List loots on a player](#list-loots-on-a-player)
    + [List Absences on a player](#list-absences-on-a-player)
    + [Delete a raid](#delete-a-raid)
    + [Delete a loot](#delete-a-loot)

<small><i><a href='http://ecotrust-canada.github.io/markdown-toc/'>Table of contents generated with markdown-toc</a></i></small>


## Introduction
There are two points of view : the player and the guild officer :
* **Player** : A player is a discord user linked to a player in the GuildOps Database. He can create an absence and get info about himself.
* **Guild Officer** : A guild officer is a discord user with the role "Staff" in the guild. He has access to various tools around loot, strikes, fails, etc.

We encourage to dispatch players and guild officers in different discord channels.

## Player actions

### Link a player to a discord user

It will link the discord user to the player specified. It outputs the player id.
It is required to perform some actions, as create absence or get info about myself.

```shell
/guildops-player-link name: milowenn

You are now linked to this player : 
Name : milowenn
ID : 902837533056499713
Discord ID : milowenn
```

* Name should be a string without space. If there is uppercase, it will be converted to lowercase.
* Name should be the name of a player already created.

### Create an absence

It will create an absence for the player linked to the discord user. It outputs the date of the absence.


```shell
/guildops-absence-create from: 30/09/23 to:31/09/23

Absence(s) created for :
* 30-09-2023
```
* Date must be in format : dd/mm/yy
* to is optional. If not specified, it will generate an absence for the date specified in from.
* Date must be in the future or today.

### Delete an absence

It will delete the absence specified on a date or a range of date. It outputs the success of the deletion. 

If you need to get info about your absences, you can use `/guildops-player-info`.

```shell
/guildops-absence-delete from: 30/09/23 to:31/09/23

Absence(s) deleted for :
* 30-09-2023
```
* Date must be in format : dd/mm/yy or `dd/mm/yy au dd/mm/yy` for a range of date.
* Absences must have been created by `/guildops-absence-create` or it will fail.

### Get info about myself

It will get info about the player linked to the discord user. It outputs the player info.
It outputs only the info of the current season.
**Info is only show to you and no one else.**

```shell
/guildops-player-info

Name : milowenn
ID : 902837533056499713
Discord ID : 271946692805263371
Loots Count: 
mythic | 1 loots
Strikes (1) : 
25/09/2023 | why not | DF/S2 | 903072156068708353
Absences (3) : 
27/09/23 | mythic | tmp
30/09/23 | mythic | example
30/09/23 | mythic | exa mple
Loots (1) : 
mythic | 27/09/23 | loot1
Fails (1) : 
28/09/2023 | p3
```
* Your discord must be linked to a player created by `/guildops-player-create` - by a guild officer - and linked by `/guildops-player-link` yourself.


---

## Guild Officer actions

Strikes, Fails and Loots should only be used by guild officers.
* **Strike** : A strike is a warning for a player. It can be used for bad behavior, bad performance, etc.
* **Fail** : A fail is a fail on a boss. It can be used to track the fails of a player across the raids.
* **Loot** : A loot is a loot given to a player. It can be used to track the loots of a player across the raids.

### Create a raid <a name="introduction"></a>

It will create a raid with the name, date and difficulty specified. It outputs the raid id.

```shell
/guildops-raid-create name: example date: 30/09/23 difficulté: Mythic

Raid 904508030127702017 créé avec succès
```
* Difficulty should be : Normal, Heroic, Mythic
* Date should be in format : dd/mm/yy
* Name should be a string

### Create a player

It will create a player with the name specified. It outputs the player id.

```shell
/guildops-player-create name: example

Joueur example créé avec succès : ID 904508329427959809
```
* Name should be a string without space. If there is uppercase, it will be converted to lowercase.


### Get info about a player

It will get info about the player specified. It outputs the player info.

```shell
/guildops-player-get name: milowenn
Name : milowenn
ID : 902837533056499713
Discord ID : 271946692805263371
Loots Count: 
mythic | 1 loots
Strikes (1) : 
25/09/2023 | why not | DF/S2 | 903072156068708353
Absences (3) : 
27/09/23 | mythic | tmp
30/09/23 | mythic | example
30/09/23 | mythic | exa mple
Loots (1) : 
mythic | 27/09/23 | loot1
Fails (1) : 
28/09/2023 | p3
```
* Player must be created by `/guildops-player-create`.
* If there is not Loots, Fails, Absences or Strikes, it does not output them.


### Get info about a raid
Not implemented yet

### List raids
Not implemented yet

### Create a strike

It will create a strike for the player specified. It outputs the strike id.

```shell
/guildops-strike-create name: milowenn reason: example of strike

Strike créé avec succès
```

### List strikes on a player

It will list the strikes on the player specified. It outputs the strikes.

```shell
/guildops-strike-list name: milowenn

Strikes de milowenn (2) :
25/09/2023 | why not
30/09/2023 | test
```

### Delete a strike

It will delete the strike specified. To get the strike id, you can use `/guildops-strike-list`.

```shell
/guildops-strike-del id: 903072156068708353

Strike supprimé avec succès
```

### Create a fail

It will create a fail for the player specified. It outputs the fail id.

```shell
/guildops-fail-create name: milowenn reason: Erreur P3 Sarkareth date: 30/09/23

Fail créé avec succès
```
* Date must be in format : dd/mm/yy
* Date should be a date of a raid created by `/guildops-raid-create`
* Name should be a string without space. If there is uppercase, it will be converted to lowercase.
* Name should be the name of a player already created.
* Reason should be a string

### List fails on a player or on a raid

It will list the fails on the player or on the raid specified. It outputs the fails.

```shell
/guildops-fail-list-player name: milowenn

Fails de milowenn (2) :
28-09-2023 - p3 - 903072156068708353
30-09-2023 - Erreur P3 Sarkareth - 903072156068708353
```
```shell
/guildops-fail-list-raid date: 30/09/23

Fails du 30/09/2023 (1) :
milowenn - Erreur P3 Sarkareth - 903072156068708353
```

* Date must be in format : dd/mm/yy
* Date should be a date of a raid created by `/guildops-raid-create`
* Name should be a string without space. If there is uppercase, it will be converted to lowercase.
* Name should be the name of a player already created.

### Delete a fail

It will delete the fail specified. To get the fail id, you can use `/guildops-fail-list-player name: <player_name>`.

```shell
/guildops-fail-del id: 904435308715671553

Fail supprimé avec succès
```

* Id should be the id of a fail created by `/guildops-fail-create`

### Attribute a loot

It attribute a loot to a player. 

```shell
/guildops-loot-attribute loot-name: example object raid-id: 64546465464 player-name: milowenn

Loot attribué avec succès
```
* Loot-name should be a string.
* Raid-id should be the id of a raid created by `/guildops-raid-create`
* Player-name should be the name of a player already created by `/guildops-player-create`.

### Select a player to attribute a loot

It takes a list of players and randomly pick a player with the lesser loots on the same difficulty and give a name. It does attribute directly the loot.

```shell
/guildops-loot-select difficulty: mythic players: milowenn, example, example2

Le joueur milowenn a été sélectionné pour recevoir le loot
```
* Difficulty should be : Normal, Heroic, Mythic
* Players should be a list of players separated by a comma. If there is uppercase, it will be converted to lowercase.
* Players should be the name of a player already created by `/guildops-player-create`.


### List loots on a player

It will list the loots on the player specified. It outputs the loots.

```shell
/guildops-loot-list player-name: milowenn

Tous les loots de milowenn:
  loot1 2023-09-27 00:00:00 +0000 UTC mythic
```

* Player-name should be the name of a player already created by `/guildops-player-create`.

### List Absences on a date

It will list the absences on the date specified. It outputs the absences.

```shell
/guildops-absence-list date: 28/09/23

Absence(s) pour le 28-09-2023 :
* milowenn
```

* Date must be in format : dd/mm/yy
* Date should be a date of a raid created by `/guildops-raid-create`

###  Delete a raid
**Warning : it will delete all the loots, fails, strikes and absences of the raid.**

It will delete the raid specified. To get the raid id, you can use `/guildops-raid-list`.

```shell
/guildops-raid-delete id:465465465465465465

Raid supprimé avec succès
```
* Id should be the id of a raid created by `/guildops-raid-create`

### Delete a loot 

It will delete the loot specified. To get the loot id, you can use `/guildops-loot-list`.

```shell
/guildops-loot-delete id:465465465465465465

Loot supprimé avec succès
```

