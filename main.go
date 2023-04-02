package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type Course struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Weight 	   int    `json:"weight"`
}

type UserCourse struct {
	UserID    int     `json:"user_id"`
	CourseID  int     `json:"course_id"`
	Grade 	  float64 `json:"grade"`
}

type CourseWithGrade struct {
	ID         int     `json:"id"`
	Name       string  `json:"name"`
	Weight     int     `json:"weight"`
	Grade 	   float64 `json:"grade"`
}


var db *sql.DB

func main() {

	err:= godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	
	db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	router := mux.NewRouter()
	createTableifNotExists()
	router.HandleFunc("/users", getUsers).Methods("GET")
	router.HandleFunc("/users", createUser).Methods("POST")
	router.HandleFunc("/users/{id}", getUserCourses).Methods("GET")
	router.HandleFunc("/users/{id}", updateUser).Methods("PUT")
	router.HandleFunc("/users/{id}", deleteUser).Methods("DELETE")
	router.HandleFunc("/users/{id}/courses", createUserCourses).Methods("POST")
	router.HandleFunc("/users/{id}/gradeaverage", getUserGradeAverage).Methods("GET")

	router.HandleFunc("/courses", getCourses).Methods("GET")
	router.HandleFunc("/courses", createCourse).Methods("POST")
	router.HandleFunc("/courses/{id}", updateCourse).Methods("PUT")
	router.HandleFunc("/courses/{id}", deleteCourse).Methods("DELETE")

	log.Fatal(http.ListenAndServe("0.0.0.0:8080", router))
}

func createTableifNotExists() {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS users (id SERIAL PRIMARY KEY, username VARCHAR(255), email VARCHAR(255))")
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS courses (id SERIAL PRIMARY KEY, name VARCHAR(255), weight INTEGER)")
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS user_courses (user_id INTEGER REFERENCES users(id), course_id INTEGER REFERENCES courses(id), grade FLOAT)")
	if err != nil {
		log.Fatal(err)
	}
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, username, email FROM users")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var users []User

	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Username, &user.Email)
		if err != nil {
			log.Fatal(err)
		}
		users = append(users, user)
	}

	err = json.NewEncoder(w).Encode(users)
	if err != nil {
		log.Fatal(err)
	}
}

func createUser(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := db.Prepare("INSERT INTO users(username, email) VALUES($1, $2) RETURNING id")
	if err != nil {
		log.Fatal(err)
	}

	err = stmt.QueryRow(user.Username, user.Email).Scan(&user.ID)
	if err != nil {
		log.Fatal(err)
	}

	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		log.Fatal(err)
	}
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Fatal(err)
	}

	var user User
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := db.Prepare("UPDATE users SET username = $1, email = $2 WHERE id = $3")
	if err != nil {
		log.Fatal(err)
	}

	_, err = stmt.Exec(user.Username, user.Email, id)
	if err != nil {
		log.Fatal(err)
	}

	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		log.Fatal(err)
	}
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := db.Prepare("DELETE FROM users WHERE id = $1")
	if err != nil {
		log.Fatal(err)
	}

	_, err = stmt.Exec(id)
	if err != nil {
		log.Fatal(err)
	}
}

func getCourses(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, name, weight FROM courses")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var courses []Course

	for rows.Next() {
		var course Course
		err := rows.Scan(&course.ID, &course.Name, &course.Weight)
		if err != nil {
			log.Fatal(err)
		}
		courses = append(courses, course)
	}

	err = json.NewEncoder(w).Encode(courses)
	if err != nil {
		log.Fatal(err)
	}
}

func createCourse(w http.ResponseWriter, r *http.Request) {
	var course Course
	err := json.NewDecoder(r.Body).Decode(&course)
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := db.Prepare("INSERT INTO courses(name, weight) VALUES($1, $2) RETURNING id")
	if err != nil {
		log.Fatal(err)
	}

	err = stmt.QueryRow(course.Name, course.Weight).Scan(&course.ID)
	if err != nil {
		log.Fatal(err)
	}

	err = json.NewEncoder(w).Encode(course)
	if err != nil {
		log.Fatal(err)
	}
}

func updateCourse(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Fatal(err)
	}

	var course Course
	err = json.NewDecoder(r.Body).Decode(&course)
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := db.Prepare("UPDATE courses SET name = $1, weight = $2 WHERE id = $3")
	if err != nil {
		log.Fatal(err)
	}

	_, err = stmt.Exec(course.Name, course.Weight, id)
	if err != nil {
		log.Fatal(err)
	}

	err = json.NewEncoder(w).Encode(course)
	if err != nil {
		log.Fatal(err)
	}
}

func deleteCourse(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := db.Prepare("DELETE FROM courses WHERE id = $1")
	if err != nil {
		log.Fatal(err)
	}

	_, err = stmt.Exec(id)
	if err != nil {
		log.Fatal(err)
	}
}

func getUserCourses(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Fatal(err)
	}

	rows, err := db.Query("SELECT c.id, c.name, c.weight, uc.grade FROM courses c INNER JOIN user_courses uc ON c.id = uc.course_id WHERE uc.user_id = $1; ", id)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var courseWithGrade []CourseWithGrade

	for rows.Next() {
		var course Course
		var user_course UserCourse
		err := rows.Scan(&course.ID, &course.Name, &course.Weight, &user_course.Grade)
		if err != nil {
			log.Fatal(err)
		}
		courseWithGrade = append(courseWithGrade, CourseWithGrade{ID: course.ID, Name: course.Name, Weight: course.Weight, Grade: user_course.Grade})
	}

	err = json.NewEncoder(w).Encode(courseWithGrade)
	if err != nil {
		log.Fatal(err)
	}
}

func createUserCourses(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Fatal(err)
	}

	var user_courses UserCourse
	err = json.NewDecoder(r.Body).Decode(&user_courses)
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := db.Prepare("INSERT INTO user_courses(user_id, course_id, grade) VALUES($1, $2, $3)")
	if err != nil {
		log.Fatal(err)
	}

	_, err = stmt.Exec(id, user_courses.CourseID, user_courses.Grade)
	if err != nil {
		log.Fatal(err)
	}

	err = json.NewEncoder(w).Encode(user_courses)
	if err != nil {
		log.Fatal(err)
	}
}

func getUserGradeAverage(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Fatal(err)
	}

	rows, err := db.Query("SELECT SUM(c.weight * uc.grade) / SUM(c.weight) FROM courses c INNER JOIN user_courses uc ON c.id = uc.course_id WHERE uc.user_id = $1; ", id)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var grade float64

	for rows.Next() {
		err := rows.Scan(&grade)
		if err != nil {
			log.Fatal(err)
		}
	}

	err = json.NewEncoder(w).Encode(grade)
	if err != nil {
		log.Fatal(err)
	}
}