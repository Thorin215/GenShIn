package sql

func Migrate() error {
	err := MigrateMetadata(DB)
	if err != nil {
		panic(err)
	}
	return nil
}
