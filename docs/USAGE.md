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
Name : chibrousse
Discord Name : milowenn
```

**Requirements:**
* Name should be a string without space. If there is uppercase, it will be converted to lowercase.
* Name should be the name of a player already created.

**Errors :**
* If the player does not exist.

  ``` Error while linking player: player not found```
* If Discord user is already linked to a player.

  ``` Error while linking player: discord account already linked to player name chibrousse. Contact Staff for modification```

### Create an absence 

It will create an absence for the player linked to the discord user. It outputs the date of the absence.


```shell
/guildops-absence-create from: 30/09/23 to:31/09/23

Absence(s) created for :
* Sat 07/10/23
* Sun 08/10/23
```
* Date must be in format : dd/mm/yy
* to is optional. If not specified, it will generate an absence for the date specified in from.
* Date must be in the future or today.
* If no raids are found for the date specified, it will show 
  ``` 
  No absence created or deleted, there is no raids on this range
  ```
  
**Errors :**
* If the date is malformed

  ``` Error while parsing date:parsing time "08/10/23A": extra text: "A"```
* If the date is in the past

  ``` Error while creating absence: date cannot be in the past```
* if to is before from

  ``` Error while creating absence: endDate is before startDate```

### Delete an absence

It will delete the absence specified on a date or a range of date. It outputs the success of the deletion. 

If you need to get info about your absences, you can use `/guildops-player-info`.

```shell
/guildops-absence-delete from: 07/10/23 to: 09/10/23

Absence(s) deleted for :
* Sat 07/10/23
* Sun 08/10/23
* Mon 09/10/23
```

**Requirements:**
* Date must be in format : dd/mm/yy or `dd/mm/yy au dd/mm/yy` for a range of date.
* Absences must have been created by `/guildops-absence-create` or it will show 
  ``` 
  No absence created or deleted, there is no raids on this range
  ```
**Errors :**
* If the date is malformed

  ``` Error while parsing date:parsing time "08/10/23A": extra text: "A"```
* If the date is in the past

  ``` You can't create or delete an absence in the past```

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
**Requirements:**
* Your discord must be linked to a player created by `/guildops-player-create` - by a guild officer - and linked by `/guildops-player-link` yourself.

**Errors :**
* If you are linked to no player or the player does not exist.

  ```Error while getting player infos: didn't find a player linked to this discord user named milowenn```

---

## Guild Officer actions

Strikes, Fails and Loots should only be used by guild officers.
* **Strike** : A strike is a warning for a player. It can be used for bad behavior, bad performance, etc.
* **Fail** : A fail is a fail on a boss. It can be used to track the fails of a player across the raids.
* **Loot** : A loot is a loot given to a player. It can be used to track the loots of a player across the raids.

### Create a raid <a name="introduction"></a>

It will create a raid with the name, date and difficulty specified. It outputs the raid id.

```shell
/guildops-raid-create name: example date: 30/09/23 difficulty: Mythic

Raid successfully created with ID 906348395984977921
```
* Difficulty should be : Normal, Heroic, Mythic.
* Date should be in format : dd/mm/yy.
* Name should be a string, from 1 to 12 characters.

**Errors :**
* If the raid already exists : same date on same difficulty.

  ``` Error while creating raid: raid already exists```
* If the player name is not compose only by letters.

  ``` Error while creating player: name must only contain letters```
* If the raid name is more than 12 characters or less than 1 character.

  ``` Error while creating raid: name must be between 1 and 12 characters```
* If the raid difficulty is not normal, heroic or mythic.

  ``` Error while creating raid: difficulty must be normal, heroic, or mythic```
* If the date is malformed

  ``` Error while creating raid: parsing time "01/101/23" as "02/01/06": cannot parse "1/23" as "/"```
### Create a player

It will create a player with the name specified. It outputs the player id.

```shell
/guildops-player-create name: example

Player a created successfully: ID 904508329427959809
```

**Requirements:**

* Name should be a string without space from 1 to 12 characters. If there is uppercase, it will be converted to lowercase.

**Errors :**
* If the player already exists.

  ``` Player chibrousse already exists```
* If the player name is not compose only by letters.

  ``` Error while creating player: name must only contain letters```
* If the player name is more than 12 characters or less than 1 character.

  ``` Error while creating player: name must be between 1 and 12 characters```

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

**Requirements:**
* Player must be created by `/guildops-player-create`.
* If there is not Loots, Fails, Absences or Strikes, it does not output them.

**Errors :**
* If the player does not exist.

  ``` Error while getting player infos: player not found```

### List raids
It lists all raids from a date to another date. It no to date is specified, it will list all raids from the date specified to from frield. 

```shell
/guildops-raid-list from: 30/09/23 to: 30/10/23

Raid List:
* example Sun 01/10/23 mythic 906348395984977921
* toto Tue 03/10/23 heroic 905258487187079169
* example Sat 07/10/23 mythic 906348159701581825

/guildops-raid-list from: 30/09/25 to: 30/10/25

no raid found
```

**Requirements:**
* Date must be in format : dd/mm/yy

**Errors:**
* If the date is malformed

  ```error while list raids: parsing time "30/101/24" as "02/01/06": cannot parse "1/24" as "/"```
### Create a strike

It will create a strike for the player specified. It outputs the strike id.

```shell
/guildops-strike-create name: milowenn reason: example of strike

Strike created successfully
```
Requirements:
* Name should be a string without space. If there is uppercase, it will be converted to lowercase.
* Name should be the name of a player already created.
* Reason should be a string, can use space and special characters. Cannot exceed 100 characters.

**Errors:**
* If the player does not exist.

  ``` Error while creating strike: player does not exist```
* If the player name is not compose only by letters.

  ``` Error while creating strike: name must only contain letters```
* If the player name is more than 12 characters or less than 1 character.

  ``` Error while creating strike: name must be between 1 and 12 characters```
* If the reason is more than 255 characters or less than 1 character.
 
  ``` Error while creating strike: reason must not be longer than 255 characterss```

  ``` Error while creating strike: reason must not be empty```



### List strikes on a player

It will list the strikes on the player specified. It outputs the strikes.

```shell
/guildops-strike-list name: milowenn

Strikes of milowenn (2) :
07/10/23 | example of strike | 906355752136933377
07/10/23 | example2 of strike | 906355886024785921
```

**Requirements:**
* Name should be a string without space. If there is uppercase, it will be converted to lowercase.
* Name should be the name of a player already created.

**Errors:**
* If the player does not exist.

  ``` Error while creating strike: player does not exist```
* If the player name is not compose only by letters.

  ``` Error while creating strike: name must only contain letters```
* If the player name is more than 12 characters or less than 1 character.

  ``` Error while creating strike: name must be between 1 and 12 characters```

### Delete a strike

It will delete the strike specified. To get the strike id, you can use `/guildops-strike-list`.

```shell
/guildops-strike-delete id: 903072156068708353

Strike deleted successfully
```

**Requirements:**
* Id should be the id of a strike created by `/guildops-strike-create`

**Errors:**
* If the strike does not exist.

  ``` strike not found```
* If the id is not a number.

  ``` Error while deleting strike: parsing "a": invalid syntax```

### Create a fail

It will create a fail for the player specified. It outputs the fail id.

```shell
/guildops-fail-create name: milowenn reason: Erreur P3 Sarkareth date: 30/09/23

Fail created successfully
```

**Requirements:**
* Date must be in format : dd/mm/yy
* Date should be a date of a raid created by `/guildops-raid-create`
* Name should be a string without space. If there is uppercase, it will be converted to lowercase.
* Name should be the name of a player already created.
* Reason should be a string

**Errors:**
* If the player does not exist or string is malformed.

  ``` Error while creating fail: player not found```

* If the reason is more than 255 characters or less than 1 character.
 
  ``` Error while creating fail: reason must not be longer than 255 characterss```

  ``` Error while creating fail: reason must not be empty```
* If the date is malformed

  ``` Error while creating fail: parsing time "30/101/24" as "02/01/06": cannot parse "1/24" as "/"```
* If the date is not a date of a raid created by `/guildops-raid-create`

  ``` Error while creating fail: raid not found```

### List fails on a player or on a raid

It will list the fails on the player or on the raid specified.

```shell
/guildops-fail-list-player name: milowenn

Fails de milowenn (2) :
28-09-2023 - p3 - 903072156068708353
30-09-2023 - Erreur P3 Sarkareth - 903072156068708353

/guildops-fail-list-player name: milowenn # with no fails

no fail found for milowenn
```
```shell
/guildops-fail-list-raid date: 30/09/23

Fails for 30/09/2023 (1) :
milowenn - Erreur P3 Sarkareth - 903072156068708353

/guildops-fail-list-raid date: 30/09/23 # with no fails

no fail found for 30/09/2023
```

**Requirements:**
* Date must be in format : dd/mm/yy
* Date should be a date of a raid created by `/guildops-raid-create`
* Name should be a string without space. If there is uppercase, it will be converted to lowercase.
* Name should be the name of a player already created.

**Errors:**
* If the player does not exist or string is malformed.

  ``` Fail to list fails on player```
* If the date is malformed

  ``` Error while parsing date:parsing time "08/10/23A": extra text: "A"```

### Delete a fail

It will delete the fail specified. To get the fail id, you can use `/guildops-fail-list-player name: <player_name>`.

```shell
/guildops-fail-delete id: 904435308715671553

Fail successfully deleted
```

**Requirements:**
* Id should be the id of a fail created by `/guildops-fail-create`

**Errors:**
* If the fail does not exist.

  ``` Error while deleting fail: fail not found from pg database```

### Attribute a loot

It attribute a loot to a player. 

```shell
/guildops-loot-attribute loot-name: example object raid-date: 03/10/23 player-name: milowenn

Loot successfully attributed
```

**Requirements:**
* Loot-name should be a string.
* Date should be the date of a raid created by `/guildops-raid-create`
* Player-name should be the name of a player already created by `/guildops-player-create`.

**Errors:**
* If the player does not exist.

  ``` Error while proceeding loot attribution: no player found```

* If the loot name is more than 20 characters or less than 1 character.
 
  ``` Error while creating loot: loot name is too long```

  ``` Error while creating loot: name must not be empty```
* If the date is malformed

  ``` invalid date```
* If the date is not a date of a raid created by `/guildops-raid-create`

  ``` Error while creating loot: raid not found```
* If the loot already exists.

  ```Error while proceeding loot attribution: database - CreateLoot - loot already exists```

### Select a player to attribute a loot

It takes a list of players and randomly pick a player with the lesser loots on the same difficulty and give a name. It does attribute directly the loot.

```shell
/guildops-loot-selector player-list: chibrousse,milowenn,prism difficulty: Mythic

milowenn have been selected to receive the loot
```
**Requirements:**
* Difficulty should be : Normal, Heroic, Mythic
* Players should be a list of players separated by a comma. If there is uppercase, it will be converted to lowercase.
* Players should be the name of a player already created by `/guildops-player-create`.

**Errors:**
* One of the players doesnt exist

    ``` Error while searching a player to attribute loot: no player found in database for milowenn```
* Difficulty is not normal, heroic or mythic

  ``` Error while searching a player to attribute loot: difficulty not valid. Must be Normal, Heroic or Mythic```

### List loots on a player

It will list the loots on the player specified. It outputs the loots.

```shell
/guildops-loot-list-on-player player-name:milowenn

All loots of milowenn:
* Test 01/10/23 mythic 123456789

/guildops-loot-list-on-player player-name:milowenn # with no loots

no loot for milowenn
```

**Requirements:**
* Player-name should be the name of a player already created by `/guildops-player-create`.

**Errors:**
* If the player does not exist.

  ``` Error while creating loot: no player found```

### List loots on a raid 

It will list the loots on the raid specified. It outputs the loots.

```shell
/guildops-loot-list-on-raid date: 01/10/23

All loots of 01/10/23:
* LootName milowenn mythic 123456789

/guildops-loot-list-on-raid date: 01/10/23 # with no loots

no loot for 08/10/23
```

**Requirements:**
* Date must be in format : dd/mm/yy and should be a date of a raid created by `/guildops-raid-create`


**Errors:**
* If the date is malformed

  ``` invalid date```
* If the date is not a date of a raid created by `/guildops-raid-create`

  ``` Error while listing loot for raid: raid not found```

### List Absences on a date

It will list the absences on the date specified. It outputs the absences.

```shell
/guildops-absence-list date: 28/09/23

08/10/23 absences :
* chibrousse

/guildops-absence-list date: 28/09/25 # with no absences

No absence for 28/09/23
```

**Requirements:**
* Date must be in format : dd/mm/yy
* Date should be a date of a raid created by `/guildops-raid-create`

**Errors:**
* If the date is malformed

  ``` Error while parsing date:parsing time "08/10/23A": extra text: "A"```
* If the date is not a date of a raid created by `/guildops-raid-create`

  ``` No absence for 09/10/23```

###  Delete a raid
**Warning : it will delete all the loots, fails, strikes and absences of the raid.**

It will delete the raid specified. To get the raid id, you can use `/guildops-raid-list`.

```shell
/guildops-raid-delete id:906348395984977921

Raid with ID 906348395984977921 successfully deleted
```

**Requirements:**
* Id should be the id of a raid created by `/guildops-raid-create`

**Errors:**
* If the raid does not exist.

  ```Error while deleting raid: database - DeleteRaid - raid not found```

* If the id is not a number.

  ```Error while deleting raid: parsing "A": invalid syntax```

### Delete a loot

It will delete the loot specified. To get the loot id, you can use `/guildops-loot-list-on-raid` or `/guildops-loot-list-on-player`.

```shell
/guildops-loot-delete id:465465465465465465

Loot successfully deleted
```
**Requirements:**
* Id should be the id of a loot created by `/guildops-loot-attribute`

**Errors:**
* If the loot does not exist.

  ``` Error while deleting loot: database - DeleteLoot - loot not found```
* If the id is not a number.

  ```id format is invalid```

### Delete a player

It will delete the player specified with a name. It removes all the loots, fails and absences of the player.

```shell
/guildops-player-delete name:milowenn

Player milowenn deleted successfully
```

**Requirements:**
* Player must exist

**Errors:**
* If the player does not exist or name malformed.

  ``` Error while deleting player: player not found```

### Delete or Create Absence

It creates or delete an absence for a player. 

```shell
/guildops-admin-absence-delete name:milowenn date: 28/09/23
```
```shell
/guildops-admin-absence-create name:milowenn date: 28/09/23
```
**Requirements:**
* Player must exist and have been created by `/guildops-player-create`
* Date must be in format : dd/mm/yy
* To is optionnal. If not specified, it will generate an absence for the date specified in from.
* Date cannot be in the past

**Errors:**
* If the player does not exist or name malformed.

  ``` Error while deleting player: player not found```

* If the date is malformed

  ``` Error while parsing date:parsing time "08/10/23A": extra text: "A"```

* If the date is in the past

  ``` Error while creating absence: date cannot be in the past```

### Create a range of raids

It creates a range of raids with the name, date and difficulty specified. It outputs the raid id.

```shell
/guildops-raid-create-multiple from: 30/09/23 to: 30/10/23 difficulty: Mythic weekdays: Monday, Wendesday

TODO
```

**Requirements:**
* Difficulty should be : Normal, Heroic, Mythic.
* Date must be in format : dd/mm/yy.
* Weekdays should be a list of weekdays separated by a comma. If there is uppercase, it will be converted to lowercase.
* to must be equal or after from.

**Errors :**
* If weekdays is malformed

    ```week days must be one of: Monday, Tuesday, Wednesday, Thursday, Friday, Saturday, Sunday```
* if to is before from 

    ```error while creating multiple raids: endDate is before startDate```
* if difficulty is incorrect 

    ```difficulty must be one of: Normal, Heroic, Mythic```
* if all raids already exists 

    ```no raid created```
