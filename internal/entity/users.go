package entity

type Users struct {
	Id        int `db:"id"`
	UserId    int `db:"user_id"`
	SegmentId int `db:"segment_id"`
}
