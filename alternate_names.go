package geonames

import "strconv"

const alternateNamesURL = `alternateNames.zip`

type AlternateName struct {
	Id              int    // alternateNameId   : the id of this alternate name, int
	GeonameId       int    // geonameid         : geonameId referring to id in table 'geoname', int
	IsoLanguage     string // isolanguage       : iso 639 language code 2- or 3-characters; 4-characters 'post' for postal codes and 'iata','icao' and faac for airport codes, fr_1793 for French Revolution names,  abbr for abbreviation, link for a website, varchar(7)
	Name            string // alternate name    : alternate name or name variant, varchar(200)
	IsPreferredName bool   // isPreferredName   : '1', if this alternate name is an official/preferred name
	IsShortName     bool   // isShortName       : '1', if this is a short name like 'California' for 'State of California'
	IsColloquial    bool   // isColloquial      : '1', if this alternate name is a colloquial or slang term
	IsHistoric      bool   // isHistoric        : '1', if this alternate name is historic and was used in the past
}

func AlternateNames() ([]*AlternateName, error) {
	var err error
	var result []*AlternateName

	zipped, err := httpGet(geonamesURL + alternateNamesURL)
	if err != nil {
		return nil, err
	}

	files, err := unzip(zipped)
	if err != nil {
		return nil, err
	}

	data, err := getZipData(files, "alternateNames.txt")
	if err != nil {
		return nil, err
	}

	parse(data, 0, func(raw [][]byte) bool {
		if len(raw) != 8 {
			return true
		}

		if string(raw[2]) == "link" {
			return true
		}

		id, _ := strconv.Atoi(string(raw[0]))
		geonameId, _ := strconv.Atoi(string(raw[1]))
		boolTrue := "1"

		result = append(result, &AlternateName{
			Id:              id,
			GeonameId:       geonameId,
			IsoLanguage:     string(raw[2]),
			Name:            string(raw[3]),
			IsPreferredName: string(raw[4]) == boolTrue,
			IsShortName:     string(raw[5]) == boolTrue,
			IsColloquial:    string(raw[6]) == boolTrue,
			IsHistoric:      string(raw[7]) == boolTrue,
		})

		return true
	})

	return result, nil
}