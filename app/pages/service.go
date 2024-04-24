package pages

import "time"

type PageService interface {
	Index() ([]Page, error)
	Create(req CreatePageRequest) (CreatePageResponse, error)
	Read(req ReadPageRequest) (ReadPageResponse, error)
	List() ([]Page, error)
	Update(req UpdatePageRequest) (UpdatePageResponse, error)
	Delete(req DeletePageRequest) (DeletePageResponse, error)
}

type SimplePageService struct {
	pages  []Page
	lastID int
}

func NewSimplePageService() *SimplePageService {
	return &SimplePageService{}
}

// @Summary      	Index
// @Description		Lists all pages
// @Tags			Pages
// @Accept			json
// @Produce			json
// @Param			Authorization					header		string			true	"Authorization Key(e.g Bearer key)"
// @Success			200								{array}	ListPageRequest
// @Router			/pages	[GET]
func (s *SimplePageService) Index() ([]Page, error) {
	return s.pages, nil
}

// @Summary      	Create
// @Description		Validates user id and title. If they are up to standard a new page will be created. The created pages ID will be returned.
// @Tags			Pages
// @Accept			json
// @Produce			json
// @Param			Authorization					header		string			true	"Authorization Key(e.g Bearer key)"
// @Param			CreatePageRequest					body		CreatePageRequest	true	"CreatePageRequest"
// @Success			200								{object}	CreatePageResponse
// @Router			/pages	[POST]
func (s *SimplePageService) Create(req CreatePageRequest) (CreatePageResponse, error) {
	s.lastID++
	newPage := Page{
		ID:        s.lastID,
		Title:     req.Title,
		Content:   req.Content,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Status:    "draft",
		AuthorID:  req.AuthorID,
	}
	s.pages = append(s.pages, newPage)
	return CreatePageResponse{ID: newPage.ID, Status: "created"}, nil
}

// @Summary      	Read
// @Description		Fetches a page by ID
// @Tags			Pages
// @Accept			json
// @Produce			json
// @Param			Authorization					header		string			true	"Authorization Key(e.g Bearer key)"
// @Param			id								path		int				true	"Page ID"
// @Success			200								{array}	ReadPageRequest
// @Router			/pages/{id}	[GET]
func (s *SimplePageService) Read(req ReadPageRequest) (ReadPageResponse, error) {
	for _, page := range s.pages {
		if page.ID == req.ID {
			return ReadPageResponse{Page: page}, nil
		}
	}
	return ReadPageResponse{}, nil
}

// @Summary      	Update
// @Description		Updates a page
// @Tags			Pages
// @Accept			json
// @Produce			json
// @Param			Authorization					header		string			true	"Authorization Key(e.g Bearer key)"
// @Param			id								path			int					true	"Page ID"
// @Param			UpdatePageRequest				body			UpdatePageRequest	true	"UpdatePageRequest"
// @Success			200								{object}		UpdatePageResponse
// @Router			/pages/{id}	[PUT]
func (s *SimplePageService) Update(req UpdatePageRequest) (UpdatePageResponse, error) {
	for i, page := range s.pages {
		if page.ID == req.ID {
			s.pages[i] = Page{
				ID:          req.ID,
				Title:       req.Title,
				Content:     req.Content,
				PublishedAt: req.PublishedAt,
				Status:      req.Status,
				AuthorID:    req.AuthorID,
			}
			return UpdatePageResponse{ID: req.ID, Status: "updated"}, nil
		}
	}
	return UpdatePageResponse{}, nil
}

// @Summary      	Delete
// @Description		Deletes a page by ID
// @Tags			Pages
// @Accept			json
// @Produce			json
// @Param			Authorization					header		string			true	"Authorization Key(e.g Bearer key)"
// @Param			id								path		int				true	"Page ID"
// @Success			200								{object}	DeletePageResponse
// @Router			/pages/{id}	[DELETE]

func (s *SimplePageService) Delete(req DeletePageRequest) (DeletePageResponse, error) {
	for i, page := range s.pages {
		if page.ID == req.ID {
			s.pages = append(s.pages[:i], s.pages[i+1:]...)
			return DeletePageResponse{ID: req.ID, Status: "deleted"}, nil
		}
	}
	return DeletePageResponse{}, nil
}
