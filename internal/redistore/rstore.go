package redistore

import (
	"bytes"
	"context"
	"encoding/base32"
	"encoding/binary"
	"encoding/gob"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/redis/go-redis/v9"
)

const (
	EXIST_TIME = 2 * 60 * 60
	KEY_PREFIX = "session:"
)

type keygenFunc func() string

type SessionSerializer interface {
	Serialize(s *sessions.Session) ([]byte, error)
	Deserialize(b []byte, s *sessions.Session) error
}

type GobSerializer struct{}

func (gs GobSerializer) Serialize(s *sessions.Session) ([]byte, error) {
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(s.Values)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (gs GobSerializer) Deserialize(b []byte, s *sessions.Session) error {
	dec := gob.NewDecoder(bytes.NewBuffer(b))
	return dec.Decode(&s.Values)
}

func randomKeyGen() string {
	key := make([]byte, 24)
	// for order -- timestamp
	ts := uint64(time.Now().UnixNano())
	binary.BigEndian.PutUint64(key[:8], ts)
	uid := uuid.New()
	copy(key[8:], uid[:])
	enc := base32.StdEncoding.WithPadding(base32.NoPadding)
	return enc.EncodeToString(key)
}

type Redistore struct {
	client     *redis.Client
	options    sessions.Options
	keyPrefix  string
	keyGen     keygenFunc
	serializer SessionSerializer
}

func (rs *Redistore) SetClient(c *redis.Client) {
	rs.client = c
}

func (rs *Redistore) SetOptions(opts sessions.Options) {
	rs.options = opts
}

func (rs *Redistore) SetKeyPrefix(keyPrefix string) {
	rs.keyPrefix = keyPrefix
}

func (rs *Redistore) SetKeyGen(keyGen keygenFunc) {
	rs.keyGen = keyGen
}

func (rs *Redistore) SetSerializer(s SessionSerializer) {
	rs.serializer = s
}

func NewRedistore(ctx context.Context, client *redis.Client) *Redistore {
	opts := sessions.Options{
		Path:   "/",
		MaxAge: EXIST_TIME,
	}
	rs := &Redistore{
		client:     client,
		options:    opts,
		keyPrefix:  KEY_PREFIX,
		keyGen:     randomKeyGen,
		serializer: GobSerializer{},
	}
	if rs.client.Ping(ctx).Err() != nil {
		log.Fatalf("[init] -- (internal/redisStore) status: fail")
	}
	return rs
}

// 从redis中取出 session,反序列化
func (rs *Redistore) load(ctx context.Context, s *sessions.Session) error {
	cmdStr := rs.client.Get(ctx, rs.keyPrefix+s.ID)
	if cmdStr.Err() != nil {
		return cmdStr.Err()
	}
	b, err := cmdStr.Bytes()
	if err != nil {
		return err
	}
	return rs.serializer.Deserialize(b, s)
}

// 尝试从 store 获取 session
func (rs *Redistore) Get(r *http.Request, name string) (*sessions.Session, error) {
	return sessions.GetRegistry(r).Get(rs, name)
}

// 不存在 session 则创建新 session
// 已存在 session 取出 session
func (rs *Redistore) New(r *http.Request, name string) (*sessions.Session, error) {
	session := sessions.NewSession(rs, name)
	opts := rs.options
	session.Options = &opts
	session.IsNew = true

	// get sessionId from cookie(request)
	cookie, err := r.Cookie(name)
	if err != nil {
		return session, nil
	}
	session.ID = cookie.Value

	// 使用 cookie 中的 session id 尝试校验 redis 中的是否存在
	err = rs.load(r.Context(), session)
	if err == redis.Nil {
		// session key 不存在或者过期，需要重新执行Save()
		err = nil
	} else {
		// 出现其他错误或者不出错，session key 存在
		session.IsNew = false
	}
	return session, err
}

func (rs *Redistore) delete(ctx context.Context, s *sessions.Session) error {
	return rs.client.Del(ctx, rs.keyPrefix+s.ID).Err()
}

func (rs *Redistore) save(ctx context.Context, s *sessions.Session) error {
	b, err := rs.serializer.Serialize(s)
	if err != nil {
		return err
	}
	return rs.client.Set(ctx, rs.keyPrefix+s.ID, b, time.Duration(s.Options.MaxAge)*time.Second).Err()
}

func (rs *Redistore) Save(r *http.Request, w http.ResponseWriter, s *sessions.Session) error {
	if s.Options.MaxAge <= 0 {
		if err := rs.delete(r.Context(), s); err != nil {
			return err
		}
		http.SetCookie(w, sessions.NewCookie(s.Name(), "", s.Options))
		return nil
	}

	if s.ID == "" {
		s.ID = rs.keyGen()
	}
	if err := rs.save(r.Context(), s); err != nil {
		return err
	}
	http.SetCookie(w, sessions.NewCookie(s.Name(), s.ID, s.Options))
	return nil
}
