package main_test

//func TestSignInHandler_Success(t *testing.T) {
//	payload := app.UserRequest{
//		Login:    "testuser",
//		Password: "testpassword",
//	}
//	payloadBytes, _ := json.Marshal(payload)
//	req, err := http.NewRequest("POST", "/sign_in", bytes.NewReader(payloadBytes))
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	rr := httptest.NewRecorder()
//	handler := http.HandlerFunc(app.SignIn)
//
//	handler.ServeHTTP(rr, req)
//
//	assert.Equal(t, http.StatusSeeOther, rr.Code)
//
//}
//
//func TestSignInHandler_Failure(t *testing.T) {
//	payload := app.UserRequest{
//		Login:    "nonexistentuser",
//		Password: "wrongpassword",
//	}
//	payloadBytes, _ := json.Marshal(payload)
//	req, err := http.NewRequest("POST", "/sign_in", bytes.NewReader(payloadBytes))
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	rr := httptest.NewRecorder()
//	handler := http.HandlerFunc(app.SignIn)
//
//	handler.ServeHTTP(rr, req)
//
//	assert.Equal(t, http.StatusInternalServerError, rr.Code)
//}
//
//func TestSignUpHandler_Success(t *testing.T) {
//	payload := app.UserRequest{
//		Login:    "newuser",
//		Password: "newpassword",
//	}
//	payloadBytes, _ := json.Marshal(payload)
//	req, err := http.NewRequest("POST", "/sign_up", bytes.NewReader(payloadBytes))
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	rr := httptest.NewRecorder()
//	handler := http.HandlerFunc(app.SignUp)
//
//	handler.ServeHTTP(rr, req)
//
//	assert.Equal(t, http.StatusOK, rr.Code)
//
//}
//
//func TestSignUpHandler_Failure(t *testing.T) {
//
//	payload := app.UserRequest{
//		Login:    "testuser",
//		Password: "testpassword",
//	}
//	payloadBytes, _ := json.Marshal(payload)
//	req, err := http.NewRequest("POST", "/sign_up", bytes.NewReader(payloadBytes))
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	rr := httptest.NewRecorder()
//	handler := http.HandlerFunc(app.SignUp)
//
//	handler.ServeHTTP(rr, req)
//
//	assert.Equal(t, http.StatusInternalServerError, rr.Code)
//
//}

//func TestMainPageHandler_Authenticated(t *testing.T) {
//	req, err := http.NewRequest("GET", "/", nil)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	session := &app.SessionMock{
//		Values: map[interface{}]interface{}{
//			"authenticated": true,
//		},
//	}
//	rr := httptest.NewRecorder()
//	req = req.WithContext(app.SessionContext(req.Context(), "session", session))
//
//	handler := http.HandlerFunc(app.MainPage)
//
//	handler.ServeHTTP(rr, req)
//
//	assert.Equal(t, http.StatusOK, rr.Code)
//
//}

//func TestMainPageHandler_Unauthenticated(t *testing.T) {
//	req, err := http.NewRequest("GET", "/", nil)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	session := &app.SessionMock{
//		Values: map[interface{}]interface{}{
//			"authenticated": false,
//		},
//	}
//	rr := httptest.NewRecorder()
//	req = req.WithContext(app.SessionContext(req.Context(), "session", session))
//
//	handler := http.HandlerFunc(app.MainPage)
//
//	handler.ServeHTTP(rr, req)
//
//	assert.Equal(t, http.StatusSeeOther, rr.Code)
//
//}
