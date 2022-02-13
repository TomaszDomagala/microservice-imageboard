package thread

type UserID = string
type CommentID = int
type ThreadID = int
type BoardID int

func (c Comment) addChild(id CommentID) Comment {
	c.Children = append(c.Children, id)
	return c
}

type Comment struct {
	Body     string      `json:"body,omitempty"`
	Author   UserID      `json:"author,omitempty"`
	Id       CommentID   `json:"Id,omitempty"`
	Children []CommentID `json:"children,omitempty"`
}
