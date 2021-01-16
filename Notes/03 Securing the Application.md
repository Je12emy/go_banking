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