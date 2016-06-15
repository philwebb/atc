package migrations

import "github.com/BurntSushi/migration"

func ReplaceSizeWithSizeInKB(tx migration.LimitedTx) error {
	_, err := tx.Exec(`
    ALTER TABLE volumes ADD COLUMN size_in_kilobytes bigint;
	`)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		UPDATE volumes SET size_in_kilobytes = size / 1000;
	`)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
    ALTER TABLE volumes DROP COLUMN size;
	`)

	return err
}
