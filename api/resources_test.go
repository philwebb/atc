package api_test

import (
	"errors"
	"io/ioutil"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/concourse/atc"
)

var _ = Describe("Resources API", func() {
	Describe("GET /api/v1/resources", func() {
		var response *http.Response

		JustBeforeEach(func() {
			var err error

			response, err = client.Get(server.URL + "/api/v1/resources")
			Ω(err).ShouldNot(HaveOccurred())
		})

		Context("when getting the resource config succeeds", func() {
			BeforeEach(func() {
				configDB.GetConfigReturns(atc.Config{
					Groups: []atc.GroupConfig{
						{
							Name:      "group-1",
							Resources: []string{"resource-1"},
						},
						{
							Name:      "group-2",
							Resources: []string{"resource-1", "resource-2"},
						},
					},

					Resources: []atc.ResourceConfig{
						{Name: "resource-1", Type: "type-1", Hidden: true},
						{Name: "resource-2", Type: "type-2"},
						{Name: "resource-3", Type: "type-3"},
					},
				}, nil)
			})

			It("returns 200 OK", func() {
				Ω(response.StatusCode).Should(Equal(http.StatusOK))
			})

			It("returns each resource", func() {
				body, err := ioutil.ReadAll(response.Body)
				Ω(err).ShouldNot(HaveOccurred())

				Ω(body).Should(MatchJSON(`[
					{
						"name": "resource-1",
						"type": "type-1",
						"groups": ["group-1", "group-2"],
						"url": "/resources/resource-1",
						"hidden": true
					},
					{
						"name": "resource-2",
						"type": "type-2",
						"groups": ["group-2"],
						"url": "/resources/resource-2",
						"hidden": false
					},
					{
						"name": "resource-3",
						"type": "type-3",
						"groups": [],
						"url": "/resources/resource-3",
						"hidden": false
					}
				]`))
			})
		})

		Context("when getting the resource config fails", func() {
			BeforeEach(func() {
				configDB.GetConfigReturns(atc.Config{}, errors.New("oh no!"))
			})

			It("returns 500", func() {
				Ω(response.StatusCode).Should(Equal(http.StatusInternalServerError))
			})
		})
	})
})