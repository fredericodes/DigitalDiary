package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/go-redis/redis"

	"diarynote/models"
)

func (srv *Server) HandleListDiaryJournals(w http.ResponseWriter, r *http.Request) {
	role, err := ReadUserRole(r)
	if err != nil {
		srv.Logger.Error().Println(err.Error())

		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(fmt.Sprintf(`{"error": "%v"}`, errNotAuthorizedToViewContent)))
	}
	if isAuthorized := IsAuthorized(HandleListDiaryJournals, *role); !isAuthorized {
		srv.Logger.Error().Println(err.Error())

		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(fmt.Sprintf(`{"error": "%v"}`, errNotAuthorizedToViewContent)))
	}

	dateParam := r.URL.Query().Get("date")
	reqBody, _ := ioutil.ReadAll(r.Body)
	srv.Logger.Info().Println(requestLogInfo(r, reqBody))

	uuidPtr, err := ReadUserId(r)
	if err != nil {
		srv.Logger.Error().Println(err)

		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, errServerProcessingErr)))
		return
	}

	uuid := *uuidPtr
	diary, err := srv.RedisDb.ReadDiaryExists(uuid)
	if err == redis.Nil {
		srv.Logger.Error().Println(err)

		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(fmt.Sprintf(`{"error": "%v"}`, err.Error())))
		return
	}

	if dateParam != "" {
		var diaryStruct models.Diary
		err = json.Unmarshal([]byte(diary), &diaryStruct)
		if err != nil {
			srv.Logger.Error().Println(err)

			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, errServerProcessingErr)))
			return
		}

		var journalEntries = diaryStruct.Entries
		entry := journalEntries[dateParam]
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf(`{"entry": "%v"}`, entry)))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(diary))
}

func (srv *Server) HandleCreateOrUpdateDiaryJournals(w http.ResponseWriter, r *http.Request) {
	role, err := ReadUserRole(r)
	if err != nil {
		srv.Logger.Error().Println(err.Error())

		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(fmt.Sprintf(`{"error": "%v"}`, errNotAuthorizedToViewContent)))
	}
	if isAuthorized := IsAuthorized(HandleListDiaryJournals, *role); !isAuthorized {
		srv.Logger.Error().Println(err.Error())

		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(fmt.Sprintf(`{"error": "%v"}`, errNotAuthorizedToViewContent)))
	}

	reqBody, _ := ioutil.ReadAll(r.Body)
	srv.Logger.Info().Println(requestLogInfo(r, reqBody))

	uuidPtr, err := ReadUserId(r)
	if err != nil {
		srv.Logger.Error().Println(err)

		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, errServerProcessingErr)))
		return
	}

	var form models.JournalEntry
	if err := json.Unmarshal(reqBody, &form); err != nil {
		srv.Logger.Error().Println(err)

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, errFormDecodingErr)))
		return
	}

	uuid := *uuidPtr
	diary, err := srv.RedisDb.ReadDiaryExists(uuid)
	if err != redis.Nil { // key doesn't exists
		var diaryStruct models.Diary
		diaryJson := []byte(diary)
		if err := json.Unmarshal(diaryJson, &diaryStruct); err != nil {
			srv.Logger.Error().Println(err)

			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, errServerProcessingErr)))
			return
		}
		entries := diaryStruct.Entries
		entries[form.Date] = form.Content

		diaryStruct.Entries = entries
		diaryJson, err := json.Marshal(&diaryStruct)
		if err != nil {
			srv.Logger.Error().Println(err)

			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, errServerProcessingErr)))
			return
		}

		if err := srv.RedisDb.CreateDiary(uuid, diaryJson); err != nil {
			srv.Logger.Error().Println(err)

			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, errServerProcessingErr)))
			return
		}
	} else {
		var newDiary models.Diary
		entries := make(map[string]string)
		entries[form.Date] = form.Content
		newDiary.ToModel(uuid, entries)
		diaryJson, err := json.Marshal(&newDiary)
		if err != nil {
			srv.Logger.Error().Println(err)

			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, errServerProcessingErr)))
			return
		}

		if err := srv.RedisDb.CreateDiary(uuid, diaryJson); err != nil {
			srv.Logger.Error().Println(err)

			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, errServerProcessingErr)))
			return
		}
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf(`{"created": "%v"}`, form)))
}
