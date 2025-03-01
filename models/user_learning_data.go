package models

type Topic struct {
	Topic_id string `bson:"topic_id"`
	Title    string `bson:"title"`
	Duration string `bson:"duration"`
	Locked   bool   `bson:"locked"`
	Done     bool   `bson:"completed"`
}

type Module struct {
	Module_id string  `bson:"module_id"`
	Completed bool    `bson:"completed"`
	Locked    bool    `bson:"locked"`
	Topics    []Topic `bson:"lessons"`
}

type ActiveCourse struct {
	Course_id string   `bson:"course_id"`
	Modules   []Module `bson:"modules"`
}

type SuspendedCourse struct {
	ActiveCourse
}

type SavedCourses struct {
	Course_id string `bson:"course_id"`
}

type User_learning_data struct {
	Uid               string            `bson:"uid"`
	Active_courses    []ActiveCourse    `bson:"active_courses"`
	Saved_courses     []SavedCourses    `bson:"saved_courses"`
	Suspended_courses []SuspendedCourse `bson:"suspended_courses"`
	Communities       []string          `bson:"communities"`
	TimeStamps
}
