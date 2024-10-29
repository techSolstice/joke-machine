package main

import (
  "encoding/json"
  "fmt"
  "net/http"
  "strings"

  "github.com/gin-gonic/gin"
)

type NameResponse struct {
  FirstName string `json:"first_name"`
  LastName string `json:"last_name"`
}

type Joke struct {
  Id string `json:"id"`
  JokeString string `json:"joke"`
  Categories map[string]interface{} `json:"categories"`
}

type JokeResponse struct {
  Type int `json:"type"`
  Value Joke `json:"value"`
}

func main() {

  // TODO: Support non-STDOUT logging.  Identify logging, metrics, and alerting solutions

  // Setup the router
	router := gin.Default()
  router.Use(ErrorHandler)
  router.Use(gin.Logger())
  router.Use(gin.Recovery())
	router.LoadHTMLGlob("templates/*")

  // Routes
  router.GET("/health", func(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
      "message": "Healthy",
    })
  })

  router.GET("/joke", func(c *gin.Context) {

    var nameResponse NameResponse
    var jokeResponse JokeResponse

    callExternal(c, "https://names.mcquay.me/api/v0/", &nameResponse);
    // TODO: Perhaps settle on a default like "Jane Doe" if the above call fails

    // TODO: Perform some encoding in case names are returned like "Jimmy John Doe" (space in name) or "Jimmy-Jane Doe" (dash in name) or non-ASCII names including Unicode names 
    url := fmt.Sprintf("http://joke.loc8u.com:8888/joke?limitTo=nerdy&firstName=%s&lastName=%s", nameResponse.FirstName, nameResponse.LastName)
    callExternal(c, url, &jokeResponse);
    // TODO: Cache a set of responses which will also help with error handling

		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"quote": jokeResponse.Value.JokeString,
		})
  })


  router.GET("/", func(c *gin.Context) {

    var nameResponse NameResponse
    var jokeResponse JokeResponse

    callExternal(c, "https://names.mcquay.me/api/v0/", &nameResponse);
    // TODO: Perhaps settle on a default like "Jane Doe" if the above call fails

    // TODO: Perform some encoding in case names are returned like "Jimmy John Doe" (space in name) or "Jimmy-Jane Doe" (dash in name) or non-ASCII names including Unicode names 
    url := fmt.Sprintf("http://joke.loc8u.com:8888/joke?limitTo=nerdy&firstName=%s&lastName=%s", nameResponse.FirstName, nameResponse.LastName)
    callExternal(c, url, &jokeResponse);
    // TODO: Cache a set of responses which will also help with error handling

    c.String(http.StatusOK, jokeResponse.Value.JokeString)
  })

  // TODO: Change depending on environment (local, Production)
  router.Run("localhost:5000")
}

// Reduce URL Call boilerplate
// TODO: Perform some assertions and sanity check response for both security and graceful handling
func callExternal(c *gin.Context, url string, target interface{}) error {
  resp, err := http.Get(url)
  if nil != err {
      return c.Error(err)
  }
  defer resp.Body.Close()

  if http.StatusOK != resp.StatusCode {
      return c.Error(fmt.Errorf("Unexpected Error: %s", resp.Status))
  }

  return json.NewDecoder(resp.Body).Decode(target)
}

func ErrorHandler(c *gin.Context) {
  if len(c.Errors) > 0 {
    var errors []string

    c.Next()

    for _, err := range c.Errors {
        errors = append(errors, err.Error())
    }

    serializedErrors := strings.Join(errors, "\n")

    c.JSON(http.StatusInternalServerError, "Joke Machine experienced an unexpected error.\n" + serializedErrors)
  }
}
