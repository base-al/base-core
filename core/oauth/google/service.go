package oauth

import (
	"base/s3"
	"context"
	"fmt"
	"strconv"
	"time"

	userService "base/core/users"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2/log"
	"github.com/lib/pq"
	"golang.org/x/oauth2"
	goauth2 "google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
	"gorm.io/gorm"
)

const (
	GoogleOAuthProviderName = string("Google OAuth 2.0.")
	OAuthVersion2           = string("2.0")
	GoogleProviderName      = string("google")
)

type oAuthAccountApi struct {
	db           *gorm.DB
	logger       log.AllLogger
	oauth2Config *oauth2.Config
	s3c          *s3.S3Client
	secretKey    string
	baseHostUrl  string
}

type OAuthAccountAPI interface {
	GoogleSignUp() (res *OAuthResponse, err error)
	GoogleSignUpCallback(req *OAuthCallbackRequest) (res *CallbackResponse, err error)
	GoogleSignIn() (res *OAuthResponse, err error)
	GoogleSignInCallback(req *OAuthCallbackRequest) (res *CallbackResponse, err error)
	UserOAuthAccounts(req *FindUserOAuthAccount) (res *UserAccountResponse, err error)
	GoogleAccountConnect(req *GoogleAccountConnectRequest) (res *OAuthResponse, err error)
	GoogleAccountConnectCallback(req *GoogleAccountConnectCallbackRequest) (res *StatusResponse, err error)
}

func NewOAuthAccountAPI(db *gorm.DB, logger log.AllLogger, oauth2Config *oauth2.Config, s3c *s3.S3Client, secretKey string, baseHostUrl string) OAuthAccountAPI {
	return &oAuthAccountApi{
		db:           db,
		logger:       logger,
		oauth2Config: oauth2Config,
		s3c:          s3c,
		secretKey:    secretKey,
		baseHostUrl:  baseHostUrl,
	}
}

// @Summary      	GoogleSignUp
// @Description		Will return the Google OAuth2.0 redirect URL for sign up.
// @Tags			OAuthAccount
// @Produce			json
// @Success			200								{object}		OAuthResponse
// @Router			/oauth/google/signup			[GET]
func (s *oAuthAccountApi) GoogleSignUp() (res *OAuthResponse, err error) {
	url := s.oauth2Config.AuthCodeURL("", oauth2.AccessTypeOffline, oauth2.SetAuthURLParam("redirect_uri", s.baseHostUrl+"/api/oauth/google/signup/callback"), oauth2.SetAuthURLParam("access_type", "offline"), oauth2.SetAuthURLParam("approval_prompt", "force"))

	return &OAuthResponse{
		RedirectURL: url,
	}, nil
}

// @Summary      	GoogleSignUpCallback
// @Description		This API endpoint is normally automatically called from google redirect of GoogleSignUp endpoint in UI(it requires the code from google to process the process of signing up), and returns the OAuth provider and token.
// @Tags			OAuthAccount
// @Produce			json
// @Param			code							query		string			true	"Code"
// @Success			200								{object}	CallbackResponse
// @Router			/oauth/google/signup/callback	[GET]
func (s oAuthAccountApi) GoogleSignUpCallback(req *OAuthCallbackRequest) (res *CallbackResponse, err error) {
	if req.Code == "" {
		return nil, fmt.Errorf("code is empty")
	}

	s.oauth2Config.RedirectURL = s.baseHostUrl + "/api/oauth/google/signup/callback"
	tok, err := s.oauth2Config.Exchange(context.Background(), req.Code)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve token from web: %v", err)
	}

	ctx := context.Background()
	client := s.oauth2Config.Client(ctx, tok)
	oauth2Service, err := goauth2.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve oauth2 Client %v", err)
	}

	userinfo, err := oauth2Service.Userinfo.Get().Do()
	if err != nil {
		return nil, fmt.Errorf("unable to get userinfo: %v", err)
	}

	// Check if email exists
	var user User
	_ = s.db.Where("email = ?", userinfo.Email).First(&user)
	if user.ID > 0 {
		return nil, fmt.Errorf("email is already used")
	}

	var oAuthAcc OAuthAccount
	_ = s.db.Where("provider = ? AND email = ?", GoogleOAuthProviderName, userinfo.Email)
	if oAuthAcc.ID > 0 {
		return nil, fmt.Errorf("email is already used 2")
	}

	user.Email = userinfo.Email
	user.Username = &userinfo.Email
	user.FirstName = userinfo.GivenName
	user.LastName = userinfo.FamilyName
	user.VerifiedEmail = true
	user.Active = true

	result := s.db.Omit("UpdatedAt").Create(&user)
	if result.Error != nil {
		s.logger.Errorf("func: GoogleSignUpCallback, operation: s.db.Omit('UpdatedAt').Create(&user), err: %s", result.Error.Error())
		return nil, result.Error
	}
	s.db.Model(User{Email: userinfo.Email}).First(&user)

	// Create directory in users.env bucket for this user
	usrsBckName := s.s3c.UsersBucketName()
	err = s.s3c.NewBucketFolder(usrsBckName, fmt.Sprint(user.ID))
	if err != nil {
		s.logger.Errorf("func: GoogleSignUpCallback, operation: s.s3c.NewBucketFolder, err: %s", err.Error())
		return nil, err
	}

	oAuthAcc.UserID = int(user.ID)
	oAuthAcc.Email = userinfo.Email
	oAuthAcc.AccessToken = tok.AccessToken
	oAuthAcc.RefreshToken = tok.RefreshToken
	oAuthAcc.Active = true
	oAuthAcc.Scopes = pq.StringArray([]string{"https://www.googleapis.com/auth/userinfo.profile", "https://www.googleapis.com/auth/userinfo.email"})
	oAuthAcc.Provider = GoogleOAuthProviderName
	oAuthAcc.OAuthVersion = OAuthVersion2
	oAuthAcc.TokenExpiry = tok.Expiry

	result = s.db.Save(&oAuthAcc)
	if result.Error != nil {
		return nil, result.Error
	}

	// Create a new token
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["oid"] = oAuthAcc.ID
	claims["email"] = user.Email
	claims["provider"] = GoogleOAuthProviderName
	claims["exp"] = time.Now().Add(time.Minute * 5).Unix()

	// Sign the token with the secret key
	t, err := token.SignedString([]byte(s.secretKey))
	if err != nil {
		fmt.Println("Error signing token:", err)
		return nil, err
	}

	return &CallbackResponse{
		Provider: GoogleProviderName,
		Token:    t,
	}, nil
}

// @Summary      	GoogleSignIn
// @Description		Will return the Google OAuth2.0 redirect URL for sign in.
// @Tags			OAuthAccount
// @Produce			json
// @Success			200								{object}		OAuthResponse
// @Router			/oauth/google					[GET]
func (s *oAuthAccountApi) GoogleSignIn() (res *OAuthResponse, err error) {
	url := s.oauth2Config.AuthCodeURL("", oauth2.AccessTypeOffline,
		oauth2.SetAuthURLParam("redirect_uri", s.baseHostUrl+"/api/oauth/google/callback"), oauth2.SetAuthURLParam("access_type", "offline"), oauth2.SetAuthURLParam("approval_prompt", "force"))

	return &OAuthResponse{
		RedirectURL: url,
	}, nil
}

// @Summary      	GoogleSignInCallback
// @Description		This API is normally automatically called from google redirect of GoogleSignIn endpoint in UI(it requires the code from google to process the process of signing in), and returns the OAuth provider and token.
// @Tags			OAuthAccount
// @Produce			json
// @Param			code							query		string			true	"Code"
// @Success			200								{object}	CallbackResponse
// @Router			/oauth/google/callback			[GET]
func (s *oAuthAccountApi) GoogleSignInCallback(req *OAuthCallbackRequest) (res *CallbackResponse, err error) {
	if req.Code == "" {
		return nil, fmt.Errorf("code is empty")
	}

	s.oauth2Config.RedirectURL = s.baseHostUrl + "/api/oauth/google/callback"
	tok, err := s.oauth2Config.Exchange(context.Background(), req.Code)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve token from web: %v", err)
	}

	ctx := context.Background()
	client := s.oauth2Config.Client(ctx, tok)
	oauth2Service, err := goauth2.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve oauth2 Client %v", err)
	}

	userinfo, err := oauth2Service.Userinfo.Get().Do()
	if err != nil {
		return nil, fmt.Errorf("unable to get userinfo: %v", err)
	}

	var oAuthAcc OAuthAccount
	result := s.db.Where("provider = ? AND email = ?", GoogleOAuthProviderName, userinfo.Email).First(&oAuthAcc)
	if result.Error != nil {
		return nil, fmt.Errorf("account doesn't exist, login with your email and password or sign up with Google.")
	}

	// Create a new token
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["oid"] = oAuthAcc.ID
	claims["email"] = userinfo.Email
	claims["provider"] = GoogleOAuthProviderName
	claims["exp"] = time.Now().Add(time.Minute * 5).Unix()

	// Sign the token with the secret key
	t, err := token.SignedString([]byte(s.secretKey))
	if err != nil {
		fmt.Println("Error signing token:", err)
		return nil, err
	}

	return &CallbackResponse{
		Provider: GoogleProviderName,
		Token:    t,
	}, nil
}

// @Summary      	UserOAuthAccounts
// @Description		Validates user id, will query DB in oauth accounts and returns records.
// @Tags			OAuthAccount
// @Produce			json
// @Param			Authorization					header		string			true	"Authorization Key(e.g Bearer key)"
// @Success			200								{array}		UserAccountResponse
// @Router			/oauth/accounts					[GET]
func (s *oAuthAccountApi) UserOAuthAccounts(req *FindUserOAuthAccount) (res *UserAccountResponse, err error) {
	if req.UserID == 0 {
		return nil, fmt.Errorf("user id is empty")
	}

	r := &UserAccountResponse{}
	var user userService.User
	result := s.db.Where("id = ?", req.UserID).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	if user.Password != "" {
		r.PasswordSet = true
	}

	scopes := []string{"https://www.googleapis.com/auth/userinfo.profile"}

	oAuthAccs := []*UserOAuthAccount{}
	result = s.db.Model(&OAuthAccount{}).Where("user_id = ? AND ?=ANY(scopes)", req.UserID, scopes).Scan(&oAuthAccs)
	if result.Error != nil {
		return nil, result.Error
	}

	r.OAuthAccounts = oAuthAccs

	return r, nil
}

func (s *oAuthAccountApi) GoogleAccountConnect(req *GoogleAccountConnectRequest) (res *OAuthResponse, err error) {
	state := req.BaseUserID
	redirectUrl := fmt.Sprintf("%s/api/oauth/google/connect/callback", s.baseHostUrl)
	url := s.oauth2Config.AuthCodeURL(fmt.Sprint(state), oauth2.AccessTypeOffline, oauth2.SetAuthURLParam("redirect_uri", redirectUrl), oauth2.SetAuthURLParam("access_type", "offline"), oauth2.SetAuthURLParam("approval_prompt", "force"))

	return &OAuthResponse{
		RedirectURL: url,
	}, nil
}

func (s *oAuthAccountApi) GoogleAccountConnectCallback(req *GoogleAccountConnectCallbackRequest) (res *StatusResponse, err error) {
	if req.Code == "" {
		return nil, fmt.Errorf("code is empty")
	}
	if req.State == "" {
		return nil, fmt.Errorf("state is empty")
	}

	s.oauth2Config.RedirectURL = s.baseHostUrl + "/api/oauth/google/connect/callback"
	tok, err := s.oauth2Config.Exchange(context.Background(), req.Code)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve token from web: %v", err)
	}

	ctx := context.Background()
	client := s.oauth2Config.Client(ctx, tok)
	fmt.Println("client", client)
	oauth2Service, err := goauth2.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve oauth2 Client %v", err)
	}

	userinfo, err := oauth2Service.Userinfo.Get().Do()
	if err != nil {
		return nil, fmt.Errorf("unable to get userinfo: %v", err)
	}

	baseUserId, err := strconv.Atoi(req.State)
	if err != nil {
		return nil, fmt.Errorf("state is not a valid integer: %v", err)
	}

	// Check if email exists
	var user User
	result := s.db.Where("id = ?", baseUserId).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	var oAuthAcc OAuthAccount
	_ = s.db.Where("provider = ? AND email = ?", GoogleOAuthProviderName, userinfo.Email).First(&oAuthAcc)
	if oAuthAcc.ID != 0 {
		for _, sc := range oAuthAcc.Scopes {
			if sc == "https://www.googleapis.com/auth/userinfo.profile" {
				return nil, fmt.Errorf("google email is already used")
			}
		}
		oAuthAcc.Scopes = append(oAuthAcc.Scopes, "https://www.googleapis.com/auth/userinfo.profile")
		oAuthAcc.Scopes = append(oAuthAcc.Scopes, "https://www.googleapis.com/auth/userinfo.email")

		result := s.db.Save(&oAuthAcc)
		if result.Error != nil {
			return nil, result.Error
		}

		return &StatusResponse{
			Status: true,
		}, nil
	}

	oAuthAcc.UserID = int(user.ID)
	oAuthAcc.Email = userinfo.Email
	oAuthAcc.AccessToken = tok.AccessToken
	oAuthAcc.RefreshToken = tok.RefreshToken
	oAuthAcc.Active = true
	oAuthAcc.Scopes = pq.StringArray([]string{"https://www.googleapis.com/auth/userinfo.profile", "https://www.googleapis.com/auth/userinfo.email"})
	oAuthAcc.Provider = GoogleOAuthProviderName
	oAuthAcc.OAuthVersion = OAuthVersion2
	oAuthAcc.TokenExpiry = tok.Expiry

	result = s.db.Save(&oAuthAcc)
	if result.Error != nil {
		return nil, result.Error
	}

	return &StatusResponse{
		Status: true,
	}, nil
}
