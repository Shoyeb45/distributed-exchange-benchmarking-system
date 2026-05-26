package validatormiddleware

import (
	"context"
	"encoding/json"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	errormiddleware "github.com/Shoyeb45/fast-docs/api/middleware/error-middleware"
	"github.com/Shoyeb45/fast-docs/pkg/apierr"
	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func init() {
    validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
        for _, tag := range []string{"json", "query", "param", "header"} {
            if name := strings.SplitN(fld.Tag.Get(tag), ",", 2)[0]; name != "" && name != "-" {
                return name
            }
        }
        return fld.Name
    })
}

type validateStoreKey struct{}

type validateStore map[reflect.Type]any

func storeFromCtx(ctx context.Context) validateStore {
    if s, ok := ctx.Value(validateStoreKey{}).(validateStore); ok {
        return s
    }
    return validateStore{}
}

// From pulls a validated value from context — call this in your handler.
func From[T any](r *http.Request) T {
    store := storeFromCtx(r.Context())

    if v, ok := store[reflect.TypeOf((*T)(nil)).Elem()]; ok {
        return v.(T)
    }
    var zero T
    return zero
}

// binding is a single extract+validate unit for one type+source pair
type binding struct {
    extract func(r *http.Request, store validateStore) error
}

// FromBody returns a binding descriptor for JSON body → T
func FromBody[T any]() binding {
    return binding{
        extract: func(r *http.Request, store validateStore) error {
            var dst T
            if err := json.NewDecoder(r.Body).Decode(&dst); err != nil {
                return apierr.NewBadRequest("malformed JSON body")
            }
            if err := validate.Struct(dst); err != nil {
                return apierr.NewValidationError("validation failed", toFieldErrors(err)...)
            }
            store[reflect.TypeOf(dst)] = dst
            return nil
        },
    }
}

// FromQuery returns a binding descriptor for query params → T
func FromQuery[T any]() binding {
    return binding{
        extract: func(r *http.Request, store validateStore) error {
            var dst T
            if err := mapToStruct(r.URL.Query(), "query", &dst); err != nil {
                return err
            }
            if err := validate.Struct(dst); err != nil {
                return apierr.NewValidationError("validation failed", toFieldErrors(err)...)
            }
            store[reflect.TypeOf(dst)] = dst
            return nil
        },
    }
}

// FromParams returns a binding descriptor for chi URL params → T
func FromParams[T any]() binding {
    return binding{
        extract: func(r *http.Request, store validateStore) error {
            var dst T
            chiCtx := chi.RouteContext(r.Context())
            params := map[string][]string{}
            if chiCtx != nil {
                for i, key := range chiCtx.URLParams.Keys {
                    params[key] = []string{chiCtx.URLParams.Values[i]}
                }
            }
            if err := mapToStruct(params, "param", &dst); err != nil {
                return err
            }
            if err := validate.Struct(dst); err != nil {
                return apierr.NewValidationError("validation failed", toFieldErrors(err)...)
            }
            store[reflect.TypeOf(dst)] = dst
            return nil
        },
    }
}

// FromHeaders returns a binding descriptor for request headers → T
func FromHeaders[T any]() binding {
    return binding{
        extract: func(r *http.Request, store validateStore) error {
            var dst T
            if err := mapToStruct(map[string][]string(r.Header), "header", &dst); err != nil {
                return err
            }
            if err := validate.Struct(dst); err != nil {
                return apierr.NewValidationError("validation failed", toFieldErrors(err)...)
            }
            store[reflect.TypeOf(dst)] = dst
            return nil
        },
    }
}

// Bind accepts any number of binding descriptors and returns a chi middleware.
// Validation errors short-circuit and flow into ErrorHandler.
func Bind(bindings ...binding) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            store := validateStore{}

            for _, b := range bindings {
                if err := b.extract(r, store); err != nil {
                    // Reuse ErrorHandler's writer path directly
                    errormiddleware.HandleError(w, r, err)
                    return
                }
            }

            ctx := context.WithValue(r.Context(), validateStoreKey{}, store)
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}

// --- Helpers ---

func mapToStruct[T any](data map[string][]string, tagName string, dst *T) error {
    rv := reflect.ValueOf(dst).Elem()
    rt := rv.Type()

    for i := range rt.NumField() {
        field := rt.Field(i)
        tag := strings.SplitN(field.Tag.Get(tagName), ",", 2)[0]
        if tag == "" || tag == "-" {
            continue
        }
        values, ok := data[tag]
        if !ok || len(values) == 0 {
            continue
        }
        if err := setField(rv.Field(i), field, values); err != nil {
            return apierr.NewBadRequest("invalid value for field: " + tag)
        }
    }
    return nil
}

func setField(fv reflect.Value, field reflect.StructField, values []string) error {
    val := values[0]
    switch fv.Kind() {
    case reflect.String:
        fv.SetString(val)
    case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
        n, err := strconv.ParseInt(val, 10, 64)
        if err != nil {
            return err
        }
        fv.SetInt(n)
    case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
        n, err := strconv.ParseUint(val, 10, 64)
        if err != nil {
            return err
        }
        fv.SetUint(n)
    case reflect.Float32, reflect.Float64:
        n, err := strconv.ParseFloat(val, 64)
        if err != nil {
            return err
        }
        fv.SetFloat(n)
    case reflect.Bool:
        b, err := strconv.ParseBool(val)
        if err != nil {
            return err
        }
        fv.SetBool(b)
    case reflect.Slice:
        if field.Type.Elem().Kind() == reflect.String {
            fv.Set(reflect.ValueOf(values))
        }
    }
    return nil
}

func toFieldErrors(err error) []apierr.FieldError {
    var out []apierr.FieldError
    for _, e := range err.(validator.ValidationErrors) {
        out = append(out, apierr.FieldError{
            Field:   e.Field(),
            Message: humanMessage(e),
        })
    }
    return out
}

func humanMessage(e validator.FieldError) string {
    switch e.Tag() {
    case "required":
        return "this field is required"
    case "email":
        return "must be a valid email address"
    case "min":
        return "must be at least " + e.Param() + " characters"
    case "max":
        return "must be at most " + e.Param() + " characters"
    case "gte":
        return "must be greater than or equal to " + e.Param()
    case "lte":
        return "must be less than or equal to " + e.Param()
    case "oneof":
        return "must be one of: " + e.Param()
    case "url":
        return "must be a valid URL"
    case "uuid4":
        return "must be a valid UUID"
    default:
        return "failed validation: " + e.Tag()
    }
}