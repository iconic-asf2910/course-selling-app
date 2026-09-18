package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"github.com/iconic-asf2910/course-selling-app/database"
	"github.com/iconic-asf2910/course-selling-app/handlers"
	"github.com/iconic-asf2910/course-selling-app/middleware"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Println(".env file not found")
	}

	mongoURI := os.Getenv("MONGO_URI")

	if mongoURI == "" {
		log.Fatal("MONGO_URI is not set")
	}

	db := database.Connect(mongoURI)

	fmt.Println("Database:", db.Name())

	http.HandleFunc("POST /user/signup", handlers.Signup(db))

	http.HandleFunc("POST /user/login", handlers.Login(db))

	http.HandleFunc("GET /courses", handlers.GetCourses(db))

	http.Handle(
		"POST /user/purchase/{courseID}",
		middleware.Auth("user")(
			handlers.PurchaseCourse(db),
		),
	)

	http.Handle(
		"GET /user/purchased-courses",
		middleware.Auth("user")(
			handlers.GetPurchasedCourses(db),
		),
	)

	http.HandleFunc("POST /admin/signup", handlers.AdminSignup(db))

	http.HandleFunc("POST /admin/login", handlers.AdminLogin(db))

	http.Handle(
		"POST /admin/course",
		middleware.Auth("admin")(
			handlers.CreateCourse(db),
		),
	)

	http.Handle(
		"DELETE /admin/course/{courseID}",
		middleware.Auth("admin")(
			handlers.DeleteCourse(db),
		),
	)

	http.Handle(
		"POST /admin/course/{courseID}/content",
		middleware.Auth("admin")(
			handlers.AddCourseContent(db),
		),
	)

	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Course Selling API"))
	})

	fmt.Println("Server running on :8080")

	err = http.ListenAndServe(":8080", middleware.CORS(http.DefaultServeMux))

	if err != nil {
		log.Fatal(err)
	}
}
