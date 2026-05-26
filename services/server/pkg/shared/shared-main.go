// This file contains the util type/functions/struct that can be used across the codebase
package shared

import "net/http"

type HandlerFunc func(w http.ResponseWriter, r *http.Request) error