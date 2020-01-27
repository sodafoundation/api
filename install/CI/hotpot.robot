*** Settings ***
Library        hotpot.Hotpot  WITH NAME  CLI

*** Variables ***
${volume_id}   ""
${fileshare_id}   ""
${file_profile_id}   ""
${block_profile_id}   ""

*** Test Cases ***
1st Hotpot Test
   CLI.Hello World

Version List
   CLI.Version List

Dock List
   CLI.Dock List

Profile List
   ${out}  ${err} =  CLI.Profile List
   Log  Profile output is ${out}

Pool List
   CLI.Pool List

Volume List
   CLI.Volume List

Fileshare List
   CLI.Fileshare List

Profile Create Block
   ${block_profile_id} =  CLI.Profile Create Block
   Set Global Variable   ${block_profile_id}
   Log  The block_profile_id is ${block_profile_id}

Profile Create File
   ${file_profile_id} =  CLI.Profile Create File
   Set Global Variable   ${file_profile_id}
   Log  The file_profile_id is ${file_profile_id}

Volume Create
   ${volume_id} =  CLI.Volume Create
   Set Global Variable   ${volume_id}
   Log  The volume_id is ${volume_id}

Check Volume Status Available
   CLI.Check volume status available  ${volume_id}

Volume Delete
   Log  The volume_id is ${volume_id}
   CLI.Volume Delete  ${volume_id}

Fileshare Create
   ${fileshare_id} =  CLI.Fileshare Create
   Set Global Variable   ${fileshare_id}
   Log  The volume_id is ${fileshare_id}

Check FileShare Status Available
   CLI.Check fileshare status available  ${fileshare_id}

Fileshare Delete
    Log  The fileshare_id is ${fileshare_id}
    CLI.Fileshare Delete  ${fileshare_id}

Profile Delete File
   CLI.Profile Delete File  ${file_profile_id}

Profile Delete Block
   CLI.Profile Delete Block  ${block_profile_id}
