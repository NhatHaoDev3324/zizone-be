package constant

type SuccessStatus int

const (
	StatusOK        SuccessStatus = 200
	StatusCreated   SuccessStatus = 201
	StatusAccepted  SuccessStatus = 202
	StatusNoContent SuccessStatus = 204
)

type FailStatus int

const (
	StatusBadRequest   FailStatus = 400
	StatusUnauthorized FailStatus = 401
	StatusForbidden    FailStatus = 403
	StatusNotFound     FailStatus = 404
	StatusConflict     FailStatus = 409
)
