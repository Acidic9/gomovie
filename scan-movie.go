package gomovie

import (
	"strings"
	"net/http"
	"io/ioutil"
	"net/url"
	"errors"
	"log"
	"fmt"
)

func GetEmbedURL(title string) (string, error) {
	sites, err := googleSearch("Watch " + title + " Online")
	if err != nil {
		return "", err
	}

	for _, url := range sites {
		splittedDomain := strings.Split(url, `//`)
		if len(splittedDomain) > 1 {
			splittedDomain = strings.Split(splittedDomain[1], `/`)
		}
		splittedDomain = strings.Split(splittedDomain[0], `.`)
		domain := splittedDomain[len(splittedDomain)-2] + `.` + splittedDomain[len(splittedDomain)-1]
		if strings.ToLower(domain) == "watchfree.to" {
			embedURL, err := WatchfreeTo(url)
			if err != nil {
				log.Println("watchfree.to:", err)
				continue
			}
			return embedURL, nil
		}
	}

	for _, url := range sites {
		splittedDomain := strings.Split(url, `//`)
		if len(splittedDomain) > 1 {
			splittedDomain = strings.Split(splittedDomain[1], `/`)
		}
		fmt.Println(splittedDomain)
		splittedDomain = strings.Split(splittedDomain[0], `.`)
		fmt.Println(splittedDomain)
		if len(splittedDomain) < 2 {
			log.Println("Domain error", splittedDomain)
			return "", errors.New("Domain error")
		}
		domain := splittedDomain[len(splittedDomain)-2] + `.` + splittedDomain[len(splittedDomain)-1]
		switch strings.ToLower(domain) {
		case "putlocker.is":
			embedURL, err := PutlockerIs(url)
			if err != nil {
				log.Println("putlocker.is:", err)
				continue
			}
			return embedURL, nil
		case "putlockerr.io":
			embedURL, err := PutlockerrIo(url)
			if err != nil {
				log.Println("putlockerr.io:", err)
				continue
			}
			return embedURL, nil
		case "putlockerr.co":
			embedURL, err := PutlockerrIo(url)
			if err != nil {
				log.Println("putlockerr.co:", err)
				continue
			}
			return embedURL, nil
		}
	}

	return "", errors.New("Unable to find movie")
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

// PutlockerIs returns the url of the embedded video in
// the url provided.
func PutlockerIs(url string) (string, error) {
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
		panic(err)
		return "", err
	}

	fmt.Printf("!%#v", doitSection)

	embedURL, err := StringBetween(doitSection, `document.write(doit('`, `'));`)
	if err != nil {
		panic(err)
		return "", err
	}

	fmt.Println(embedURL)
	return DecryptPutlocker(embedURL), nil
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
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return StringBetween(strings.ToLower(string(body)), `<iframe src="`, `"`)
}

// WatchfreeTo returns the url of the embedded video in
// the url provided.
func WatchfreeTo(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return StringBetween(strings.ToLower(string(body)), `var locations = ["`, `"`)
}

// googleSearch searches a query to google.com and returns all
// the website url's on the first page in a slice of string.
func googleSearch(query string) (results []string, err error) {
	results = make([]string, 0, 10)
	query = strings.Replace(query, " ", "%20", -1)
	resp, err := http.Get("https://www.google.com.au/search?q="+query)
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
	if strings.Index(str, start) == -1 || strings.Index(str, end) == -1 {
		return "", errors.New("String does not include start/end as substring.")
	}
	str = str[len(start)+strings.Index(str, start):]
	return str[:strings.Index(str, end)], nil
}