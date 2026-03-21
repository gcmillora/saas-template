package provider

import (
	storage_go "github.com/supabase-community/storage-go"
)

const StorageBucket = "assets"

func NewSupabaseStorageClient(env *EnvProvider) *storage_go.Client {
	return storage_go.NewClient(
		env.SupabaseUrl()+"/storage/v1",
		env.SupabaseKey(),
		nil,
	)
}
