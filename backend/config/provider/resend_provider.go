package provider

import (
	"github.com/resend/resend-go/v2"
)

func NewResendProvider(env *EnvProvider) *resend.Client {
	return resend.NewClient(env.ResendApiKey())
}
