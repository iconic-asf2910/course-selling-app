package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/iconic-asf2910/course-selling-app/middleware"
	"github.com/iconic-asf2910/course-selling-app/models"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"golang.org/x/crypto/bcrypt"
)

func AdminSignup(db *mongo.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var admin models.Admin

		if err := json.NewDecoder(r.Body).Decode(&admin); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		var existing models.Admin

		err := db.Collection("admins").
			FindOne(ctx, bson.M{"email": admin.Email}).
			Decode(&existing)

		if err == nil {
			http.Error(w, "Admin already exists", http.StatusConflict)
			return
		}

		if err != mongo.ErrNoDocuments {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		hashed, err := bcrypt.GenerateFromPassword(
			[]byte(admin.Password),
			bcrypt.DefaultCost,
		)

		if err != nil {
			http.Error(w, "Could not hash password", http.StatusInternalServerError)
			return
		}

		admin.Password = string(hashed)
		admin.AdminID = int(time.Now().UnixNano() / int64(time.Millisecond))

		_, err = db.Collection("admins").InsertOne(ctx, admin)

		if err != nil {
			http.Error(w, "Could not create admin", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)

		json.NewEncoder(w).Encode(map[string]string{
			"message": "Admin created successfully",
		})
	}
}

func AdminLogin(db *mongo.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var data struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		var admin models.Admin

		err := db.Collection("admins").
			FindOne(ctx, bson.M{"email": data.Email}).
			Decode(&admin)

		if err != nil {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}

		err = bcrypt.CompareHashAndPassword(
			[]byte(admin.Password),
			[]byte(data.Password),
		)

		if err != nil {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}

		token, err := middleware.GenerateToken(admin.AdminID, "admin")

		if err != nil {
			http.Error(w, "Could not create token", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]string{
			"token": token,
		})
	}
}

func CreateCourse(db *mongo.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		adminID := r.Context().Value(middleware.UserIDKey).(int)

		var course models.Course

		if err := json.NewDecoder(r.Body).Decode(&course); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		course.CourseID = int(time.Now().UnixNano() / int64(time.Millisecond))
		course.InstructorID = adminID

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		_, err := db.Collection("courses").InsertOne(ctx, course)

		if err != nil {
			http.Error(w, "Could not create course", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)

		json.NewEncoder(w).Encode(course)
	}
}

func DeleteCourse(db *mongo.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		courseID := r.PathValue("courseID")

		var id int

		_, err := fmt.Sscanf(courseID, "%d", &id)

		if err != nil {
			http.Error(w, "Invalid course ID", http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		result, err := db.Collection("courses").
			DeleteOne(ctx, bson.M{"courseId": id})

		if err != nil {
			http.Error(w, "Could not delete course", http.StatusInternalServerError)
			return
		}

		if result.DeletedCount == 0 {
			http.Error(w, "Course not found", http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(map[string]string{
			"message": "Course deleted successfully",
		})
	}
}

func AddCourseContent(db *mongo.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		courseID := r.PathValue("courseID")

		var id int

		_, err := fmt.Sscanf(courseID, "%d", &id)

		if err != nil {
			http.Error(w, "Invalid course ID", http.StatusBadRequest)
			return
		}

		var data struct {
			Content string `json:"content"`
		}

		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		result, err := db.Collection("courses").UpdateOne(
			ctx,
			bson.M{"courseId": id},
			bson.M{
				"$set": bson.M{
					"content": data.Content,
				},
			},
		)

		if err != nil {
			http.Error(w, "Could not add content", http.StatusInternalServerError)
			return
		}

		if result.MatchedCount == 0 {
			http.Error(w, "Course not found", http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(map[string]string{
			"message": "Course content added successfully",
		})
	}
}
