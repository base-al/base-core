package orgs

import (
	"base-core/core/subscriptions"

	"fmt"
	"regexp"
	"strings"

	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

type orgApi struct {
	db     *gorm.DB
	logger log.AllLogger
}

type OrgAPI interface {
	Add(req *OrgRequest) (res *OrgResponse, err error)
	FindMyOrgs(req *FindOrgRequest) (res []*OrgWithRole, err error)
}

func NewOrgAPI(db *gorm.DB, logger log.AllLogger) OrgAPI {
	return &orgApi{
		db:     db,
		logger: logger,
	}
}

func (s orgApi) Add(req *OrgRequest) (res *OrgResponse, err error) {
	if req.UserID == 0 {
		return nil, fmt.Errorf("User ID is empty")
	}
	req.Name = strings.TrimSpace(req.Name)
	req.Size = strings.TrimSpace(req.Size)
	if req.Name == "" {
		return nil, fmt.Errorf("Name is empty")
	}
	if req.Size == "" {
		return nil, fmt.Errorf("Size is empty")
	}

	var org Org
	// Escape special chars and replace spaces with hyphens
	orgSlug := regexp.MustCompile(`[^a-zA-Z0-9 ]+`).ReplaceAllString(strings.ToLower(req.Name), "")
	orgSlug = strings.ReplaceAll(strings.TrimSpace(orgSlug), " ", "-")
	// Check slug duplicate
	s.db.Where("slug = ?", orgSlug).Or("name = ?", req.Name).First(&org)
	if org.ID != 0 {
		return nil, fmt.Errorf("Organization name is already used")
	}

	// Create new features and subscription records
	var trialSubscription subscriptions.Subscription
	trialSubscription.Name = "Trial"
	trialSubDesc := "Platform Trial"
	trialSubscription.Description = &trialSubDesc
	trialSubscription.DurationType = "day"
	trialSubscription.DurationTime = 14
	trialSubscription.Features = append(trialSubscription.Features,
		subscriptions.Feature{
			Key:   "MembersLimit",
			Value: "100",
		})
	result := s.db.Omit("UpdatedAt").Create(&trialSubscription)
	if result.Error != nil {
		return nil, result.Error
	}

	// Save new org record
	org.Name = req.Name
	org.Size = req.Size
	org.Slug = orgSlug
	org.TradeName = req.TradeName
	org.UIN = req.UIN
	org.VAT = req.VAT
	org.SubscriptionID = trialSubscription.ID

	result = s.db.Omit("UpdatedAt").Create(&org)
	if result.Error != nil {
		return nil, result.Error
	}

	// Save relationship
	var usrOrgRole UserOrgRole
	usrOrgRole.OrgID = org.ID
	usrOrgRole.UserID = req.UserID
	usrOrgRole.RoleID = 1 // TODO: handle dynamically
	usrOrgRole.Status = "active"
	result = s.db.Create(&usrOrgRole)
	if result.Error != nil {
		return nil, result.Error
	}

	result = s.db.Save(org)
	if result.Error != nil {
		return nil, result.Error
	}

	// Create new bucket for this new organization
	// err = s.s3c.NewOrgBucket(org.ID)
	// if err != nil {
	// 	return nil, err
	// }
	// bucketOrgName := s.s3c.OrgBucketName(org.ID)
	// // Create new bucket folder for settings
	// err = s.s3c.NewBucketFolder(*bucketOrgName, "settings")
	// if err != nil {
	// 	return nil, err
	// }
	// err = s.s3c.PutOrgSettingsPublicPolicy(org.ID)
	// if err != nil {
	// 	return nil, err
	// }

	return &OrgResponse{
		ID:      org.ID,
		Name:    org.Name,
		OrgSlug: orgSlug,
	}, nil
}

func (s orgApi) FindMyOrgs(req *FindOrgRequest) (res []*OrgWithRole, err error) {
	if req.UserID == 0 {
		return nil, fmt.Errorf("User ID is empty")
	}

	rows, err := s.db.Table("orgs").
		Select("id", "name", "slug", "role_id").
		Joins("LEFT JOIN user_org_roles ON user_org_roles.org_id=orgs.id").
		Where("user_org_roles.user_id = ?", req.UserID).Order("id ASC").Rows()
	var orgRoles []*OrgWithRole
	for rows.Next() {
		onr := &OrgWithRole{}
		err := rows.Scan(&onr.OrgID, &onr.Name, &onr.Slug, &onr.RoleID)
		if err != nil {
			fmt.Println("err", err)
		}
		orgRoles = append(orgRoles, onr)
	}

	return orgRoles, nil
}
