package response

type VersionOneBaseResponse struct {
	Code    int         `json:"code"`
	Message interface{} `json:"message"`
}

func NewVersionOneBaseResponse(code int, msg interface{}) *VersionOneBaseResponse {
	v1Response := new(VersionOneBaseResponse)
	v1Response.Code = code
	v1Response.Message = msg

	return v1Response
}

func (r *VersionOneBaseResponse) Response() interface{} {
	return r
}
