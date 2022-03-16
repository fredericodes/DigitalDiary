package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"api/model"
	"api/util/auth"
	"api/util/configs"
)

func (srv *Server) HandleLogin(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)

	var form model.LoginForm
	if err := json.Unmarshal(reqBody, &form); err != nil {

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, errFormDecodingErr)))
		return
	}

	if len(form.Username) == 0 || len(form.Password) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(`{"error": "%v"}`, errInvalidCredentials)))
		return
	}

	tx := srv.DB.TxBegin()
	user, _ := tx.ReadUserByUsername(form.Username)
	if user == nil {
		tx.Rollback()

		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, errUserDoesntExistsErr)))
		return
	}

	match, err := auth.ComparePassword(form.Password, user.Password)
	if !match || err != nil {
		tx.Rollback()

		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(fmt.Sprintf(`{"error": "%v"}`, errInvalidCredentials)))
		return
	}

	token, err := GenerateToken(user.Id)
	if err != nil {
		tx.Rollback()

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error generating JWT token: " + err.Error()))
		return
	}

	tx.Commit()

	w.Header().Set("Authorization", "Bearer "+token)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`{"token": "%v"}`, token)))
	return
}

func (srv *Server) HandleRegister(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)

	var form model.RegisterForm
	if err := json.Unmarshal(reqBody, &form); err != nil {

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, errFormDecodingErr)))
		return
	}

	passwordConfigs := configs.ReadPasswordEnvConfigs()
	hash, err := auth.GeneratePassword(passwordConfigs, form.Password)
	if err != nil {

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, errPasswordHashGenFailure)))
		return
	}

	tx := srv.DB.TxBegin()

	user := form.ToModel(hash)
	userFromDb, _ := tx.ReadUserByUsername(form.Username)
	if userFromDb != nil {
		tx.Rollback()

		w.WriteHeader(http.StatusConflict)
		w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, errUserEmailExistsErr)))
		return
	}

	if err := tx.CreateUser(user); err != nil {
		tx.Rollback()

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, errServerProcessingErr)))
		return
	}

	tx.Commit()

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf(`{"created": "%s"}`, user.Username)))
}
