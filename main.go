package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

func main() {
	getver := flag.Bool("v", false, "version")
	steamId := flag.String("steamid", "", "Steam app ID")
	flag.Parse()

	if *getver {
		fmt.Println(APP_NAME, VERSION, "(", GH_LINK, ")")
		return
	}

	fmt.Println("Running", APP_NAME, VERSION, "(", GH_LINK, ")")

	var err error = nil
	gameId := ""

	if len(*steamId) != 0 {
		gameId = *steamId
	} else {
		// Ask for input from the user
		for len(gameId) == 0 || err != nil {
			print("Insert the Steam app ID: ")
			gameId, err = TakeInput()
			if err != nil {
				fmt.Println(err)
			}
		}
	}

	fmt.Println("Fetching game app details...")

	var gameJson []byte
	gameJson, err = ParseGame(gameId, false)
	if err != nil {
		fmt.Println(err)
		return
	}

	game, err := UnmarshalGame(gameJson)
	if err != nil {
		fmt.Printf("An error occurred while attempting to unmarshal the JSON... (%s)\n", err)
		return
	} else if !game.Success {
		fmt.Println("The app ID provided does not exist or does not have a Store page...")
		return
	}

	var output string = game_article_template

	fmt.Println("* [2/26] Adding app cover")
	output = strings.ReplaceAll(output, "$game_title$", SanitiseName(game.Data.Name, true))

	fmt.Println("* [3/26] Adding app developers")
	var developers string = ""
	for _, developer := range game.Data.Developers {
		developers += strings.ReplaceAll(developer_template, "$developer$", SanitiseName(developer, false))
	}
	output = strings.ReplaceAll(output, "$developers$", developers)

	fmt.Println("* [4/26] Adding app publishers")
	var publishers string = ""
	for _, publisher := range game.Data.Publishers {
		if len(game.Data.Publishers) == 1 {
			skip := false
			for _, developer := range game.Data.Developers {
				if developer == publisher {
					skip = true
					break
				}
			}
			if skip {
				continue
			}
		}
		publishers += strings.ReplaceAll(publisher_template, "$publisher$", SanitiseName(publisher, false))
	}
	output = strings.ReplaceAll(output, "$publishers$", publishers)

	fmt.Println("* [5/26] Adding app release date")

	date := ""
	if game.HasSteamGenre(EarlyAccess) {
		date += "EA"
	} else if game.Data.ReleaseDate.ComingSoon {
		if success, _ := IsDate(game.Data.ReleaseDate.Date); success {
			date += ParseDate(game.Data.ReleaseDate.Date)
		} else {
			date += "TBA"
		}
	} else {
		date += ParseDate(game.Data.ReleaseDate.Date)
	}

	var dates string = ""
	if game.Data.Platforms.Windows {
		release_date := strings.ReplaceAll(release_date_template, "$os_name$", "Windows")
		release_date = strings.ReplaceAll(release_date, "$release_date$", date)
		dates += release_date
	}

	if game.Data.Platforms.MAC {
		release_date := strings.ReplaceAll(release_date_template, "$os_name$", "OS X")
		release_date = strings.ReplaceAll(release_date, "$release_date$", date)
		dates += release_date
	}

	if game.Data.Platforms.Linux {
		release_date := strings.ReplaceAll(release_date_template, "$os_name$", "Linux")
		release_date = strings.ReplaceAll(release_date, "$release_date$", date)
		dates += release_date
	}
	output = strings.ReplaceAll(output, "$release_dates$", dates)

	fmt.Println("* [6/26] Adding reception score")
	if game.Data.Metacritic != nil {
		meta, err := regexSubstr(game.Data.Metacritic.URL, `https://(?:www.)?metacritic.com/game/pc/([^?/]+)`)
		if err == nil {
			output = strings.ReplaceAll(output, "$metacritic_link$", meta)
			output = strings.ReplaceAll(output, "$metacritic_rating$", strconv.Itoa(game.Data.Metacritic.Score))
		} else {
			output = strings.ReplaceAll(output, "$metacritic_link$", "link")
			output = strings.ReplaceAll(output, "$metacritic_rating$", "rating")
		}
	} else if val, ok := game.Data.Ratings["Metascore"]; ok {
		meta, err := regexSubstr(val.URL, `https://(?:www.)?metacritic.com/game/pc/([^?/]+)`)
		if err == nil {
			output = strings.ReplaceAll(output, "$metacritic_link$", meta)
			output = strings.ReplaceAll(output, "$metacritic_rating$", strconv.Itoa(val.Score))
		} else {
			output = strings.ReplaceAll(output, "$metacritic_link$", "link")
			output = strings.ReplaceAll(output, "$metacritic_rating$", "rating")
		}
	} else {
		output = strings.ReplaceAll(output, "$metacritic_link$", "link")
		output = strings.ReplaceAll(output, "$metacritic_rating$", "rating")
	}

	if val, ok := game.Data.Ratings["OpenCritic"]; ok {
		output = strings.ReplaceAll(output, "$opencritic_link$", strings.TrimPrefix(val.URL, "https://opencritic.com/game/"))
		output = strings.ReplaceAll(output, "$opencritic_rating$", strconv.Itoa(val.Score))
	} else {
		output = strings.ReplaceAll(output, "$opencritic_link$", "link")
		output = strings.ReplaceAll(output, "$opencritic_rating$", "rating")
	}

	output = strings.ReplaceAll(output, "$igdb_link$", "link")
	output = strings.ReplaceAll(output, "$igdb_rating$", "rating")

	if game.Data.IsFree {
		fmt.Println("* [7/26] Game is F2P")
		output = strings.ReplaceAll(output, "$taxonomy_monetization$", "Free-to-play")
	} else {
		fmt.Println("* [7/26] Game is not F2P")
		output = strings.ReplaceAll(output, "$taxonomy_monetization$", "One-time game purchase")
	}

	fmt.Println("* [8/26] Taxonomy...")
	if !game.HasCategory(InAppPurchases) {
		output = strings.ReplaceAll(output, "$taxonomy_microtransactions$", "None")
	} else {
		output = strings.ReplaceAll(output, "$taxonomy_microtransactions$", "")
	}

	modes := ""
	if game.HasCategory(Singleplayer) {
		modes += "Singleplayer, "
	}

	if game.HasCategory(Multiplayer) {
		modes += "Multiplayer, "
	}

	modes = strings.TrimSuffix(modes, ", ")
	output = strings.ReplaceAll(output, "$taxonomy_modes$", modes)

	output = strings.ReplaceAll(output, "$taxonomy_pacing$", game.Data.Pacing)

	output = strings.ReplaceAll(output, "$taxonomy_perspectives$", game.Data.Perspectives)

	output = strings.ReplaceAll(output, "$taxonomy_controls$", game.Data.Controls)

	output = strings.ReplaceAll(output, "$taxonomy_genres$", game.Data.Genres)

	output = strings.ReplaceAll(output, "$taxonomy_sports$", game.Data.Sports)

	output = strings.ReplaceAll(output, "$taxonomy_vehicles$", game.Data.Vehicles)

	output = strings.ReplaceAll(output, "$taxonomy_artstyles$", game.Data.ArtStyles)

	output = strings.ReplaceAll(output, "$taxonomy_themes$", game.Data.Themes)

	if len(game.Data.Franchise) != 0 {
		output = strings.ReplaceAll(output, "$taxonomy_series$", game.Data.Franchise)
	} else {
		output = strings.ReplaceAll(output, "$taxonomy_series$", "")
	}

	output = strings.ReplaceAll(output, "$steam_appid$", gameId)

	if game.Data.Dlc != nil {
		sort.Slice(game.Data.Dlc, func(i, j int) bool { return game.Data.Dlc[i] < game.Data.Dlc[j] })
		var dlcs string = ""
		for _, v := range game.Data.Dlc {
			dlcs += fmt.Sprintf("%v, ", v)
		}
		dlcs = strings.TrimSuffix(dlcs, ", ")
		output = strings.ReplaceAll(output, "$steam_appid_side$", dlcs)
	} else {
		output = strings.ReplaceAll(output, "$steam_appid_side$", "")
	}

	if game.Data.Website != nil {
		output = strings.ReplaceAll(output, "$official_site$", *game.Data.Website)
	} else {
		output = strings.ReplaceAll(output, "$official_site$", game.Data.SupportInfo.URL)
	}

	fmt.Println("* [9/26] Processing introduction...")

	fmt.Println("* [10/26] Processing Availability!")

	platforms := ""
	if game.Data.Platforms.Windows {
		platforms += "Windows, "
	}
	if game.Data.Platforms.MAC {
		platforms += "OS X, "
	}
	if game.Data.Platforms.Linux {
		platforms += "Linux, "
	}

	platforms = strings.TrimSuffix(platforms, ", ")
	var editions []string
	trimPrice := regexp.MustCompile(`(\$[0-9.]+( USD)?)`)

	for _, v := range game.Data.PackageGroups {
		diplayType, _ := GetInt(v.DisplayType)
		if diplayType == 1 {
			continue
		}

		for _, sub := range v.Subs {
			edition := RemoveTags(sub.OptionText, "")
			fmt.Println(edition)
			edition = strings.ReplaceAll(edition, SanitiseName(game.Data.Name, true), "")
			edition = strings.TrimSpace(edition)
			edition = strings.TrimPrefix(edition, ": ")
			edition = strings.TrimPrefix(edition, "- ")
			edition = trimPrice.ReplaceAllLiteralString(edition, "")
			edition = strings.TrimSpace(edition)
			edition = strings.TrimSuffix(edition, " -")

			if len(edition) != 0 {
				editions = append(editions, "'''"+edition+"'''")
			}
		}
	}

	editionList := ""
	for i := 0; i < len(editions); i++ {
		editionList += editions[i]
		if i == len(editions)-2 {
			editionList += " and "
		} else {
			editionList += ", "
		}
	}

	if len(editionList) != 0 {
		editionList = strings.TrimSuffix(editionList, ", ")
		editionList += " also available"
	}

	rows := ""
	availability_row := availability_template
	availability_row = strings.ReplaceAll(availability_row, "$store_name$", "Steam")
	availability_row = strings.ReplaceAll(availability_row, "$store_id$", gameId)
	availability_row = strings.ReplaceAll(availability_row, "$store_drm$", "Steam")
	availability_row = strings.ReplaceAll(availability_row, "$store_editions$", editionList)
	availability_row = strings.ReplaceAll(availability_row, "$store_oses$", platforms)

	if len(game.Data.Packages) == 0 {
		availability_row = strings.ReplaceAll(availability_row, "$store_unavailable$", "| unavailable ")
	} else {
		availability_row = strings.ReplaceAll(availability_row, "$store_unavailable$", "")
	}

	rows += availability_row

	for store, data := range game.Data.Stores {
		availability_row = availability_template
		availability_row = strings.ReplaceAll(availability_row, "$store_name$", store)
		availability_row = strings.ReplaceAll(availability_row, "$store_id$", data.URL)
		availability_row = strings.ReplaceAll(availability_row, "$store_drm$", "DRM")
		availability_row = strings.ReplaceAll(availability_row, "$store_editions$", editionList)
		availability_row = strings.ReplaceAll(availability_row, "$store_oses$", data.Platforms)
		availability_row = strings.ReplaceAll(availability_row, "$store_unavailable$", "")
		rows += availability_row
	}

	output = strings.ReplaceAll(output, "$availability_rows$", rows)

	// Third party account check
	availability_info := ""
	if len(game.Data.ExternalAccountNotice) != 0 {
		availability_info += fmt.Sprintf("\n{{ii}} Requires 3rd-Party Account: %s", game.Data.ExternalAccountNotice)
	}

	// DRM check
	drms := ""
	if strings.Contains(game.Data.DRMNotice, "Denuvo") {
		drms += "{{DRM|Denuvo}}, "
	}
	// if strings.Contains(game.Data.PCRequirements["minimum"].(string), "VMProtect") {
	// 	drms += "{{DRM|VMProtect}}, "
	// }

	drms = strings.TrimSuffix(drms, ", ")
	if len(drms) == 0 {
		drms += game.Data.DRMNotice
		if len(drms) != 0 {
			availability_info += fmt.Sprintf("\n{{ii}} All versions require %s.", drms)
		}
	}

	if len(editionList) > 1 {
		availability_info += "\n\n===Version differences===\n{{ii}} "
		availability_info += editionList
	}

	availability_info += "\n\n<!-- PAGE GENERATED BY STEAM2PCGW -->"

	output = strings.ReplaceAll(output, "$availability_info$", availability_info)

	fmt.Println("* [11/26] Processing Monetization!")

	if game.Data.IsFree {
		output = strings.ReplaceAll(output, "$monetization_free_to_play$", "The game has such monetization.")
	} else {
		output = strings.ReplaceAll(output, "$monetization_free_to_play$", "")
	}

	if !game.Data.IsFree {
		output = strings.ReplaceAll(output, "$monetization_purchase$", "The game requires an upfront purchase to access.")
	} else {
		output = strings.ReplaceAll(output, "$monetization_purchase$", "")
	}

	fmt.Println("* [12/26] Processing Microtransactions!")

	if !game.HasCategory(InAppPurchases) {
		output = strings.ReplaceAll(output, "$microtransactions_none$", "None")
	} else {
		output = strings.ReplaceAll(output, "$microtransactions_none$", "")
	}

	fmt.Println("* [13/26] Processing DLCs!")
	var addedDlcs int = 0
	if game.Data.Dlc != nil {
		dlcs := ""

		for i, dlcsteamid := range game.Data.Dlc {
			// TODO dlcs
			fmt.Printf("  Processing DLC [%d/%d] steamid %d... ", i+1, len(game.Data.Dlc), dlcsteamid)
			var dlcJson []byte
			dlcJson, err = ParseGame(strconv.FormatInt(dlcsteamid, 10), true)
			if err != nil {
				fmt.Println(err)
				continue
			}

			dlc, err := UnmarshalGame(dlcJson)
			if err != nil {
				fmt.Printf("An error occurred while attempting to unmarshal the JSON... (%s)\n", err)
				continue
			} else if !dlc.Success {
				fmt.Println("DLC Steam ID provided does not exist or does not have a Store page...")
				continue
			}

			var dlcname = dlc.Data.Name
			dlcname = strings.TrimPrefix(dlcname, game.Data.Name) // Removing game name
			dlcname = strings.TrimLeft(dlcname, " -:")
			dlcname = strings.TrimPrefix(dlcname, "DLC") // "DLC" can be left sometimes, removing it too
			dlcname = strings.TrimLeft(dlcname, " -:")
			dlcname = strings.ReplaceAll(dlcname, "=", "&#61;") // Replacing unsupported symbols

			if strings.Contains(strings.ToLower(dlcname), "soundtrack") {
				fmt.Printf("  Skipping soundtrack '%s'\n", dlc.Data.Name)
				continue
			}

			if strings.Contains(strings.ToLower(dlcname), "artwork") ||
				strings.Contains(strings.ToLower(dlcname), "concept art") ||
				strings.Contains(strings.ToLower(dlcname), "digital art") {
				fmt.Printf("  Skipping artwork '%s'\n", dlc.Data.Name)
				continue
			}

			fmt.Printf(" Adding DLC '%s'\n", dlcname)

			var oses string = ""
			if game.Data.Platforms.Windows {
				oses += "Windows, "
			}
			if game.Data.Platforms.MAC {
				oses += "OS X, "
			}
			if game.Data.Platforms.Linux {
				oses += "Linux, "
			}
			oses = strings.TrimSuffix(oses, ", ")
			dlc_row := dlc_row_template
			dlc_row = strings.ReplaceAll(dlc_row, "$dlc_name$", dlcname)
			dlc_row = strings.ReplaceAll(dlc_row, "$dlc_oses$", oses)

			dlcs += dlc_row
			addedDlcs++
		}

		fmt.Printf("  Added %d DLCs\n", addedDlcs)
		// insert comment if no DLCs were actually added
		if addedDlcs == 0 {
			output = strings.ReplaceAll(output, "$dlcs$", "")
		} else {
			dlc_table := strings.ReplaceAll(dlcs_template, "$dlc_rows$", dlcs)
			output = strings.ReplaceAll(output, "$dlcs$", dlc_table)
		}
	} else {
		output = strings.ReplaceAll(output, "$dlcs$", "")
	}

	fmt.Println("* [14/26] Processing Config File Location!")

	config_rows := ""
	if game.Data.Platforms.Windows {
		config_rows += strings.ReplaceAll(game_data_config_template, "$os_name$", "Windows")
	}
	if game.Data.Platforms.MAC {
		config_rows += strings.ReplaceAll(game_data_config_template, "$os_name$", "OS X")
	}
	if game.Data.Platforms.Linux {
		config_rows += strings.ReplaceAll(game_data_config_template, "$os_name$", "Linux")
	}
	output = strings.ReplaceAll(output, "$game_data_config_rows$", config_rows)

	if game.Data.Platforms.Linux {
		output = strings.ReplaceAll(output, "$xdg$", xdg_template)
	} else {
		output = strings.ReplaceAll(output, "$xdg$", "")
	}

	fmt.Println("* [15/26] Processing Save Game Location!")
	saves_rows := ""
	if game.Data.Platforms.Windows {
		saves_rows += strings.ReplaceAll(game_data_saves_template, "$os_name$", "Windows")
	}
	if game.Data.Platforms.MAC {
		saves_rows += strings.ReplaceAll(game_data_saves_template, "$os_name$", "OS X")
	}
	if game.Data.Platforms.Linux {
		saves_rows += strings.ReplaceAll(game_data_saves_template, "$os_name$", "Linux")
	}
	output = strings.ReplaceAll(output, "$game_data_saves_rows$", saves_rows)

	fmt.Println("* [16/26] Processing Save Game Sync!")

	// Game has steam cloud, then we can just add it
	// Otherwise, we can check if the game is out yet or not
	// to determine whether we should add `unknown` or `false`
	if game.HasCategory(SteamCloud) {
		output = strings.ReplaceAll(output, "$steam_cloud$", "true")
	} else {
		if game.Data.ReleaseDate.ComingSoon {
			output = strings.ReplaceAll(output, "$steam_cloud$", "unknown")
		} else {
			output = strings.ReplaceAll(output, "$steam_cloud$", "false")
		}
	}

	// TODO: Scan the description to search for widescreen, ray tracing etc support
	fmt.Println("* [17/26] Processing Video!")

	fmt.Println("* [18/26] Processing Input!")

	controller := false
	if game.Data.ControllerSupport != nil {
		controller = true
	}

	output = strings.ReplaceAll(output, "$controller_support$", strconv.FormatBool(controller))
	if controller {
		if *game.Data.ControllerSupport == "full" {
			output = strings.ReplaceAll(output, "$full_controller$", "true")
		} else {
			output = strings.ReplaceAll(output, "$full_controller$", "false")
		}
	} else {
		output = strings.ReplaceAll(output, "$full_controller$", "unknown")
	}

	if game.HasCategory(TrackedControllerSupport) {
		output = strings.ReplaceAll(output, "$tracked_motion_controllers$", "true")
	} else {
		output = strings.ReplaceAll(output, "$tracked_motion_controllers$", "unknown")
	}

	fmt.Println("* [19/26] Processing Audio!")

	game.ProcessLanguages()

	// Better leave this field unknown since it can be n/a or always on instead of true
	output = strings.ReplaceAll(output, "$subtitles$", "unknown") // strconv.FormatBool(game.Data.Subtitles)

	fmt.Println("* [20/26] Processing Languages!")

	orderedLangauges := make([]string, 0, len(game.Data.Languages))
	for key := range game.Data.Languages {
		sanitisedKey := key
		orderedLangauges = append(orderedLangauges, sanitisedKey)
	}

	sort.Strings(orderedLangauges)

	// find English and swap it to be the first language instead...

	if orderedLangauges[0] != "English" {
		// Only swap English to the first language if it isn't already...

		foundIndex := 0
		for i := 1; i < len(orderedLangauges); i++ {
			if orderedLangauges[i] == "English" {
				foundIndex = i
				break
			}
		}

		if foundIndex != 0 {
			for i := foundIndex; i > 0; i-- {
				orderedLangauges[i] = orderedLangauges[i-1]
			}
			orderedLangauges[0] = "English"
		}
	}

	langs := ""
	for _, key := range orderedLangauges {
		langs += game.FormatLanguage(key)
	}

	output = strings.ReplaceAll(output, "$languages$", langs)

	fmt.Println("* [21/26] Processing Network!")

	if game.HasCategory(Multiplayer) {
		network := network_template
		network = strings.ReplaceAll(network, "$local_play$",
			strconv.FormatBool(game.HasCategory(LocalMultiPlayer) || game.HasCategory(LocalCoOp)))
		network = strings.ReplaceAll(network, "$lan_play$",
			strconv.FormatBool(game.HasCategory(CoOp)))
		network = strings.ReplaceAll(network, "$online_play$",
			strconv.FormatBool(game.HasCategory(OnlineMultiPlayer) || game.HasCategory(OnlineCoOp)))
		output = strings.ReplaceAll(output, "$network$", network)
	} else {
		output = strings.ReplaceAll(output, "$network$", "")
	}

	fmt.Println("* [22/26] Processing VR!")
	if game.HasCategory(VRSupport) || game.HasCategory(VRSupported) || game.HasCategory(VROnly) {
		vr := vr_template
		vr = strings.ReplaceAll(vr, "$vr_only$", strconv.FormatBool(game.HasCategory(VROnly)))
		output = strings.ReplaceAll(output, "$vrsupport$", vr)
	} else {
		output = strings.ReplaceAll(output, "$vrsupport$", "")
	}

	fmt.Println("* [23/26] Processing API!")

	// D3D is empty for now. Better not rely on any info from system requirements.
	output = strings.ReplaceAll(output, "$d3d_versions$", "") // game.FindDirectX()
	output = strings.ReplaceAll(output, "$windows_x32$", GetExeBit(true, "windows", game.Data.Platforms, game.Data.PCRequirements))
	output = strings.ReplaceAll(output, "$windows_x64$", GetExeBit(false, "windows", game.Data.Platforms, game.Data.PCRequirements))
	output = strings.ReplaceAll(output, "$osx_x32$", GetExeBit(true, "mac", game.Data.Platforms, game.Data.MACRequirements))
	output = strings.ReplaceAll(output, "$osx_x64$", GetExeBit(false, "mac", game.Data.Platforms, game.Data.MACRequirements))
	output = strings.ReplaceAll(output, "$linux_x32$", GetExeBit(true, "linux", game.Data.Platforms, game.Data.LinuxRequirements))
	output = strings.ReplaceAll(output, "$linux_x64$", GetExeBit(false, "linux", game.Data.Platforms, game.Data.LinuxRequirements))

	fmt.Println("* [24/26] Processing Middleware!")

	fmt.Println("* [25/26] Processing System Requirements!")

	output = strings.ReplaceAll(output, "$system_requirements$", game.OutputSpecs())

	fmt.Println("* [26/26] Processing References!")

	outputFile, err := os.Create(fmt.Sprintf("output/%s.txt", gameId))
	if err != nil {
		fmt.Println("Failed to create the output file.")
		return
	}

	outputFile.WriteString(output)

	outputFile.Close()

	println(fmt.Sprintf("Successfully parsed information for game: '%s'", SanitiseName(game.Data.Name, true)))
}
