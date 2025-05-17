-- is_skipped カラムを is_aborted に名前変更
ALTER TABLE rounds RENAME COLUMN is_skipped TO is_aborted;
