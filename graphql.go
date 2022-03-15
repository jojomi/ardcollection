package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/tidwall/gjson"
)

type Endpoint struct {
	Query     string
	Variables map[string]interface{}
}

var queryProgramSet = `query ProgramSetEpisodesQuery($id: ID!, $offset: Int!, $count: Int!) {
	result: programSet(id: $id) {
	  items(
		offset: $offset
		first: $count
		orderBy: PUBLISH_DATE_DESC
		filter: { isPublished: { equalTo: true } }
	  ) {
		pageInfo {
		  hasNextPage
		  endCursor
		}
		nodes {
		  id
		  title
		  publishDate
		  summary
		  duration
		  path
		  image {
			url
			url1X1
			description
			attribution
		  }
		  programSet {
			id
			title
			path
			publicationService {
			  title
			  genre
			  path
			  organizationName
			}
		  }
		  audios {
			url
			downloadUrl
			allowDownload
		  }
		}
	  }
	}
  }`
var queryEditorialCollection = `query EpisodesQuery($id: ID!, $offset: Int!, $limit: Int!) {
	result: editorialCollection(id: $id, offset: $offset, limit: $limit) {
	  items {
		pageInfo {
		  hasNextPage
		}
		nodes {
		  id
		  title
		  publishDate
		  duration
		  path
		  image {
			url
			url1X1
			description
			attribution
		  }
		  programSet {
			id
			title
			publicationService {
			  title
			  organizationName
			}
		  }
		  audios {
			url
			downloadUrl
			allowDownload
		  }
		}
	  }
	}
  }
  `

func getEpisodes(env EnvRoot) ([]Episode, error) {
	endpoints := []Endpoint{
		{
			Query: queryProgramSet,
			Variables: map[string]interface{}{
				"id":     env.ID,
				"offset": 0,
				"count":  env.MaxCount,
			},
		}, {
			Query: queryEditorialCollection,
			Variables: map[string]interface{}{
				"id":     env.ID,
				"offset": 0,
				"limit":  env.MaxCount,
			},
		},
	}

	results := make([]Episode, 0)
	var (
		body string
		err  error
	)
	for _, endpoint := range endpoints {
		body, err = getGraphqlResult(endpoint.Query, endpoint.Variables, env.Verbose)

		if err == nil && gjson.Get(body, "data.result").Raw != "null" {
			break
		}
	}
	if body == "" {
		return results, errors.New("could not retrieve ARD data")
	}

	for i := 0; i < env.MaxCount; i++ {
		title := gjson.Get(body, "data.result.items.nodes."+strconv.Itoa(i)+".title")
		if title.String() == "" {
			continue
		}
		mp3 := gjson.Get(body, "data.result.items.nodes."+strconv.Itoa(i)+".audios.0.downloadUrl")
		if mp3.String() == "" {
			continue
		}

		results = append(results, Episode{
			Title:     title.String(),
			RemoteURL: mp3.String(),
		})
	}

	return results, nil
}

func getGraphqlResult(query string, variables map[string]interface{}, verbose bool) (string, error) {
	q := strings.ReplaceAll(query, "\n", " ")
	re := regexp.MustCompile(`[ \t]+`)
	q = re.ReplaceAllString(q, " ")

	variableJSON, err := json.Marshal(variables)
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("https://api.ardaudiothek.de/graphql?query=%s&variables=%s", url.QueryEscape(q), url.QueryEscape(string(variableJSON)))

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
