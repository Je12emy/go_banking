# Securing the Application

For authentication we wil be creating a separate server with the responsibility to generate and validate authentication tokens, this is the basic flow.

![[Auth server flow.png]]

## JWT or JSON Web Token

JWT tokens have the following structure when encoded

![[JWT Structure.png]]

The data located in a token should NOT be sensitive since it can be viewed and modified, which of course will cause it to be rejected by our auth server. 

We will be using [jwt-go](https://github.com/dgrijalva/jwt-go) to create our JWT.

Based on our current endpoints there are user which should be able to use limited resources, like a admin role and a user role. this is something we need to take care of.

## Auth Server: Login API

We will be using yet again the Hexagonal architecture. The code has already been made for us, so let's breakdown the login handler.

This is the DTO for the login request.

```go
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
```

Which is then used used to decode de incoming request `json` and then pass it into the `authService`, if no error is returned the jwt is returned in the response.

```go
func (h AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var loginRequest dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		log.Println("Error while decoding login request: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
	} else {
		token, err := h.service.Login(loginRequest)
		if err != nil {
			writeResponse(w, http.StatusUnauthorized, err.Error())
		} else {
			writeResponse(w, http.StatusOK, *token)
		}
	}
}
```

In the `authService` login function the domain's `FindById`.

```go
func (s DefaultAuthService) Login(req dto.LoginRequest) (*string, error) {
	login, err := s.repo.FindBy(req.Username, req.Password)
	if err != nil {
		return nil, err
	}
	token, err := login.GenerateToken()
	if err != nil {
		return nil, err
	}
	return token, nil
}
```

`FindBy`  returns the needed data if both the username and password match, the `group_concat` built-in function returns the account's accounts in a single column with a specified separator and finally these results are grouped by the the customer's id

If the user has not accounts (like a admin user would) these are returned as `null` thanks to the `LEFT` join.

```go
func (d AuthRepositoryDb) FindBy(username, password string) (*Login, error) {
	var login Login
	sqlVerify := `SELECT username, u.customer_id, role, group_concat(a.account_id) as account_numbers FROM users u
                  LEFT JOIN accounts a ON a.customer_id = u.customer_id
                WHERE username = ? and password = ?
                GROUP BY a.customer_id, u.username`
	err := d.client.Get(&login, sqlVerify, username, password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("invalid credentials")
		} else {
			log.Println("Error while verifying login request from database: " + err.Error())
			return nil, errors.New("Unexpected database error")
		}
	}
	return &login, nil
}
```

If data is returned the domain's Login object's `GenerateToken` function is called, here the `nullable` fields are checked where; if there are `null` fields the token for the admin account is generated, else the user's token is generated.

```go
func (l Login) GenerateToken() (*string, error) {
	var claims jwt.MapClaims
	if l.Accounts.Valid && l.CustomerId.Valid {
		claims = l.claimsForUser()
	} else {
		claims = l.claimsForAdmin()
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedTokenAsString, err := token.SignedString([]byte(HMAC_SAMPLE_SECRET))
	if err != nil {
		log.Println("Failed while signing token: " + err.Error())
		return nil, errors.New("cannot generate token")
	}
	return &signedTokenAsString, nil
}
```

These "claims"  build the JWT body.

```go
func (l Login) claimsForUser() jwt.MapClaims {
	accounts := strings.Split(l.Accounts.String, ",")
	return jwt.MapClaims{
		"customer_id": l.CustomerId.String,
		"role":        l.Role,
		"username":    l.Username,
		"accounts":    accounts,
		"exp":         time.Now().Add(TOKEN_DURATION).Unix(),
	}
}

func (l Login) claimsForAdmin() jwt.MapClaims {
	return jwt.MapClaims{
		"role":     l.Role,
		"username": l.Username,
		"exp":      time.Now().Add(TOKEN_DURATION).Unix(),
	}
}
```

After these are generated, the JWT is signed with a method and a signature constant is used as the secret which returns the whole JWT.

```go
func (l Login) GenerateToken() (*string, error) {
	// ...CODE
	
	signedTokenAsString, err := token.SignedString([]byte(HMAC_SAMPLE_SECRET))
	if err != nil {
		log.Println("Failed while signing token: " + err.Error())
		return nil, errors.New("cannot generate token")
	}
	return &signedTokenAsString, nil
}
```

This whole token is then returned to the handler, which returns it as a response.

```json
{
  "accounts": [
    "95472",
    "95473",
    "95474"
  ],
  "customer_id": "2001",
  "exp": 1610867384,
  "role": "user",
  "username": "2001"
}
```

## Auth Server: Verify API

Validation for routes will be made through a middleware, middleware is a chain of function to be executed when a route request is received.

![[Middleware.png]]

This middleware will allow us to check the validity of each JWT and stop the request's flow if needed.

### Creating Middleware

This is how we can create middleware inside our router.

```go
	// middleware
	router.Use(func(next http.Handler) http.Handler {
		// next is the handler to the next middleware
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// this is the middleware
			// before

			next.ServeHTTP(w, r)
			// next middleware
		})
	})
```

Let's create this auth middleware inside our app.

```go
func (a AuthMiddleware) authorizationHandler() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get the current route, this route is named in the router config
			currentRoute := mux.CurrentRoute(r)
			// Get all the vars
			currentRouteVars := mux.Vars(r)
			// Get the auth bearer auth header
			authHeader := r.Header.Get("Authorization")

			if authHeader != "" {
				token := getTokenFromHeader(authHeader)

				isAuthorized := a.repo.IsAuthorized(token, currentRoute.GetName(), currentRouteVars)

				if isAuthorized {
					// pass the writer and request to the next middleware or router
					next.ServeHTTP(w, r)
				} else {
					appError := errs.AppError{http.StatusForbidden, "Unauthorized"}
					writeResponse(w, appError.Code, appError.AsMessage())
				}
			} else {
				writeResponse(w, http.StatusUnauthorized, "missing token")
			}
		})
	}
}

func getTokenFromHeader(header string) string {
	/*
	   token is coming in the format as below
	   "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2NvdW50cyI6W.yI5NTQ3MCIsIjk1NDcyIiw"
	*/
	splitToken := strings.Split(header, "Bearer")
	if len(splitToken) == 2 {
		// extract the token itself
		return strings.TrimSpace(splitToken[1])
	}
	return ""
}
```

This middleware uses a `authRepository` dependency, which will prepare the request URL for the auth server and send it.

```go
func (r RemoteAuthRepository) IsAuthorized(token string, routeName string, vars map[string]string) bool {

	u := buildVerifyURL(token, routeName, vars)

	if response, err := http.Get(u); err != nil {
		fmt.Println("Error while sending..." + err.Error())
		return false
	} else {
		// Create a map with string as a key and bool as a value: ["isAuthorized": true]
		m := map[string]bool{}
		if err = json.NewDecoder(response.Body).Decode(&m); err != nil {
			logger.Error("Error while decoding response from auth server:" + err.Error())
			return false
		}
		return m["isAuthorized"]
	}
}

/*
  This will generate a url for token verification in the below format
  /auth/verify?token={token string}
              &routeName={current route name}
              &customer_id={customer id from the current route}
              &account_id={account id from current route if available}
  Sample: /auth/verify?token=aaaa.bbbb.cccc&routeName=MakeTransaction&customer_id=2000&account_id=95470
*/
func buildVerifyURL(token string, routeName string, vars map[string]string) string {
	// hardcoded auth server url
	u := url.URL{Host: "localhost:8181", Path: "/auth/verify", Scheme: "http"}
	q := u.Query()
	// add the token and route name params
	q.Add("token", token)
	q.Add("routeName", routeName)
	// for each loop in the vars map
	for k, v := range vars {
		// k is the key and v is the value
		q.Add(k, v)
	}
	u.RawQuery = q.Encode()
	return u.String()
}
```

In our router we are able to name our routes, this will allows us to restrict their access based on a matrix which specifies which end-users have access to each resource.

```go
	router.HandleFunc("/customers", ch.getAllCustomers).Methods(http.MethodGet).Name("GetAllCustomers")
	router.HandleFunc("/customers/{customer_id:[0-9]+}", ch.getCustomer).Methods(http.MethodGet).Name("GetCustomer")
	router.HandleFunc("/customers/{customer_id:[0-9]+}/account", ah.newAccount).Methods(http.MethodPost).Name("NewAccount")
	router.HandleFunc("/transaction", th.newTransaction).Methods(http.MethodPost).Name("NewTransaction")
```

### The Auth Server

In our `authHandler` we have a `Verify` function which takes care of token and role validation.

```go
/*
  Sample URL string
 http://localhost:8181/auth/verify?token=somevalidtokenstring&routeName=GetCustomer&customer_id=2000&account_id=95470
*/
func (h AuthHandler) Verify(w http.ResponseWriter, r *http.Request) {
	urlParams := make(map[string]string)

	// converting from Query to map type
	for k := range r.URL.Query() {
		urlParams[k] = r.URL.Query().Get(k)
	}
	// from the map, if the token key is not empty
	if urlParams["token"] != "" {

		isAuthorized, appError := h.service.Verify(urlParams)
		if appError != nil {
			writeResponse(w, http.StatusForbidden, notAuthorizedResponse())
		} else {
			if isAuthorized {
				writeResponse(w, http.StatusOK, authorizedResponse())
			} else {
				writeResponse(w, http.StatusForbidden, notAuthorizedResponse())
			}
		}
	} else {
		writeResponse(w, http.StatusForbidden, "missing token")
	}
}
```

It used the `Verify` function from the `authService` which verifies the token validity.

```go
/*
  Sample URL string
 http://localhost:8181/auth/verify?token=somevalidtokenstring&routeName=GetCustomer&customer_id=2000&account_id=95470
*/
func (h AuthHandler) Verify(w http.ResponseWriter, r *http.Request) {
	urlParams := make(map[string]string)

	// converting from Query to map type
	for k := range r.URL.Query() {
		urlParams[k] = r.URL.Query().Get(k)
	}
	// from the map, if the token key is not empty
	if urlParams["token"] != "" {

		isAuthorized, appError := h.service.Verify(urlParams)
		if appError != nil {
			writeResponse(w, http.StatusForbidden, notAuthorizedResponse())
		} else {
			if isAuthorized {
				writeResponse(w, http.StatusOK, authorizedResponse())
			} else {
				writeResponse(w, http.StatusForbidden, notAuthorizedResponse())
			}
		}
	} else {
		writeResponse(w, http.StatusForbidden, "missing token")
	}
}

func jwtTokenFromString(tokenString string) (*jwt.Token, error) {
	// decode the token with the signature key
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(domain.HMAC_SAMPLE_SECRET), nil
	})
	if err != nil {
		log.Println("Error while parsing token: " + err.Error())
		return nil, err
	}
	return token, nil
}
```

if the token contains the user role, the passed `customer_id` and `account_id` are checked to match

```go
func (c Claims) IsRequestVerifiedWithTokenClaims(urlParams map[string]string) bool {
	// check if the customer_id and the consulted account in the token and request are the same
	if c.CustomerId != urlParams["customer_id"] {
		return false
	}

	if !c.IsValidAccountId(urlParams["account_id"]) {
		return false
	}
	return true
}
```

if the token is valid, the roles are checked with the role domain object.

```go
type RolePermissions struct {
	rolePermissions map[string][]string
}

func (p RolePermissions) IsAuthorizedFor(role string, routeName string) bool {
	perms := p.rolePermissions[role]
	// Loop through all the allowed routes to match which the current route
	for _, r := range perms {
		if r == strings.TrimSpace(routeName) {
			return true
		}
	}
	return false
}

func GetRolePermissions() RolePermissions {
	return RolePermissions{map[string][]string{
		"admin": {"GetAllCustomers", "GetCustomer", "NewAccount", "NewTransaction"},
		"user":  {"GetCustomer", "NewTransaction"},
	}}
}
```