package astorgconv

import "strings"

func sanitizeNameDashed(name string) string {
	name = strings.ToLower(name)
	name = strings.Replace(name, " ", "-", -1)
	return name
}

func sanitizeNameUnderscored(name string) string {
	name = strings.ToLower(name)
	name = strings.Replace(name, " ", "_", -1)
	return name
}

func sanitizeNumber(number string) string {
	number = strings.Replace(number, ".", "", -1)
	return number
}

/*
	switch {
	case person.MAC != "":
		conf.MAC = person.MAC
	case person.Username != "":
		conf.MAC = "no-mac-for-" + sanitizeName(person.Username)
	case person.FullName != "":
		conf.MAC = "no-mac-for-" + sanitizeName(person.FullName)
	default:
		conf.MAC = "no-mac-for-unnamed-" + strconv.Itoa(i)
	}
*/
