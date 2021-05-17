package main

//
// The tests below should be run inside of docker compose environment:
//
// docker-compose -f integration_test.compose.yaml up -d
// docker-compose -f integration_test.compose.yaml exec test-client go test
// docker-compose -f integration_test.compose.yaml down
//

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/AntonAverchenkov/cards-http-service/internal/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type IntegrationTestSuite struct {
	suite.Suite

	endpointCards        string
	endpointCardsShuffle string
	endpointCardsDeal    string
	endpointCardsReturn  string
}

func TestIntegrationTestSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping docker-dependent test suite in short mode")
	}

	suite.Run(t, new(IntegrationTestSuite))
}

func (suite *IntegrationTestSuite) SetupSuite() {
	/* */ log.Println("IntegrationTestSuite::SetupSuite : begin")
	defer log.Println("IntegrationTestSuite::SetupSuite : end")

	var found bool

	suite.endpointCards, found = os.LookupEnv("ENDPOINT_CARDS")
	require.True(suite.T(), found)

	suite.endpointCardsShuffle, found = os.LookupEnv("ENDPOINT_CARDS_SHUFFLE")
	require.True(suite.T(), found)

	suite.endpointCardsDeal, found = os.LookupEnv("ENDPOINT_CARDS_DEAL")
	require.True(suite.T(), found)

	suite.endpointCardsReturn, found = os.LookupEnv("ENDPOINT_CARDS_RETURN")
	require.True(suite.T(), found)
}

func (suite *IntegrationTestSuite) TestCardsEndpoint() {
	/* */ log.Println("IntegrationTestSuite::TestCardsEndpoint : begin")
	defer log.Println("IntegrationTestSuite::TestCardsEndpoint : end")

	request, err := http.NewRequest("GET", suite.endpointCards, nil)
	require.NoError(suite.T(), err)

	response, err := http.DefaultClient.Do(request)
	require.NoError(suite.T(), err)
	defer response.Body.Close()

	require.Equal(suite.T(), 200, response.StatusCode)
	require.Len(suite.T(), response.Cookies(), 1)
	assert.Equal(suite.T(), sessionCookie, response.Cookies()[0].Name)

	// parse the response
	var cards []api.Card

	reponseBytes, err := ioutil.ReadAll(response.Body)
	require.NoError(suite.T(), err)

	err = json.Unmarshal(reponseBytes, &cards)
	require.NoError(suite.T(), err)

	require.Len(suite.T(), cards, 52)
	assert.Equal(
		suite.T(),
		[]api.Card{
			{Value: "ace", Suit: "clubs"},
			{Value: "two", Suit: "clubs"},
			{Value: "three", Suit: "clubs"},
			{Value: "four", Suit: "clubs"},
			{Value: "five", Suit: "clubs"},
		},
		cards[0:5],
	)
}

func (suite *IntegrationTestSuite) TestCardsShuffleEndpoint() {
	/* */ log.Println("IntegrationTestSuite::TestCardsShuffleEndpoint : begin")
	defer log.Println("IntegrationTestSuite::TestCardsShuffleEndpoint : end")

	request, err := http.NewRequest("POST", suite.endpointCardsShuffle, nil)
	require.NoError(suite.T(), err)

	response, err := http.DefaultClient.Do(request)
	require.NoError(suite.T(), err)
	defer response.Body.Close()

	require.Equal(suite.T(), 200, response.StatusCode)
	require.Len(suite.T(), response.Cookies(), 1)
	assert.Equal(suite.T(), sessionCookie, response.Cookies()[0].Name)

	// parse the response
	var cards []api.Card

	reponseBytes, err := ioutil.ReadAll(response.Body)
	require.NoError(suite.T(), err)

	err = json.Unmarshal(reponseBytes, &cards)
	require.NoError(suite.T(), err)

	require.Len(suite.T(), cards, 52)
	// it is highly unlikely for the random shuffle to result in the following:
	assert.NotEqual(
		suite.T(),
		[]api.Card{
			{Value: "ace", Suit: "clubs"},
			{Value: "two", Suit: "clubs"},
			{Value: "three", Suit: "clubs"},
			{Value: "four", Suit: "clubs"},
			{Value: "five", Suit: "clubs"},
		},
		cards[0:5],
	)
}
