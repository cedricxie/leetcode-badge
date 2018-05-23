package internal

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func fetchLeetcode(name string) (io.Reader, error) {
	resp, err := http.Get(fmt.Sprintf("https://leetcode.com/%s/", name))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(b), nil
}

func analysis(selection *goquery.Selection, data map[string]int) (err error) {
	var solvedQuestion = "Solved Question"
	var acceptedSubmission = "Accepted Submission"
	keys := []string{solvedQuestion, acceptedSubmission}
	for _, key := range keys {
		if strings.Contains(selection.Text(), key) {
			span := selection.Find("span").Text()
			span = strings.Replace(span, " ", "", -1)
			span = strings.Replace(span, "\n", "", -1)

			splited := strings.Split(span, "/")
			if len(splited) != 2 {
				return
			}

			if key == solvedQuestion {
				solved, err := strconv.Atoi(splited[0])
				if err != nil {
					return err
				}

				data["solved_question"] = solved

				all, err := strconv.Atoi(splited[1])
				if err != nil {
					return err
				}
				data["all_question"] = all
			} else if key == acceptedSubmission {
				accepted, err := strconv.Atoi(splited[0])
				if err != nil {
					return err
				}
				data["accepted_submission"] = accepted

				all, err := strconv.Atoi(splited[1])
				if err != nil {
					return err
				}
				data["all_submission"] = all
			}
			return
		}
	}

	return
}

func FetchLeetcodeData(name string) (map[string]int, error) {
	read, err := fetchLeetcode(name)
	if err != nil {
		return nil, err
	}
	doc, err := goquery.NewDocumentFromReader(read)
	if err != nil {
		return nil, err
	}

	var data = make(map[string]int)
	doc.Find(".list-group-item").Each(func(i int, selection *goquery.Selection) {
		if err2 := analysis(selection, data); err2 != nil {
			err = err2
		}
	})

	if err != nil {
		return nil, err
	}
	return data, nil
}