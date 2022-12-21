package models

type BaseRequest struct {
	ID         uint              `mapper:"id"`
	Page       int               `mapper:"page"`
	Limit      int               `mapper:"limit"`
	BodyData   interface{}       `mapper:"data"`
	QueryParam map[string]string `mapper:"query_param"`
}
