package tryout

type TryOut struct {
	ID                string     `gorm:"type:text;primaryKey"`
	CreatedAt         int64      `gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt         int64      `gorm:"column:updated_at;autoCreateTime:milli;autoUpdateTime:milli"`
	TryOutName        string     `gorm:"column:tryout_name"`
	TestType          string     `gorm:"column:tryout_type"`
	TryOutQuota       int        `gorm:"column:tryout_quota"`
	TryoutPrice       int        `json:"tryout_price" form:"tryout_price"`
	DueDate           int64      `json:"due_date" form:"tryout_price"`
	QuestionExercises []Question `gorm:"foreignKey:tryout_id"`
	Version           int64      `gorm:"column:version"`
}

func (u *TryOut) TableName() string {
	return "try_outs"
}

type Question struct {
	ID              string   `gorm:"type:text;primaryKey"`
	CreatedAt       int64    `gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt       int64    `gorm:"column:updated_at;autoCreateTime:milli;autoUpdateTime:milli"`
	TryoutID        string   `gorm:"column:tryout_id"`
	Content         string   `gorm:"column:content"`
	Weight          int      `gorm:"column:weight"`
	ChoiceExercises []Choice `gorm:"foreignKey:question_id"`
	version         int64    `gorm:"column:version"`
}

func (u *Question) TableName() string {
	return "questions"
}

type Choice struct {
	ID         string `gorm:"type:text;primaryKey"`
	CreatedAt  int64  `gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt  int64  `gorm:"column:updated_at;autoCreateTime:milli;autoUpdateTime:milli"`
	QuestionID string `gorm:"column:question_id"`
	Content    string `gorm:"column:content"`
	IsCorrect  bool   `gorm:"column:is_correct"`
	Weight     int    `gorm:"column:weight"`
	version    int64  `gorm:"column:version"`
}

func (u *Choice) TableName() string {
	return "choices"
}
