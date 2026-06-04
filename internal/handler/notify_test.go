package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandleNotify(t *testing.T) {
	cases := []struct {
		name           string
		token          string
		body           string
		channelFails   bool
		expectedStatus int
		expectedError  string
	}{
		// ── Seguridad ──────────────────────────────────────────
		{
			name:           "sin header authorization",
			token:          "",
			body:           `{"source":"test","type":"mail","receivers":["t@test.com"],"subject":"Test"}`,
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "header Authorization requerido",
		},
		{
			name:           "token incorrecto",
			token:          "Bearer token-malo",
			body:           `{"source":"test","type":"mail","receivers":["t@test.com"],"subject":"Test"}`,
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "token inválido",
		},

		// ── Validación del request ──────────────────────────────
		{
			name:           "json malformado",
			token:          authHeader(),
			body:           `{esto no es json`,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "request inválido",
		},
		{
			name:           "type inválido",
			token:          authHeader(),
			body:           `{"source":"test","type":"fax"}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "mail sin receivers",
			token:          authHeader(),
			body:           `{"source":"test","type":"mail","subject":"Test","body":"hola"}`,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "mail requiere al menos un receiver",
		},
		{
			name:           "mail sin subject",
			token:          authHeader(),
			body:           `{"source":"test","type":"mail","receivers":["t@test.com"],"body":"hola"}`,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "mail requiere subject",
		},
		{
			name:           "mail sin body ni template",
			token:          authHeader(),
			body:           `{"source":"test","type":"mail","receivers":["t@test.com"],"subject":"Test"}`,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "mail requiere body o template_id",
		},
		{
			name:           "receiver con email inválido",
			token:          authHeader(),
			body:           `{"source":"test","type":"mail","receivers":["no-es-email"],"subject":"Test","body":"hola"}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "canal no registrado",
			token:          authHeader(),
			body:           `{"source":"test","type":"slack","text":"hola","channel":"#general"}`,
			expectedStatus: http.StatusBadRequest,
		},

		// ── Errores del canal ───────────────────────────────────
		{
			name:           "error al enviar",
			token:          authHeader(),
			body:           `{"source":"test","type":"mail","receivers":["t@test.com"],"subject":"Test","body":"hola"}`,
			channelFails:   true,
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "error al enviar notificación",
		},

		// ── Happy path ──────────────────────────────────────────
		{
			name:           "request válido",
			token:          authHeader(),
			body:           `{"source":"test","type":"mail","receivers":["t@test.com"],"subject":"Test","body":"hola"}`,
			expectedStatus: http.StatusOK,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {

			mock := &mockChannel{shouldFail: tc.channelFails}
			srv := newTestServer(mock)

			req := httptest.NewRequest(http.MethodPost, "/notify", bytes.NewBufferString(tc.body))
			req.Header.Set("Content-Type", "application/json")
			if tc.token != "" {
				req.Header.Set("Authorization", tc.token)
			}

			rec := httptest.NewRecorder()

			srv.Router().ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedStatus, rec.Code, "status code incorrecto")

			if tc.expectedError != "" {
				var resp map[string]interface{}
				require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
				assert.Contains(t, resp["error"], tc.expectedError,
					"mensaje de error incorrecto")
			}

			assert.NotEmpty(t, rec.Header().Get("X-Request-ID"),
				"X-Request-ID debería estar presente en toda respuesta")
		})
	}
}
