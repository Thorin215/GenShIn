package sql

func Migrate() error {
	err := MigrateMetadata(DB)
	if err != nil {
		return err
	}

	err = MigrateUser(DB)
	if err != nil {
		return err
	}

	return nil
}
