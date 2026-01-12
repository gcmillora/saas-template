package provider

import (
	"log/slog"

	"github.com/supabase-community/supabase-go"
)

type SupabaseProvider struct {
	client *supabase.Client
}

func NewSupabaseProvider(env *EnvProvider) *SupabaseProvider {
	client, err := supabase.NewClient(env.SupabaseURL(), env.SupabaseKey(), &supabase.ClientOptions{})

	if err != nil {
		slog.Default().Error("Unable to connect to Supabase", "error", err)
		return nil
	}

	return &SupabaseProvider{client}
}

func (p *SupabaseProvider) Client() *supabase.Client {
	return p.client
}