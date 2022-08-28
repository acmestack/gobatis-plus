package constants

const SELECT_SQL = "SELECT #{columns} FROM #{tableName} WHERE #{conditions}"

const SAVE_SQL = "insert into #{tableName} (#{columns}) values (#{columnsValue})"
const SAVE_BATCH_SQL = "insert into #{tableName} (#{columns}) values (#{columnsValue}),(#{columnsValue}),(#{columnsValue})"

const DELETEBYID_SQL = "delete from #{tableName} where #{conditions}"
const DELETEBATCHIDS_SQL = "delete from #{tableName} where in #{conditions}"

const UPDATEBYID_SQL = "update #{tableName} set #{filed} where #{conditions}"
