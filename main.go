package main

import (
	"fmt"
	"github.com/alexedwards/scs/engine/memstore"
	"github.com/alexedwards/scs/session"
	"github.com/casbin/casbin"
	"github.com/zupzup/casbin-http-role-example/authorization"
	"github.com/zupzup/casbin-http-role-example/model"
	"log"
	"net/http"
	"time"
)

func main() {
	// setup casbin auth rules
	authEnforcer, err := casbin.NewEnforcerSafe("./auth_model.conf", "./policy.csv")
	if err != nil {
		log.Fatal(err)
	}

	// setup session store
	engine := memstore.New(30 * time.Minute)
	sessionManager := session.Manage(engine, session.IdleTimeout(30*time.Minute), session.Persist(true), session.Secure(true))
	users := createUsers()

	// setup routes
	mux := http.NewServeMux()
	mux.HandleFunc("/login", loginHandler(users))
	mux.HandleFunc("/logout", logoutHandler())
	mux.HandleFunc("/member/current", currentMemberHandler())
	mux.HandleFunc("/member/role", memberRoleHandler())
	mux.HandleFunc("/admin/stuff", adminHandler())

	log.Print("Server started on localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", sessionManager(authorization.Authorizer(authEnforcer, users)(mux))))

}

func createUsers() model.Users {
	users := model.Users{}
	users = append(users, model.User{ID: 1, Name: "Admin", Role: "admin"})
	users = append(users, model.User{ID: 2, Name: "Sabine", Role: "member"})
	users = append(users, model.User{ID: 3, Name: "Sepp", Role: "member"})
	return users
}

func loginHandler(users model.Users) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		name := r.PostFormValue("name")
		user, err := users.FindByName(name)
		if err != nil {
			writeError(http.StatusBadRequest, "WRONG_CREDENTIALS", w, err)
			return
		}
		// setup ession
		if err := session.RegenerateToken(r); err != nil {
			writeError(http.StatusInternalServerError, "ERROR", w, err)
			return
		}
		session.PutInt(r, "userID", user.ID)
		session.PutString(r, "role", user.Role)
		writeSuccess("SUCCESS", w)
	})
}

func logoutHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := session.Renew(r); err != nil {
			writeError(http.StatusInternalServerError, "ERROR", w, err)
			return
		}
		writeSuccess("SUCCESS", w)
	})
}

func currentMemberHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uid, err := session.GetInt(r, "userID")
		if err != nil {
			writeError(http.StatusInternalServerError, "ERROR", w, err)
			return
		}
		writeSuccess(fmt.Sprintf("User with ID: %d", uid), w)
	})
}

func memberRoleHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role, err := session.GetString(r, "role")
		if err != nil {
			writeError(http.StatusInternalServerError, "ERROR", w, err)
			return
		}
		writeSuccess(fmt.Sprintf("User with Role: %s", role), w)
	})
}

func adminHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeSuccess("I'm an Admin!", w)
	})
}

func writeError(status int, message string, w http.ResponseWriter, err error) {
	log.Print("ERROR: ", err.Error())
	w.WriteHeader(status)
	w.Write([]byte(message))
}

func writeSuccess(message string, w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(message))
}
