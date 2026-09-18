package models

type Course struct {
	CourseID     int    `json:"courseId" bson:"courseId"`
	Name         string `json:"title" bson:"title"`
	Description  string `json:"description" bson:"description"`
	Price        int    `json:"price" bson:"price"`
	InstructorID int    `json:"instructorId" bson:"instructorId"`
	Content      string `json:"content" bson:"content"`
}

type User struct {
	UserID   int    `json:"userId" bson:"userId"`
	Name     string `json:"name" bson:"name"`
	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
}

type Admin struct {
	AdminID  int    `json:"adminId" bson:"adminId"`
	Name     string `json:"name" bson:"name"`
	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
}

type Purchase struct {
	ID       int `json:"id" bson:"id"`
	UserID   int `json:"userId" bson:"userId"`
	CourseID int `json:"courseId" bson:"courseId"`
}
