package handlers

// import (
// 	"net/http"
// 	// "user-profiles/internal/http/req"
// 	// "user-profiles/internal/http/resp"
// 	//"html/template"
// 	"user-profiles/configs"
// 	//"user-profiles/internal/http/middleware"

// 	"user-profiles/internal/auth"
// 	"user-profiles/internal/users"

// 	"github.com/go-chi/chi/v5"
// )

// type AdminHandler struct {
// 	Config       *configs.Config
// 	UsersService users.AdminService
// 	AuthService  auth.AuthService
// 	JWTService   auth.JWTService
// }

// type AdminHandlerDeps struct {
// 	Config       *configs.Config
// 	AdminService users.AdminService
// 	AuthService  auth.AuthService
// 	JWTService   auth.JWTService
// }

// func NewAdminHandler(router chi.Router, deps AdminHandlerDeps) *AdminHandler {
// 	return &AdminHandler{
// 		Config:       deps.Config,
// 		AdminService: deps.AdminService,
// 		AuthService:  deps.AuthService,
// 		JWTService:   deps.JWTService,
// 	}
// }

// func (handler *AdminHandler) AdminProfilePage(w http.ResponseWriter, r *http.Request) {
// 	// err := handler.Tmpl.ExecuteTemplate(w, "profile", nil)
// 	// if err != nil {
// 	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
// 	// }
// }

// func (handler *AdminHandler) UserProfilePage(w http.ResponseWriter, r *http.Request) {
// 	// err := handler.Tmpl.ExecuteTemplate(w, "profile", nil)
// 	// if err != nil {
// 	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
// 	// }
// }

// func (handler *Handler) ViewById(w http.ResponseWriter, r *http.Request) {
// 	uId, err := getUserId(r)
// 	if err != nil {
// 		resp.Json(w, err.Error(), http.StatusUnauthorized)
// 		return
// 	}

// 	user, err := handler.Service.View(uId)
// 	if err != nil {
// 		resp.Json(w, UserNotFound, http.StatusNotFound)
// 		return
// 	}
// 	profile := &ProfileResponse{
// 		Name:  user.Name,
// 		Email: user.Email,
// 		Banned: user.Banned,
// 	}
// 	resp.Json(w, profile, http.StatusOK)
// }

// func (handler *Handler) UpdateById(w http.ResponseWriter, r *http.Request) {
// 	uId, err := getUserId(r)
// 	if err != nil {
// 		resp.Json(w, err.Error(), http.StatusUnauthorized)
// 		return
// 	}
// 	res, err := req.HandleBody[UpdateNameRequest](&w, r)
// 	if err != nil {
// 		resp.Json(w, UserNotFound, http.StatusNotFound)
// 		return
// 	}
// 	_, err = handler.Service.ChangeName(uId, res.Name)
// 	if err != nil {
// 		resp.Json(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	resp.Json(w, "Profile updated", http.StatusOK)
// }

// func (handler *Handler) DeleteById(w http.ResponseWriter, r *http.Request) {
// 	uId, err := getUserId(r)
// 	if err != nil {
// 		resp.Json(w, err.Error(), http.StatusUnauthorized)
// 		return
// 	}
// 	err = handler.Service.Delete(uId)
// 	if err != nil {
// 		resp.Json(w, DeleteFailed, http.StatusInternalServerError)
// 		return
// 	}
// 	resp.Json(w, "Profile deleted", http.StatusOK)
// }




// func (handler *Handler) UpdateByAdmin(w http.ResponseWriter, r *http.Request) {
// 	uId, err := getIdFromReq(w, r)
// 	if err != nil {
// 		resp.Json(w, err.Error(), http.StatusForbidden)
// 		return
// 	}
// 	res, err := req.HandleBody[UpdateNameRequest](&w, r)
// 	if err != nil {
// 		resp.Json(w, UserNotFound, http.StatusNotFound)
// 		return
// 	}
// 	user, err := handler.Service.ChangeName(uId, res.Name)
// 	if err != nil {
// 		resp.Json(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	profile := &ShortProfileResponseForAdmin{
// 		Id:        user.Id,
// 		Name:      user.Name,
// 	}
// 	resp.Json(w, profile, http.StatusOK)
// }

// func (handler *Handler) DeleteByAdmin(w http.ResponseWriter, r *http.Request) {
// 	uId, err := selfDeletionDetected(w, r)
// 	if err != nil {
// 		resp.Json(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}
// 	err = handler.Service.Delete(*uId)
// 	if err != nil {
// 		resp.Json(w, DeleteFailed, http.StatusInternalServerError)
// 		return
// 	}
// 	resp.Json(w, "Profile deleted", http.StatusOK)
// }



// func getIdFromReq(w http.ResponseWriter, r *http.Request) (int, error) {
// 	idStr := strings.TrimSpace(r.PathValue("id"))
// 	if idStr == "" {
// 		resp.Json(w, EmptyParams, http.StatusForbidden)
// 		return 0, errors.New(EmptyParams)
// 	}
// 	uId, err := strconv.Atoi(idStr)
// 	if err != nil {
// 		return 0, errors.New(WrongParam)
// 	}
// 	return uId, nil
// }