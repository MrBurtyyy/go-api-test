package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/MrBurtyyy/go-api-test/internal/validation"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"
	"strconv"
	"strings"
)

const (
	ContentTypeJson = "application/json"
)

var (
	httpProductAlreadyExists = ErrorResponse{StatusCode: http.StatusBadRequest, ErrorCode: "PRODUCT_EXISTS", ErrorMessage: "Product already exists"}

	v *validator.Validate
	p ProductService
)

func Init(r *chi.Mux, db *pgxpool.Pool) {
	v = validation.New()

	p = productService{db: db}

	r.Post("/product", addProduct)
	r.Get("/product/{id}", getProduct)
}

func addProduct(w http.ResponseWriter, r *http.Request) {
	var j ProductRequest
	if err := parseAndValidate(w, r, &j); err != nil {
		return
	}

	product, err := p.Add(r.Context(), &NewProduct{Name: j.Name, Price: j.Price})
	if err != nil {
		if errors.Is(err, ProductAlreadyExists) {
			jsonError(httpProductAlreadyExists, w)
			return
		}

		fmt.Printf("error in handler: %v\n", err)
		jsonError(httpServerError, w)
		return
	}

	success(http.StatusOK, ProductResponse{ID: product.ID, Name: product.Name, Price: product.Price}, w)
}

func getProduct(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		jsonError(httpMalformedRequest, w)
	}

	product, err := p.Get(r.Context(), int64(id))
	if err != nil {
		fmt.Printf("failed to fetch product: %v\n", err)
	}

	success(http.StatusOK, ProductResponse{ID: product.ID, Name: product.Name, Price: product.Price}, w)
}

func parseAndValidate(w http.ResponseWriter, r *http.Request, o any) error {
	if err := json.NewDecoder(r.Body).Decode(o); err != nil {
		var syntaxError *json.SyntaxError
		if errors.As(err, &syntaxError) {
			fmt.Println("encountered syntax error")
			jsonError(httpServerError, w)
			return err
		}
		fmt.Println("encountered malformed request error")
		jsonError(httpMalformedRequest, w)
		return err
	}

	if err := v.Struct(o); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			fmt.Println("encountered payload validation error")
			payloadValidationError(ve, w)
		} else {
			fmt.Println("encountered payload validation error")
			jsonError(httpPayloadValidationError, w)
		}
		return err
	}

	return nil
}

func success(statusCode int, response any, w http.ResponseWriter) {
	w.Header().Add("Content-Type", ContentTypeJson)
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		panic("failed to write JSON response")
	}
}

func payloadValidationError(ve validator.ValidationErrors, w http.ResponseWriter) {
	e := httpPayloadValidationError
	e.ValidationErrors = convertValidationErrors(ve)
	jsonError(e, w)
}

func convertValidationErrors(ve validator.ValidationErrors) []PropertyValidationErrorResponse {
	var e []PropertyValidationErrorResponse
	for _, v := range ve {
		e = append(e, toErrorResponse(v))
	}
	return e
}

func toErrorResponse(v validator.FieldError) PropertyValidationErrorResponse {
	s := v.Namespace()
	i := strings.IndexRune(v.Namespace(), '.')
	p := s[i+1:]
	return PropertyValidationErrorResponse{Property: p, Message: getTagErrorMessage(v.Tag(), v.Param())}
}

func getTagErrorMessage(tag string, param string) string {
	switch tag {
	case "dp":
		return fmt.Sprintf("must not be more than %v decimal places", param)
	case "required":
		return "is required"
	default:
		return "unknown jsonError"
	}
}

func jsonError(e ErrorResponse, w http.ResponseWriter) {
	w.Header().Add("Content-Type", ContentTypeJson)
	w.WriteHeader(e.StatusCode)
	if err := json.NewEncoder(w).Encode(e); err != nil {
		panic("failed to write JSON response")
	}
}
