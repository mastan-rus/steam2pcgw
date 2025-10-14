package main

const game_article_template = `{{stub}}
{{Infobox game
|cover        = $game_title$ cover.jpg
|developers   = $developers$
|publishers   = $publishers$
|engines      = 
|release dates= $release_dates$
|reception    = 
{{Infobox game/row/reception|Metacritic|$metacritic_link$|$metacritic_rating$}}
{{Infobox game/row/reception|OpenCritic|$opencritic_link$|$opencritic_rating$}}
{{Infobox game/row/reception|IGDB|$igdb_link$|$igdb_rating$}}
|taxonomy     =
{{Infobox game/row/taxonomy/monetization      | $taxonomy_monetization$ }}
{{Infobox game/row/taxonomy/microtransactions | $taxonomy_microtransactions$ }}
{{Infobox game/row/taxonomy/modes             | $taxonomy_modes$ }}
{{Infobox game/row/taxonomy/pacing            | $taxonomy_pacing$ }}
{{Infobox game/row/taxonomy/perspectives      | $taxonomy_perspectives$ }}
{{Infobox game/row/taxonomy/controls          | $taxonomy_controls$ }}
{{Infobox game/row/taxonomy/genres            | $taxonomy_genres$ }}
{{Infobox game/row/taxonomy/sports            | $taxonomy_sports$ }}
{{Infobox game/row/taxonomy/vehicles          | $taxonomy_vehicles$ }}
{{Infobox game/row/taxonomy/art styles        | $taxonomy_artstyles$ }}
{{Infobox game/row/taxonomy/themes            | $taxonomy_themes$ }}
{{Infobox game/row/taxonomy/series            | $taxonomy_series$ }}
|steam appid  = $steam_appid$
|steam appid side = $steam_appid_side$
|gogcom id    = 
|gogcom id side = 
|official site= $official_site$
|hltb         = 
|lutris       = 
|mobygames    = 
|strategywiki = 
|wikipedia    = 
|winehq       = 
|license      = commercial
}}

{{Introduction
|introduction      = 

|release history   = 

|current state     = 
}}

'''General information'''
{{mm}} [https://steamcommunity.com/app/$steam_appid$/discussions/ Steam Community Discussions]

==Availability==
{{Availability|$availability_rows$
}}$availability_info$$demo_info$

==Monetization==
{{Monetization
|ad-supported                = 
|cross-game bonus            = 
|dlc                         = 
|expansion pack              = 
|freeware                    = 
|free-to-play                = $monetization_free_to_play$
|one-time game purchase      = $monetization_purchase$
|subscription                = 
|subscription gaming service = 
}}

===Microtransactions===
{{Microtransactions
|boost               = 
|cosmetic            = 
|currency            = 
|finite spend        = 
|infinite spend      = 
|free-to-grind       = 
|loot box            = 
|none                = $microtransactions_none$
|player trading      = 
|time-limited        = 
|unlock              = 
}}
$dlcs$
==Game data==
===Configuration file(s) location===
{{Game data|$game_data_config_rows$
}}$xdg$

===Save game data location===
{{Game data|$game_data_saves_rows$
}}

===[[Glossary:Save game cloud syncing|Save game cloud syncing]]===
{{Save game cloud syncing
|epic games launcher       = 
|epic games launcher notes = 
|gog galaxy                = 
|gog galaxy notes          = 
|origin                    = 
|origin notes              = 
|steam cloud               = $steam_cloud$
|steam cloud notes         = 
|ubisoft connect           = 
|ubisoft connect notes     = 
|xbox cloud                = 
|xbox cloud notes          = 
}}

==Video==
{{Video
|wsgf link                  = 
|widescreen wsgf award      = 
|multimonitor wsgf award    = 
|ultrawidescreen wsgf award = 
|4k ultra hd wsgf award     = 
|widescreen resolution      = unknown
|widescreen resolution notes= 
|multimonitor               = unknown
|multimonitor notes         = 
|ultrawidescreen            = unknown
|ultrawidescreen notes      = 
|4k ultra hd                = unknown
|4k ultra hd notes          = 
|fov                        = unknown
|fov notes                  = 
|windowed                   = unknown
|windowed notes             = 
|borderless windowed        = unknown
|borderless windowed notes  = 
|anisotropic                = unknown
|anisotropic notes          = 
|antialiasing               = unknown
|antialiasing notes         = 
|upscaling                  = unknown
|upscaling tech             = 
|upscaling notes            = 
|framegen                   = unknown
|framegen tech              = 
|framegen notes             = 
|vsync                      = unknown
|vsync notes                = 
|60 fps                     = unknown
|60 fps notes               = 
|120 fps                    = unknown
|120 fps notes              = 
|hdr                        = unknown
|hdr notes                  = 
|ray tracing                = unknown
|ray tracing notes          = 
|color blind                = unknown
|color blind notes          = 
}}

==Input==
{{Input
|key remap                 = unknown
|key remap notes           = 
|acceleration option       = unknown
|acceleration option notes = 
|mouse sensitivity         = unknown
|mouse sensitivity notes   = 
|mouse menu                = unknown
|mouse menu notes          = 
|keyboard and mouse prompts = unknown
|keyboard and mouse prompts notes = 
|invert mouse y-axis       = unknown
|invert mouse y-axis notes = 
|touchscreen               = unknown
|touchscreen notes         = 
|controller support        = $controller_support$
|controller support notes  = 
|full controller           = $full_controller$
|full controller notes     = 
|controller remap          = unknown
|controller remap notes    = 
|controller sensitivity    = unknown
|controller sensitivity notes = 
|invert controller y-axis  = unknown
|invert controller y-axis notes = 
|xinput controllers        = unknown
|xinput controllers notes  = 
|xbox prompts              = unknown
|xbox prompts notes        = 
|impulse triggers          = unknown
|impulse triggers notes    = 
|playstation controllers   = unknown
|playstation controller models = 
|playstation controllers notes = 
|playstation prompts       = unknown
|playstation prompts notes = 
|playstation motion sensors = unknown
|playstation motion sensors modes = 
|playstation motion sensors notes = 
|light bar support         = unknown
|light bar support notes   = 
|dualsense adaptive trigger support = unknown
|dualsense adaptive trigger support notes = 
|dualsense haptics support  = unknown
|dualsense haptics support notes = 
|playstation connection modes = unknown
|playstation connection modes notes = 
|tracked motion controllers = $tracked_motion_controllers$
|tracked motion controllers notes = 
|tracked motion prompts    = unknown
|tracked motion prompts notes = 
|other controllers         = unknown
|other controllers notes   = 
|other button prompts      = unknown
|other button prompts notes= 
|controller hotplug        = unknown
|controller hotplug notes  = 
|input prompt override     = unknown
|input prompt override notes = 
|haptic feedback           = unknown
|haptic feedback notes     = 
|haptic feedback hd        = unknown
|haptic feedback hd notes  = 
|haptic feedback hd controller models = 
|digital movement supported = unknown
|digital movement supported notes = 
|simultaneous input        = unknown
|simultaneous input notes  = 
|steam input api           = unknown
|steam input api notes     = 
|steam hook input          = unknown
|steam hook input notes    = 
|steam input prompts = unknown
|steam input prompts icons = 
|steam input prompts styles = 
|steam input prompts notes = 
|steam deck prompts  = unknown
|steam deck prompts notes = 
|steam controller prompts  = unknown
|steam controller prompts notes = 
|steam input motion sensors = unknown
|steam input motion sensors modes = 
|steam input motion sensors notes = 
|steam input presets = unknown
|steam input preset notes = 
|steam cursor detection    = unknown
|steam cursor detection notes = 
}}

==Audio==
{{Audio
|separate volume           = unknown
|separate volume notes     = 
|surround sound            = unknown
|surround sound notes      = 
|subtitles                 = $subtitles$
|subtitles notes           = 
|closed captions           = unknown
|closed captions notes     = 
|mute on focus lost        = unknown
|mute on focus lost notes  = 
|eax support               = 
|eax support notes         = 
|royalty free audio        = unknown
|royalty free audio notes  = 
|red book cd audio         = 
|red book cd audio notes   = 
|general midi audio        = 
|general midi audio notes  = 
}}

{{L10n|content=$languages$
}}
$network$$vrsupport$
==Other information==
===API===
{{API
|direct3d versions      = $d3d_versions$
|direct3d notes         = 
|directdraw versions    = 
|directdraw notes       = 
|wing                   = 
|wing notes             = 
|opengl versions        = 
|opengl notes           = 
|glide versions         = 
|glide notes            = 
|software mode          = 
|software mode notes    = 
|mantle support         = 
|mantle support notes   = 
|metal support          = 
|metal support notes    = 
|vulkan versions        = 
|vulkan notes           = 
|dos modes              = 
|dos modes notes        = 
|windows 32-bit exe     = $windows_x32$
|windows 64-bit exe     = $windows_x64$
|windows arm app        = unknown
|windows exe notes      = 
|mac os x powerpc app   = unknown
|macos intel 32-bit app = $osx_x32$
|macos intel 64-bit app = $osx_x64$
|macos arm app          = unknown
|macos app notes        = 
|linux powerpc app      = unknown
|linux 32-bit executable= $linux_x32$
|linux 64-bit executable= $linux_x64$
|linux arm app          = unknown
|linux 68k app          = unknown
|linux executable notes = 
|mac os powerpc app     = unknown
|mac os 68k app         = unknown
|mac os executable notes= 
}}

===Middleware===
{{Middleware
|physics          = 
|physics notes    = 
|audio            = 
|audio notes      = 
|interface        = 
|interface notes  = 
|input            = 
|input notes      = 
|cutscenes        = 
|cutscenes notes  = 
|multiplayer      = 
|multiplayer notes= 
|anticheat        = 
|anticheat notes  = 
}}

==System requirements==$system_requirements$
{{References}}
`
const developer_template = "\n{{Infobox game/row/developer|$developer$}}"

const publisher_template = "\n{{Infobox game/row/publisher|$publisher$}}"

const release_date_template = "\n{{Infobox game/row/date|$os_name$|$release_date$}}"

const availability_template = "\n{{Availability/row| $store_name$ | $store_id$ | $store_drm$ | $store_editions$ |  | $store_oses$ $store_unavailable$}}"

const dlcs_template = `
{{DLC|$dlc_rows$
}}
`
const dlc_row_template = "\n{{DLC/row| $dlc_name$ | | $dlc_oses$ }}"

const game_data_config_template = "\n{{Game data/config|$os_name$|}}"

const game_data_saves_template = "\n{{Game data/saves|$os_name$|}}"

const xdg_template = "\n{{XDG|unknown}}"

const language_template = `
{{L10n/switch
 |language  = $language_name$
 |interface = $language_interface$
 |audio     = $language_audio$
 |subtitles = $language_subtitles$
 |notes     = 
 |fan       = 
 |ref       = 
}}`

const network_template = `
==Network==
{{Network/Multiplayer
|local play           = $local_play$
|local play players   = 
|local play modes     = 
|local play notes     = 
|lan play             = $lan_play$
|lan play players     = 
|lan play modes       = 
|lan play notes       = 
|online play          = $online_play$
|online play players  = 
|online play modes    = 
|online play notes    = 
|asynchronous         = 
|asynchronous notes   = 
|crossplay            = 
|crossplay platforms  = 
|crossplay notes      = 
}}{{Network/Connections
|matchmaking        = 
|matchmaking notes  = 
|p2p                = 
|p2p notes          = 
|dedicated          = 
|dedicated notes    = 
|self-hosting       = 
|self-hosting notes = 
|direct ip          = 
|direct ip notes    = 
}}{{Network/Ports
|tcp  = 
|udp  = 
|upnp = 
}}
`

const vr_template = `
==VR support==
{{VR support
|native 3d                   = 
|native 3d notes             = 
|nvidia 3d vision            = 
|nvidia 3d vision notes      = 
|vorpx                       = 
|vorpx modes                 = 
|vorpx notes                 = 
|vr only                     = $vr_only$
|openxr                      = 
|openxr notes                = 
|steamvr                     = unknown
|steamvr notes               = 
|oculusvr                    = unknown
|oculusvr notes              = 
|windows mixed reality       = unknown
|windows mixed reality notes = 
|osvr                        = 
|osvr notes                  = 
|forte vfx1                  = 
|forte vfx1 notes            = 
|keyboard-mouse              = unknown
|keyboard-mouse notes        = 
|body tracking               = 
|body tracking notes         = 
|hand tracking               = 
|hand tracking notes         = 
|face tracking               = 
|face tracking notes         = 
|eye tracking                = 
|eye tracking notes          = 
|tobii eye tracking          = 
|tobii eye tracking notes    = 
|trackir                     = 
|trackir notes               = 
|3rd space gaming vest       = 
|3rd space gaming vest notes = 
|novint falcon               = 
|novint falcon notes         = 
|play area seated            = 
|play area seated notes      = 
|play area standing          = 
|play area standing notes    = 
|play area room-scale        = 
|play area room-scale notes  = 
}}
`

const system_requirements_template = `
{{System requirements
|OSfamily = $os_name$

|minTGT   = 
|minOS    = $min_os_versions$
|minCPU   = $min_cpu1$
|minCPU2  = $min_cpu2$
|minRAM   = $min_ram$
|minHD    = $min_hd$
|minGPU   = $min_gpu1$
|minGPU2  = $min_gpu2$
|minVRAM  = $min_vram$
|minDX    = $min_dx$

|recTGT   = 
|recOS    = $rec_os_versions$
|recCPU   = $rec_cpu1$
|recCPU2  = $rec_cpu2$
|recRAM   = $rec_ram$
|recHD    = $rec_hd$
|recGPU   = $rec_gpu1$
|recGPU2  = $rec_gpu2$
|recVRAM  = $rec_vram$
|recDX    = $rec_dx$

<!-- Please see the Editing Guide before filling in the following -->
|alt1Title = 
|alt1TGT   = 
|alt1OS    = 
|alt1CPU   = 
|alt1CPU2  = 
|alt1RAM   = 
|alt1HD    = 
|alt1GPU   = 
|alt1GPU2  = 
|alt1GPU3  = 
|alt1VRAM  = 

|alt2Title = 
|alt2TGT   = 
|alt2OS    = 
|alt2CPU   = 
|alt2CPU2  = 
|alt2RAM   = 
|alt2HD    = 
|alt2GPU   = 
|alt2GPU2  = 
|alt2GPU3  = 
|alt2VRAM  = 
|notes     = $notes$
}}
`
