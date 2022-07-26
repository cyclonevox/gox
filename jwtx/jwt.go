package jwtx

import (
	`encoding/json`
	`reflect`
	`strings`
	`time`

	`github.com/dgrijalva/jwt-go`
)

// JWTConfig JWT配置
type JWTConfig struct {
	// 签名算法
	Method string `default:"HS256" yaml:"method" validate:"required,oneof=HS256 HS384 HS512"`
	// 密钥
	Key []byte `yaml:"key" validate:"required"`
	// 统一前缀
	Scheme string `yaml:"scheme"`
	// 有效期，单位分组
	Expiration int `default:"720" yaml:"expiration" validate:"required"`
}

// 额外的Token验证方法
// 例如只允许一个Token有效
type extraKeyFunc func(token *jwt.Token) bool

type jwtTool struct {
	// 签名算法
	method string
	// 签名密钥
	key []byte
	// Token前缀
	scheme string
	// 有效期，单位分钟
	expiration int
	// 需要是一个结构体指针
	payload interface{}
	// Token校验方法
	keyFunc jwt.Keyfunc
}

func NewJWT(config JWTConfig, payload interface{}, extra ...extraKeyFunc) *jwtTool {
	return &jwtTool{
		method:     config.Method,
		key:        config.Key,
		scheme:     config.Scheme,
		expiration: config.Expiration,
		payload:    payload,
		keyFunc: func(token *jwt.Token) (interface{}, error) {
			if token.Method.Alg() != config.Method {
				return nil, jwt.ErrInvalidKey
			}

			var (
				ok       bool
				standard *jwt.StandardClaims
			)

			if standard, ok = token.Claims.(*jwt.StandardClaims); !ok {
				return nil, jwt.ErrInvalidKey
			}

			if err := standard.Valid(); err != nil {
				return nil, err
			}

			if len(extra) != 0 {
				if !extra[0](token) {
					return nil, jwt.ErrInvalidKey
				}
			}

			return config.Key, nil
		},
	}
}

func (j *jwtTool) Sign(payload interface{}) (string, error) {
	var (
		data []byte
		err  error
	)

	if data, err = json.Marshal(payload); err != nil {
		return "", err
	}

	now := time.Now().Unix()
	token := jwt.NewWithClaims(
		jwt.GetSigningMethod(j.method),
		&jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(j.expiration) * time.Minute).Unix(),
			IssuedAt:  now,
			NotBefore: now,
			Subject:   string(data),
		},
	)

	var tokenString string
	if tokenString, err = token.SignedString(j.key); err != nil {
		return "", err
	}

	if j.scheme != "" {
		tokenString = j.scheme + " " + tokenString
	}

	return tokenString, nil
}

func (j *jwtTool) Parse(tokenString string) (*jwt.StandardClaims, map[string]interface{}, error) {
	var (
		err   error
		token *jwt.Token
	)

	if j.scheme != "" && strings.HasPrefix(tokenString, j.scheme) {
		tokenString = tokenString[len(j.scheme)+1:]
	}

	standard := new(jwt.StandardClaims)
	if token, err = jwt.ParseWithClaims(tokenString, standard, j.keyFunc); err != nil {
		return nil, nil, err
	}

	return token.Claims.(*jwt.StandardClaims), token.Header, nil
}

func (j *jwtTool) Payload(tokenString string) (interface{}, error) {
	var err error

	if j.scheme != "" && strings.HasPrefix(tokenString, j.scheme) {
		tokenString = tokenString[len(j.scheme)+1:]
	}

	standard := new(jwt.StandardClaims)
	if _, err = jwt.ParseWithClaims(tokenString, standard, j.keyFunc); err != nil {
		return nil, err
	}

	payload := reflect.New(reflect.TypeOf(j.payload).Elem()).Interface()
	if err = json.Unmarshal([]byte(standard.Subject), payload); err != nil {
		return nil, err
	}

	return payload, nil
}
