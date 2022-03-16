package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"gorm.io/gorm"

	"api/model"
)

func (srv *Server) HandleListJournal(w http.ResponseWriter, r *http.Request) {
	entryDate := r.URL.Query().Get("date")
	uuidPtr, err := ReadUserId(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, errServerProcessingErr)))
		return
	}

	if entryDate == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, errQueryParamDoesntExistsErr)))
		return
	}

	userUuid := *uuidPtr

	tx := srv.DB.TxBegin()
	journal, err := tx.ReadJournalExistsByUserIdAndEntryDate(userUuid, entryDate)
	if err == gorm.ErrRecordNotFound {
		tx.Rollback()

		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(fmt.Sprintf(`{"error": "%v"}`, gorm.ErrRecordNotFound.Error())))
		return
	}

	if err != nil {
		tx.Rollback()

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf(`{"error": "%v"}`, err.Error())))
		return
	}

	tx.Commit()

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`{"entry": "%v"}`, journal.Content)))
}

func (srv *Server) HandleCreateOrUpdateJournal(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	uuidPtr, err := ReadUserId(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, errServerProcessingErr)))
		return
	}

	var form model.JournalForm
	if err := json.Unmarshal(reqBody, &form); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, errFormDecodingErr)))
		return
	}

	if form.Date == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, errFormDecodingErr)))
		return
	}

	tx := srv.DB.TxBegin()
	userUuid := *uuidPtr

	user, _ := tx.ReadUserById(userUuid)
	if user == nil {
		tx.Rollback()

		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, errUserDoesntExistsErr)))
		return
	}

	_, err = tx.ReadJournalExistsByUserIdAndEntryDate(userUuid, form.Date)
	if err == gorm.ErrRecordNotFound {
		createJournal := form.ToCreateModel(userUuid)
		if err = tx.CreateJournal(createJournal); err != nil {
			tx.Rollback()

			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf(`{"error": "%v"}`, err.Error())))
			return
		}

		tx.Commit()

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(fmt.Sprintf(`{"created": "%v"}`, form)))
		return
	}

	if err != nil {
		tx.Rollback()

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf(`{"error": "%v"}`, err.Error())))
		return
	}

	rowsAffected, err := tx.UpdateJournalContentByUserIdAndEntryDate(userUuid, form.Date, form.Content)
	if rowsAffected == 0 || err != nil {
		tx.Rollback()

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf(`{"error": "%v"}`, err.Error())))
		return
	}

	tx.Commit()

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf(`{"created": "%v"}`, form)))
}
