package entity

type Category string

const (
	FootballCategory  = "football"
	HistoryCategory   = "history"
	GeographyCategory = "geography"
	ArtCategory       = "art"
	CinemaCategory    = "cinema"
	CultureCategory   = "culture"
	BodyCategory      = "body"
	ExerciseCategory  = "exercise"
	MathCategory      = "math"
)

func (c Category) IsValid() bool {
	switch c {
	case FootballCategory,
		HistoryCategory,
		GeographyCategory,
		ArtCategory,
		CinemaCategory,
		CultureCategory,
		BodyCategory,
		ExerciseCategory,
		MathCategory:
		return true
	default:
		return false
	}
}

func (c Category) GetCategories() []Category {
	return []Category{
		FootballCategory,
		HistoryCategory,
		GeographyCategory,
		ArtCategory,
		CinemaCategory,
		CultureCategory,
		BodyCategory,
		ExerciseCategory,
		MathCategory,
	}
}
