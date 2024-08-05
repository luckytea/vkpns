package vkpns

import (
	"crypto/tls"
	"errors"
	"net"
	"net/http"
	"time"
)

const defaultMessagingEndpoint = "https://vkpns.rustore.ru/v1/projects/%s/messages:send"

type Client struct {
	client  *http.Client
	options ClientOptions
}

type ClientOptions struct {
	ProjectID     string
	ServiceToken  string
	VKPNSEndpoint string
}

// Push - модель для запроса на отправку сообщения.
type Push struct {
	Message Message `json:"message"`
}

// Структура push-уведомления.
type Message struct {
	Token        string            `json:"token"`        // Push-токен пользователя, полученный в приложении.
	Data         map[string]string `json:"data"`         // Объект, содержащий пары "key": value.
	Notification Notification      `json:"notification"` // Базовый шаблон уведомления для использования на всех платформах.
	Android      Android           `json:"android"`      // Специальные параметры Android для сообщений.
}

type Notification struct {
	Title string `json:"title"` // Название уведомления.
	Body  string `json:"body"`  // Основной текст уведомления.
	Image string `json:"image"` // Содержит URL-адрес изображения, которое будет отображаться в уведомлении.
}

type Android struct {
	TTL          string              `json:"ttl"`          // Как долго (в секундах) сообщение должно храниться в хранилище. Пример: 3.5s.
	Notification AndroidNotification `json:"notification"` // Уведомление для отправки на устройства Android.
}

type AndroidNotification struct {
	Title           string `json:"title"`              // Название уведомления.
	Body            string `json:"body"`               // Основной текст уведомления.
	Icon            string `json:"icon"`               // Значок уведомления..
	Color           string `json:"color"`              // Цвет значка уведомления в формате #rrggbb.
	Image           string `json:"image" `             // Содержит URL-адрес изображения, которое будет отображаться в уведомлении.
	ChannelID       string `json:"channel_id"`         // Идентификатор канала уведомления.
	ClickAction     string `json:"click_action" `      // Действие, связанное с кликом пользователя по уведомлению.
	ClickActionType int    `json:"click_action_type" ` // Необязательное поле, тип click_action (значение по умолчанию 0 - click_action будет использоваться как intent action, 1 - click_action будет использоваться как deep link).
}

// Response - ответ сервиса.
type Response struct {
	Code    int    `json:"code"`    // Числовой код ошибки.
	Message string `json:"message"` // Детальное описание ошибки.
	Status  string `json:"status" ` // Код ошибки в текстовом формате.
}

const (
	ErrInvalid   = "INVALID_ARGUMENT"  // неправильно указаны параметры запроса при отправке сообщения.
	ErrInternal  = "INTERNAL"          // внутренняя ошибка сервиса.
	ErrRatelimit = "TOO_MANY_REQUESTS" // превышено количество попыток отправить сообщение.
	ErrDenied    = "PERMISSION_DENIED" // неправильно указан сервисный ключ.
	ErrNotFound  = "NOT_FOUND"         // неправильно указан push-токен пользователя.
)

var (
	ErrNoData         = errors.New("no data to access VKPNS")
	ErrNotImplemented = errors.New("method is not implemented")
	ErrNoMessage      = errors.New("нет сообщения для отправки")
)

var (
	// HTTPClientTimeout specifies a time limit for requests made by the
	// HTTPClient. The timeout includes connection time, any redirects,
	// and reading the response body.
	HTTPClientTimeout = 60 * time.Second

	// ReadIdleTimeout is the timeout after which a health check using a ping
	// frame will be carried out if no frame is received on the connection. If
	// zero, no health check is performed.
	ReadIdleTimeout = 15 * time.Second

	// TCPKeepAlive specifies the keep-alive period for an active network
	// connection. If zero, keep-alive probes are sent with a default value
	// (currently 15 seconds).
	TCPKeepAlive = 15 * time.Second

	// TLSDialTimeout is the maximum amount of time a dial will wait for a connect
	// to complete.
	TLSDialTimeout = 20 * time.Second
)

// DialTLS is the default dial function for creating TLS connections for
// non-proxied HTTPS requests.
var DialTLS = func(network, addr string, cfg *tls.Config) (net.Conn, error) {
	dialer := &net.Dialer{
		Timeout:   TLSDialTimeout,
		KeepAlive: TCPKeepAlive,
	}
	return tls.DialWithDialer(dialer, network, addr, cfg)
}
