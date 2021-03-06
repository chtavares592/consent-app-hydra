package handler

import (
	"fmt"
	"net/http"

	"github.com/ory/hydra/sdk/go/hydra/swagger"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo"
	"github.com/ory/hydra/sdk/go/hydra"
)

var store = sessions.NewCookieStore([]byte("jaaiuohtrhua"))

const sessionName = "authentication"

type Worker struct {
	Client hydra.SDK
}

type User struct {
	Name     string
	Password string
}

type Rules struct {
	ID          string   `json:"id"`
	Description string   `json:"description"`
	Subjects    []string `json:"subjects"`
	Effect      string   `json:"effect"`
}

func (w Worker) HandlerConsent(c echo.Context) error {
	consentRequestID := c.QueryParam("consent_challenge")
	if consentRequestID == "" {
		return c.JSON(http.StatusBadRequest, "Consent endpoint was called without a consent request id")

	}
	consentRequest, response, err := w.Client.GetConsentRequest(consentRequestID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "The consent request endpoint does not respond")

	} else if response.StatusCode != http.StatusOK {
		return c.JSON(http.StatusBadRequest, "Consent request endpoint")
	}

	fmt.Println(consentRequest.Subject)

	completeRequest, _, err := w.Client.AcceptConsentRequest(consentRequest.Challenge, swagger.AcceptConsentRequest{
		Session: swagger.ConsentRequestSession{
			AccessToken: map[string]interface{}{
				"test": "testing",
				"user": "claudio",
			},
		},
	})

	if err != nil {
		return c.JSON(http.StatusBadRequest, "The accept consent request endpoint encountered a network error")
	}

	return c.Redirect(http.StatusMovedPermanently, completeRequest.RedirectTo)
}

func (w Worker) HandlerLogin(c echo.Context) error {
	loginChallengeID := c.QueryParam("login_challenge")
	if loginChallengeID == "" {
		return c.JSON(http.StatusBadRequest, "Consent endpoint was called without a consent request id")
	}

	// client := &http.Client{}
	// resp, _ := client.Get("http://localhost:4466/policies/test")

	// buf := new(bytes.Buffer)
	// buf.ReadFrom(resp.Body)

	// rules := &Rules{}
	// err := json.Unmarshal(buf.Bytes(), rules)
	// if err != nil {
	// 	return c.JSON(http.StatusBadRequest, "Error parse json")
	// }

	// var val int64
	// val = -1

	// for i := range rules.Subjects {
	// 	if rules.Subjects[i] == "claudio" {
	// 		val = 0
	// 		break
	// 	}
	// }
	// if val == -1 {
	// 	return c.JSON(http.StatusBadRequest, "no rules for this member")
	// }

	request := c.Request()
	user := authenticated(request)

	if user == "" {
		recv := &User{
			Name:     "claudio",
			Password: "userpassword",
		}
		if recv.Name != "claudio" || recv.Password != "userpassword" {
			return c.JSON(http.StatusBadRequest, "User or Password incorrect")
		}

		request = c.Request()
		response := c.Response()
		session, _ := store.Get(request, sessionName)
		session.Values["user"] = "userid22"

		if err := store.Save(request, response.Writer, session); err != nil {
			return c.JSON(http.StatusBadRequest, "error to save section")

		}
	}

	loginRequest, _, err := w.Client.GetLoginRequest(loginChallengeID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Error get login request")
	}

	completedRequest, _, err := w.Client.AcceptLoginRequest(loginRequest.Challenge, swagger.AcceptLoginRequest{
		Subject:     user,
		RememberFor: 0,
	})

	if err != nil {
		return c.JSON(http.StatusBadRequest, "Error accept login request")
	}

	return c.Redirect(http.StatusMovedPermanently, completedRequest.RedirectTo)
}

func authenticated(r *http.Request) string {
	session, _ := store.Get(r, sessionName)
	if u, ok := session.Values["user"]; !ok {
		return ""
	} else if user, ok := u.(string); !ok {
		return ""
	} else {
		return user
	}
}
