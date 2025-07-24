package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/html"
)

func GetInt(v interface{}) (int, error) {
	switch v := v.(type) {
	case float64:
		return int(v), nil
	case string:
		c, err := strconv.Atoi(v)
		if err != nil {
			return 0, err
		}
		return c, nil
	default:
		return 0, fmt.Errorf("conversion to int from %T not supported", v)
	}
}

func (req *Requirement) UnmarshalJSON(data []byte) error {
	if string(data) == `""` || string(data) == `{}` || string(data) == `[]` {
		return nil
	}

	type requirement Requirement
	return json.Unmarshal(data, (*requirement)(req))
}

func UnmarshalGame(data []byte) (result Game, err error) {
	var tempResult map[string]Game
	err = json.Unmarshal(data, &tempResult)
	if err != nil {
		return
	}

	var gameId string
	for key := range tempResult {
		gameId = key
		break
	}

	result = Game(tempResult[gameId])
	result.Data.Ratings = make(map[string]Rating)
	result.Data.Stores = make(map[string]Store)

	var scrapeData []byte
	scrapeData, err = os.ReadFile("cache/" + gameId + ".html")
	if err != nil {
		fmt.Printf("Failed to read scraped Steam page data")
	} else {
		franchiseNames := regexp.MustCompile(`<div class="dev_row">\s*<b>Franchise:</b>\s*<a href=".*">([^<]+)</a>\s*</div>`).FindStringSubmatch(string(scrapeData))
		if len(franchiseNames) > 1 {
			franchiseName := RemoveTags(html.UnescapeString(franchiseNames[0]), "")
			franchiseName = strings.ReplaceAll(franchiseName, "Franchise:", "")
			franchiseName = strings.TrimSpace(franchiseName)
			result.SetFranchise(franchiseName)
		}

		dirtyTags := regexp.MustCompile(`<a href=".+" class="app_tag" style=".+">\s+(.+)\s+<\/a>{1,}`).FindAllStringSubmatch(string(scrapeData), 50)
		var appTags []string
		for _, tag := range dirtyTags {
			cleanTag := html.UnescapeString(tag[1])
			cleanTag = strings.Replace(cleanTag, "+", "", 1)
			cleanTag = strings.Replace(cleanTag, "Point & Click", "Point and Select", 1)
			cleanTag = strings.TrimSpace(cleanTag)

			appTags = append(appTags, cleanTag)
		}

		result.SetPacing(appTags)
		result.SetPerspective(appTags)
		result.SetControls(appTags)
		result.SetGenres(appTags)
		result.SetSports(appTags)
		result.SetVehicles(appTags)
		result.SetArtStyles(appTags)
		result.SetThemes(appTags)
	}

	// Is There Any Deals
	response, optionalErr := makeRequest(fmt.Sprintf("https://isthereanydeal.com/steam/app/%s", gameId))
	if optionalErr = checkRequest(response, optionalErr); optionalErr == nil {
		defer response.Body.Close()
		body, _ := parseResponseToBody(response)
		htmlString := string(body)

		result.parseReviews(htmlString)
		result.parseAvailability(htmlString)
	} else {
		fmt.Println("Failed to scrape IsThereAnyDeals page...")
	}

	return
}

func makeRequest(url string) (*http.Response, error) {
	client := &http.Client{}
	getData := strings.NewReader("")
	req, _ := http.NewRequest("GET", url, getData)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Safari/537.36")
	req.Header.Set("Cookie", "birthtime=0; max-age=315360000;")
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func doesCacheExistOrLatest(fileName string) bool {
	fi, err := os.Stat(fileName)
	return err == nil && time.Since(fi.ModTime()).Hours() < (7*24)
}

func createCache(gameId string, apiBody []byte, scrapeBody []byte) (err error) {
	err = os.WriteFile("cache/"+gameId+".json", apiBody, 0777)
	if len(scrapeBody) != 0 {
		os.WriteFile("cache/"+gameId+".html", scrapeBody, 0777)
	}
	return
}

func checkRequest(response *http.Response, err error) error {
	if err != nil {
		fmt.Printf("Failed to connect (error: %s)\n", err)
	} else if response.StatusCode != http.StatusOK {
		fmt.Printf("Failed to connect to the '%v'... (HTTP code: %d)\n", response.Request.URL, response.StatusCode)
		err = errors.New("status code not OK")
	}

	return err
}

func parseResponseToBody(response *http.Response) (body []byte, err error) {
	body, err = io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("An error occurred while attempting to parse the response body...")
	}
	return
}

func fetchGame(gameId string) (err error) {
	var response *http.Response
	var apiBody []byte
	var scrapeBody []byte

	response, err = makeRequest(fmt.Sprintf("%s%s%s&cc=us", API_LINK, gameId, LOCALE))
	if err = checkRequest(response, err); err != nil {
		return
	}
	defer response.Body.Close()
	apiBody, err = parseResponseToBody(response)
	if err != nil {
		return
	}

	// TODO
	// optionalResponse, optionalErr := makeRequest(fmt.Sprintf("https://cdn.cloudflare.steamstatic.com/steam/apps/%s/library_600x900_2x.jpg", gameId))
	// if optionalErr = checkRequest(response, optionalErr); optionalErr == nil {
	// 	defer optionalResponse.Body.Close()
	// 	scrapeBody, _ = parseResponseToBody(optionalResponse)
	// 	file, optionalErr := os.Create(gameId + ".jpg")
	// 	if optionalErr == nil {
	// 		defer file.Close()
	// 		_, optionalErr = io.Copy(file, optionalResponse.Body)
	// 		if optionalErr == nil {
	// 			fmt.Println("Downloaded game cover!")
	// 		}
	// 	}
	// }
	// if optionalErr != nil {
	// 	fmt.Println("Game cover download failed")
	// }

	optionalResponse, optionalErr := makeRequest(fmt.Sprintf("https://store.steampowered.com/app/%s/?cc=us%s", gameId, LOCALE))
	if optionalErr = checkRequest(response, optionalErr); optionalErr == nil {
		defer optionalResponse.Body.Close()
		scrapeBody, _ = parseResponseToBody(optionalResponse)
	} else {
		fmt.Println("Failed to scrape Steam Store page...")
	}

	err = createCache(gameId, apiBody, scrapeBody)
	if err != nil {
		fmt.Println("Failed to create the cache, but continuing the process...")
	} else {
		fmt.Println("Cached!")
	}

	return err
}

func ParseGame(gameId string, shouldsleep bool) (body []byte, err error) {
	os.Mkdir("cache", 0777)
	os.Mkdir("output", 0777)

	fileName := fmt.Sprintf("cache/%s.json", gameId)

	if doesCacheExistOrLatest(fileName) {
		fmt.Println("Found cache...")
		body, err = os.ReadFile(fileName)
		return
	}

	fmt.Println("Did not find game cache or cache is older than 7 days...")

	// Optional 1-second sleep to not query DLC pages too often
	if shouldsleep {
		time.Sleep(time.Second)
	}

	err = fetchGame(gameId)
	if err == nil {
		body, err = os.ReadFile(fileName)
	}

	return body, err
}

func TakeInput() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')

	// For Windows and Linux
	text = strings.TrimSuffix(text, "\n")
	// For Windows
	text = strings.TrimSuffix(text, "\r")

	if len(text) == 0 {
		return "", errors.New("invalid input")
	}

	return text, nil
}

func GetExeBit(is32 bool, platform string, platforms Platforms, requirements Requirement) string {
	value := "unknown"
	if requirements == nil {
		return value
	}

	if (platform == "windows" && platforms.Windows) || (platform == "mac" && platforms.MAC) || (platform == "linux" && platforms.Linux) {
		var sanitised = strings.ToLower(requirements["minimum"].(string))
		sanitised = RemoveTags(sanitised, "\n")

		if strings.Contains(sanitised, "Requires a 64-bit processor and operating system") {
			if is32 {
				value = "false"
			} else {
				value = "true"
			}
		} else if strings.Contains(sanitised, "32/64") {
			value = "true"
		} else {
			if requirements["recommended"] != nil {
				sanitised = strings.ToLower(requirements["recommended"].(string))
				sanitised = RemoveTags(sanitised, "\n")
			}
			ramFinder := regexp.MustCompile(`memory:(\d+) gb`)
			ramFound := ramFinder.FindStringSubmatch(sanitised)

			var ram = 0
			if len(ramFound) != 0 {
				ram, _ = strconv.Atoi(ramFound[1])
				ram *= 1000
			} else {
				ramFinder = regexp.MustCompile(`memory:(\d+) mb`)
				ramFound = ramFinder.FindStringSubmatch(sanitised)
				if len(ramFound) != 0 {
					ram, _ = strconv.Atoi(ramFound[1])
				}
			}

			if is32 && (strings.Contains(sanitised, "64-bit") || strings.Contains(sanitised, "64 bit") || ram > 4096) {
				value = "false"
			} else {
				value = "true"
			}
		}
	}

	fmt.Printf("* [23/26] %s (32-bit: %v): %s\n", platform, is32, value)

	return value
}

func RemoveTags(input, replacement string) string {
	noTag, _ := regexp.Compile(`(<[^>]*>)+`)
	output := noTag.ReplaceAllLiteralString(input, replacement)
	output = strings.ReplaceAll(output, "\n ", "")
	return output
}

func formatCitation(note string) string {
	return `{{cn|` + note + `}}`
}

func (game *Game) parseAvailability(htmlString string) {
	doc, _ := html.Parse(strings.NewReader(htmlString))

	var f func(n *html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "table" {
			for _, a := range n.Attr {
				if a.Key != "class" || a.Val != "t-st3 priceTable" {
					continue
				}

				for c := n.FirstChild; c != nil; c = c.NextSibling {
					if c.Type != html.ElementNode || c.Data != "tbody" {
						continue
					}

					for d := c.FirstChild; d != nil; d = d.NextSibling {
						if d.Type != html.ElementNode || d.Data != "tr" {
							continue
						}

						for e := d.FirstChild; e != nil; e = e.NextSibling {
							if e.Type != html.ElementNode || e.Data != "td" {
								continue
							}

							for _, f := range e.Attr {
								if f.Key != "class" || f.Val != "priceTable__shop" {
									continue
								}

								for g := e.FirstChild; g != nil; g = g.NextSibling {
									if g.Type != html.ElementNode || g.Data != "a" {
										continue
									}

									var store, platform, url string
									// cut, current, lowest, regular

									for _, h := range g.Attr {
										if h.Key == "href" {
											url = h.Val
										}
									}

									for i := g.FirstChild; i != nil; i = i.NextSibling {
										if i.Type == html.TextNode {
											store = i.Data
										}
									}

									e = e.NextSibling
									for _, h := range e.Attr {
										if h.Key == "class" && h.Val == "priceTable__platforms" {
											for i := e.FirstChild; i != nil; i = i.NextSibling {
												if i.Type == html.TextNode {
													platform = i.Data
												}
											}
										}
									}

									// e = e.NextSibling
									// for _, h := range e.Attr {
									// 	if h.Key == "class" && h.Val == "priceTable__cut t-st3__num" {
									// 		for i := e.FirstChild; i != nil; i = i.NextSibling {
									// 			if i.Type == html.TextNode {
									// 				cut = i.Data
									// 			}
									// 		}
									// 	}
									// }

									// e = e.NextSibling
									// for _, h := range e.Attr {
									// 	if h.Key == "class" && h.Val == "priceTable__new t-st3__price s-low g-low" {
									// 		for i := e.FirstChild; i != nil; i = i.NextSibling {
									// 			if i.Type == html.TextNode {
									// 				current = i.Data
									// 			}
									// 		}
									// 	}
									// }

									// e = e.NextSibling
									// for _, h := range e.Attr {
									// 	if h.Key == "class" && h.Val == "priceTable__low t-st3__price s-low g-low" {
									// 		for i := e.FirstChild; i != nil; i = i.NextSibling {
									// 			if i.Type == html.TextNode {
									// 				lowest = i.Data
									// 			}
									// 		}
									// 	}
									// }

									// e = e.NextSibling
									// for _, h := range e.Attr {
									// 	if h.Key == "class" && h.Val == "priceTable__old t-st3__price" {
									// 		for i := e.FirstChild; i != nil; i = i.NextSibling {
									// 			if i.Type == html.TextNode {
									// 				regular = i.Data
									// 			}
									// 		}
									// 	}
									// }

									game.AddStore(store, platform, url)
									// fmt.Printf("Store: %s, Platforms: %s, Price Cut: %s, Current: %s, Lowest: %s, Regular: %s\n", store, platform, cut, current, lowest, regular)
								}
							}
						}
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
}

func regexSubstr(input string, format string) (string, error) {
	infoRe := regexp.MustCompile(format)
	match := infoRe.FindSubmatch([]byte(input))
	if match == nil {
		return "", errors.New("not found")
	}

	return string(match[1]), nil
}

func (game *Game) parseReviews(htmlString string) {
	body, _ := regexSubstr(htmlString, `var page = (.+);`)
	var allBody []interface{}
	json.Unmarshal([]byte(body), &allBody)

	if allBody == nil {
		return
	}

	mainBody := allBody[1]

	details := mainBody.(map[string]interface{})["detail"]
	if details == nil {
		return
	}
	reviews := details.(map[string]interface{})["reviews"]
	if reviews == nil {
		return
	}

	reviewList := reviews.([]interface{})

	for _, v := range reviewList {
		nn := v.(map[string]interface{})
		count := nn["count"].(float64)

		if count <= 4 {
			continue
		}

		positive := (int64)(nn["positive"].(float64))
		source := nn["source"].(string)
		url := nn["url"].(string)
		game.AddRating(source, strconv.FormatInt((positive), 10), url)
	}
}

func (game *Game) AddStore(name, platforms, link string) {
	validStores := []ValidStore{
		{
			ScrapeName:  "Blizzard",
			DisplayName: "Battle.net",
			LinkToStrip: "https://us.shop.battle.net/en-us/product/",
		},
		{
			ScrapeName:  "Discord",
			DisplayName: "Discord",
			LinkToStrip: "https://discordapp.com/store/skus/",
		},
		{
			ScrapeName:  "Epic Game Store",
			DisplayName: "Epic Games Store",
			LinkToStrip: "https://www.epicgames.com/store/en-US/product/",
		},
		{
			ScrapeName:  "GamersGate",
			DisplayName: "GamersGate",
			LinkToStrip: "https://www.gamersgate.com/product/",
		},
		{
			ScrapeName:  "GamesPlanet UK",
			DisplayName: "GamesPlanet",
			LinkToStrip: "https://uk.gamesplanet.com/game/",
		},
		{
			ScrapeName:  "GOG.com",
			DisplayName: "GOG.com",
			LinkToStrip: "https://www.gog.com/game/",
		},
		{
			ScrapeName:  "GreenManGaming",
			DisplayName: "Green Man Gaming",
			LinkToStrip: "https://www.greenmangaming.com/games/",
		},
		{
			ScrapeName:  "Humble Store",
			DisplayName: "Humble",
			LinkToStrip: "https://www.humblebundle.com/store/",
		},
		{
			ScrapeName:  "Itch.io",
			DisplayName: "Itch.io",
			LinkToStrip: "",
		},
		{
			ScrapeName:  "Origin",
			DisplayName: "Origin",
			LinkToStrip: "https://www.origin.com/store/",
		},
	}

	key := -1
	for k, v := range validStores {
		if v.ScrapeName != name {
			continue
		}

		key = k
		name = validStores[key].DisplayName
		break
	}

	if key == -1 {
		return
	}

	sanitised := platforms
	sanitised = strings.Replace(sanitised, "Win", "Windows", 1)
	sanitised = strings.Replace(sanitised, "Mac", "OS X", 1)

	game.Data.Stores[name] = Store{
		Platforms: sanitised,
		URL:       strings.TrimPrefix(link, validStores[key].LinkToStrip),
	}
}

func (game *Game) AddRating(name, scoreString, link string) {
	if _, ok := game.Data.Ratings[name]; ok {
		return
	}

	score, _ := strconv.Atoi(scoreString)

	game.Data.Ratings[name] = Rating{
		Score: score,
		URL:   link,
	}
}

func (game *Game) FindDirectX() string {
	if game.Data.PCRequirements == nil {
		return ""
	}

	if len(game.Data.PCRequirements["minimum"].(string)) == 0 {
		return ""
	}

	retVal := RemoveTags(game.Data.PCRequirements["minimum"].(string), "\n")
	dxRegex := regexp.MustCompile(`DirectX:(.+)\n`)
	version := dxRegex.FindString(retVal)
	if len(version) != 0 {
		version = strings.Trim(version, "DirectX")
		version = strings.Trim(version, "Version")
		version = strings.TrimSpace(version)

		// Games with "DirectX 10" are using D3D11 API
		if version == "10" {
			version = "11"
		}

		version += formatCitation("This has been extracted from the game's store page using Steam2PCGW and needs to be confirmed.")
	}

	return version
}

func splitChars(r rune) bool {
	return r == '/' || r == '／' || r == ','
}

func fixMemSize(input string) string {
	sizeRegEx := regexp.MustCompile(`(?i)\d[kmg]b`)
	l := sizeRegEx.FindStringIndex(input)
	if l != nil {
		input = input[:l[0]+1] + " " + input[l[0]+1:]
	}
	return input
}

func ProcessSpecs(input string, isMin bool) SysRequirements {
	result := SysRequirements{}

	input = strings.ReplaceAll(input, "only", "")
	input = strings.ReplaceAll(input, " or greater", "")
	input = strings.ReplaceAll(input, " or better", "")
	input = strings.ReplaceAll(input, " or later", "")
	input = strings.ReplaceAll(input, " or higher", "")
	input = strings.ReplaceAll(input, " or lower", "")
	input = strings.ReplaceAll(input, " or equivalent", "")
	input = strings.ReplaceAll(input, "Ghz", "GHz")
	input = strings.ReplaceAll(input, "ghz", "GHz")
	input = strings.ReplaceAll(input, "GHz", " GHz")
	input = strings.ReplaceAll(input, "  GHz", " GHz")
	input = strings.ReplaceAll(input, "®", "")
	input = strings.ReplaceAll(input, "©", "")
	lines := strings.Split(RemoveTags(input, "\n"), "\n")

	for _, line := range lines {
		if len(line) == 0 {
			continue
		}

		defRegEx := regexp.MustCompile(`([a-zA-Z *]+):(.+)$`)
		defs := defRegEx.FindStringSubmatch(line)
		if len(defs) != 3 {
			continue
		}
		param := defs[2]
		switch defs[1] {
		case "OS":
			fallthrough
		case "OS *":
			param = regexp.MustCompile(`(?i)\(?(32|64)[- ]?bits?\)?`).ReplaceAllString(param, "")
			param = regexp.MustCompile(`(?i)microsoft`).ReplaceAllString(param, "")
			param = regexp.MustCompile(`(?i)windows`).ReplaceAllString(param, "")
			param = regexp.MustCompile(`(?i)os x`).ReplaceAllString(param, "")
			param = regexp.MustCompile(`(?i)macos`).ReplaceAllString(param, "")
			param = regexp.MustCompile(`(?i)versions?`).ReplaceAllString(param, "")
			param = regexp.MustCompile(`(?i)\(r\)`).ReplaceAllString(param, "")
			param = regexp.MustCompile(`(?i)\bor `).ReplaceAllString(param, ", ")
			param = strings.ReplaceAll(param, ", ,", ",")
			param = strings.ReplaceAll(param, ",,", ",")
			param = strings.ReplaceAll(param, "/", ",")
			param = strings.ReplaceAll(param, " ,", ",")
			param = strings.ReplaceAll(param, ",", ", ")
			param = strings.ReplaceAll(param, ",  ", ", ")
			param = strings.Trim(param, " .,")
			result.OS = param
		case "Processor":
			param = regexp.MustCompile(`(?i)processor`).ReplaceAllString(param, "")
			param = regexp.MustCompile(`(?i)\(?(32|64)[- ]?bits?\)?`).ReplaceAllString(param, "")
			param = regexp.MustCompile(`(?i)cpu`).ReplaceAllString(param, "")
			param = regexp.MustCompile(`(?i)\bor `).ReplaceAllString(param, ", ")
			param = strings.ReplaceAll(param, ", ,", ",")
			param = strings.ReplaceAll(param, ",,", ",")
			result.CPU = strings.FieldsFunc(param, splitChars)
			for i := 0; i < len(result.CPU); i++ {
				result.CPU[i] = strings.Trim(result.CPU[i], " .,")
			}
		case "Memory":
			param = regexp.MustCompile(`(?i)ram`).ReplaceAllString(param, "")
			param = fixMemSize(param)
			result.RAM = strings.Trim(param, " ")
		case "Graphics":
			fallthrough
		case "Video Card":
			param = regexp.MustCompile(`(?i)\bor `).ReplaceAllString(param, ", ")
			param = regexp.MustCompile(`(?i)\band `).ReplaceAllString(param, ", ")
			param = strings.ReplaceAll(param, ", ,", ",")
			param = strings.ReplaceAll(param, ",,", ",")

			gpus := []string{}
			fields := strings.FieldsFunc(param, splitChars)
			for i := 0; i < len(fields); i++ {
				ogl := regexp.MustCompile(`OpenGL\s*([0-9.]+)`).FindStringSubmatch(fields[i])
				if len(ogl) == 2 {
					result.OGL = ogl[1]
				} else {
					p := strings.ReplaceAll(fields[i], "NVIDIA", "Nvidia")
					p = strings.ReplaceAll(p, "Amd", "AMD")
					p = strings.Trim(p, " .,")
					gpus = append(gpus, p)
				}
			}
			result.GPU = gpus
		case "DirectX":
			param = regexp.MustCompile(`(?i)version`).ReplaceAllString(param, "")
			param = regexp.MustCompile(`(?i)directx`).ReplaceAllString(param, "")
			result.DX = strings.Trim(param, " ")
		case "Storage":
			fallthrough
		case "Hard Drive":
			fallthrough
		case "Hard Disk Space":
			param = regexp.MustCompile(`(?i)available space`).ReplaceAllString(param, "")
			param = regexp.MustCompile(`(?i)hd space`).ReplaceAllString(param, "")
			param = regexp.MustCompile(`(?i)free space`).ReplaceAllString(param, "")
			param = fixMemSize(param)
			result.HD = strings.Trim(param, " ")
		case "Other Requirements":
			result.Other = param
		case "Additional Notes":
			result.Notes = param
		default:
			// fmt.Println("Unprocessed requirement ", defs[1], " - ", param)
		}

	}
	return result
}

func processSysRequirements(input string, specs SysRequirements, level string) string {
	output := input

	mkItem := func(op string) string {
		return "$" + level + "_" + op + "$"
	}

	output = strings.ReplaceAll(output, mkItem("os_versions"), specs.OS)

	if len(specs.CPU) >= 1 {
		output = strings.ReplaceAll(output, mkItem("cpu1"), specs.CPU[0])
	} else {
		output = strings.ReplaceAll(output, mkItem("cpu1"), "")
	}

	if len(specs.CPU) >= 2 {
		output = strings.ReplaceAll(output, mkItem("cpu2"), specs.CPU[1])
	} else {
		output = strings.ReplaceAll(output, mkItem("cpu2"), "")
	}

	output = strings.ReplaceAll(output, mkItem("ram"), specs.RAM)
	output = strings.ReplaceAll(output, mkItem("hd"), specs.HD)
	if len(specs.GPU) >= 1 {
		output = strings.ReplaceAll(output, mkItem("gpu1"), specs.GPU[0])
	} else {
		output = strings.ReplaceAll(output, mkItem("gpu1"), "")
	}

	if len(specs.GPU) >= 2 {
		gpus := specs.GPU[1]
		if len(specs.GPU) >= 3 { // if GPU3 is present insert it right after GPU2
			gpus += "\n|" + level + "GPU3  = " + specs.GPU[2]
		}
		output = strings.ReplaceAll(output, mkItem("gpu2"), gpus)
	} else {
		output = strings.ReplaceAll(output, mkItem("gpu2"), "")
	}

	vram := specs.VRAM
	if len(specs.OGL) != 0 { // if OGL is present insert it after VRAM
		vram += "\n|" + level + "OGL   = " + specs.OGL
	}
	output = strings.ReplaceAll(output, mkItem("vram"), vram)

	dx := specs.DX
	if len(specs.Other) != 0 { // if other is present insert it after DX
		dx += "\n|" + level + "other = " + specs.Other
	}
	output = strings.ReplaceAll(output, mkItem("dx"), dx)

	output = strings.ReplaceAll(output, "$notes$", specs.Notes)

	return output
}

func CleanRecommended(min SysRequirements, rec SysRequirements) SysRequirements {
	if min.OS == rec.OS {
		rec.OS = ""
	}

	if len(min.CPU) == len(rec.CPU) {
		for i := 0; i < len(min.CPU); i++ {
			if min.CPU[i] == rec.CPU[i] {
				rec.CPU[i] = ""
			}
		}
	}

	if min.RAM == rec.RAM {
		rec.RAM = ""
	}

	if min.HD == rec.HD {
		rec.HD = ""
	}

	if len(min.GPU) == len(rec.GPU) {
		for i := 0; i < len(min.GPU); i++ {
			if min.GPU[i] == rec.GPU[i] {
				rec.GPU[i] = ""
			}
		}
	}

	if min.VRAM == rec.VRAM {
		rec.VRAM = ""
	}

	if min.OGL == rec.OGL {
		rec.OGL = ""
	}

	if min.DX == rec.DX {
		rec.DX = ""
	}

	if min.Other == rec.Other {
		rec.Other = ""
	}

	if min.Notes == rec.Notes {
		rec.Notes = ""
	}

	return rec
}

func outputPlaform(isMin bool, isRec bool, minspecs string, recspecs string, osname string) string {
	var output string = system_requirements_template
	var reqsmin SysRequirements
	var reqsrec SysRequirements

	output = strings.ReplaceAll(output, "$os_name$", osname)

	if isMin {
		reqsmin = ProcessSpecs(minspecs, true)
		output = processSysRequirements(output, reqsmin, "min")
	} else {
		output = processSysRequirements(output, SysRequirements{}, "min")
	}

	// Handle recommended specs
	if isRec {
		reqsrec = ProcessSpecs(recspecs, false)
		reqsrec = CleanRecommended(reqsmin, reqsrec)
		output = processSysRequirements(output, reqsrec, "rec")
	} else {
		output = processSysRequirements(output, SysRequirements{}, "rec")
	}

	return output
}

func (game *Game) OutputSpecs() string {
	var output string = ""
	var minStr = ""
	var recStr = ""

	if game.Data.Platforms.Windows {
		isMin := game.Data.PCRequirements != nil && game.Data.PCRequirements["minimum"] != nil
		isRec := game.Data.PCRequirements != nil && game.Data.PCRequirements["recommended"] != nil

		if isMin {
			minStr = game.Data.PCRequirements["minimum"].(string)
		}
		if isRec {
			recStr = game.Data.PCRequirements["recommended"].(string)
		}
		output += outputPlaform(isMin, isRec, minStr, recStr, "Windows")
	}

	if game.Data.Platforms.MAC {
		isMin := game.Data.MACRequirements != nil && game.Data.MACRequirements["minimum"] != nil
		isRec := game.Data.MACRequirements != nil && game.Data.MACRequirements["recommended"] != nil
		if isMin {
			minStr = game.Data.MACRequirements["minimum"].(string)
		}
		if isRec {
			recStr = game.Data.MACRequirements["recommended"].(string)
		}
		output += outputPlaform(isMin, isRec, minStr, recStr, "OS X")
	}

	if game.Data.Platforms.Linux {
		isMin := game.Data.LinuxRequirements != nil && game.Data.LinuxRequirements["minimum"] != nil
		isRec := game.Data.LinuxRequirements != nil && game.Data.LinuxRequirements["recommended"] != nil
		if isMin {
			minStr = game.Data.LinuxRequirements["minimum"].(string)
		}
		if isRec {
			recStr = game.Data.LinuxRequirements["recommended"].(string)
		}
		output += outputPlaform(isMin, isRec, minStr, recStr, "Linux")
	}

	return output
}

func (game *Game) addLanguage(name string, ui, audio, subtitles bool) {

	switch name {
	case "Simplified Chinese":
		name = "Chinese Simplified"
	case "Traditional Chinese":
		name = "Chinese Traditional"
	}

	game.Data.Languages[name] = LanguageData{
		UI:        ui,
		Audio:     audio,
		Subtitles: subtitles,
	}

	// Any one language should at least have subtitles
	if subtitles {
		game.Data.Subtitles = true
	}
}

func (game *Game) ProcessLanguages() {
	var language string

	game.Data.Languages = make(Language)

	input := game.Data.SupportedLanguages
	input = strings.Replace(input, "<br><strong>*</strong>languages with full audio support", "", 1)
	input = strings.ReplaceAll(input, ", ", "\n")
	input = strings.ReplaceAll(input, "<strong>", "")
	input = strings.ReplaceAll(input, "</strong>", "")

	for i := 0; i < len(input); i++ {
		// fmt.Printf("[ProcessLanguages] '%c' char found (language: '%s')\n", input[i], language)
		if rune(input[i]) == '\n' {
			// New line, new language!

			if len(language) != 0 {
				game.addLanguage(language, true, false, true)
				// fmt.Printf("[ProcessLanguages] %s added (\\n found)\n", language)
			}

			language = ""
			continue
		}

		// Found * this means that it has complete support
		if input[i] == '*' {
			game.addLanguage(language, true, true, true)
			// fmt.Printf("[ProcessLanguages] %s added (* found)\n", language)

			language = ""
			continue
		}

		// Append that language string
		language += string(input[i])
	}

	if len(language) != 0 {
		game.addLanguage(language, true, false, true)
	}
}

func IsDate(date string) (bool, []string) {
	dateRe := regexp.MustCompile(`(\d+) (\w+), (\d+)`)
	tokens := dateRe.FindStringSubmatch(date)
	return (len(tokens) != 0), tokens
}

func ParseDate(date string) (output string) {
	success, tokens := IsDate(date)
	if success {
		output = fmt.Sprintf("%s %s %s", tokens[2], tokens[1], tokens[3])
	}
	return output
}

func (game *Game) FormatLanguage(language string) string {
	sanitisedLanguage := language

	switch sanitisedLanguage {
	case "Spanish - Spain":
		sanitisedLanguage = "Spanish"
	case "Spanish - Latin America":
		sanitisedLanguage = "Latin American Spanish"
	case "Portuguese - Brazil":
		sanitisedLanguage = "Brazilian Portuguese"
	case "Chinese Simplified":
		sanitisedLanguage = "Simplified Chinese"
	case "Chinese Traditional":
		sanitisedLanguage = "Traditional Chinese"
	}

	output := language_template
	output = strings.ReplaceAll(output, "$language_name$", sanitisedLanguage)
	output = strings.ReplaceAll(output, "$language_interface$", strconv.FormatBool(game.Data.Languages[language].UI))
	output = strings.ReplaceAll(output, "$language_audio$", strconv.FormatBool(game.Data.Languages[language].Audio))
	output = strings.ReplaceAll(output, "$language_subtitles$", strconv.FormatBool(game.Data.Languages[language].Subtitles))
	return output
}

func SanitiseName(name string, title bool) string {
	name = strings.ReplaceAll(name, "™", "")
	name = strings.ReplaceAll(name, "®", "")
	name = strings.ReplaceAll(name, "©", "")
	name = strings.ReplaceAll(name, ":", "")

	if !title {
		// game titles can have LLC
		name = strings.ReplaceAll(name, " LLC", "")
	}
	return name
}

func (game *Game) HasCategory(category CategoryId) bool {
	for _, v := range game.Data.Categories {
		if CategoryId(v.ID) == category {
			return true
		}
	}
	return false
}

func (game *Game) HasSteamGenre(genre GenreId) bool {
	for _, v := range game.Data.SteamGenres {
		id, _ := strconv.Atoi(v.ID)
		if GenreId(id) == genre {
			return true
		}
	}
	return false
}

func (game *Game) SetFranchise(name string) {
	game.Data.Franchise = name
}

func (game *Game) SetPacing(tags []string) {
	var output string
	pacing := []string{
		"Continuous turn-based",
		"Persistent",
		"Real-time",
		"Relaxed",
		"Turn-based"}
	for _, pace := range pacing {
		for _, tag := range tags {
			if strings.Contains(strings.ToLower(tag), strings.ToLower(pace)) {
				output += pace + ", "
				break
			}
		}
	}

	if len(output) == 0 {
		output += "Real-time"
	} else {
		output = strings.TrimSuffix(output, ", ")
		output = strings.TrimSpace(output)
	}
	game.Data.Pacing = output
}

func (game *Game) SetPerspective(tags []string) {
	var output string
	perspectives := []string{
		"Audio-based",
		"Bird's-eye view",
		"Cinematic camera",
		"First-person",
		"Flip screen",
		"Free-roaming camera",
		"Isometric",
		"Scrolling",
		"Side view",
		"Text-based",
		"Third-person",
		"Top-down view"}
	for _, perspective := range perspectives {
		for _, tag := range tags {
			if strings.Contains(strings.ToLower(tag), strings.ToLower(perspective)) {
				output += perspective + ", "
				break
			}
		}
	}
	output = strings.TrimSuffix(output, ", ")
	output = strings.TrimSpace(output)
	game.Data.Perspectives = output
}

func (game *Game) SetControls(tags []string) {
	var output string
	controls := []string{
		"Direct control",
		"Gestures",
		"Menu-based",
		"Multiple select",
		"Point and select",
		"Text input",
		"Voice control"}
	for _, control := range controls {
		for _, tag := range tags {
			if strings.Contains(strings.ToLower(tag), strings.ToLower(control)) {
				output += control + ", "
				break
			}
		}
	}

	if len(output) == 0 {
		output += "Direct control"
	} else {
		output = strings.TrimSuffix(output, ", ")
		output = strings.TrimSpace(output)
	}
	game.Data.Controls = output
}

func (game *Game) SetGenres(tags []string) {
	var output string
	genres := []string{
		"4X",
		"Action",
		"Adventure",
		"Arcade",
		"ARPG",
		"Artillery",
		"Battle royale",
		"Board",
		"Brawler",
		"Building",
		"Business",
		"Card/tile",
		"CCG",
		"Chess",
		"Clicker",
		"Dating",
		"Driving",
		"Educational",
		"Endless runner",
		"Exploration",
		"Falling block",
		"Fighting",
		"FPS",
		"Gambling/casino",
		"Hack and slash",
		"Hidden object",
		"Hunting",
		"Idle",
		"Immersive sim",
		"Interactive book",
		"JRPG",
		"Life sim",
		"Mental training",
		"Metroidvania",
		"Mini-games",
		"MMO",
		"MMORPG",
		"Music/rhythm",
		"Open world",
		"Paddle",
		"Party game",
		"Pinball",
		"Platform",
		"Puzzle",
		"Quick time events",
		"Racing",
		"Rail shooter",
		"Roguelike",
		"Rolling ball",
		"RPG",
		"RTS",
		"Sandbox",
		"Shooter",
		"Simulation",
		"Sports",
		"Stealth",
		"Strategy",
		"Survival",
		"Survival horror",
		"Tactical RPG",
		"Tactical shooter",
		"TBS",
		"Text adventure",
		"Tile matching",
		"Time management",
		"Tower defense",
		"TPS",
		"Tricks",
		"Trivia/quiz",
		"Vehicle combat",
		"Vehicle simulator",
		"Visual novel",
		"Wargame",
		"Word"}
	for _, genre := range genres {
		for _, tag := range tags {
			if strings.Contains(strings.ToLower(tag), strings.ToLower(genre)) {
				output += genre + ", "
				break
			}
		}
	}
	output = strings.TrimSuffix(output, ", ")
	output = strings.TrimSpace(output)
	game.Data.Genres = output
}

func (game *Game) SetSports(tags []string) {
	var output string
	sports := []string{
		"American football",
		"Australian football",
		"Baseball",
		"Basketball",
		"Bowling",
		"Boxing",
		"Cricket",
		"Darts/target shooting",
		"Dodgeball",
		"Extreme sports",
		"Fictional sport",
		"Fishing",
		"Football (Soccer)",
		"Golf",
		"Handball",
		"Hockey",
		"Horse",
		"Lacrosse",
		"Martial arts",
		"Mixed sports",
		"Paintball",
		"Parachuting",
		"Pool or snooker",
		"Racquetball/squash",
		"Rugby",
		"Sailing/boating",
		"Skateboarding",
		"Skating",
		"Snowboarding or skiing",
		"Surfing",
		"Table tennis",
		"Tennis",
		"Volleyball",
		"Water sports",
		"Wrestling"}
	for _, sport := range sports {
		for _, tag := range tags {
			if strings.Contains(strings.ToLower(tag), strings.ToLower(sport)) {
				output += sport + ", "
				break
			}
		}
	}

	output = strings.TrimSuffix(output, ", ")
	output = strings.TrimSpace(output)
	game.Data.Sports = output
}

func (game *Game) SetVehicles(tags []string) {
	var output string
	vehicles := []string{
		"Automobile",
		"Bicycle",
		"Bus",
		"Flight",
		"Helicopter",
		"Hovercraft",
		"Industrial",
		"Motorcycle",
		"Naval/watercraft",
		"Off-roading",
		"Robot",
		"Self-propelled artillery",
		"Space flight",
		"Street racing",
		"Tank",
		"Track racing",
		"Train",
		"Transport",
		"Truck"}
	for _, vehicle := range vehicles {
		for _, tag := range tags {
			if strings.Contains(strings.ToLower(tag), strings.ToLower(vehicle)) {
				output += vehicle + ", "
				break
			}
		}
	}

	output = strings.TrimSuffix(output, ", ")
	output = strings.TrimSpace(output)
	game.Data.Vehicles = output
}

func (game *Game) SetArtStyles(tags []string) {
	var output string
	artStyles := []string{
		"Abstract",
		"Anime",
		"Cartoon",
		"Cel-shaded",
		"Comic book",
		"Digitized",
		"FMV",
		"Live action",
		"Pixel art",
		"Pre-rendered graphics",
		"Realistic",
		"Stylized",
		"Vector art",
		"Video backdrop",
		"Voxel art"}
	for _, artStyle := range artStyles {
		for _, tag := range tags {
			if strings.Contains(strings.ToLower(tag), strings.ToLower(artStyle)) {
				output += artStyle + ", "
				break
			}
		}
	}

	if len(output) == 0 {
		output = "Realistic"
	} else {
		output = strings.TrimSuffix(output, ", ")
		output = strings.TrimSpace(output)
	}
	game.Data.ArtStyles = output
}

func (game *Game) SetThemes(tags []string) {
	var output string
	themes := []string{
		"Adult",
		"Africa",
		"Amusement park",
		"Antarctica",
		"Arctic",
		"Asia",
		"China",
		"Classical",
		"Cold War",
		"Comedy",
		"Contemporary",
		"Cyberpunk",
		"Dark",
		"Detective/mystery",
		"Eastern Europe",
		"Egypt",
		"Europe",
		"Fantasy",
		"Healthcare",
		"Historical",
		"Horror",
		"Industrial Age",
		"Interwar",
		"Japan",
		"LGBTQ",
		"Lovecraftian",
		"Medieval",
		"Middle East",
		"North America",
		"Oceania",
		"Piracy",
		"Post-apocalyptic",
		"Pre-Columbian Americas",
		"Prehistoric",
		"Renaissance",
		"Romance",
		"Sci-fi",
		"South America",
		"Space",
		"Steampunk",
		"Supernatural",
		"Victorian",
		"Western",
		"World War I",
		"World War II",
		"Zombies"}
	for _, theme := range themes {
		for _, tag := range tags {
			if strings.Contains(strings.ToLower(tag), strings.ToLower(theme)) {
				output += theme + ", "
				break
			}
		}
	}
	output = strings.TrimSuffix(output, ", ")
	game.Data.Themes = output
}
