-- is_aborted カラムを is_skipped に戻す
ALTER TABLE rounds RENAME COLUMN is_aborted TO is_skipped;
