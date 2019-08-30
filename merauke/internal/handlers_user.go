package internal

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

// GetUserByID a method to get user given userID params in URL
func (h *Handler) GetUserByID(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	userID, err := strconv.ParseInt(param.ByName("userID"), 10, 64)
	if err != nil {
		log.Printf("[internal][GetUserById] fail to convert user_id into int :%+v", err)
		renderJSON(w, []byte(`
		{
			"message":"Gaboleh nakal :)"
		}
		`), http.StatusBadRequest)
		return
	}
	// TODO: implement this. Query = SELECT * FROM users WHERE id = <userID>
	query := fmt.Sprintf("SELECT id, COALESCE(name,'-') FROM users WHERE id=$1")

	rows, err := h.DB.Query(query, userID)
	if err != nil {
		log.Printf("[internal][GetUserById] fail to select user user_id:%s :%+v\n",
			param.ByName("userID"), err)
		return
	}
	var users []User
	for rows.Next() {
		user := User{}
		err := rows.Scan(&user.ID, &user.Name)
		if err != nil {
			log.Printf("[internal][GetUserById] fail to scan into variable user user_id:%s :%+v\n",
				param.ByName("userID"), err)
			return
		}
		users = append(users, user)
	}
	bytes, err := json.Marshal(users)
	if err != nil {
		log.Printf("[internal][GetUserById] fail to convert array of users into byte :%+v\n",
			err)
		return
	}
	renderJSON(w, bytes, http.StatusOK)

}

// InsertUser a function to insert user data (id, name) to DB
func (h *Handler) InsertUser(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	// TODO: implement this. Query = INSERT INTO users (id, name) VALUES (<userID>, '<name>')
	// read json body
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		renderJSON(w, []byte(`
			message: "Fail to read body"
			`), http.StatusBadRequest)
		return
	}
	// parse json body
	var user User
	err = json.Unmarshal(body, &user)
	if err != nil {
		log.Printf("[internal][InsertUser] fail to convert json into array :%+v\n",
			err)
		return
	}
	// executing insert query
	query := fmt.Sprintf("INSERT INTO users (id,name) VALUES (%d,'%s') ", user.ID, user.Name)
	_, err = h.DB.Query(query)
	if err != nil {
		log.Printf("[internal][InsertUser] fail to create user user_id:%d :%+v\n",
			user.ID, err)
		return
	}
	renderJSON(w, []byte(`
	{
		status:"success",
		message:"Insert User Successfully"
	}
	`), http.StatusOK)

}

// EditUserByID a function to change user data (name) in DB with given params (id, name)
func (h *Handler) EditUserByID(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	// TODO: implement this. Query = UPDATE users SET name = '<name>' WHERE id = <userID>
	// read json body
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		renderJSON(w, []byte(`
			message: "Fail to read body"
			`), http.StatusBadRequest)
		return
	}
	// parse json body
	var user User
	err = json.Unmarshal(body, &user)
	if err != nil {
		log.Printf("[internal][EditUserByID] fail to convert json into array :%+v\n",
			err)
		return
	}
	// executing insert query
	query := fmt.Sprintf("UPDATE users SET name='%s' WHERE id='%s' ", user.Name, param.ByName("userID"))
	_, err = h.DB.Query(query)
	if err != nil {
		log.Printf("[internal][EditUserByID] fail to update user user_id:%s :%+v\n",
			param.ByName("userID"), err)
		return
	}
	renderJSON(w, []byte(`
	{
		status:"success",
		message:"Insert User Successfully"
	}
	`), http.StatusOK)
}

// DeleteUserByID a function to remove user data from DB given the userID
func (h *Handler) DeleteUserByID(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	// TODO: implement this. Query = DELETE FROM users WHERE id = <userID>
	userID := param.ByName("userID")
	query := fmt.Sprintf("DELETE FROM users WHERE id=%s", userID)
	_, err := h.DB.Exec(query)
	if err != nil {
		log.Printf("[internal][DeleteUserByID] fail to delete user user_id:%s :%+v\n",
			userID, err)
		return
	}
	renderJSON(w, []byte(`
	{
		status:"success",
		message:"User Deleted Successfully"
	}
	`), http.StatusOK)
}
