package gomovie

import (
	"strings"
	"net/http"
	"io/ioutil"
	"net/url"
	"errors"
	"log"
	"strconv"
)

func GetEmbedURL(title string) ([]string, error) {
	if title == "" {
		return []string{}, errors.New("The title argument must not be empty")
	}

	domainBlacklist := []string{"videoweed.es"}

	locations := make([]string, 0, 32)

	sites, err := googleSearch("Watch " + title + " Online Putlocker", 3)
	if err != nil {
		return locations, err
	}
	for _, url := range sites {
		domain := DomainFromURL(url)
		if strings.ToLower(domain) == "watchfree.to" {
			locations, err := WatchfreeTo(url)
			if err != nil {
				log.Println("watchfree.to:", err)
				continue
			}
			for k, location := range locations {
				if checkBlacklist(domainBlacklist, location) {
					locations = append(locations[:k], locations[k+1:]...)
					log.Println("Blacklisted Domain", location)
					continue
				}
			}
			if len(locations) <= 0 {
				continue
			}
			locations = append(locations, locations...)
		}
	}

	for _, url := range sites {
		domain := DomainFromURL(url)
		switch strings.ToLower(domain) {
		case "putlocker.is":
			embedURL, err := PutlockerIs(url)
			if err != nil {
				continue
			}
			if checkBlacklist(domainBlacklist, embedURL) {
				log.Println("Blacklisted Domain", embedURL)
				continue
			}
			locations = append(locations, embedURL)
		case "putlockerr.io":
			embedURL, err := PutlockerrIo(url)
			if err != nil {
				continue
			}
			if checkBlacklist(domainBlacklist, embedURL) {
				log.Println("Blacklisted Domain", embedURL)
				continue
			}
			locations = append(locations, embedURL)
		case "putlockerr.co":
			embedURL, err := PutlockerrIo(url)
			if err != nil {
				continue
			}
			if checkBlacklist(domainBlacklist, embedURL) {
				log.Println("Blacklisted Domain", embedURL)
				continue
			}
			locations = append(locations, embedURL)
		}
	}

	if len(locations) <= 0 {
		return locations, errors.New("Unable to find movie")
	}

	return locations, nil
}

func checkBlacklist(blacklist []string, domain string) bool {
	domain = DomainFromURL(domain)
	for _, url := range blacklist {
		if strings.ToLower(url) == strings.ToLower(domain) {
			return true
		}
	}
	return false
}

func GetIMDBTitle(id string) (string, error) {
	resp, err := http.Get("http://www.imdb.com/title/" + id)
	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return StringBetween(strings.ToLower(string(body)), `<h1 itemprop="name" class="">`, `&nbsp;`)
}

func DomainFromURL(url string) string {
	splittedDomain := strings.Split(url, `//`)
	if len(splittedDomain) > 1 {
		splittedDomain = strings.Split(splittedDomain[1], `/`)
	}
	splittedDomain = strings.Split(splittedDomain[0], `.`)
	if len(splittedDomain) < 2 {
		return splittedDomain[0]
	}
	return splittedDomain[len(splittedDomain)-2] + `.` + splittedDomain[len(splittedDomain)-1]
}

// PutlockerIs returns the url of the embedded video in
// the url provided.
func PutlockerIs(url string) (string, error) {
	if url == "" {
		return "", errors.New("The url argument must not be empty")
	}
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	doitSection, err := StringBetween(string(body), `<div class="video">`, `<font color="red">`)
	if err != nil {
		return "", err
	}
	embedURL, err := StringBetween(doitSection, `document.write(doit('`, `'));`)
	if err != nil {
		return "", err
	}
	embedURL, err = StringBetween(strings.ToLower(DecryptPutlocker(embedURL)), `<iframe src="`, `"`)
	return strings.Replace(embedURL, `\`, "", -1), err
}

// PutlockerhdCo returns the url of the embedded video in
// the url provided.
/*func PutlockerhdCo(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return StringBetween(string(body), `<IFRAME SRC="`, `"`)
}*/

// PutlockerIo returns the url of the embedded video in
// the url provided.
func PutlockerrIo(url string) (string, error) {
	if url == "" {
		return "", errors.New("The url argument must not be empty")
	}
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	embedURL, err := StringBetween(strings.ToLower(string(body)), `<iframe src="`, `"`)
	return strings.Replace(embedURL, `\`, "", -1), err
}

// WatchfreeTo returns the url of the embedded video in
// the url provided.
func WatchfreeTo(url string) ([]string, error) {
	if url == "" {
		return []string{}, errors.New("The url argument must not be empty")
	}
	resp, err := http.Get(url)
	if err != nil {
		return []string{}, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []string{}, err
	}
	locations, err := StringBetween(strings.ToLower(string(body)), `var locations = [`, `];`)
	locationSlice := strings.Split(locations, `,`)
	for i, loc := range locationSlice {
		locationSlice[i] = strings.Replace(strings.Replace(loc, `\`, "", -1), `"`, "", -1)
	}
	return locationSlice, nil
}

// googleSearch searches a query to google.com and returns all
// the website url's on the first page in a slice of string.
func googleSearch(query string, pages int) (results []string, err error) {
	results = make([]string, 0, 10)
	query = strings.Replace(query, " ", "%20", -1)

	if pages == 0 {
		pages = 1
	}

	for i := 0; i <= 1; i++ {
		resp, err := http.Get("https://www.google.com.au/search?q="+query+"&start="+strconv.Itoa(i*10))
		if err != nil {
			return results, err
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return results, err
		}
		sites := strings.Split(string(body), `class="g"`)
		for _, site := range sites {
			if strings.Index(site, `<a href="`) != -1 {
				site = site[strings.Index(site, `<a href="`) + len(`<a href="`):]
				if site[:len(`/url?q=`)] == `/url?q=` {
					if strings.Index(site, `">`) != -1 {
						site = site[len(`/url?q=`):strings.Index(site, `">`)]
						if strings.Index(site, `&`) != -1 {
							site = site[:strings.Index(site, `&`)]
							site, err := url.QueryUnescape(site)
							if err != nil {
								return results, err
							}
							results = append(results, site)
						}
					}
				}
			}
		}
	}

	return results, nil
}

// isPutlockerOnline checks weather a putlocker url's response
// and returns a bool representing weather it is accessable.
// An error is returned if a process failed during the process.
func isPutlockerOnline(url string) (bool, error) {
	resp, err := http.Get(url)
	if err != nil {
		return false, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	str, err := StringBetween(string(body), `<title>`, `</title>`)
	if err != nil {
		return false, err
	}

	if strings.Index(strings.ToLower(str), strings.ToLower(`Website is offline`)) != -1 {
		return false, nil
	}

	return true, nil
}

// urlIsPutlock returns a bool showing weather
// the provided URL is located at putlocker;
// no matter the final url extension.
func urlIsPutlocker(url string) bool {
	splitted := strings.Split(url, ".")

	for i, section := range splitted {
		if strings.Index(section, "/") != -1 && i > 0 {
			if strings.ToLower(splitted[i-1]) == "putlocker" {
				return true
			}
		}
	}

	return false
}

// stringBetween returns a substring located between the first occurrence of
// both the provided start and end strings. An error will be returned if
// str does not include both start and end as a substring.
func StringBetween(str, start, end string) (string, error) {
	if strings.Index(str, start) == -1 && strings.Index(str, end) == -1 {
		return str, errors.New("String does not include start/end as substring.")
	}
	str = str[len(start)+strings.Index(str, start):]
	return str[:strings.Index(str, end)], nil
}