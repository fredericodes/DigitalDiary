package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"diarynote/models"
	"diarynote/util/auth"
	"diarynote/util/configs"
)

func (srv *Server) HandleLogin(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	srv.Logger.Info().Println(requestLogInfo(r, reqBody))

	var form models.LoginForm
	if err := json.Unmarshal(reqBody, &form); err != nil {
		srv.Logger.Error().Println(err)

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, errFormDecodingErr)))
		return
	}

	if len(form.Email) == 0 || len(form.Password) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(`{"error": "%v"}`, errInvalidCredentials)))
		return
	}

	userKeyValue, err := srv.RedisDb.ReadUserExists(form.Email)
	if err != nil {
		srv.Logger.Error().Println(fmt.Sprintf(`%v %v`, form.Email, errUserEmailDoesntExistsErr))

		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, errUserEmailDoesntExistsErr)))
		return
	}

	var user models.User
	jsonUser := []byte(userKeyValue)
	if err := json.Unmarshal(jsonUser, &user); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, errServerProcessingErr)))
		return
	}

	match, err := auth.ComparePassword(form.Password, user.Password)
	if !match || err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(fmt.Sprintf(`{"error": "%v"}`, errInvalidCredentials)))
		return
	} else {
		token, err := GenerateToken(user.Id, PermissionTypeUser)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error generating JWT token: " + err.Error()))
			return
		} else {
			w.Header().Set("Authorization", "Bearer "+token)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(fmt.Sprintf(`{"token": "%v"}`, token)))
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}

func (srv *Server) HandleRegister(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	srv.Logger.Info().Println(requestLogInfo(r, reqBody))

	var form models.RegisterForm
	if err := json.Unmarshal(reqBody, &form); err != nil {
		srv.Logger.Error().Println(err)

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, errFormDecodingErr)))
		return
	}

	passwordConfigs := configs.ReadPasswordEnvConfigs()
	hash, err := auth.GeneratePassword(passwordConfigs, form.Password)
	if err != nil {
		srv.Logger.Error().Println(err)

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, errPasswordHashGenFailure)))
		return
	}

	user := form.ToModel(hash)
	jsonData, err := json.Marshal(user)
	if err != nil {
		srv.Logger.Error().Println(err)

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, errServerProcessingErr)))
		return
	}

	_, err = srv.RedisDb.ReadUserExists(form.Email)
	if err == nil {
		srv.Logger.Error().Println(fmt.Sprintf(`%v %v`, form.Email, errUserEmailExistsErr))

		w.WriteHeader(http.StatusConflict)
		w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, errUserEmailExistsErr)))
		return
	}

	if err := srv.RedisDb.CreateUser(user.Email, jsonData); err != nil {
		srv.Logger.Error().Println(err)

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, errServerProcessingErr)))
		return
	}

	srv.Logger.Info().Println(fmt.Sprintf(`"created user": "%v"`, string(jsonData)))
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf(`{"created": "%s"}`, string(jsonData))))
}

func (srv *Server) HandleChangePassword(w http.ResponseWriter, r *http.Request) {

}
