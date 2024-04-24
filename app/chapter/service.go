package chapter

import "time"

type ChapterService interface {
	Index() ([]Chapter, error)
	Create(req CreateChapterRequest) (CreateChapterResponse, error)
	Read(req ReadChapterRequest) (ReadChapterResponse, error)
	List() ([]Chapter, error)
	Update(req UpdateChapterRequest) (UpdateChapterResponse, error)
	Delete(req DeleteChapterRequest) (DeleteChapterResponse, error)
}

type SimpleChapterService struct {
	chapters  []Chapter
	lastID int
}

func NewSimpleChapterService() *SimpleChapterService {
	return &SimpleChapterService{}
}

// @Summary      	Index
// @Description		Lists all chapters
// @Tags			Chapters
// @Accept			json
// @Produce			json
// @Param			Authorization					header		string			true	"Authorization Key(e.g Bearer key)"
// @Success			200								{array}	ListChapterRequest
// @Router			/chapters	[GET]
func (s *SimpleChapterService) Index() ([]Chapter, error) {
	return s.chapters, nil
}

// @Summary      	Create
// @Description		Validates user id and title. If they are up to standard a new chapter will be created. The created chapters ID will be returned.
// @Tags			Chapters
// @Accept			json
// @Produce			json
// @Param			Authorization					header		string			true	"Authorization Key(e.g Bearer key)"
// @Param			CreateChapterRequest					body		CreateChapterRequest	true	"CreateChapterRequest"
// @Success			200								{object}	CreateChapterResponse
// @Router			/chapters	[POST]
func (s *SimpleChapterService) Create(req CreateChapterRequest) (CreateChapterResponse, error) {
	s.lastID++
	newChapter := Chapter{
		ID:        s.lastID,
		Title:     req.Title,
		Content:   req.Content,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Status:    "draft",
		AuthorID:  req.AuthorID,
	}
	s.chapters = append(s.chapters, newChapter)
	return CreateChapterResponse{ID: newChapter.ID, Status: "created"}, nil
}

// @Summary      	Read
// @Description		Fetches a chapter by ID
// @Tags			Chapters
// @Accept			json
// @Produce			json
// @Param			Authorization					header		string			true	"Authorization Key(e.g Bearer key)"
// @Param			id								path		int				true	"Chapter ID"
// @Success			200								{array}	ReadChapterRequest
// @Router			/chapters/{id}	[GET]
func (s *SimpleChapterService) Read(req ReadChapterRequest) (ReadChapterResponse, error) {
	for _, chapter := range s.chapters {
		if chapter.ID == req.ID {
			return ReadChapterResponse{ Chapter: chapter }, nil
		}
	}
	return ReadChapterResponse{}, nil
}

// @Summary      	Update
// @Description		Updates a chapter
// @Tags			Chapters
// @Accept			json
// @Produce			json
// @Param			Authorization					header		string			true	"Authorization Key(e.g Bearer key)"
// @Param			id								path			int					true	"Chapter ID"
// @Param			UpdateChapterRequest				body			UpdateChapterRequest	true	"UpdateChapterRequest"
// @Success			200								{object}		UpdateChapterResponse
// @Router			/chapters/{id}	[PUT]
func (s *SimpleChapterService) Update(req UpdateChapterRequest) (UpdateChapterResponse, error) {
	for i, chapter := range s.chapters {
		if chapter.ID == req.ID {
			s.chapters[i] = Chapter{
				ID:          req.ID,
				Title:       req.Title,
				Content:     req.Content,
				PublishedAt: req.PublishedAt,
				Status:      req.Status,
				AuthorID:    req.AuthorID,
			}
			return UpdateChapterResponse{ID: req.ID, Status: "updated"}, nil
		}
	}
	return UpdateChapterResponse{}, nil
}

// @Summary      	Delete
// @Description		Deletes a chapter by ID
// @Tags			Chapters
// @Accept			json
// @Produce			json
// @Param			Authorization					header		string			true	"Authorization Key(e.g Bearer key)"
// @Param			id								path		int				true	"Chapter ID"
// @Success			200								{object}	DeleteChapterResponse
// @Router			/chapters/{id}	[DELETE]
func (s *SimpleChapterService) Delete(req DeleteChapterRequest) (DeleteChapterResponse, error) {
	for i, chapter := range s.chapters {
		if chapter.ID == req.ID {
			s.chapters = append(s.chapters[:i], s.chapters[i+1:]...)
			return DeleteChapterResponse{ID: req.ID, Status: "deleted"}, nil
		}
	}
	return DeleteChapterResponse{}, nil
}
